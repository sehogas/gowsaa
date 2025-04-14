package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"github.com/sehogas/gowsaa/afip"
	"github.com/sehogas/gowsaa/cmd/api/docs"
	"github.com/sehogas/gowsaa/internal/middleware"
	"github.com/sehogas/gowsaa/internal/util"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

var (
	Version string = "development"

	Wscoem      *afip.Wscoem
	Wscoemcons  *afip.Wscoemcons
	Wsgestabref *afip.Wsgestabref

	validate *validator.Validate
)

//	@title			API Rest Json proxy a los servicios web SOAP ARCA
//	@version		1.0
//	@description	Esta API Json Rest actua como proxy SOAP a los servicios web de ARCA.
//	@termsOfService	http://swagger.io/terms/

// @contact.name	Sebastian Hogas
// @contact.email	sehogas@gmail.com
func main() {
	godotenv.Load()

	environment := afip.TESTING
	if strings.ToLower(strings.TrimSpace(os.Getenv("PROD"))) == "true" {
		environment = afip.PRODUCTION
	}

	cuit, err := strconv.ParseInt(os.Getenv("CUIT"), 10, 64)
	if err != nil {
		log.Fatalln("variable de entorno CUIT faltante o no numérica")
	}

	if os.Getenv("PRIVATE_KEY_FILE") == "" {
		log.Fatalln("variable de entorno PRIVATE_KEY_FILE faltante")
	}

	if os.Getenv("CERTIFICATE_FILE") == "" {
		log.Fatalln("falta variable de entorno CERTIFICATE_FILE")
	}

	if len(os.Getenv("WSCOEM_TIPO_AGENTE")) != 4 {
		log.Fatalln("falta variable de entorno WSCOEM_TIPO_AGENTE")
	}

	if len(os.Getenv("WSCOEM_ROL")) != 4 {
		log.Fatalln("variable de entorno WSCOEM_ROL faltante o inválida")
	}

	if len(os.Getenv("WGESTABREF_TIPO_AGENTE")) != 4 {
		log.Fatalln("falta variable de entorno WGESTABREF_TIPO_AGENTE")
	}

	if len(os.Getenv("WGESTABREF_ROL")) != 4 {
		log.Fatalln("variable de entorno WGESTABREF_ROL faltante o inválida")
	}

	port := 3000

	if os.Getenv("PORT") != "" {
		port, err = strconv.Atoi(os.Getenv("PORT"))
		if err != nil {
			log.Fatalln("variable de entorno PORT no numérica. ", err)
		}
	}

	Wscoem, err = afip.NewWscoem(environment, cuit, os.Getenv("WSCOEM_TIPO_AGENTE"), os.Getenv("WSCOEM_ROL"))
	if err != nil {
		log.Fatalln(err)
	}

	Wscoemcons, err = afip.NewWscoemcons(environment, cuit, os.Getenv("WSCOEM_TIPO_AGENTE"), os.Getenv("WSCOEM_ROL"))
	if err != nil {
		log.Fatalln(err)
	}

	Wsgestabref, err = afip.Newgestabref(environment, cuit, os.Getenv("WGESTABREF_TIPO_AGENTE"), os.Getenv("WGESTABREF_ROL"))
	if err != nil {
		log.Fatalln(err)
	}

	/* API Rest */
	validate = validator.New(validator.WithRequiredStructEnabled())

	router := http.NewServeMux()

	docs.SwaggerInfo.Host = os.Getenv("HOST")
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	router.HandleFunc("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	))

	coem := http.NewServeMux()
	coem.HandleFunc("/dummy", DummyCoemHandler)
	coem.HandleFunc("POST /registrar-caratula", RegistrarCaratulaHandler)
	coem.HandleFunc("PUT /rectificar-caratula", RectificarCaratulaHandler)
	coem.HandleFunc("DELETE /anular-caratula", AnularCaratulaHandler)
	coem.HandleFunc("PUT /solicitar-cambio-buque", SolicitarCambioBuqueHandler)
	coem.HandleFunc("PUT /solicitar-cambio-fechas", SolicitarCambioFechasHandler)
	coem.HandleFunc("PUT /solicitar-cambio-lot", SolicitarCambioLOTHandler)
	coem.HandleFunc("POST /registrar-coem", RegistrarCOEMHandler)
	coem.HandleFunc("PUT /rectificar-coem", RectificarCOEMHandler)
	coem.HandleFunc("POST /cerrar-coem", CerrarCOEMHandler)
	coem.HandleFunc("DELETE /anular-coem", AnularCOEMHandler)
	coem.HandleFunc("POST /solicitar-anulacion-coem", SolicitarAnulacionCOEMHandler)
	coem.HandleFunc("POST /solicitar-no-abordo", SolicitarNoABordoHandler)
	coem.HandleFunc("POST /solicitar-cierre-carga-conto-bulto", SolicitarCierreCargaContoBultoHandler)
	coem.HandleFunc("POST /solicitar-cierre-carga-granel", SolicitarCierreCargaGranelHandler)

	coemcons := http.NewServeMux()
	coemcons.HandleFunc("/dummy", DummyCoemconsHandler)
	coemcons.HandleFunc("/obtener-consulta-estados-coem", ObtenerConsultaEstadosCOEMHandler)
	coemcons.HandleFunc("/obtener-consulta-no-abordo", ObtenerConsultaNoAbordoHandler)
	coemcons.HandleFunc("/obtener-consulta-solicitudes", ObtenerConsultaSolicitudesHandler)

	gestabref := http.NewServeMux()
	gestabref.HandleFunc("/dummy", DummyGesTabRefHandler)
	gestabref.HandleFunc("/consultar-fecha-ult-act", ConsultarFechaUltActHandler)
	gestabref.HandleFunc("/lista-arancel", ListaArancelHandler)
	gestabref.HandleFunc("/lista-descripcion", ListaDescripcionHandler)
	gestabref.HandleFunc("/lista-descripcion-decodificacion", ListaDescripcionDecodificacionHandler)
	gestabref.HandleFunc("/lista-empresas", ListaEmpresasHandler)
	gestabref.HandleFunc("/lista-lugares-operativos", ListaLugaresOperativosHandler)
	gestabref.HandleFunc("/lista-paises-aduanas", ListaPaisesAduanasHandler)
	gestabref.HandleFunc("/lista-tablas-referencia", ListaTablasReferenciaHandler)

	v1 := http.NewServeMux()
	v1.HandleFunc("/info", InfoHandler)
	v1.Handle("/coem/", http.StripPrefix("/coem", coem))
	v1.Handle("/coemcons/", http.StripPrefix("/coemcons", coemcons))
	v1.Handle("/gestabref/", http.StripPrefix("/gestabref", gestabref))
	router.Handle("/api/v1/", http.StripPrefix("/api/v1", v1))

	middlewareCors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		Debug:            false,
	})

	stack := middleware.CreateStack(
		middlewareCors.Handler,
		middleware.Logging,
	)

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      stack(router),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	go util.GracefulShutdown(server)

	log.Printf("Starting server on port %v", server.Addr)
	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("http server error: %s", err))
	}

	/*
		err = wscoem.Dummy()
		if err != nil {
			log.Fatalln(err)
		}

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
