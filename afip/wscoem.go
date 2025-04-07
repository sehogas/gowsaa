package afip

import (
	"crypto/tls"
	"encoding/xml"
	"fmt"

	"github.com/hooklift/gowsdl/soap"
	"github.com/sehogas/gowsaa/ws/wscoem"
)

type Wscoem struct {
	environment Environment
	url         string
	loginTicket *LoginTicket
	auth        *wscoem.WSAutenticacionEmpresa
}

func NewWscoem(environment Environment, loginTicket *LoginTicket, tipoAgente, rol string) (*Wscoem, error) {
	var url string
	if environment == PRODUCTION {
		url = URLWSCOEMProduction
	} else {
		url = URLWSCOEMTesting
	}

	return &Wscoem{
		environment: environment,
		url:         url,
		loginTicket: loginTicket,
		auth: &wscoem.WSAutenticacionEmpresa{
			WSAutenticacion: &wscoem.WSAutenticacion{
				Token: loginTicket.Token,
				Sign:  loginTicket.Sign,
			},
			CuitEmpresaConectada: loginTicket.Cuit,
			TipoAgente:           tipoAgente,
			Rol:                  rol,
		},
	}, nil
}

func (ws *Wscoem) Dummy() error {

	request := &wscoem.Dummy{}

	requestXml, err := xml.MarshalIndent(request, " ", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(requestXml))

	client := soap.NewClient(ws.url, soap.WithTLS(&tls.Config{InsecureSkipVerify: true}))
	service := wscoem.NewWgescomunicacionembarqueSoap(client)

	response, err := service.Dummy(request)
	if err != nil {
		return err
	}

	if ws.environment == TESTING {
		fmt.Println("DUMMY AppServer: ", response.DummyResult.AppServer)
		fmt.Println("DUMMY AuthServer: ", response.DummyResult.AuthServer)
		fmt.Println("DUMMY DbServer: ", response.DummyResult.DbServer)
	}

	return nil
}

func (ws *Wscoem) RegistrarCaratula(params *RegistrarCaratulaParams) error {
	request := &wscoem.RegistrarCaratula{
		ArgWSAutenticacionEmpresa: ws.auth,
		ArgRegistrarCaratula: &wscoem.RegistrarCaratulaRequest{
			Caratula: &wscoem.Caratula{
				IdentificadorBuque:    params.IdentificadorBuque,
				NombreMedioTransporte: params.NombreMedioTransporte,
				CodigoAduana:          params.CodigoAduana,
				CodigoLugarOperativo:  params.CodigoLugarOperativo,
				FechaArribo:           soap.CreateXsdDateTime(params.FechaEstimadaArribo, true),
				FechaZarpada:          soap.CreateXsdDateTime(params.FechaEstimadaZarpada, true),
				Via:                   params.Via,
				NumeroViaje:           params.NumeroViaje,
				PuertoDestino:         params.PuertoDestino,
				Itinerario:            nil,
			},
		},
	}

	requestXml, err := xml.MarshalIndent(request, " ", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(requestXml))

	client := soap.NewClient(ws.url, soap.WithTLS(&tls.Config{InsecureSkipVerify: true}))
	service := wscoem.NewWgescomunicacionembarqueSoap(client)

	response, err := service.RegistrarCaratula(request)
	if err != nil {
		return err
	}

	responseXml, err := xml.MarshalIndent(response, " ", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(responseXml))

	return nil
}
