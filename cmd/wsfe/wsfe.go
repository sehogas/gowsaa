package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/sehogas/gowsaa/afip"
)

func main() {
	godotenv.Load("wsfe.TA")

	environment := afip.TESTING
	if os.Getenv("PROD") == "True" {
		environment = afip.PRODUCTION
	}

	data := fmt.Sprintf("SERVICE=%s\nTOKEN=%s\nSIGN=%s\nEXPIRATION=%s\n",
		os.Getenv("SERVICE"),
		os.Getenv("TOKEN"),
		os.Getenv("SIGN"),
		os.Getenv("EXPIRATION"))

	fmt.Println(data)
	fmt.Println(environment)
}
