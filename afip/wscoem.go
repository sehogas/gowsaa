package afip

import (
	"crypto/tls"
	"encoding/xml"
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/hooklift/gowsdl/soap"
	validador "github.com/sehogas/gowsaa/internal/util/validator"
	"github.com/sehogas/gowsaa/ws/wscoem"
)

type Wscoem struct {
	serviceName string
	environment Environment
	url         string
	cuit        int64
	tipoAgente  string
	rol         string
	validate    *validator.Validate
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
		validate:    validator.New(validator.WithRequiredStructEnabled()),
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

func (ws *Wscoem) Dummy() (appServer, authServer, DbServer string, err error) {
	request := &wscoem.Dummy{}
	PrintlnAsXML(request)

	client := soap.NewClient(ws.url, soap.WithTLS(&tls.Config{InsecureSkipVerify: true}))
	service := wscoem.NewWgescomunicacionembarqueSoap(client)

	response, err := service.Dummy(request)
	if err != nil {
		return "", "", "", err
	}

	PrintlnAsXML(response)

	if response.DummyResult != nil {
		return response.DummyResult.AppServer, response.DummyResult.AuthServer, response.DummyResult.DbServer, nil
	}

	return "", "", "", nil
}

func (ws *Wscoem) RegistrarCaratula(params *CaratulaParams) (string, error) {
	if err := ws.validate.Struct(params); err != nil {
		return "", errors.New(strings.Join(validador.ToErrResponse(err).Errors, ", "))
	}

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
				FechaArribo:           soap.CreateXsdDateTime(params.FechaArribo, true),
				FechaZarpada:          soap.CreateXsdDateTime(params.FechaZarpada, true),
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

func (ws *Wscoem) AnularCaratula(params *IdentificadorCaraturaParams) (string, error) {
	if err := ws.validate.Struct(params); err != nil {
		return "", errors.New(strings.Join(validador.ToErrResponse(err).Errors, ", "))
	}

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
			IdentificadorCaratula: params.IdentificadorCaratula,
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

func (ws *Wscoem) RectificarCaratula(params *RectificarCaratulaParams) (string, error) {
	if err := ws.validate.Struct(params); err != nil {
		return "", errors.New(strings.Join(validador.ToErrResponse(err).Errors, ", "))
	}

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
			IdentificadorCaratula: params.IdentificadorCaratula,
			Caratula: &wscoem.Caratula{
				IdentificadorBuque:    params.IdentificadorBuque,
				NombreMedioTransporte: params.NombreMedioTransporte,
				CodigoAduana:          params.CodigoAduana,
				CodigoLugarOperativo:  params.CodigoLugarOperativo,
				FechaArribo:           soap.CreateXsdDateTime(params.FechaArribo, true),
				FechaZarpada:          soap.CreateXsdDateTime(params.FechaZarpada, true),
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

func (ws *Wscoem) SolicitarCambioBuque(params *CambioBuqueParams) (string, error) {
	if err := ws.validate.Struct(params); err != nil {
		return "", errors.New(strings.Join(validador.ToErrResponse(err).Errors, ", "))
	}

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
			IdentificadorCaratula: params.IdentificadorCaratula,
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

func (ws *Wscoem) SolicitarCambioFechas(params *CambioFechasParams) (string, error) {
	if err := ws.validate.Struct(params); err != nil {
		return "", errors.New(strings.Join(validador.ToErrResponse(err).Errors, ", "))
	}

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
			IdentificadorCaratula: params.IdentificadorCaratula,
			FechaArribo:           soap.CreateXsdDateTime(params.FechaArribo, true),
			FechaZarpada:          soap.CreateXsdDateTime(params.FechaZarpada, true),
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

func (ws *Wscoem) SolicitarCambioLOT(params *CambioLOTParams) (string, error) {
	if err := ws.validate.Struct(params); err != nil {
		return "", errors.New(strings.Join(validador.ToErrResponse(err).Errors, ", "))
	}

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
			IdentificadorCaratula: params.IdentificadorCaratula,
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

func (ws *Wscoem) RegistrarCOEM(params *RegistrarCOEMParams) (string, error) {
	if err := ws.validate.Struct(params); err != nil {
		return "", errors.New(strings.Join(validador.ToErrResponse(err).Errors, ", "))
	}

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
			IdentificadorCaratula: params.IdentificadorCaratula,
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

func (ws *Wscoem) RectificarCOEM(params *RectificarCOEMParams) (string, error) {
	if err := ws.validate.Struct(params); err != nil {
		return "", errors.New(strings.Join(validador.ToErrResponse(err).Errors, ", "))
	}

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
			IdentificadorCaratula: params.IdentificadorCaratula,
			IdentificadorCOEM:     params.IdentificadorCOEM,
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

func (ws *Wscoem) CerrarCOEM(params *IdentificadorCOEMParams) (string, error) {
	if err := ws.validate.Struct(params); err != nil {
		return "", errors.New(strings.Join(validador.ToErrResponse(err).Errors, ", "))
	}

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
			IdentificadorCaratula: params.IdentificadorCaratula,
			IdentificadorCOEM:     params.IdentificadorCOEM,
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

func (ws *Wscoem) AnularCOEM(params *IdentificadorCOEMParams) (string, error) {
	if err := ws.validate.Struct(params); err != nil {
		return "", errors.New(strings.Join(validador.ToErrResponse(err).Errors, ", "))
	}

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
			IdentificadorCaratula: params.IdentificadorCaratula,
			IdentificadorCOEM:     params.IdentificadorCOEM,
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

func (ws *Wscoem) SolicitarAnulacionCOEM(params *IdentificadorCOEMParams) (string, error) {
	if err := ws.validate.Struct(params); err != nil {
		return "", errors.New(strings.Join(validador.ToErrResponse(err).Errors, ", "))
	}

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
			IdentificadorCaratula: params.IdentificadorCaratula,
			IdentificadorCOEM:     params.IdentificadorCOEM,
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

func (ws *Wscoem) SolicitarNoABordo(params *SolicitarNoABordoParams) (string, error) {
	if err := ws.validate.Struct(params); err != nil {
		return "", errors.New(strings.Join(validador.ToErrResponse(err).Errors, ", "))
	}

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
			IdentificadorCaratula: params.IdentificadorCaratula,
			IdentificadorCOEM:     params.IdentificadorCOEM,
			CodigoMotivo:          params.CodigoMotivo,
			DescripcionMotivo:     params.DescripcionMotivo,
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

func (ws *Wscoem) SolicitarCierreCargaContoBulto(params *SolicitarCierreCargaContoBultoParams) (string, error) {
	if err := ws.validate.Struct(params); err != nil {
		return "", errors.New(strings.Join(validador.ToErrResponse(err).Errors, ", "))
	}

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
			IdentificadorCaratula: params.IdentificadorCaratula,
			FechaZarpada:          soap.CreateXsdDateTime(params.FechaZarpada, true),
			NumeroViaje:           params.NumeroViaje,
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

func (ws *Wscoem) SolicitarCierreCargaGranel(params *SolicitarCierreCargaGranelParams) (string, error) {
	if err := ws.validate.Struct(params); err != nil {
		return "", errors.New(strings.Join(validador.ToErrResponse(err).Errors, ", "))
	}

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
			IdentificadorCaratula: params.IdentificadorCaratula,
			FechaZarpada:          soap.CreateXsdDateTime(params.FechaZarpada, true),
			NumeroViaje:           params.NumeroViaje,
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
