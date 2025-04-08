package afip

import "time"

type CaratulaParams struct {
	IdentificadorBuque    string    `json:"identificador_buque,omitempty"`
	NombreMedioTransporte string    `json:"nombre_medio_transporte"`
	CodigoAduana          string    `json:"codigo_aduana"`
	CodigoLugarOperativo  string    `json:"codigo_lugar_operativo"`
	FechaEstimadaArribo   time.Time `json:"fecha_estimada_arribo"`
	FechaEstimadaZarpada  time.Time `json:"fecha_estimada_zarpada"`
	Via                   string    `json:"via,omitempty"`
	NumeroViaje           string    `json:"numero_viaje,omitempty"`
	PuertoDestino         string    `json:"puerto_destino,omitempty"`
	Itinerario            []string  `json:"itinerario,omitempty"`
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
