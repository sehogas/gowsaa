package afip

import (
	"crypto/tls"
	"encoding/xml"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hooklift/gowsdl/soap"
	"github.com/sehogas/gowsaa/ws/wscoem"
)

type Wscoem struct {
	serviceName string
	environment Environment
	url         string
	cuit        int64
	tipoAgente  string
	rol         string
}

func NewWscoem(environment Environment, cuit int64, tipoAgente, rol string) (*Wscoem, error) {
	var url string
	if environment == PRODUCTION {
		url = URLWSCOEMProduction
	} else {
		url = URLWSCOEMTesting
	}

	return &Wscoem{
		serviceName: "wgescomunicacionembarque",
		environment: environment,
		url:         url,
		cuit:        cuit,
		tipoAgente:  tipoAgente,
		rol:         rol,
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
	ticket, err := GetTA(ws.environment, ws.serviceName, ws.cuit)
	if err != nil {
		return "", err
	}
	request := &wscoem.RegistrarCaratula{
		ArgWSAutenticacionEmpresa: &wscoem.WSAutenticacionEmpresa{
			WSAutenticacion: &wscoem.WSAutenticacion{
				Token: ticket.Token,
				Sign:  ticket.Sign,
			},
			CuitEmpresaConectada: ws.cuit,
			TipoAgente:           ws.tipoAgente,
			Rol:                  ws.rol,
		},
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

func (ws *Wscoem) AnularCaratula(identificadorCaratula string) (string, error) {
	ticket, err := GetTA(ws.environment, ws.serviceName, ws.cuit)
	if err != nil {
		return "", err
	}

	request := &wscoem.AnularCaratula{
		ArgWSAutenticacionEmpresa: &wscoem.WSAutenticacionEmpresa{
			WSAutenticacion: &wscoem.WSAutenticacion{
				Token: ticket.Token,
				Sign:  ticket.Sign,
			},
			CuitEmpresaConectada: ws.cuit,
			TipoAgente:           ws.tipoAgente,
			Rol:                  ws.rol,
		},

		ArgAnularCaratula: &wscoem.AnularCaratulaRequest{
			IdentificadorCaratula: identificadorCaratula,
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

func (ws *Wscoem) RectificarCaratula(identificadorCaratula string, params *CaratulaParams) (string, error) {
	ticket, err := GetTA(ws.environment, ws.serviceName, ws.cuit)
	if err != nil {
		return "", err
	}

	request := &wscoem.RectificarCaratula{
		ArgWSAutenticacionEmpresa: &wscoem.WSAutenticacionEmpresa{
			WSAutenticacion: &wscoem.WSAutenticacion{
				Token: ticket.Token,
				Sign:  ticket.Sign,
			},
			CuitEmpresaConectada: ws.cuit,
			TipoAgente:           ws.tipoAgente,
			Rol:                  ws.rol,
		},
		ArgRectificarCaratula: &wscoem.RectificarCaratulaRequest{
			IdentificadorCaratula: identificadorCaratula,
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

func (ws *Wscoem) SolicitarCambioBuque(identificadorCaratula string, params *CambioBuqueParams) (string, error) {
	ticket, err := GetTA(ws.environment, ws.serviceName, ws.cuit)
	if err != nil {
		return "", err
	}

	request := &wscoem.SolicitarCambioBuque{
		ArgWSAutenticacionEmpresa: &wscoem.WSAutenticacionEmpresa{
			WSAutenticacion: &wscoem.WSAutenticacion{
				Token: ticket.Token,
				Sign:  ticket.Sign,
			},
			CuitEmpresaConectada: ws.cuit,
			TipoAgente:           ws.tipoAgente,
			Rol:                  ws.rol,
		},
		ArgSolicitarCambioBuque: &wscoem.SolicitarCambioBuqueRequest{
			IdentificadorCaratula: identificadorCaratula,
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

func (ws *Wscoem) SolicitarCambioFechas(identificadorCaratula string, params *CambioFechasParams) (string, error) {
	ticket, err := GetTA(ws.environment, ws.serviceName, ws.cuit)
	if err != nil {
		return "", err
	}

	request := &wscoem.SolicitarCambioFechas{
		ArgWSAutenticacionEmpresa: &wscoem.WSAutenticacionEmpresa{
			WSAutenticacion: &wscoem.WSAutenticacion{
				Token: ticket.Token,
				Sign:  ticket.Sign,
			},
			CuitEmpresaConectada: ws.cuit,
			TipoAgente:           ws.tipoAgente,
			Rol:                  ws.rol,
		},
		ArgSolicitarCambioFechas: &wscoem.SolicitarCambioFechasRequest{
			IdentificadorCaratula: identificadorCaratula,
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

func (ws *Wscoem) SolicitarCambioLOT(identificadorCaratula string, params *CambioLOTParams) (string, error) {
	ticket, err := GetTA(ws.environment, ws.serviceName, ws.cuit)
	if err != nil {
		return "", err
	}

	request := &wscoem.SolicitarCambioLOT{
		ArgWSAutenticacionEmpresa: &wscoem.WSAutenticacionEmpresa{
			WSAutenticacion: &wscoem.WSAutenticacion{
				Token: ticket.Token,
				Sign:  ticket.Sign,
			},
			CuitEmpresaConectada: ws.cuit,
			TipoAgente:           ws.tipoAgente,
			Rol:                  ws.rol,
		},
		ArgSolicitarCambioLOT: &wscoem.SolicitarCambioLOTRequest{
			IdentificadorCaratula: identificadorCaratula,
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

func (ws *Wscoem) RegistrarCOEM(identificadorCaratula string, params *COEMParams) (string, error) {
	ticket, err := GetTA(ws.environment, ws.serviceName, ws.cuit)
	if err != nil {
		return "", err
	}

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
		ArgWSAutenticacionEmpresa: &wscoem.WSAutenticacionEmpresa{
			WSAutenticacion: &wscoem.WSAutenticacion{
				Token: ticket.Token,
				Sign:  ticket.Sign,
			},
			CuitEmpresaConectada: ws.cuit,
			TipoAgente:           ws.tipoAgente,
			Rol:                  ws.rol,
		},
		ArgRegistrarCOEM: &wscoem.RegistrarCOEMRequest{
			IdentificadorCaratula: identificadorCaratula,
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

func (ws *Wscoem) RectificarCOEM(identificadorCaratula string, identificadorCOEM string, params *COEMParams) (string, error) {
	ticket, err := GetTA(ws.environment, ws.serviceName, ws.cuit)
	if err != nil {
		return "", err
	}

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

	request := &wscoem.RectificarCOEM{
		ArgWSAutenticacionEmpresa: &wscoem.WSAutenticacionEmpresa{
			WSAutenticacion: &wscoem.WSAutenticacion{
				Token: ticket.Token,
				Sign:  ticket.Sign,
			},
			CuitEmpresaConectada: ws.cuit,
			TipoAgente:           ws.tipoAgente,
			Rol:                  ws.rol,
		},
		ArgRectificarCOEM: &wscoem.RectificarCOEMRequest{
			IdentificadorCaratula: identificadorCaratula,
			IdentificadorCOEM:     identificadorCOEM,
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

	response, err := service.RectificarCOEM(request)
	if err != nil {
		return "", err
	}
	PrintlnAsXML(response)

	var errs []error
	result := ""
	for _, e := range response.RectificarCOEMResult.ListaErrores.DetalleError {
		if *e.Codigo != 0 {
			errs = append(errs, fmt.Errorf("%d - %s - %s", *e.Codigo, e.Descripcion, e.DescripcionAdicional))
		} else {
			result = strings.TrimSpace(strings.Replace(e.DescripcionAdicional, "Identificador:", "", -1))
		}
	}

	return result, errors.Join(errs...)
}

func (ws *Wscoem) CerrarCOEM(identificadorCaratula string, identificadorCOEM string) (string, error) {
	ticket, err := GetTA(ws.environment, ws.serviceName, ws.cuit)
	if err != nil {
		return "", err
	}

	request := &wscoem.CerrarCOEM{
		ArgWSAutenticacionEmpresa: &wscoem.WSAutenticacionEmpresa{
			WSAutenticacion: &wscoem.WSAutenticacion{
				Token: ticket.Token,
				Sign:  ticket.Sign,
			},
			CuitEmpresaConectada: ws.cuit,
			TipoAgente:           ws.tipoAgente,
			Rol:                  ws.rol,
		},
		ArgCerrarCOEM: &wscoem.CerrarCOEMRequest{
			IdentificadorCaratula: identificadorCaratula,
			IdentificadorCOEM:     identificadorCOEM,
		},
	}
	PrintlnAsXML(request)

	client := soap.NewClient(ws.url, soap.WithTLS(&tls.Config{InsecureSkipVerify: true}))
	service := wscoem.NewWgescomunicacionembarqueSoap(client)

	response, err := service.CerrarCOEM(request)
	if err != nil {
		return "", err
	}
	PrintlnAsXML(response)

	var errs []error
	result := ""
	for _, e := range response.CerrarCOEMResult.ListaErrores.DetalleError {
		if *e.Codigo != 0 {
			errs = append(errs, fmt.Errorf("%d - %s - %s", *e.Codigo, e.Descripcion, e.DescripcionAdicional))
		} else {
			result = strings.TrimSpace(strings.Replace(e.DescripcionAdicional, "Identificador:", "", -1))
		}
	}

	return result, errors.Join(errs...)
}

func (ws *Wscoem) AnularCOEM(identificadorCaratula string, identificadorCOEM string) (string, error) {
	ticket, err := GetTA(ws.environment, ws.serviceName, ws.cuit)
	if err != nil {
		return "", err
	}

	request := &wscoem.AnularCOEM{
		ArgWSAutenticacionEmpresa: &wscoem.WSAutenticacionEmpresa{
			WSAutenticacion: &wscoem.WSAutenticacion{
				Token: ticket.Token,
				Sign:  ticket.Sign,
			},
			CuitEmpresaConectada: ws.cuit,
			TipoAgente:           ws.tipoAgente,
			Rol:                  ws.rol,
		},
		ArgAnularCOEM: &wscoem.AnularCOEMRequest{
			IdentificadorCaratula: identificadorCaratula,
			IdentificadorCOEM:     identificadorCOEM,
		},
	}
	PrintlnAsXML(request)

	client := soap.NewClient(ws.url, soap.WithTLS(&tls.Config{InsecureSkipVerify: true}))
	service := wscoem.NewWgescomunicacionembarqueSoap(client)

	response, err := service.AnularCOEM(request)
	if err != nil {
		return "", err
	}
	PrintlnAsXML(response)

	var errs []error
	result := ""
	for _, e := range response.AnularCOEMResult.ListaErrores.DetalleError {
		if *e.Codigo != 0 {
			errs = append(errs, fmt.Errorf("%d - %s - %s", *e.Codigo, e.Descripcion, e.DescripcionAdicional))
		} else {
			result = strings.TrimSpace(strings.Replace(e.DescripcionAdicional, "Identificador:", "", -1))
		}
	}

	return result, errors.Join(errs...)
}

func (ws *Wscoem) SolicitarAnulacionCOEM(identificadorCaratula string, identificadorCOEM string) (string, error) {
	ticket, err := GetTA(ws.environment, ws.serviceName, ws.cuit)
	if err != nil {
		return "", err
	}

	request := &wscoem.SolicitarAnulacionCOEM{
		ArgWSAutenticacionEmpresa: &wscoem.WSAutenticacionEmpresa{
			WSAutenticacion: &wscoem.WSAutenticacion{
				Token: ticket.Token,
				Sign:  ticket.Sign,
			},
			CuitEmpresaConectada: ws.cuit,
			TipoAgente:           ws.tipoAgente,
			Rol:                  ws.rol,
		},
		ArgSolicitarAnulacionCOEM: &wscoem.SolicitarAnulacionCOEMRequest{
			IdentificadorCaratula: identificadorCaratula,
			IdentificadorCOEM:     identificadorCOEM,
		},
	}
	PrintlnAsXML(request)

	client := soap.NewClient(ws.url, soap.WithTLS(&tls.Config{InsecureSkipVerify: true}))
	service := wscoem.NewWgescomunicacionembarqueSoap(client)

	response, err := service.SolicitarAnulacionCOEM(request)
	if err != nil {
		return "", err
	}
	PrintlnAsXML(response)

	var errs []error
	result := ""
	for _, e := range response.SolicitarAnulacionCOEMResult.ListaErrores.DetalleError {
		if *e.Codigo != 0 {
			errs = append(errs, fmt.Errorf("%d - %s - %s", *e.Codigo, e.Descripcion, e.DescripcionAdicional))
		} else {
			result = strings.TrimSpace(strings.Replace(e.DescripcionAdicional, "Identificador:", "", -1))
		}
	}

	return result, errors.Join(errs...)
}

func (ws *Wscoem) SolicitarNoABordo(identificadorCaratula string, identificadorCOEM string, codigoMotivo, descripcionMotivo string, params *NoABordoParams) (string, error) {
	ticket, err := GetTA(ws.environment, ws.serviceName, ws.cuit)
	if err != nil {
		return "", err
	}

	var contenedoresCarga []*wscoem.Contenedor
	var contenedoresVacios []*wscoem.Contenedor
	var declaraciones []*wscoem.Declaracion

	if params.ContenedoresCarga != nil {
		for _, v := range *params.ContenedoresCarga {
			contenedoresCarga = append(contenedoresCarga, &wscoem.Contenedor{
				IdentificadorContenedor: v.IdentificadorContenedor,
			})
		}
	}

	if params.ContenedoresVacios != nil {
		for _, v := range *params.ContenedoresVacios {
			contenedoresVacios = append(contenedoresVacios, &wscoem.Contenedor{
				IdentificadorContenedor: v.IdentificadorContenedor,
			})
		}
	}

	if params.MercaderiasSueltas != nil {
		for _, v := range *params.MercaderiasSueltas {
			declaraciones = append(declaraciones, &wscoem.Declaracion{
				IdentificadorDeclaracion: v.IdentificadorDeclaracion,
			})
		}
	}

	request := &wscoem.SolicitarNoABordo{
		ArgWSAutenticacionEmpresa: &wscoem.WSAutenticacionEmpresa{
			WSAutenticacion: &wscoem.WSAutenticacion{
				Token: ticket.Token,
				Sign:  ticket.Sign,
			},
			CuitEmpresaConectada: ws.cuit,
			TipoAgente:           ws.tipoAgente,
			Rol:                  ws.rol,
		},
		ArgSolicitarNoABordo: &wscoem.SolicitarNoABordoRequest{
			IdentificadorCaratula: identificadorCaratula,
			IdentificadorCOEM:     identificadorCOEM,
			CodigoMotivo:          codigoMotivo,
			DescripcionMotivo:     descripcionMotivo,
			IdentificadoresContenedoresVacios: &wscoem.ArrayOfContenedor{
				Contenedor: contenedoresVacios,
			},
			IdentificadoresContenedorCarga: &wscoem.ArrayOfContenedor{
				Contenedor: contenedoresCarga,
			},
			IdentificadoresDeclaracionesMercaderiaSuelta: &wscoem.ArrayOfDeclaracion{
				Declaracion: declaraciones,
			},
		},
	}
	PrintlnAsXML(request)

	client := soap.NewClient(ws.url, soap.WithTLS(&tls.Config{InsecureSkipVerify: true}))
	service := wscoem.NewWgescomunicacionembarqueSoap(client)

	response, err := service.SolicitarNoABordo(request)
	if err != nil {
		return "", err
	}
	PrintlnAsXML(response)

	var errs []error
	result := ""
	for _, e := range response.SolicitarNoABordoResult.ListaErrores.DetalleError {
		if *e.Codigo != 0 {
			errs = append(errs, fmt.Errorf("%d - %s - %s", *e.Codigo, e.Descripcion, e.DescripcionAdicional))
		} else {
			result = strings.TrimSpace(strings.Replace(e.DescripcionAdicional, "Identificador:", "", -1))
		}
	}

	return result, errors.Join(errs...)
}

func (ws *Wscoem) SolicitarCierreCargaContoBulto(identificadorCaratula string, fechaZarpada time.Time, numViaje string, params *CierreCargaContoBultoParams) (string, error) {
	ticket, err := GetTA(ws.environment, ws.serviceName, ws.cuit)
	if err != nil {
		return "", err
	}

	var declaraciones []*wscoem.DeclaracionCont

	if params.Declaraciones != nil {
		for _, v := range *params.Declaraciones {
			declaraciones = append(declaraciones, &wscoem.DeclaracionCont{
				IdentificadorDeclaracion: v.IdentificadorDeclaracion,
				FechaEmbarque:            soap.CreateXsdDateTime(v.FechaEmbarque, true),
			})
		}
	}

	request := &wscoem.SolicitarCierreCargaContoBulto{
		ArgWSAutenticacionEmpresa: &wscoem.WSAutenticacionEmpresa{
			WSAutenticacion: &wscoem.WSAutenticacion{
				Token: ticket.Token,
				Sign:  ticket.Sign,
			},
			CuitEmpresaConectada: ws.cuit,
			TipoAgente:           ws.tipoAgente,
			Rol:                  ws.rol,
		},
		ArgSolicitarCierreCargaContoBulto: &wscoem.SolicitarCierreCargaContoBultoRequest{
			IdentificadorCaratula: identificadorCaratula,
			FechaZarpada:          soap.CreateXsdDateTime(fechaZarpada, true),
			NumeroViaje:           numViaje,
			Declaraciones: &wscoem.ArrayOfDeclaracionCont{
				DeclaracionCont: declaraciones,
			},
		},
	}
	PrintlnAsXML(request)

	client := soap.NewClient(ws.url, soap.WithTLS(&tls.Config{InsecureSkipVerify: true}))
	service := wscoem.NewWgescomunicacionembarqueSoap(client)

	response, err := service.SolicitarCierreCargaContoBulto(request)
	if err != nil {
		return "", err
	}
	PrintlnAsXML(response)

	var errs []error
	result := ""
	for _, e := range response.SolicitarCierreCargaContoBultoResult.ListaErrores.DetalleError {
		if *e.Codigo != 0 {
			errs = append(errs, fmt.Errorf("%d - %s - %s", *e.Codigo, e.Descripcion, e.DescripcionAdicional))
		} else {
			result = strings.TrimSpace(strings.Replace(e.DescripcionAdicional, "Identificador:", "", -1))
		}
	}

	return result, errors.Join(errs...)
}

func (ws *Wscoem) SolicitarCierreCargaGranel(identificadorCaratula string, fechaZarpada time.Time, numViaje string, params *CierreCargaGranelParams) (string, error) {
	ticket, err := GetTA(ws.environment, ws.serviceName, ws.cuit)
	if err != nil {
		return "", err
	}

	var declaracionesCoemGranel []*wscoem.CoemGranel

	if params.DeclaracionesCOEMGranel != nil {
		for _, v := range *params.DeclaracionesCOEMGranel {
			var declaracionesGranel []*wscoem.DeclaracionGranel
			for _, d := range *params.DeclaracionesCOEMGranel {
				for _, g := range *d.DeclaracionesGranel {
					var items []*wscoem.Item
					for _, i := range *g.Items {
						items = append(items, &wscoem.Item{
							NumeroItem:   i.NumeroItem,
							CantidadReal: i.CantidadReal,
						})
					}
					declaracionesGranel = append(declaracionesGranel, &wscoem.DeclaracionGranel{
						IdentificadorDeclaracion:    g.IdentificadorDeclaracion,
						FechaEmbarque:               soap.CreateXsdDateTime(g.FechaEmbarque, true),
						IdentificadorCierreCumplido: g.IdentificadorCierreCumplido,
						Items: &wscoem.ArrayOfItem{
							Item: items,
						},
					})
				}
			}
			declaracionesCoemGranel = append(declaracionesCoemGranel, &wscoem.CoemGranel{
				IdentificadorCoem: v.IdentificadorCOEM,
				Declaraciones: &wscoem.ArrayOfDeclaracionGranel{
					DeclaracionGranel: declaracionesGranel,
				},
			})
		}
	}

	request := &wscoem.SolicitarCierreCargaGranel{
		ArgWSAutenticacionEmpresa: &wscoem.WSAutenticacionEmpresa{
			WSAutenticacion: &wscoem.WSAutenticacion{
				Token: ticket.Token,
				Sign:  ticket.Sign,
			},
			CuitEmpresaConectada: ws.cuit,
			TipoAgente:           ws.tipoAgente,
			Rol:                  ws.rol,
		},
		ArgSolicitarCierreCargaGranel: &wscoem.SolicitarCierreCargaGranelRequest{
			IdentificadorCaratula: identificadorCaratula,
			FechaZarpada:          soap.CreateXsdDateTime(fechaZarpada, true),
			NumeroViaje:           numViaje,
			Coems: &wscoem.ArrayOfCoemGranel{
				CoemGranel: declaracionesCoemGranel,
			},
		},
	}
	PrintlnAsXML(request)

	client := soap.NewClient(ws.url, soap.WithTLS(&tls.Config{InsecureSkipVerify: true}))
	service := wscoem.NewWgescomunicacionembarqueSoap(client)

	response, err := service.SolicitarCierreCargaGranel(request)
	if err != nil {
		return "", err
	}
	PrintlnAsXML(response)

	var errs []error
	result := ""
	for _, e := range response.SolicitarCierreCargaGranelResult.ListaErrores.DetalleError {
		if *e.Codigo != 0 {
			errs = append(errs, fmt.Errorf("%d - %s - %s", *e.Codigo, e.Descripcion, e.DescripcionAdicional))
		} else {
			result = strings.TrimSpace(strings.Replace(e.DescripcionAdicional, "Identificador:", "", -1))
		}
	}

	return result, errors.Join(errs...)
}
