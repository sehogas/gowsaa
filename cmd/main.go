package main

import (
	"log"
	"os"

	"github.com/sehogas/gowsaa/afip"
)

func main() {
	log.Println("Inicio del proceso...")

	log.Printf("archivo p12 => %s\n", os.Getenv("AfipP12"))
	log.Printf("password    => %s\n", os.Getenv("AfipP12Pass"))

	// Creo instancia de servicios afip
	cliente := afip.Create(afip.TESTING)

	// Requerido para los servicios que requieren autenticación
	err := cliente.SetFileP12(os.Getenv("AfipP12"), os.Getenv("AfipP12Pass"))

	if err != nil {
		log.Fatalf("Fin de proceso con ERROR.\n\n%s\n", err.Error())
	}

	token, sign, expiration, err := cliente.GetLoginTicket("wgesstockdepositosfiscales")

	if err != nil {
		log.Fatalf("Fin de proceso con ERROR.\n\n%s\n", err.Error())
	}

	log.Printf("TOKEN: %s\n\n", token)
	log.Printf("SIGN: %s\n\n", sign)
	log.Printf("EXPIRATION: %s\n\n", expiration)

	log.Println("Fin del proceso")

}
