package main

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/sehogas/gowsaa/afip"
)

var cuit int64

func main() {
	godotenv.Load()

	environment := afip.TESTING
	if strings.ToLower(strings.TrimSpace(os.Getenv("PROD"))) == "true" {
		environment = afip.PRODUCTION
	}

	cuit, err := strconv.ParseInt(os.Getenv("CUIT"), 10, 64)
	if err != nil {
		log.Fatalln("Falta o es errónea la variable de entorno CUIT")
	}

	if os.Getenv("PRIVATE_KEY_FILE") == "" {
		log.Fatalln("Falta variable de entorno PRIVATE_KEY_FILE")
	}

	if os.Getenv("CERTIFICATE_FILE") == "" {
		log.Fatalln("Falta variable de entorno CERTIFICATE_FILE")
	}

	wsfe, err := afip.NewWsfe(environment, cuit)
	if err != nil {
		log.Fatalln(err)
	}

	cbteNro, err := wsfe.FEUltimoComprobanteEmitido(1, 6) //PtoVta=1, TpoCbte=6 (Factura B)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("CbteNro: %d\n", cbteNro)
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
