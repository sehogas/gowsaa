package main

import (
	"log"
	"os"

	"github.com/sehogas/gowsaa/afip"
)

func main() {
	log.Printf("Inicio del proceso...\n\n")

	// Creo instancia de afip
	cliente := afip.Create(afip.TESTING)

	// Requerido para los servicios que necesitan autenticación de AFIP
	err := cliente.SetFileP12(os.Getenv("AfipP12"), os.Getenv("AfipP12Pass"))
	if err != nil {
		log.Fatal(err)
	}

	// Llamo al webservice wsaa para obtener el token, sign y expiration
	token, sign, expiration, err := cliente.GetLoginTicket("wgesstockdepositosfiscales")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("\n\tTOKEN: %s\n\tSIGN:%s\n\tEXPIRATION:%s\n", token, sign, expiration)

	log.Printf("\nFin del proceso\n")
}
