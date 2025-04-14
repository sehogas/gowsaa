package afip

import (
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/hooklift/gowsdl/soap"
	"github.com/sehogas/gowsaa/ws/wgestabref"
)

type Wsgestabref struct {
	serviceName string
	environment Environment
	url         string
	cuit        int64
	tipoAgente  string
	rol         string
	validate    *validator.Validate
}

func Newgestabref(environment Environment, cuit int64, tipoAgente, rol string) (*Wsgestabref, error) {
	var url string
	if environment == PRODUCTION {
		url = URLWGESTABREFProduction
	} else {
		url = URLWGESTABREFTesting
	}

	return &Wsgestabref{
		serviceName: "wGesTabRef",
		environment: environment,
		url:         url,
		cuit:        cuit,
		tipoAgente:  tipoAgente,
		rol:         rol,
		validate:    validator.New(validator.WithRequiredStructEnabled()),
	}, nil
}

func (ws *Wsgestabref) PrintlnAsXML(obj interface{}) {
	if ws.environment == TESTING {
		data, err := xml.MarshalIndent(obj, " ", "  ")
		if err == nil {
			fmt.Println(string(data))
		}
	}
}

func (ws *Wsgestabref) Dummy() (appServer, authServer, DbServer string, err error) {
	request := &wgestabref.Dummy{}
	PrintlnAsXML(request)

	client := soap.NewClient(ws.url, soap.WithTLS(&tls.Config{InsecureSkipVerify: true}))
	service := wgestabref.NewWgesTabRefSoap(client)

	response, err := service.Dummy(request)
	if err != nil {
		return "", "", "", err
	}

	PrintlnAsXML(response)

	if response.DummyResult != nil {
		return response.DummyResult.Appserver, response.DummyResult.Authserver, response.DummyResult.Dbserver, nil
	}

	return "", "", "", nil
}

func (ws *Wsgestabref) ConsultarFechaUltAct(argNombreTabla string) (*time.Time, error) {
	ticket, err := GetTA(ws.environment, ws.serviceName, ws.cuit)
	if err != nil {
		return nil, err
	}

	request := &wgestabref.ConsultarFechaUltAct{
		Autentica: &wgestabref.Autenticacion{
			Cuit:       strconv.FormatInt(ws.cuit, 10),
			TipoAgente: ws.tipoAgente,
			Rol:        ws.rol,
			Token:      &ticket.Token,
			Sign:       &ticket.Sign,
		},
		IdReferencia: argNombreTabla,
	}

	PrintlnAsXML(request)

	client := soap.NewClient(ws.url, soap.WithTLS(&tls.Config{InsecureSkipVerify: true}))
	service := wgestabref.NewWgesTabRefSoap(client)

	response, err := service.ConsultarFechaUltAct(request)
	if err != nil {
		return nil, err
	}
	PrintlnAsXML(response)

	var errServ error
	var fecUltAct time.Time
	if response.ConsultarFechaUltActResult != nil {
		if response.ConsultarFechaUltActResult.CodError != 0 {
			errServ = fmt.Errorf("%d - %s", response.ConsultarFechaUltActResult.CodError, response.ConsultarFechaUltActResult.InfoAdicional)
		}

		fecUltAct = response.ConsultarFechaUltActResult.Fecha.ToGoTime()
	}

	return &fecUltAct, errServ
}

func (ws *Wsgestabref) ListaArancel(argNombreTabla string) ([]*wgestabref.Opcion, error) {
	ticket, err := GetTA(ws.environment, ws.serviceName, ws.cuit)
	if err != nil {
		return nil, err
	}

	request := &wgestabref.ListaArancel{
		Autentica: &wgestabref.Autenticacion{
			Cuit:       strconv.FormatInt(ws.cuit, 10),
			TipoAgente: ws.tipoAgente,
			Rol:        ws.rol,
			Token:      &ticket.Token,
			Sign:       &ticket.Sign,
		},
		IdReferencia: argNombreTabla,
	}

	PrintlnAsXML(request)

	client := soap.NewClient(ws.url, soap.WithTLS(&tls.Config{InsecureSkipVerify: true}))
	service := wgestabref.NewWgesTabRefSoap(client)

	response, err := service.ListaArancel(request)
	if err != nil {
		return nil, err
	}
	PrintlnAsXML(response)

	var errServ error
	var lista []*wgestabref.Opcion = make([]*wgestabref.Opcion, 0)
	if response.ListaArancelResult != nil {
		if response.ListaArancelResult.CodError != 0 {
			errServ = fmt.Errorf("%d - %s", response.ListaArancelResult.CodError, response.ListaArancelResult.InfoAdicional)
		}
	}
	if response.ListaArancelResult.Opciones != nil {
		lista = response.ListaArancelResult.Opciones.Opcion
	}

	return lista, errServ
}

