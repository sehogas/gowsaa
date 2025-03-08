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
	"github.com/sehogas/gowsaa/afip/wsaa"
)

const URLWSAATesting string = "https://wsaahomo.afip.gov.ar/ws/services/LoginCms?WSDL"
const URLWSAAProduction string = "https://wsaa.afip.gov.ar/ws/services/LoginCms?WSDL"

type Environment int

const (
	TESTING Environment = iota
	PRODUCTION
)

type Afip struct {
	environment     Environment
	urlWsaa         string
	tickets         map[string]*LoginTicketResponse
	privateKeyFile  string
	certificateFile string
}

type LoginTicket struct {
	ServiceName    string
	Token          string
	Sign           string
	ExpirationTime time.Time
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

func NewClient(environment Environment, privateKeyFile string, certificateFile string) (*Afip, error) {
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
	return &Afip{
		environment:     environment,
		urlWsaa:         url,
		tickets:         make(map[string]*LoginTicketResponse),
		privateKeyFile:  privateKeyFile,
		certificateFile: certificateFile,
	}, nil
}

// GetLoginTicket devuelve el ticket de acceso afip correspondiente al servicio pasado por parámetro.
func (c *Afip) GetLoginTicket(serviceName string) (*LoginTicket, error) {
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

	return &LoginTicket{
		ServiceName:    serviceName,
		Token:          c.tickets[serviceName].Credentials.Token,
		Sign:           c.tickets[serviceName].Credentials.Sign,
		ExpirationTime: expirationTime,
	}, nil
}
