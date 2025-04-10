package afip

import (
	"crypto/tls"
	"encoding/xml"
	"errors"
	"fmt"

	"github.com/hooklift/gowsdl/soap"
	"github.com/sehogas/gowsaa/ws/wscoemcons"
)

type Wscoemcons struct {
	serviceName string
	environment Environment
	url         string
	cuit        int64
	tipoAgente  string
	rol         string
}

func NewWscoemcons(environment Environment, cuit int64, tipoAgente, rol string) (*Wscoemcons, error) {
	var url string
	if environment == PRODUCTION {
		url = URLWSCOEMConsProduction
	} else {
		url = URLWSCOEMConsTesting
	}

	return &Wscoemcons{
		serviceName: "wconscomunicacionembarque",
		environment: environment,
		cuit:        cuit,
		url:         url,
		tipoAgente:  tipoAgente,
		rol:         rol,
	}, nil
}

func (ws *Wscoemcons) PrintlnAsXML(obj interface{}) {
	if ws.environment == TESTING {
		data, err := xml.MarshalIndent(obj, " ", "  ")
		if err == nil {
			fmt.Println(string(data))
		}
	}
}

func (ws *Wscoemcons) Dummy() error {
	request := &wscoemcons.Dummy{}
	PrintlnAsXML(request)

	client := soap.NewClient(ws.url, soap.WithTLS(&tls.Config{InsecureSkipVerify: true}))
	service := wscoemcons.NewWconscomunicacionembarqueSoap(client)

	response, err := service.Dummy(request)
	if err != nil {
		return err
	}
	PrintlnAsXML(response)

	return nil
}

func (ws *Wscoemcons) ObtenerConsultaEstadosCOEM(identificadorCaratula string) ([]*ConsultaEstadoCOEM, error) {
	ticket, err := GetTA(ws.environment, ws.serviceName, ws.cuit)
	if err != nil {
		return nil, err
	}

	request := &wscoemcons.ObtenerConsultaEstadosCOEM{
		ArgWSAutenticacionEmpresa: &wscoemcons.WSAutenticacionEmpresa{
			CuitEmpresaConectada: ws.cuit,
			TipoAgente:           ws.tipoAgente,
			Rol:                  ws.rol,
			WSAutenticacion: &wscoemcons.WSAutenticacion{
				Token: ticket.Token,
				Sign:  ticket.Sign,
			},
		},
		IdentificadorCabecera: identificadorCaratula,
	}

	PrintlnAsXML(request)

	client := soap.NewClient(ws.url, soap.WithTLS(&tls.Config{InsecureSkipVerify: true}))
	service := wscoemcons.NewWconscomunicacionembarqueSoap(client)

	response, err := service.ObtenerConsultaEstadosCOEM(request)
	if err != nil {
		return nil, err
	}
	PrintlnAsXML(response)

	var errs []error
	for _, e := range response.ObtenerConsultaEstadosCOEMResult.Errores.ErrorEjecucion {
		errs = append(errs, fmt.Errorf("%s - %s", e.Codigo, e.Descripcion))
	}

	var listado []*ConsultaEstadoCOEM
	for _, r := range response.ObtenerConsultaEstadosCOEMResult.Resultado.Listado.ConsultaEstadoCOEM {
		listado = append(listado, &ConsultaEstadoCOEM{
			IdentificadorCOEM: r.IdentificadorCOEM,
			CuitAlta:          r.CuitAlta,
			Motivo:            r.Motivo,
			FechaEstado:       r.FechaEstado.ToGoTime(),
			Estado:            r.Estado,
			CODE:              r.CODE,
		})
	}

	return listado, errors.Join(errs...)
}

