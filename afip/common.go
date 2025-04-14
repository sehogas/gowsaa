package afip

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

const URLWSAATesting string = "https://wsaahomo.afip.gov.ar/ws/services/LoginCms?WSDL"
const URLWSAAProduction string = "https://wsaa.afip.gov.ar/ws/services/LoginCms?WSDL"

const URLWSFETesting string = "https://wswhomo.afip.gov.ar/wsfev1/service.asmx?WSDL"
const URLWSFEProduction string = "https://servicios1.afip.gov.ar/wsfev1/service.asmx?WSDL"

const URLWSCOEMTesting string = "https://wsaduhomoext.afip.gob.ar/DIAV2/wgescomunicacionembarque/wgescomunicacionembarque.asmx?WSDL"
const URLWSCOEMProduction string = "https://webservicesadu.afip.gob.ar/DIAV2/wgescomunicacionembarque/wgescomunicacionembarque.asmx?WSDL"

const URLWSCOEMConsTesting string = "https://wsaduhomoext.afip.gob.ar/DIAV2/wconscomunicacionembarque/wconscomunicacionembarque.asmx?WSDL"
const URLWSCOEMConsProduction string = "https://webservicesadu.afip.gob.ar/DIAV2/wconscomunicacionembarque/wconscomunicacionembarque.asmx?WSDL"

const URLWGESTABREFTesting string = "https://testdia.afip.gob.ar/Dia/ws/wgesTabRef/wgesTabRef.asmx?WSDL"
const URLWGESTABREFProduction string = "https://servicios1.afip.gov.ar/Dia/ws/wgesTabRef/wgesTabRef.asmx?WSDL"

type Environment int

const (
	TESTING Environment = iota
	PRODUCTION
)

// func GenTA(environment Environment, serviceName string, cuit int64) error {
// 	wsaa, err := NewWsaa(environment,
// 		os.Getenv("PRIVATE_KEY_FILE"),
// 		os.Getenv("CERTIFICATE_FILE"),
// 		cuit)
// 	if err != nil {
// 		return err
// 	}

// 	loginTicket, err := wsaa.GetLoginTicket(serviceName)
// 	if err != nil {
// 		return err
// 	}

// 	os.Setenv("CUIT", strconv.FormatInt(loginTicket.Cuit, 10))
// 	os.Setenv("TOKEN", loginTicket.Token)
// 	os.Setenv("SIGN", loginTicket.Sign)
// 	os.Setenv("EXPIRATION", loginTicket.ExpirationTime.Format(time.RFC3339))

//		return nil
//	}
var tickets map[string]*LoginTicket

func GenerarTA(environment Environment, serviceName string, cuit int64) (*LoginTicket, error) {
	wsaa, err := NewWsaa(environment,
		os.Getenv("PRIVATE_KEY_FILE"),
		os.Getenv("CERTIFICATE_FILE"),
		cuit)
	if err != nil {
		return nil, err
	}

	return wsaa.GetLoginTicket(serviceName)
}

func GrabarTA(serviceName string, ticket *LoginTicket) error {
	fileName := fmt.Sprintf("data/%s.TA", serviceName)
	prefix := strings.TrimSpace(strings.ToUpper(serviceName))
	log.Printf("Almacenando ticket de acceso para  [%s] en [%s]...\n", serviceName, fileName)
	data := fmt.Sprintf("%s_CUIT=%s\n%s_TOKEN=%s\n%s_SIGN=%s\n%s_EXPIRATION=%s\n",
		prefix, strconv.FormatInt(ticket.Cuit, 10),
		prefix, ticket.Token,
		prefix, ticket.Sign,
		prefix, ticket.ExpirationTime.Format(time.RFC3339))
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(data)
	if err != nil {
		return err
	}
	log.Printf("Ticket de acceso almacenado para [%s]\n", serviceName)
	return nil
}

func GetTA(environment Environment, serviceName string, cuit int64) (*LoginTicket, error) {
	if tickets == nil {
		tickets = make(map[string]*LoginTicket)
	}
	ticketMemory, exist := tickets[serviceName]
	if exist {
		if ticketMemory.ExpirationTime.After(time.Now()) {
			return ticketMemory, nil
		}
	}

	var ticket *LoginTicket
	err := godotenv.Load(fmt.Sprintf("data/%s.TA", serviceName))
	if err != nil {
		log.Printf("Intentando genera ticket de acceso para [%s]...\n", serviceName)
		ticket, err = GenerarTA(environment, serviceName, cuit)
		if err != nil {
			return nil, err
		}
		log.Printf("Ticket de acceso emitido para [%s]\n", serviceName)
		tickets[serviceName] = ticket
		if err = GrabarTA(serviceName, ticket); err != nil {
			return ticket, err
		}
	} else {
		prefix := strings.TrimSpace(strings.ToUpper(serviceName))
		expirationTime, err := time.Parse(time.RFC3339, os.Getenv(fmt.Sprintf("%s_EXPIRATION", prefix)))
		if err != nil {
			return nil, err
		}
		ticket = &LoginTicket{
			ServiceName:    serviceName,
			Token:          os.Getenv(fmt.Sprintf("%s_TOKEN", prefix)),
			Sign:           os.Getenv(fmt.Sprintf("%s_SIGN", prefix)),
			ExpirationTime: expirationTime,
			Cuit:           cuit,
		}
	}

	if ticket.ExpirationTime.Before(time.Now()) {
		log.Printf("Renovando ticket de acceso para [%s]...\n", serviceName)
		ticket, err = GenerarTA(environment, serviceName, cuit)
		if err != nil {
			return nil, err
		}
		log.Printf("Ticket de acceso emitido para [%s]\n", serviceName)
		if err = GrabarTA(serviceName, ticket); err != nil {
			tickets[serviceName] = ticket
			return ticket, err
		}
	}

	tickets[serviceName] = ticket
	return ticket, nil
}
