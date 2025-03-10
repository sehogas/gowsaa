package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/sehogas/gowsaa/afip"
)

func main() {
	godotenv.Load()

	service := flag.String("service", "wsfe", "Nombre del servicio afip para el cual se quiere obtener el ticket de acceso. Default: wsfe")
	flag.Parse()

	environment := afip.TESTING
	if os.Getenv("PROD") == "True" {
		environment = afip.PRODUCTION
	}

	if environment == afip.TESTING {
		// leo si hay un ticket de prueba
		if _, err := os.Stat(fmt.Sprintf("%s.TA", *service)); err == nil {
			fmt.Println("Ya existe un ticket de prueba para el servicio", *service)
			os.Exit(-1)
		}
	}

	cuit, err := strconv.ParseInt(os.Getenv("CUIT"), 10, 64)
	if err != nil {
		fmt.Println("Falta o es errónea la variable de entorno CUIT")
		os.Exit(-1)
	}

	wsaa, err := afip.NewWsaa(environment,
		os.Getenv("PRIVATE_KEY_FILE"),
		os.Getenv("CERTIFICATE_FILE"),
		cuit)
	if err != nil {
		panic(err)
	}

	_, err = wsaa.GetLoginTicket(*service)
	if err != nil {
		panic(err)
	}
}