func (ws *Wsgestabref) ListaDescripcion(argNombreTabla string) ([]*wgestabref.Descripcion, error) {
	ticket, err := GetTA(ws.environment, ws.serviceName, ws.cuit)
	if err != nil {
		return nil, err
	}

	request := &wgestabref.ListaDescripcion{
		Autentica: &wgestabref.Autenticacion{
			Cuit:       strconv.FormatInt(ws.cuit, 10),
			TipoAgente: ws.tipoAgente,
			Rol:        ws.rol,
			Token:      &ticket.Token,
			Sign:       &ticket.Sign,
		},
		IdReferencia: argNombreTabla,
	}

	PrintlnAsXML(request)

	client := soap.NewClient(ws.url, soap.WithTLS(&tls.Config{InsecureSkipVerify: true}))
	service := wgestabref.NewWgesTabRefSoap(client)

	response, err := service.ListaDescripcion(request)
	if err != nil {
		return nil, err
	}
	PrintlnAsXML(response)

	var errServ error
	var lista []*wgestabref.Descripcion = make([]*wgestabref.Descripcion, 0)
	if response.ListaDescripcionResult != nil {
		if response.ListaDescripcionResult.CodError != 0 {
			errServ = fmt.Errorf("%d - %s", response.ListaDescripcionResult.CodError, response.ListaDescripcionResult.InfoAdicional)
		}
	}
	if response.ListaDescripcionResult.Descripciones != nil {
		lista = response.ListaDescripcionResult.Descripciones.Descripcion
	}
	return lista, errServ
}

func (ws *Wsgestabref) ListaDescripcionDecodificacion(argNombreTabla string) ([]*wgestabref.DescripcionCodificacion, error) {
	ticket, err := GetTA(ws.environment, ws.serviceName, ws.cuit)
	if err != nil {
		return nil, err
	}

	request := &wgestabref.ListaDescripcionDecodificacion{
		Autentica: &wgestabref.Autenticacion{
			Cuit:       strconv.FormatInt(ws.cuit, 10),
			TipoAgente: ws.tipoAgente,
			Rol:        ws.rol,
			Token:      &ticket.Token,
			Sign:       &ticket.Sign,
		},
		IdReferencia: argNombreTabla,
	}

	PrintlnAsXML(request)

	client := soap.NewClient(ws.url, soap.WithTLS(&tls.Config{InsecureSkipVerify: true}))
	service := wgestabref.NewWgesTabRefSoap(client)

	response, err := service.ListaDescripcionDecodificacion(request)
	if err != nil {
		return nil, err
	}
	PrintlnAsXML(response)

	var errServ error
	var lista []*wgestabref.DescripcionCodificacion = make([]*wgestabref.DescripcionCodificacion, 0)
	if response.ListaDescripcionDecodificacionResult != nil {
		if response.ListaDescripcionDecodificacionResult.CodError != 0 {
			errServ = fmt.Errorf("%d - %s", response.ListaDescripcionDecodificacionResult.CodError, response.ListaDescripcionDecodificacionResult.InfoAdicional)
		}
	}
	if response.ListaDescripcionDecodificacionResult.DescripcionesCodificaciones != nil {
		lista = response.ListaDescripcionDecodificacionResult.DescripcionesCodificaciones.DescripcionCodificacion
	}
	return lista, errServ
}

func (ws *Wsgestabref) ListaEmpresas(argNombreTabla string) ([]*wgestabref.Empresa, error) {
	ticket, err := GetTA(ws.environment, ws.serviceName, ws.cuit)
	if err != nil {
		return nil, err
	}

	request := &wgestabref.ListaEmpresas{
		Autentica: &wgestabref.Autenticacion{
			Cuit:       strconv.FormatInt(ws.cuit, 10),
			TipoAgente: ws.tipoAgente,
			Rol:        ws.rol,
			Token:      &ticket.Token,
			Sign:       &ticket.Sign,
		},
		IdReferencia: argNombreTabla,
	}

	PrintlnAsXML(request)

	client := soap.NewClient(ws.url, soap.WithTLS(&tls.Config{InsecureSkipVerify: true}))
	service := wgestabref.NewWgesTabRefSoap(client)

	response, err := service.ListaEmpresas(request)
	if err != nil {
		return nil, err
	}
	PrintlnAsXML(response)

	var errServ error
	var lista []*wgestabref.Empresa = make([]*wgestabref.Empresa, 0)
	if response.ListaEmpresasResult != nil {
		if response.ListaEmpresasResult.CodError != 0 {
			errServ = fmt.Errorf("%d - %s", response.ListaEmpresasResult.CodError, response.ListaEmpresasResult.InfoAdicional)
		}
	}
	if response.ListaEmpresasResult.Empresas != nil {
		lista = response.ListaEmpresasResult.Empresas.Empresa
	}
	return lista, errServ
}

