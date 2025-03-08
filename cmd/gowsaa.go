package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/sehogas/gowsaa/afip"
)

func main() {
	godotenv.Load()

	client, err := afip.NewClient(afip.TESTING, os.Getenv("PRIVATE_KEY_FILE"), os.Getenv("CERTIFICATE_FILE"))
	if err != nil {
		log.Fatalln(err)
	}

	loginTicket, err := client.GetLoginTicket("wsfe")
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("\n\tSERVICE: %s\n\tTOKEN: %s\n\tSIGN: %s\n\tEXPIRATION: %s\n",
		loginTicket.ServiceName,
		loginTicket.Token,
		loginTicket.Sign,
		loginTicket.ExpirationTime)

}
