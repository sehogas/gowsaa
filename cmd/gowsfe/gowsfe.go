package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/sehogas/gowsaa/afip"
)

var cuit int64

func main() {
	serviceName := "wsfe"

	godotenv.Load()

	environment := afip.TESTING
	if strings.ToLower(strings.TrimSpace(os.Getenv("PROD"))) == "true" {
		environment = afip.PRODUCTION
	}

	cuit, err := strconv.ParseInt(os.Getenv("CUIT"), 10, 64)
	if err != nil {
		fmt.Println("Falta o es errónea la variable de entorno CUIT")
		os.Exit(-1)
	}

	err = godotenv.Load(fmt.Sprintf("%s.TA", serviceName))
	if err != nil {
		err = GenTA(environment, serviceName)
		if err != nil {
			panic(err)
		}
	}

	expirationTime, err := time.Parse(time.RFC3339, os.Getenv("EXPIRATION"))
	if err != nil {
		panic(err)
	}

	if expirationTime.Before(time.Now()) {
		fmt.Println("TA expirado. Renovando ticket de acceso...")
		err = GenTA(environment, serviceName)
		if err != nil {
			panic(err)
		}
	}

	data := fmt.Sprintf("CUIT=%s\nTOKEN=%s\nSIGN=%s\nEXPIRATION=%s\n",
		os.Getenv("CUIT"),
		os.Getenv("TOKEN"),
		os.Getenv("SIGN"),
		os.Getenv("EXPIRATION"))

	fmt.Println(data)

	expirationTime, err = time.Parse(time.RFC3339, os.Getenv("EXPIRATION"))
	if err != nil {
		panic(err)
	}

	wsfe, err := afip.NewWsfe(environment, &afip.LoginTicket{
		ServiceName:    serviceName,
		Token:          os.Getenv("TOKEN"),
		Sign:           os.Getenv("SIGN"),
		ExpirationTime: expirationTime,
		Cuit:           cuit,
	})
	if err != nil {
		panic(err)
	}

	err = wsfe.FEUltimoComprobanteEmitido(1, 6) //PtoVta=1, TpoCbte=6 (Factura B)
	if err != nil {
		fmt.Println(err)
	}

}

func GenTA(environment afip.Environment, serviceName string) error {
	wsaa, err := afip.NewWsaa(environment,
		os.Getenv("PRIVATE_KEY_FILE"),
		os.Getenv("CERTIFICATE_FILE"),
		cuit)
	if err != nil {
		return err
	}

	loginTicket, err := wsaa.GetLoginTicket(serviceName)
	if err != nil {
		return err
	}

	os.Setenv("CUIT", strconv.FormatInt(loginTicket.Cuit, 10))
	os.Setenv("TOKEN", loginTicket.Token)
	os.Setenv("SIGN", loginTicket.Sign)
	os.Setenv("EXPIRATION", loginTicket.ExpirationTime.Format(time.RFC3339))

	return nil
}
