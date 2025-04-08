package afip

import (
	"crypto/tls"
	"encoding/xml"
	"errors"
	"fmt"
	"strings"

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

func (ws *Wscoem) PrintlnAsXML(obj interface{}) {
	if ws.environment == TESTING {
		data, err := xml.MarshalIndent(obj, " ", "  ")
		if err == nil {
			fmt.Println(string(data))
		}
	}
}

func (ws *Wscoem) Dummy() error {

	request := &wscoem.Dummy{}
	PrintlnAsXML(request)

	client := soap.NewClient(ws.url, soap.WithTLS(&tls.Config{InsecureSkipVerify: true}))
	service := wscoem.NewWgescomunicacionembarqueSoap(client)

	response, err := service.Dummy(request)
	if err != nil {
		return err
	}
	PrintlnAsXML(response)

	return nil
}

func (ws *Wscoem) RegistrarCaratula(params *CaratulaParams) (string, error) {
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
	PrintlnAsXML(request)

	client := soap.NewClient(ws.url, soap.WithTLS(&tls.Config{InsecureSkipVerify: true}))
	service := wscoem.NewWgescomunicacionembarqueSoap(client)

	response, err := service.RegistrarCaratula(request)
	if err != nil {
		return "", err
	}
	PrintlnAsXML(response)

	var errs []error
	identificador := ""
	for _, e := range response.RegistrarCaratulaResult.ListaErrores.DetalleError {
		if *e.Codigo != 0 {
			errs = append(errs, fmt.Errorf("%d - %s - %s", *e.Codigo, e.Descripcion, e.DescripcionAdicional))
		} else {
			identificador = strings.TrimSpace(strings.Replace(e.DescripcionAdicional, "Identificador:", "", -1))
		}
	}

	return identificador, errors.Join(errs...)
}

func (ws *Wscoem) AnularCaratula(IdentificadorCaratula string) (string, error) {
	request := &wscoem.AnularCaratula{
		ArgWSAutenticacionEmpresa: ws.auth,
		ArgAnularCaratula: &wscoem.AnularCaratulaRequest{
			IdentificadorCaratula: IdentificadorCaratula,
		},
	}
	PrintlnAsXML(request)

	client := soap.NewClient(ws.url, soap.WithTLS(&tls.Config{InsecureSkipVerify: true}))
	service := wscoem.NewWgescomunicacionembarqueSoap(client)

	response, err := service.AnularCaratula(request)
	if err != nil {
		return "", err
	}
	PrintlnAsXML(response)

	var errs []error
	result := ""
	for _, e := range response.AnularCaratulaResult.ListaErrores.DetalleError {
		if *e.Codigo != 0 {
			errs = append(errs, fmt.Errorf("%d - %s - %s", *e.Codigo, e.Descripcion, e.DescripcionAdicional))
		} else {
			result = strings.TrimSpace(strings.Replace(e.DescripcionAdicional, "Identificador:", "", -1))
		}
	}
	return result, errors.Join(errs...)
}

func (ws *Wscoem) RectificarCaratula(IdentificadorCaratula string, params *CaratulaParams) (string, error) {
	request := &wscoem.RectificarCaratula{
		ArgWSAutenticacionEmpresa: ws.auth,
		ArgRectificarCaratula: &wscoem.RectificarCaratulaRequest{
			IdentificadorCaratula: IdentificadorCaratula,
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
	PrintlnAsXML(request)

	client := soap.NewClient(ws.url, soap.WithTLS(&tls.Config{InsecureSkipVerify: true}))
	service := wscoem.NewWgescomunicacionembarqueSoap(client)

	response, err := service.RectificarCaratula(request)
	if err != nil {
		return "", err
	}
	PrintlnAsXML(response)

	var errs []error
	result := ""
	for _, e := range response.RectificarCaratulaResult.ListaErrores.DetalleError {
		if *e.Codigo != 0 {
			errs = append(errs, fmt.Errorf("%d - %s - %s", *e.Codigo, e.Descripcion, e.DescripcionAdicional))
		} else {
			result = strings.TrimSpace(strings.Replace(e.DescripcionAdicional, "Identificador:", "", -1))
		}
	}

	return result, errors.Join(errs...)
}

func (ws *Wscoem) SolicitarCambioBuque(IdentificadorCaratula string, params *CambioBuqueParams) (string, error) {
	request := &wscoem.SolicitarCambioBuque{
		ArgWSAutenticacionEmpresa: ws.auth,
		ArgSolicitarCambioBuque: &wscoem.SolicitarCambioBuqueRequest{
			IdentificadorCaratula: IdentificadorCaratula,
			IdentificadorBuque:    params.IdentificadorBuque,
			NombreMedioTransporte: params.NombreMedioTransporte,
		},
	}
	PrintlnAsXML(request)

	client := soap.NewClient(ws.url, soap.WithTLS(&tls.Config{InsecureSkipVerify: true}))
	service := wscoem.NewWgescomunicacionembarqueSoap(client)

	response, err := service.SolicitarCambioBuque(request)
	if err != nil {
		return "", err
	}
	PrintlnAsXML(response)

	var errs []error
	result := ""
	for _, e := range response.SolicitarCambioBuqueResult.ListaErrores.DetalleError {
		if *e.Codigo != 0 {
			errs = append(errs, fmt.Errorf("%d - %s - %s", *e.Codigo, e.Descripcion, e.DescripcionAdicional))
		} else {
			result = strings.TrimSpace(strings.Replace(e.DescripcionAdicional, "Identificador:", "", -1))
		}
	}

	return result, errors.Join(errs...)
}

func (ws *Wscoem) SolicitarCambioFechas(IdentificadorCaratula string, params *CambioFechasParams) (string, error) {
	request := &wscoem.SolicitarCambioFechas{
		ArgWSAutenticacionEmpresa: ws.auth,
		ArgSolicitarCambioFechas: &wscoem.SolicitarCambioFechasRequest{
			IdentificadorCaratula: IdentificadorCaratula,
			FechaArribo:           soap.CreateXsdDateTime(params.FechaEstimadaArribo, true),
			FechaZarpada:          soap.CreateXsdDateTime(params.FechaEstimadaZarpada, true),
			CodigoMotivo:          params.CodigoMotivo,
			DescripcionMotivo:     params.DescripcionMotivo,
		},
	}
	PrintlnAsXML(request)

	if ws.environment == TESTING {
		requestXml, err := xml.MarshalIndent(request, " ", "  ")
		if err != nil {
			return "", err
		}
		fmt.Println(string(requestXml))
	}

	client := soap.NewClient(ws.url, soap.WithTLS(&tls.Config{InsecureSkipVerify: true}))
	service := wscoem.NewWgescomunicacionembarqueSoap(client)

	response, err := service.SolicitarCambioFechas(request)
	if err != nil {
		return "", err
	}
	PrintlnAsXML(response)

	var errs []error
	result := ""
	for _, e := range response.SolicitarCambioFechasResult.ListaErrores.DetalleError {
		if *e.Codigo != 0 {
			errs = append(errs, fmt.Errorf("%d - %s - %s", *e.Codigo, e.Descripcion, e.DescripcionAdicional))
		} else {
			result = strings.TrimSpace(strings.Replace(e.DescripcionAdicional, "Identificador:", "", -1))
		}
	}

	return result, errors.Join(errs...)
}

func (ws *Wscoem) SolicitarCambioLOT(IdentificadorCaratula string, params *CambioLOTParams) (string, error) {
	request := &wscoem.SolicitarCambioLOT{
		ArgWSAutenticacionEmpresa: ws.auth,
		ArgSolicitarCambioLOT: &wscoem.SolicitarCambioLOTRequest{
			IdentificadorCaratula: IdentificadorCaratula,
			CodigoLugarOperativo:  params.CodigoLugarOperativo,
		},
	}
	PrintlnAsXML(request)

	client := soap.NewClient(ws.url, soap.WithTLS(&tls.Config{InsecureSkipVerify: true}))
	service := wscoem.NewWgescomunicacionembarqueSoap(client)

	response, err := service.SolicitarCambioLOT(request)
	if err != nil {
		return "", err
	}
	PrintlnAsXML(response)

	var errs []error
	result := ""
	for _, e := range response.SolicitarCambioLOTResult.ListaErrores.DetalleError {
		if *e.Codigo != 0 {
			errs = append(errs, fmt.Errorf("%d - %s - %s", *e.Codigo, e.Descripcion, e.DescripcionAdicional))
		} else {
			result = strings.TrimSpace(strings.Replace(e.DescripcionAdicional, "Identificador:", "", -1))
		}
	}

	return result, errors.Join(errs...)
}

func (ws *Wscoem) RegistrarCOEM(IdentificadorCaratula string, params *COEMParams) (string, error) {
	var contenedoresCarga []*wscoem.ContenedorCarga
	var contenedoresVacios []*wscoem.ContenedorVacio
	var mercaderiasSueltas []*wscoem.MercaderiaSuelta

	if params.ContenedoresCarga != nil {
		for _, v := range *params.ContenedoresCarga {
			var precintos []*wscoem.Precinto
			for _, p := range v.Precintos {
				precintos = append(precintos, &wscoem.Precinto{
					IdentificadorPrecinto: p,
				})
			}
			var declaraciones []*wscoem.Declaracion
			for _, d := range v.Declaraciones {
				declaraciones = append(declaraciones, &wscoem.Declaracion{
					IdentificadorDeclaracion: d,
				})
			}
			contenedoresCarga = append(contenedoresCarga, &wscoem.ContenedorCarga{
				IdentificadorContenedor: v.IdentificadorContenedor,
				CuitATA:                 v.CuitATA,
				Tipo:                    v.Tipo,
				Peso:                    v.Peso,
				Precintos: &wscoem.ArrayOfPrecinto{
					Precinto: precintos,
				},
				Declaraciones: &wscoem.ArrayOfDeclaracion{
					Declaracion: declaraciones,
				},
			})
		}
	}

	if params.ContenedoresVacios != nil {
		for _, v := range *params.ContenedoresVacios {
			contenedoresVacios = append(contenedoresVacios, &wscoem.ContenedorVacio{
				IdentificadorContenedor: v.IdentificadorContenedor,
				Tipo:                    v.Tipo,
				CuitATA:                 v.CuitATA,
				CodigoPais:              v.CodigoPais,
			})
		}
	}

	if params.MercaderiasSueltas != nil {
		for _, v := range *params.MercaderiasSueltas {
			var embalajes []*wscoem.Embalaje
			for _, e := range *v.Embalajes {
				embalajes = append(embalajes, &wscoem.Embalaje{
					CodigoEmbalaje: e.CodigoEmbalaje,
					Peso:           e.Peso,
					CantidadBultos: e.CantidadBultos,
				})
			}
			mercaderiasSueltas = append(mercaderiasSueltas, &wscoem.MercaderiaSuelta{
				IdentificadorDeclaracion: v.IdentificadorDeclaracion,
				CuitATA:                  v.CuitATA,
				Embalajes: &wscoem.ArrayOfEmbalaje{
					Embalaje: embalajes,
				},
			})
		}
	}

	request := &wscoem.RegistrarCOEM{
		ArgWSAutenticacionEmpresa: ws.auth,
		ArgRegistrarCOEM: &wscoem.RegistrarCOEMRequest{
			IdentificadorCaratula: IdentificadorCaratula,
			Coem: &wscoem.Coem{
				ContenedoresConCarga: &wscoem.ArrayOfContenedorCarga{
					ContenedorCarga: contenedoresCarga,
				},
				ContenedoresVacios: &wscoem.ArrayOfContenedorVacio{
					ContenedorVacio: contenedoresVacios,
				},
				MercaderiasSueltas: &wscoem.ArrayOfMercaderiaSuelta{
					MercaderiaSuelta: mercaderiasSueltas,
				},
			},
		},
	}
	PrintlnAsXML(request)

	client := soap.NewClient(ws.url, soap.WithTLS(&tls.Config{InsecureSkipVerify: true}))
	service := wscoem.NewWgescomunicacionembarqueSoap(client)

	response, err := service.RegistrarCOEM(request)
	if err != nil {
		return "", err
	}
	PrintlnAsXML(response)

	var errs []error
	result := ""
	for _, e := range response.RegistrarCOEMResult.ListaErrores.DetalleError {
		if *e.Codigo != 0 {
			errs = append(errs, fmt.Errorf("%d - %s - %s", *e.Codigo, e.Descripcion, e.DescripcionAdicional))
		} else {
			result = strings.TrimSpace(strings.Replace(e.DescripcionAdicional, "Identificador:", "", -1))
		}
	}

	return result, errors.Join(errs...)
}
