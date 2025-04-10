package main

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/sehogas/gowsaa/afip"
)

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

	if len(os.Getenv("WSCOEM_TIPO_AGENTE")) != 4 {
		log.Fatalln("Falta o es errónea la variable de entorno WSCOEM_TIPO_AGENTE")
	}

	if len(os.Getenv("WSCOEM_ROL")) != 4 {
		log.Fatalln("Falta o es errónea la variable de entorno WSCOEM_ROL")
	}

	wscoemcons, err := afip.NewWscoemcons(environment, cuit, os.Getenv("WSCOEM_TIPO_AGENTE"), os.Getenv("WSCOEM_ROL"))
	if err != nil {
		log.Fatalln(err)
	}

	err = wscoemcons.Dummy()
	if err != nil {
		log.Fatalln(err)
	}
}
