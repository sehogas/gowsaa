package afip

import "time"

type CaratulaParams struct {
	IdentificadorBuque    string    `json:"IdentificadorBuque,omitempty"`
	NombreMedioTransporte string    `json:"NombreMedioTransporte" validate:"required"`
	CodigoAduana          string    `json:"CodigoAduana" validate:"required"`
	CodigoLugarOperativo  string    `json:"CodigoLugarOperativo" validate:"required"`
	FechaArribo           time.Time `json:"FechaArribo" validate:"required"`
	FechaZarpada          time.Time `json:"FechaZarpada" validate:"required"`
	Via                   string    `json:"Via,omitempty" validate:"required"`
	NumeroViaje           string    `json:"NumeroViaje,omitempty"`
	PuertoDestino         string    `json:"PuertoDestino,omitempty"`
	Itinerario            []string  `json:"itinerario,omitempty"`
}

type RectificarCaratulaParams struct {
	IdentificadorCaratula string `json:"IdentificadorCaratula,omitempty" validate:"required"`
	*CaratulaParams
}

type CambioBuqueParams struct {
	IdentificadorBuque    string `json:"identificador_buque,omitempty"`
	NombreMedioTransporte string `json:"nombre_medio_transporte,omitempty"`
}

type CambioFechasParams struct {
	FechaEstimadaArribo  time.Time `json:"fecha_estimada_arribo"`
	FechaEstimadaZarpada time.Time `json:"fecha_estimada_zarpada"`
	CodigoMotivo         string    `json:"codigo_motivo"`
	DescripcionMotivo    string    `json:"descripcion_motivo,omitempty"`
}

type CambioLOTParams struct {
	CodigoLugarOperativo string `json:"codigo_lugar_operativo"`
}

type ContenedorCarga struct {
	IdentificadorContenedor string   `json:"identificador_contenedor"`
	CuitATA                 string   `json:"cuit_ata,omitempty"`
	Tipo                    string   `json:"tipo"`
	Peso                    float64  `json:"peso"`
	Precintos               []string `json:"precintos"`
	Declaraciones           []string `json:"declaraciones,omitempty"`
}

type ContenedorVacio struct {
	IdentificadorContenedor string `json:"identificador_contenedor"`
	Tipo                    string `json:"tipo"`
	CuitATA                 string `json:"cuit_ata,omitempty"`
	CodigoPais              string `json:"codigo_pais"`
}

type Embalaje struct {
	CodigoEmbalaje string  `json:"codigo_embalaje"`
	Peso           float64 `json:"peso"`
	CantidadBultos int32   `json:"cantidad_bultos"`
}

type MercaderiaSuelta struct {
	IdentificadorDeclaracion string      `json:"identificador_declaracion"`
	CuitATA                  string      `json:"cuit_ata,omitempty"`
	Embalajes                *[]Embalaje `json:"embalajes"`
}

type COEMParams struct {
	ContenedoresCarga  *[]ContenedorCarga  `json:"contenedores_carga,omitempty"`
	ContenedoresVacios *[]ContenedorVacio  `json:"contenedores_vacios,omitempty"`
	MercaderiasSueltas *[]MercaderiaSuelta `json:"mercaderias_sueltas,omitempty"`
}

type ContenedorNoABordo struct {
	IdentificadorContenedor string `json:"identificador_contenedor"`
}

type DeclaracionNoABordo struct {
	IdentificadorDeclaracion string `json:"identificador_declaracion"`
}

type NoABordoParams struct {
	ContenedoresCarga  *[]ContenedorNoABordo  `json:"contenedores_carga,omitempty"`
	ContenedoresVacios *[]ContenedorNoABordo  `json:"contenedores_vacios,omitempty"`
	MercaderiasSueltas *[]DeclaracionNoABordo `json:"declaraciones,omitempty"`
}

type DeclaracionCont struct {
	IdentificadorDeclaracion string    `json:"identificador_declaracion"`
	FechaEmbarque            time.Time `json:"fecha_embarque"`
}

type CierreCargaContoBultoParams struct {
	Declaraciones *[]DeclaracionCont `json:"declaraciones"`
}

type ItemGranel struct {
	NumeroItem   int32   `json:"numero_item"`
	CantidadReal float64 `json:"cantidad_real"`
}

type DeclaracionGranel struct {
	IdentificadorDeclaracion    string        `json:"identificador_declaracion"`
	FechaEmbarque               time.Time     `json:"fecha_embarque"`
	IdentificadorCierreCumplido string        `json:"identificador_cierre_cumplido"`
	Items                       *[]ItemGranel `json:"items"`
}

type DeclaracionCOEMGranel struct {
	IdentificadorCOEM   string               `json:"identificador_coem"`
	DeclaracionesGranel *[]DeclaracionGranel `json:"declaraciones_granel"`
}

type CierreCargaGranelParams struct {
	DeclaracionesCOEMGranel *[]DeclaracionCOEMGranel `json:"declaraciones_coem_granel"`
}

type ConsultaEstadoCOEM struct {
	IdentificadorCOEM string    `json:"identificador_coem"`
	CuitAlta          string    `json:"cuit_alta"`
	Motivo            string    `json:"motivo"`
	FechaEstado       time.Time `json:"fecha_estado"`
	Estado            string    `json:"estado"`
	CODE              string    `json:"code"`
}

type ConsultaNoAbordo struct {
	IdentificadorCACE   string    `json:"identificador_cace"`
	IdentificadorCOEM   string    `json:"identificador_coem"`
	Tipo                string    `json:"tipo"`
	Contenedor          string    `json:"contenedor"`
	Destinacion         string    `json:"destinacion"`
	Cuit                string    `json:"cuit"`
	MotivoNoAbordo      string    `json:"motivo_no_abordo"`
	FechaNoAbordo       time.Time `json:"fecha_no_abordo"`
	TipoNoAbordo        string    `json:"tipo_no_abordo"`
	DescripcionNoAbordo string    `json:"descripcion_no_abordo"`
}

type ConsultaSolicitud struct {
	IdentificadorCACE string    `json:"Identificador_cace"`
	NumeroSolicitud   string    `json:"numero_solicitud"`
	IdentificadorCOEM string    `json:"identificador_coem"`
	Estado            string    `json:"estado"`
	FechaEstado       time.Time `json:"fecha_estado"`
	TipoSolicitud     string    `json:"tipo_solicitud"`
}
