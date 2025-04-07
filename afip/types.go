package afip

import "time"

type RegistrarCaratulaParams struct {
	IdentificadorBuque    string    `json:"identificador_buque,omitempty"`
	NombreMedioTransporte string    `json:"NombreMedioTransporte"`
	CodigoAduana          string    `json:"codigo_aduana"`
	CodigoLugarOperativo  string    `json:"codigo_lugar_operativo"`
	FechaEstimadaArribo   time.Time `json:"fecha_estimada_arribo"`
	FechaEstimadaZarpada  time.Time `json:"fecha_estimada_zarpada"`
	Via                   string    `json:"via,omitempty"`
	NumeroViaje           string    `json:"numero_viaje,omitempty"`
	PuertoDestino         string    `json:"puerto_destino,omitempty"`
	Itinerario            []string  `json:"itinerario,omitempty"`
}
