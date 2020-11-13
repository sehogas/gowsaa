package afip

import (
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/sehogas/gowsaa/soap"
)

// URLWSAATesting ... wsdl de wsaa en ambiente de homolagación de afip
const URLWSAATesting string = "https://wsaahomo.afip.gov.ar/ws/services/LoginCms?WSDL"

// URLWSAAProduction ... wsdl de wsaa en ambiente de producción de afip
const URLWSAAProduction string = "https://wsaa.afip.gov.ar/ws/services/LoginCms?WSDL"

// Ambiente es un tipo de dato
type Ambiente int

// Constantes de ambiente
const (
	TESTING Ambiente = iota
	PRODUCTION
)

// Afip es la estructura global del paquete
type Afip struct {
	ambiente Ambiente
	urlWsaa  string
	p12      string
	password string
	tickets  map[string]*LoginTicketResponse
}

// LoginTicket es una estructura que representa un ticket de un servicio de afip
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

// Create crea un objeto cliente para acceder a los servicios web de afip
func Create(ambiente Ambiente) *Afip {

	var url string
	if ambiente == PRODUCTION {
		url = URLWSAAProduction
	} else {
		url = URLWSAATesting
	}

	return &Afip{ambiente: ambiente, urlWsaa: url, tickets: make(map[string]*LoginTicketResponse)}
}

// SetFileP12 especifica el archivo .p12 que se utilizará para autenticar contra los servicios de afip
func (c *Afip) SetFileP12(p12, password string) error {

	if strings.TrimSpace(p12) == "" {
		return errors.New("SetFileP12(): Se requiere el parámetro p12")
	}

	if _, err := os.Stat(p12); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("SetFileP12(): No existe el archivo %s", p12)
		} else {
			return fmt.Errorf("SetFileP12(): Error verificando archivo .12. %s", err.Error())
		}
	}

	c.p12 = strings.TrimSpace(p12)
	c.password = strings.TrimSpace(password)

	return nil
}

// GetLoginTicket devuelve información del ticket de acceso afip correspondiente al servicio pasado por parámetro.
func (c *Afip) GetLoginTicket(serviceName string) (token string, sign string, expiration string, err error) {

	var renovar bool = true

	ticket, _ := c.tickets[serviceName]

	if ticket != nil {
		expTime, err := time.Parse(time.RFC3339, ticket.Header.ExpirationTime)

		if err != nil {
			return "", "", "", fmt.Errorf("GetLoginTicket(): Error parseando fecha de expiración del ticket anterior. %s", err)
		}

		renovar = time.Now().After(expTime)
	}

	if renovar {

		expiration := time.Now().Add(10 * time.Minute)
		generationTime := fmt.Sprintf("%s", time.Now().Add(-10*time.Minute).Format(time.RFC3339))
		expirationTime := fmt.Sprintf("%s", expiration.Format(time.RFC3339))

		// Decodifico archivo con certificado y clave privada
		certificate, privateKey, err := decodePkcs12(c.p12, c.password)

		if err != nil {
			return "", "", "", fmt.Errorf("GetLoginTicket(): %s", err)
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

		// Armo XML
		loginTicketRequestXML, err := xml.MarshalIndent(loginTicketRequest, " ", "  ")

		if err != nil {
			return "", "", "", fmt.Errorf("GetLoginTicket(): Error armando login ticket request XML. %s", err)
		}

		content := []byte(string(loginTicketRequestXML))

		// Creo CMS (Cryptographic Message Syntax)
		cms, err := encodeCMS(content, certificate, privateKey)

		if err != nil {
			return "", "", "", fmt.Errorf("GetLoginTicket(): %s", err)
		}

		cmsBase64 := base64.StdEncoding.EncodeToString(cms)

		//Llamo al servicio de autenticación afip wssa
		conexion := soap.NewClient(c.urlWsaa)
		service := NewLoginCMS(conexion)

		request := LoginCms{In0: cmsBase64}

		if c.ambiente == TESTING {
			requestXML, _ := xml.MarshalIndent(request, " ", "  ")
			log.Printf("REQUEST XML:\n%s\n\n", xml.Header+string(requestXML))
		}

		responseXML, err := service.LoginCms(&request)

		if err != nil {
			return "", "", "", fmt.Errorf("GetLoginTicket(): %s", err)
		}

		if c.ambiente == TESTING {
			log.Printf("RESPONSE XML:\n%s\n\n", responseXML)
		}

		response := LoginTicketResponse{}

		if err := xml.Unmarshal([]byte(responseXML.LoginCmsReturn), &response); err != nil {
			return "", "", "", fmt.Errorf("GetLoginTicket(): Error desarmando respuesta XML. %s", err)
		}

		// Almaceno ticket de respuesta (porque no se puede llamar nuevamente al servicio hasta dentro de 10 minutos,
		// hay que seguir usando el ticket actual. El vencimiento de los ticket de afip suele ser de 12 horas)
		c.tickets[serviceName] = &response

	}

	return c.tickets[serviceName].Credentials.Token, c.tickets[serviceName].Credentials.Sign, c.tickets[serviceName].Header.ExpirationTime, nil
}
