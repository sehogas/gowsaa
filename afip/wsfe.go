package afip

import (
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/hooklift/gowsdl/soap"
	"github.com/sehogas/gowsaa/ws/wsfe"
)

type Wsfe struct {
	environment Environment
	url         string
	loginTicket *LoginTicket
	auth        *wsfe.FEAuthRequest
}

func NewWsfe(environment Environment, loginTicket *LoginTicket) (*Wsfe, error) {
	var url string
	if environment == PRODUCTION {
		url = URLWSAAProduction
	} else {
		url = URLWSAATesting
	}

	return &Wsfe{
		environment: environment,
		url:         url,
		loginTicket: loginTicket,
		auth: &wsfe.FEAuthRequest{
			Token: loginTicket.Token,
			Sign:  loginTicket.Sign,
			Cuit:  loginTicket.Cuit,
		},
	}, nil
}

func (ws *Wsfe) FEUltimoComprobanteEmitido(ptoVta int32, cbteTipo int32) error {
	request := &wsfe.FECompUltimoAutorizado{
		Auth: &wsfe.FEAuthRequest{
			Token: ws.loginTicket.Token,
			Sign:  ws.loginTicket.Sign,
			Cuit:  ws.loginTicket.Cuit},
		PtoVta:   ptoVta,
		CbteTipo: cbteTipo,
	}

	requestXml, err := xml.MarshalIndent(request, " ", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(requestXml))

	conexion := soap.NewClient(ws.url, soap.WithTLS(&tls.Config{InsecureSkipVerify: true}))
	service := wsfe.NewServiceSoap(conexion)

	responseXml, err := service.FECompUltimoAutorizado(request)
	if err != nil {
		response := soap.SOAPEnvelopeResponse{
			Body: soap.SOAPBodyResponse{
				Content: &soap.SOAPFault{},
				Fault:   &soap.SOAPFault{},
			},
		}
		if err := xml.Unmarshal([]byte(err.Error()[strings.Index(err.Error(), "<soapenv:Envelope"):]), &response); err != nil {
			return err
		}
		return fmt.Errorf("%s", response.Body.Fault.String)

	}

	fmt.Println(responseXml.FECompUltimoAutorizadoResult.CbteNro)

	/* 	response := wsfe.FERecuperaLastCbteResponse{}
	   	tmp := responseXml.FECompUltimoAutorizadoResult
	   	if err := xml.Unmarshal(tmp, &response); err != nil {
	   		return err
	   	}  */

	fmt.Println("respuesta:", responseXml)

	return nil
}

func ObtenerCAE() error {

	return nil

}
