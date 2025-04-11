package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"github.com/sehogas/gowsaa/afip"
	"github.com/sehogas/gowsaa/internal/middleware"
	"github.com/sehogas/gowsaa/internal/util"
)

var (
	Version string = "development"

	Wscoemcons *afip.Wscoemcons
)

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
		log.Fatalln("variable de entorno CERTIFICATE_FILE faltante")
	}

	if len(os.Getenv("WSCOEM_TIPO_AGENTE")) != 4 {
		log.Fatalln("variable de entorno WSCOEM_TIPO_AGENTE faltante o inválida")
	}

	if len(os.Getenv("WSCOEM_ROL")) != 4 {
		log.Fatalln("variable de entorno WSCOEM_ROL faltante o inválida")
	}

	port := 3000

	if os.Getenv("PORT") != "" {
		port, err = strconv.Atoi(os.Getenv("PORT"))
		if err != nil {
			log.Fatalln("variable de entorno PORT no numérica. ", err)
		}
	}

	Wscoemcons, err = afip.NewWscoemcons(environment, cuit, os.Getenv("WSCOEM_TIPO_AGENTE"), os.Getenv("WSCOEM_ROL"))
	if err != nil {
		log.Fatalln(err)
	}

	/* API Rest */
	router := http.NewServeMux()
	router.HandleFunc("/dummy", DummyHandler)
	router.HandleFunc("/obtener-consulta-estados-coem", ObtenerConsultaEstadosCOEMHandler)
	router.HandleFunc("/obtener-consulta-no-abordo", ObtenerConsultaNoAbordoHandler)
	router.HandleFunc("/obtener-consulta-solicitudes", ObtenerConsultaSolicitudesHandler)

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
}
