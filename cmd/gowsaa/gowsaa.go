package main

import (
	"flag"
	"fmt"
	"os"

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
			os.Exit(0)
		}
	}

	client, err := afip.NewClient(environment, os.Getenv("PRIVATE_KEY_FILE"), os.Getenv("CERTIFICATE_FILE"))
	if err != nil {
		panic(err)
	}

	loginTicket, err := client.GetLoginTicket(*service)
	if err != nil {
		panic(err)
	}

	data := fmt.Sprintf("SERVICE=%s\nTOKEN=%s\nSIGN=%s\nEXPIRATION=%s\n",
		loginTicket.ServiceName,
		loginTicket.Token,
		loginTicket.Sign,
		loginTicket.ExpirationTime)

	f, err := os.Create(fmt.Sprintf("%s.TA", *service))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	l, err := f.WriteString(data)
	if err != nil {
		panic(err)
	}

	fmt.Println(l, "bytes written successfully")
}
