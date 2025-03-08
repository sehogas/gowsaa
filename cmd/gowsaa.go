package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/sehogas/gowsaa/afip"
)

func main() {
	godotenv.Load()

	environment := afip.TESTING
	if os.Getenv("PROD") == "True" {
		environment = afip.PRODUCTION
	}

	client, err := afip.NewClient(environment, os.Getenv("PRIVATE_KEY_FILE"), os.Getenv("CERTIFICATE_FILE"))
	if err != nil {
		panic(err)
	}

	loginTicket, err := client.GetLoginTicket("wsfe")
	if err != nil {
		panic(err)
	}

	fmt.Printf("SERVICE=%s\nTOKEN=%s\nSIGN=%s\nEXPIRATION=%s\n",
		loginTicket.ServiceName,
		loginTicket.Token,
		loginTicket.Sign,
		loginTicket.ExpirationTime)

}
