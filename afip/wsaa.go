package afip

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/hooklift/gowsdl/soap"
	"github.com/sehogas/gowsaa/ws/wsaa"
)

type Wsaa struct {
	environment     Environment
	urlWsaa         string
	tickets         map[string]*LoginTicketResponse
	privateKeyFile  string
	certificateFile string
	cuit            int64
}

type LoginTicket struct {
	ServiceName    string
	Token          string
	Sign           string
	ExpirationTime time.Time
	Cuit           int64
}

// HeaderLoginTicket es la cabecera de la estructura de request y response
type HeaderLoginTicket struct {
	Source         string `xml:"source,omitempty"`
	Destination    string `xml:"destination,omitempty"`
	UniqueID       uint32 `xml:"uniqueId,omitempty"`
	GenerationTime string `xml:"generationTime,omitempty"`
	ExpirationTime string `xml:"expirationTime,omitempty"`
}

// LoginTicketRequest es la estructura general del request
type LoginTicketRequest struct {
	XMLName xml.Name           `xml:"loginTicketRequest"`
	Version string             `xml:"version,attr"`
	Header  *HeaderLoginTicket `xml:"header,omitempty"`
	Service string             `xml:"service,omitempty"`
}

// Credentials es la estructura que devuelve el response con la info principal
type Credentials struct {
	Token string `xml:"token,omitempty"`
	Sign  string `xml:"sign,omitempty"`
}

// LoginTicketResponse ...
type LoginTicketResponse struct {
	XMLName     xml.Name           `xml:"loginTicketResponse"`
	Header      *HeaderLoginTicket `xml:"header,omitempty"`
	Credentials *Credentials       `xml:"credentials,omitempty"`
}

func NewWsaa(environment Environment, privateKeyFile string, certificateFile string, cuit int64) (*Wsaa, error) {
	var url string
	if environment == PRODUCTION {
		url = URLWSAAProduction
	} else {
		url = URLWSAATesting
	}

	if _, err := os.Stat(privateKeyFile); err != nil {
		return nil, fmt.Errorf("privateKeyFile: %s", err)
	}
	if _, err := os.Stat(certificateFile); err != nil {
		return nil, fmt.Errorf("certificateFile: %s", err)
	}
	return &Wsaa{
		environment:     environment,
		urlWsaa:         url,
		tickets:         make(map[string]*LoginTicketResponse),
		privateKeyFile:  privateKeyFile,
		certificateFile: certificateFile,
		cuit:            cuit,
	}, nil
}

// GetLoginTicket devuelve el ticket de acceso afip correspondiente al servicio pasado por parámetro.
func (c *Wsaa) GetLoginTicket(serviceName string) (*LoginTicket, error) {
	var renovar bool = true

	ticket := c.tickets[serviceName]
	if ticket != nil {
		expTime, err := time.Parse(time.RFC3339, ticket.Header.ExpirationTime)
		if err != nil {
			return nil, fmt.Errorf("GetLoginTicket: Error leyendo fecha de expiración del ticket. %s", err)
		}
		renovar = time.Now().After(expTime)
	}

	if renovar {
		expiration := time.Now().Add(10 * time.Minute)
		generationTime := time.Now().Add(-10 * time.Minute).Format(time.RFC3339)
		expirationTime := expiration.Format(time.RFC3339)

		privateKey, err := readPrivateKey(c.privateKeyFile)
		if err != nil {
			return nil, fmt.Errorf("GetLoginTicket: Error leyendo clave privada: %s", err)
		}

		certificate, err := readCertificate(c.certificateFile)
		if err != nil {
			return nil, fmt.Errorf("GetLoginTicket: Error leyendo certificado: %s", err)
		}

		// Armo estructura request
		loginTicketRequest := LoginTicketRequest{
			Version: "1.0",
			Header: &HeaderLoginTicket{
				UniqueID:       1,
				GenerationTime: generationTime,
				ExpirationTime: expirationTime,
			},
			Service: serviceName,
		}

		loginTicketRequestXML, err := xml.MarshalIndent(loginTicketRequest, " ", "  ")
		if err != nil {
			return nil, fmt.Errorf("GetLoginTicket: Error armando request XML. %s", err)
		}
		content := []byte(string(loginTicketRequestXML))

		// Creo CMS (Cryptographic Message Syntax)
		cms, err := encodeCMS(content, certificate, privateKey.(*rsa.PrivateKey))
		if err != nil {
			return nil, fmt.Errorf("GetLoginTicket: Error creando CMS: %s", err)
		}

		// Convierto CMS a base64
		cmsBase64 := base64.StdEncoding.EncodeToString(cms)

		// Armo conexión SOAP y solicitud
		conexion := soap.NewClient(c.urlWsaa)
		service := wsaa.NewLoginCMS(conexion)

		request := wsaa.LoginCms{In0: cmsBase64}

		// Llamo al servicio de autenticación afip wssa
		responseXML, err := service.LoginCms(&request)
		if err != nil {
			response := soap.SOAPEnvelopeResponse{
				Body: soap.SOAPBodyResponse{
					Content: &soap.SOAPFault{},
					Fault:   &soap.SOAPFault{},
				},
			}
			if err := xml.Unmarshal([]byte(err.Error()[strings.Index(err.Error(), "<soapenv:Envelope"):]), &response); err != nil {
				return nil, fmt.Errorf("GetLoginTicket: Error desarmando respuesta XML. %s", err)
			}
			return nil, fmt.Errorf("GetLoginTicket: Error del servicio: %s", response.Body.Fault.String)
		}

		// Desarmo respuesta XML
		response := LoginTicketResponse{}
		if err := xml.Unmarshal([]byte(responseXML.LoginCmsReturn), &response); err != nil {
			return nil, fmt.Errorf("GetLoginTicket: Error desarmando respuesta XML. %s", err)
		}

		// Almaceno ticket de respuesta (porque no se puede llamar nuevamente al servicio hasta dentro de 10 minutos,
		// hay que seguir usando el ticket actual. El vencimiento de los ticket de afip suele ser de 12 horas)
		c.tickets[serviceName] = &response
	}

	expirationTime, err := time.Parse(time.RFC3339, c.tickets[serviceName].Header.ExpirationTime)
	if err != nil {
		return nil, fmt.Errorf("GetLoginTicket: Error leyendo fecha de expiración del ticket. %s", err)
	}

	loginTicket := &LoginTicket{
		ServiceName:    serviceName,
		Token:          c.tickets[serviceName].Credentials.Token,
		Sign:           c.tickets[serviceName].Credentials.Sign,
		ExpirationTime: expirationTime,
		Cuit:           c.cuit,
	}

	if c.environment == TESTING {
		// grabo en local el archivo temporal con el ticket de acceso

		data := fmt.Sprintf("CUIT=%s\nTOKEN=%s\nSIGN=%s\nEXPIRATION=%s\n",
			loginTicket.Cuit,
			loginTicket.Token,
			loginTicket.Sign,
			loginTicket.ExpirationTime.Format(time.RFC3339))

		fileName := fmt.Sprintf("%s.TA", serviceName)
		f, err := os.Create(fileName)
		if err != nil {
			return nil, fmt.Errorf("GetLoginTicket: No se pudo crear el archivo %s: %s", fileName, err)
		}
		defer f.Close()

		_, err = f.WriteString(data)
		if err != nil {
			return nil, fmt.Errorf("GetLoginTicket: No se pudo escribir en el archivo %s: %s", fileName, err)
		}

	}

	return loginTicket, nil
}