func (ws *Wsgestabref) ListaListaLugaresOperativos(argNombreTabla string) ([]*wgestabref.LugarOperativo, error) {
	ticket, err := GetTA(ws.environment, ws.serviceName, ws.cuit)
	if err != nil {
		return nil, err
	}

	request := &wgestabref.ListaLugaresOperativos{
		Autentica: &wgestabref.Autenticacion{
			Cuit:       strconv.FormatInt(ws.cuit, 10),
			TipoAgente: ws.tipoAgente,
			Rol:        ws.rol,
			Token:      &ticket.Token,
			Sign:       &ticket.Sign,
		},
		IdReferencia: argNombreTabla,
	}

	PrintlnAsXML(request)

	client := soap.NewClient(ws.url, soap.WithTLS(&tls.Config{InsecureSkipVerify: true}))
	service := wgestabref.NewWgesTabRefSoap(client)

	response, err := service.ListaLugaresOperativos(request)
	if err != nil {
		return nil, err
	}
	PrintlnAsXML(response)

	var errServ error
	var lista []*wgestabref.LugarOperativo = make([]*wgestabref.LugarOperativo, 0)
	if response.ListaLugaresOperativosResult != nil {
		if response.ListaLugaresOperativosResult.CodError != 0 {
			errServ = fmt.Errorf("%d - %s", response.ListaLugaresOperativosResult.CodError, response.ListaLugaresOperativosResult.InfoAdicional)
		}
	}
	if response.ListaLugaresOperativosResult.LugaresOperativos != nil {
		lista = response.ListaLugaresOperativosResult.LugaresOperativos.LugarOperativo
	}
	return lista, errServ
}

func (ws *Wsgestabref) ListaPaisesAduanas(argNombreTabla string) ([]*wgestabref.PaisAduana, error) {
	ticket, err := GetTA(ws.environment, ws.serviceName, ws.cuit)
	if err != nil {
		return nil, err
	}

	request := &wgestabref.ListaPaisesAduanas{
		Autentica: &wgestabref.Autenticacion{
			Cuit:       strconv.FormatInt(ws.cuit, 10),
			TipoAgente: ws.tipoAgente,
			Rol:        ws.rol,
			Token:      &ticket.Token,
			Sign:       &ticket.Sign,
		},
		IdReferencia: argNombreTabla,
	}

	PrintlnAsXML(request)

	client := soap.NewClient(ws.url, soap.WithTLS(&tls.Config{InsecureSkipVerify: true}))
	service := wgestabref.NewWgesTabRefSoap(client)

	response, err := service.ListaPaisesAduanas(request)
	if err != nil {
		return nil, err
	}
	PrintlnAsXML(response)

	var errServ error
	var lista []*wgestabref.PaisAduana = make([]*wgestabref.PaisAduana, 0)
	if response.ListaPaisesAduanasResult != nil {
		if response.ListaPaisesAduanasResult.CodError != 0 {
			errServ = fmt.Errorf("%d - %s", response.ListaPaisesAduanasResult.CodError, response.ListaPaisesAduanasResult.InfoAdicional)
		}
	}
	if response.ListaPaisesAduanasResult.PaisesAduanas != nil {
		lista = response.ListaPaisesAduanasResult.PaisesAduanas.PaisAduana
	}
	return lista, errServ
}

func (ws *Wsgestabref) ListaTablasReferencia() ([]*wgestabref.TablaReferencia, error) {
	ticket, err := GetTA(ws.environment, ws.serviceName, ws.cuit)
	if err != nil {
		return nil, err
	}

	request := &wgestabref.ListaTablasReferencia{
		Autentica: &wgestabref.Autenticacion{
			Cuit:       strconv.FormatInt(ws.cuit, 10),
			TipoAgente: ws.tipoAgente,
			Rol:        ws.rol,
			Token:      &ticket.Token,
			Sign:       &ticket.Sign,
		},
	}

	PrintlnAsXML(request)

	client := soap.NewClient(ws.url, soap.WithTLS(&tls.Config{InsecureSkipVerify: true}))
	service := wgestabref.NewWgesTabRefSoap(client)

	response, err := service.ListaTablasReferencia(request)
	if err != nil {
		return nil, err
	}
	PrintlnAsXML(response)

	var errServ error
	var lista []*wgestabref.TablaReferencia
	if response.ListaTablasReferenciaResult != nil {
		if response.ListaTablasReferenciaResult.CodError != 0 {
			errServ = fmt.Errorf("%d - %s", response.ListaTablasReferenciaResult.CodError, response.ListaTablasReferenciaResult.InfoAdicional)
		}
	}
	if response.ListaTablasReferenciaResult.TablasReferencia != nil {
		lista = response.ListaTablasReferenciaResult.TablasReferencia.TablaReferencia
	}
	return lista, errServ
}
