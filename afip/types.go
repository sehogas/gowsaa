package afip

import "time"

type PuertoDestinoParams struct {
	PuertoDestino string `json:"PuertoDestino,omitempty" validate:"len=5"`
}

type CaratulaParams struct {
	IdentificadorBuque    string    `json:"IdentificadorBuque,omitempty" validate:"lte=20,startswith=IMO"`
	NombreMedioTransporte string    `json:"NombreMedioTransporte" validate:"required,lte=200"`
	CodigoAduana          string    `json:"CodigoAduana" validate:"required,len=3"`
	CodigoLugarOperativo  string    `json:"CodigoLugarOperativo" validate:"required,len=5"`
	FechaArribo           time.Time `json:"FechaArribo" validate:"required"`
	FechaZarpada          time.Time `json:"FechaZarpada" validate:"required"`
	Via                   string    `json:"Via" validate:"required,len=1"`
	NumeroViaje           string    `json:"NumeroViaje,omitempty" validate:"lte=16"`
	*PuertoDestinoParams
	Itinerario []PuertoDestinoParams `json:"Itinerario,omitempty"`
}

type IdentificadorCaraturaParams struct {
	IdentificadorCaratula string `json:"IdentificadorCaratula" validate:"required,len=16"`
}

type RectificarCaratulaParams struct {
	*IdentificadorCaraturaParams
	*CaratulaParams
}

type CambioBuqueParams struct {
	*IdentificadorCaraturaParams
	IdentificadorBuque    string `json:"IdentificadorBuque,omitempty" validate:"lte=20"`
	NombreMedioTransporte string `json:"NombreMedioTransporte,omitempty" validate:"lte=200"`
}

type CambioFechasParams struct {
	*IdentificadorCaraturaParams
	FechaArribo       time.Time `json:"FechaArribo" validate:"required,datetime"`
	FechaZarpada      time.Time `json:"FechaZarpada" validate:"required,datetime"`
	CodigoMotivo      string    `json:"CodigoMotivo" validate:"required,lte=2"`
	DescripcionMotivo string    `json:"DescripcionMotivo,omitempty" validate:"lte=200"`
}

type CambioLOTParams struct {
	*IdentificadorCaraturaParams
	CodigoLugarOperativo string `json:"CodigoLugarOperativo" validate:"required,lte=5"`
}

type IdentificadorContenedorParams struct {
	IdentificadorContenedor string `json:"IdentificadorContenedor" validate:"required,len=11"`
}

type ContenedorCarga struct {
	*IdentificadorContenedorParams
	CuitATA       string   `json:"CuitATA,omitempty"`
	Tipo          string   `json:"Tipo"`
	Peso          float64  `json:"Peso"`
	Precintos     []string `json:"Precintos"`
	Declaraciones []string `json:"Declaraciones,omitempty"`
}

type ContenedorVacio struct {
	*IdentificadorContenedorParams
	Tipo       string `json:"Tipo"`
	CuitATA    string `json:"CuitATA,omitempty"`
	CodigoPais string `json:"CodigoPais"`
}

type Embalaje struct {
	CodigoEmbalaje string  `json:"CodigoEmbalaje"`
	Peso           float64 `json:"Peso"`
	CantidadBultos int32   `json:"CantidadBultos"`
}

type MercaderiaSuelta struct {
	IdentificadorDeclaracion string      `json:"IdentificadorDeclaracion"`
	CuitATA                  string      `json:"CuitATA,omitempty"`
	Embalajes                *[]Embalaje `json:"Embalajes"`
}

type COEMParams struct {
	ContenedoresCarga  *[]ContenedorCarga  `json:"ContenedoresConCarga,omitempty"`
	ContenedoresVacios *[]ContenedorVacio  `json:"ContenedoresVacios,omitempty"`
	MercaderiasSueltas *[]MercaderiaSuelta `json:"MercaderiasSueltas,omitempty"`
}

type IdentificadorCOEMParams struct {
	*IdentificadorCaraturaParams
	IdentificadorCOEM string `json:"IdentificadorCOEM" validate:"required,len=16"`
}

type RegistrarCOEMParams struct {
	*IdentificadorCaraturaParams
	*COEMParams
}

type RectificarCOEMParams struct {
	*IdentificadorCOEMParams
	*COEMParams
}

