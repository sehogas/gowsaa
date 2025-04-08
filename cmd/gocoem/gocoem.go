package main

import (
	"fmt"
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
	serviceName := "wgescomunicacionembarque"

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

	wscoem, err := afip.NewWscoem(environment, &afip.LoginTicket{
		ServiceName:    serviceName,
		Token:          os.Getenv("TOKEN"),
		Sign:           os.Getenv("SIGN"),
		ExpirationTime: expirationTime,
		Cuit:           cuit,
	}, os.Getenv("WSCOEM_TIPO_AGENTE"), os.Getenv("WSCOEM_ROL"))
	if err != nil {
		panic(err)
	}

	err = wscoem.Dummy()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	/*
		identificadorAnulado, err := wscoem.AnularCaratula("25067EMBA000279P")
		if identificadorAnulado != "" {
			log.Printf("Carátula anulada: %s\n", identificadorAnulado)
			if err != nil {
				log.Println(err) //Imprimo los warnings
			}
		} else {
			log.Println(err) //Imprimo los errores
			os.Exit(1)
		}
	*/

	/*
		identificador, err := wscoem.RegistrarCaratula(&afip.CaratulaParams{
			IdentificadorBuque:    "IMO9262871",
			NombreMedioTransporte: "ARGENTINO II",
			CodigoAduana:          "067",
			CodigoLugarOperativo:  "10056",
			FechaEstimadaArribo:   time.Now().AddDate(0, 0, 1).Local(),
			FechaEstimadaZarpada:  time.Now().AddDate(0, 0, 3).Local(),
			Via:                   "8", //ACUATICA --SIEMPRE 8
			NumeroViaje:           "",
			PuertoDestino:         "ARBUE",
			Itinerario:            nil,
		})
		if identificador != "" {
			log.Printf("Carátula registrada: %s\n", identificador)
			if err != nil {
				log.Println(err) //Imprimo Warnings
			}
		} else {
			log.Println(err) //Imprimo errores
			os.Exit(1)
		}
	*/

	/*
		identificador, err := wscoem.RectificarCaratula("25067EMBA000281X", &afip.CaratulaParams{
			IdentificadorBuque:    "IMO9262871",
			NombreMedioTransporte: "ARGENTINO II",
			CodigoAduana:          "067",
			CodigoLugarOperativo:  "10056",
			FechaEstimadaArribo:   time.Now().AddDate(0, 0, 1).Local(),
			FechaEstimadaZarpada:  time.Now().AddDate(0, 0, 3).Local(),
			Via:                   "8", //ACUATICA --SIEMPRE 8
			NumeroViaje:           "",
			PuertoDestino:         "ARBUE",
			Itinerario:            nil,
		})
		if identificador != "" {
			log.Printf("Carátula rectificada: %s\n", identificador)
			if err != nil {
				log.Println(err) //Imprimo Warnings
			}
		} else {
			log.Println(err) //Imprimo errores
			os.Exit(1)
		}
	*/

	/*
		identificador, err := wscoem.SolicitarCambioBuque("25067EMBA000281X", &afip.CambioBuqueParams{
			IdentificadorBuque:    "IMO9262871",
			NombreMedioTransporte: "ARGENTINO II",
		})
		if identificador != "" {
			log.Printf("Cambio de buque OK para carátula: %s\n", identificador)
			if err != nil {
				log.Println(err) //Imprimo Warnings
			}
		} else {
			log.Println(err) //Imprimo errores
			os.Exit(1)
		}
	*/
	/*
		identificador, err := wscoem.SolicitarCambioFechas("25067EMBA000281X", &afip.CambioFechasParams{
			FechaEstimadaArribo:  time.Now().AddDate(0, 0, 1).Local(),
			FechaEstimadaZarpada: time.Now().AddDate(0, 0, 3).Local(),
			CodigoMotivo:         "1",
			DescripcionMotivo:    "Se averio",
		})
		if identificador != "" {
			log.Printf("Cambio de fechas OK para carátula: %s\n", identificador)
			if err != nil {
				log.Println(err) //Imprimo Warnings
			}
		} else {
			log.Println(err) //Imprimo errores
			os.Exit(1)
		}
	*/
	/*
		identificador, err := wscoem.SolicitarCambioLOT("25067EMBA000281X", &afip.CambioLOTParams{
			CodigoLugarOperativo: "10057",
		})
		if identificador != "" {
			log.Printf("Cambio de LOT OK para carátula: %s\n", identificador)
			if err != nil {
				log.Println(err) //Imprimo Warnings
			}
		} else {
			log.Println(err) //Imprimo errores
			os.Exit(1)
		}*/

	/*
		plan, err := os.ReadFile("./cmd/gocoem/data_test/coem.json") // filename is the JSON file to read
		if err != nil {
			log.Println("Error abriendo archivo", err)
		}
		var datos afip.COEMParams
		err = json.Unmarshal(plan, &datos)
		if err != nil {
			log.Println("Cannot unmarshal the json ", err)
		}
		afip.PrintlnAsJSON(datos)

		identificadorCOEM, err := wscoem.RegistrarCOEM("25067EMBA000281X", &datos)
		if identificadorCOEM != "" {
			log.Printf("COEM registrado: %s\n", identificadorCOEM)
			if err != nil {
				log.Println(err) //Imprimo Warnings
			}
		}
		if err != nil {
			log.Println(err) //Imprimo errores
			os.Exit(1)
		}
	*/
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