func (ws *Wscoemcons) ObtenerConsultaNoAbordo(identificadorCaratula string) ([]*ConsultaNoAbordo, error) {
	ticket, err := GetTA(ws.environment, ws.serviceName, ws.cuit)
	if err != nil {
		return nil, err
	}

	request := &wscoemcons.ObtenerConsultaNoAbordo{
		ArgWSAutenticacionEmpresa: &wscoemcons.WSAutenticacionEmpresa{
			CuitEmpresaConectada: ws.cuit,
			TipoAgente:           ws.tipoAgente,
			Rol:                  ws.rol,
			WSAutenticacion: &wscoemcons.WSAutenticacion{
				Token: ticket.Token,
				Sign:  ticket.Sign,
			},
		},
		IdentificadorCabecera: identificadorCaratula,
	}

	PrintlnAsXML(request)

	client := soap.NewClient(ws.url, soap.WithTLS(&tls.Config{InsecureSkipVerify: true}))
	service := wscoemcons.NewWconscomunicacionembarqueSoap(client)

	response, err := service.ObtenerConsultaNoAbordo(request)
	if err != nil {
		return nil, err
	}
	PrintlnAsXML(response)

	var errs []error
	for _, e := range response.ObtenerConsultaNoAbordoResult.Errores.ErrorEjecucion {
		errs = append(errs, fmt.Errorf("%s - %s", e.Codigo, e.Descripcion))
	}

	var listado []*ConsultaNoAbordo
	for _, r := range response.ObtenerConsultaNoAbordoResult.Resultado.Listado.ConsultaNoAbordo {
		listado = append(listado, &ConsultaNoAbordo{
			IdentificadorCACE:   r.IdentificadorCACE,
			IdentificadorCOEM:   r.IdentificadorCOEM,
			Tipo:                r.Tipo,
			Contenedor:          r.Contenedor,
			Destinacion:         r.Destinacion,
			Cuit:                r.Cuit,
			MotivoNoAbordo:      r.MotivoNoAbordo,
			FechaNoAbordo:       r.FechaNoAbordo.ToGoTime(),
			TipoNoAbordo:        r.TipoNoAbordo,
			DescripcionNoAbordo: r.DescripcionMotivo,
		})
	}

	return listado, errors.Join(errs...)
}

func (ws *Wscoemcons) ObtenerConsultaSolicitudes(identificadorCaratula string) ([]*ConsultaSolicitud, error) {
	ticket, err := GetTA(ws.environment, ws.serviceName, ws.cuit)
	if err != nil {
		return nil, err
	}

	request := &wscoemcons.ObtenerConsultaSolicitudes{
		ArgWSAutenticacionEmpresa: &wscoemcons.WSAutenticacionEmpresa{
			CuitEmpresaConectada: ws.cuit,
			TipoAgente:           ws.tipoAgente,
			Rol:                  ws.rol,
			WSAutenticacion: &wscoemcons.WSAutenticacion{
				Token: ticket.Token,
				Sign:  ticket.Sign,
			},
		},
		IdentificadorCabecera: identificadorCaratula,
	}

	PrintlnAsXML(request)

	client := soap.NewClient(ws.url, soap.WithTLS(&tls.Config{InsecureSkipVerify: true}))
	service := wscoemcons.NewWconscomunicacionembarqueSoap(client)

	response, err := service.ObtenerConsultaSolicitudes(request)
	if err != nil {
		return nil, err
	}
	PrintlnAsXML(response)

	var errs []error
	for _, e := range response.ObtenerConsultaSolicitudesResult.Errores.ErrorEjecucion {
		errs = append(errs, fmt.Errorf("%s - %s", e.Codigo, e.Descripcion))
	}

	var listado []*ConsultaSolicitud
	for _, r := range response.ObtenerConsultaSolicitudesResult.Resultado.Listado.ConsultaSolicitudes {
		listado = append(listado, &ConsultaSolicitud{
			IdentificadorCACE: r.IdentificadorCACE,
			NumeroSolicitud:   r.NumeroSolicitud,
			IdentificadorCOEM: r.IdentificadorCOEM,
			Estado:            r.Estado,
			FechaEstado:       r.FechaEstado.ToGoTime(),
			TipoSolicitud:     r.TipoSolicitud,
		})
	}

	return listado, errors.Join(errs...)
}