type ContenedorNoABordo struct {
	*IdentificadorContenedorParams
}

type DeclaracionNoABordo struct {
	IdentificadorDeclaracion string `json:"IdentificadorDeclaracion"`
}

type NoABordoParams struct {
	ContenedoresCarga  *[]ContenedorNoABordo  `json:"ContenedoresCarga,omitempty"`
	ContenedoresVacios *[]ContenedorNoABordo  `json:"ContenedoresVacios,omitempty"`
	MercaderiasSueltas *[]DeclaracionNoABordo `json:"MercaderiasSueltas,omitempty"`
}

type DeclaracionCont struct {
	IdentificadorDeclaracion string    `json:"IdentificadorDeclaracion"`
	FechaEmbarque            time.Time `json:"FechaEmbarque"`
}

type CierreCargaContoBultoParams struct {
	Declaraciones *[]DeclaracionCont `json:"Declaraciones"`
}

type ItemGranel struct {
	NumeroItem   int32   `json:"NumeroItem"`
	CantidadReal float64 `json:"CantidadReal"`
}

type DeclaracionGranel struct {
	IdentificadorDeclaracion    string        `json:"IdentificadorDeclaracion"`
	FechaEmbarque               time.Time     `json:"FechaEmbarque"`
	IdentificadorCierreCumplido string        `json:"IdentificadorCierreCumplido"`
	Items                       *[]ItemGranel `json:"Items"`
}

type DeclaracionCOEMGranel struct {
	IdentificadorCOEM   string               `json:"IdentificadorCOEM"`
	DeclaracionesGranel *[]DeclaracionGranel `json:"DeclaracionesGranel"`
}

type CierreCargaGranelParams struct {
	DeclaracionesCOEMGranel *[]DeclaracionCOEMGranel `json:"DeclaracionesCOEMGranel"`
}

type SolicitarNoABordoParams struct {
	*IdentificadorCOEMParams
	CodigoMotivo      string `json:"CodigoMotivo" validate:"required,lte=2"`
	DescripcionMotivo string `json:"DescripcionMotivo,omitempty" validate:"lte=200"`
	*NoABordoParams
}

type SolicitarCierreCargaContoBultoParams struct {
	*IdentificadorCaraturaParams
	FechaZarpada time.Time `json:"FechaZarpada" validate:"required,datetime"`
	NumeroViaje  string    `json:"NumeroViaje,omitempty" validate:"len=16"`
	*CierreCargaContoBultoParams
}

type SolicitarCierreCargaGranelParams struct {
	*IdentificadorCaraturaParams
	FechaZarpada time.Time `json:"FechaZarpada" validate:"required,datetime"`
	NumeroViaje  string    `json:"NumeroViaje,omitempty" validate:"len=16"`
	*CierreCargaGranelParams
}

type ConsultaEstadoCOEM struct {
	IdentificadorCOEM string    `json:"IdentificadorCOEM"`
	CuitAlta          string    `json:"CuitAlta"`
	Motivo            string    `json:"Motivo"`
	FechaEstado       time.Time `json:"FechaEstado"`
	Estado            string    `json:"Estado"`
	CODE              string    `json:"CODE"`
}

type ConsultaNoAbordo struct {
	IdentificadorCACE   string    `json:"IdentificadorCACE"`
	IdentificadorCOEM   string    `json:"IdentificadorCOEM"`
	Tipo                string    `json:"Tipo"`
	Contenedor          string    `json:"Contenedor"`
	Destinacion         string    `json:"Destinacion"`
	Cuit                string    `json:"Cuit"`
	MotivoNoAbordo      string    `json:"MotivoNoAbordo"`
	FechaNoAbordo       time.Time `json:"FechaNoAbordo"`
	TipoNoAbordo        string    `json:"TipoNoAbordo"`
	DescripcionNoAbordo string    `json:"DescripcionNoAbordo"`
}

type ConsultaSolicitud struct {
	IdentificadorCACE string    `json:"IdentificadorCACE"`
	NumeroSolicitud   string    `json:"NumeroSolicitud"`
	IdentificadorCOEM string    `json:"IdentificadorCOEM"`
	Estado            string    `json:"Estado"`
	FechaEstado       time.Time `json:"FechaEstado"`
	TipoSolicitud     string    `json:"TipoSolicitud"`
}
