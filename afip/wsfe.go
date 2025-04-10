package afip

import (
	"crypto/tls"
	"errors"
	"fmt"

	"github.com/hooklift/gowsdl/soap"
	"github.com/sehogas/gowsaa/ws/wsfe"
)

type Wsfe struct {
	serviceName string
	environment Environment
	url         string
	cuit        int64
}

func NewWsfe(environment Environment, cuit int64) (*Wsfe, error) {
	var url string
	if environment == PRODUCTION {
		url = URLWSFEProduction
	} else {
		url = URLWSFETesting
	}

	return &Wsfe{
		serviceName: "wsfe",
		environment: environment,
		url:         url,
		cuit:        cuit,
	}, nil
}

func (ws *Wsfe) FEUltimoComprobanteEmitido(ptoVta int32, cbteTipo int32) (int32, error) {
	ticket, err := GetTA(ws.environment, ws.serviceName, ws.cuit)
	if err != nil {
		return 0, err
	}

	request := &wsfe.FECompUltimoAutorizado{
		Auth: &wsfe.FEAuthRequest{
			Token: ticket.Token,
			Sign:  ticket.Sign,
			Cuit:  ticket.Cuit},
		PtoVta:   ptoVta,
		CbteTipo: cbteTipo,
	}

	PrintlnAsXML(request)

	conexion := soap.NewClient(ws.url, soap.WithTLS(&tls.Config{InsecureSkipVerify: true}))
	service := wsfe.NewServiceSoap(conexion)

	response, err := service.FECompUltimoAutorizado(request)
	if err != nil {
		return 0, err
	}
	PrintlnAsXML(response)

	var errs []error
	if response.FECompUltimoAutorizadoResult.Events != nil {
		for _, e := range response.FECompUltimoAutorizadoResult.Events.Evt {
			errs = append(errs, fmt.Errorf("event %d - %s", e.Code, e.Msg))
		}
	}
	if response.FECompUltimoAutorizadoResult.Errors != nil {
		for _, e := range response.FECompUltimoAutorizadoResult.Errors.Err {
			errs = append(errs, fmt.Errorf("error %d - %s", e.Code, e.Msg))
		}
	} else {
		cbteNro := response.FECompUltimoAutorizadoResult.CbteNro
		return cbteNro, errors.Join(errs...)
	}

	return 0, errors.Join(errs...)
}

func ObtenerCAE() error {

	return nil

}
