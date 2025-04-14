package dto

import "time"

type DummyResponse struct {
	AppServer  string `json:"AppServer"`
	AuthServer string `json:"AuthServer"`
	DbServer   string `json:"DbServer"`
}

type FecUltActResponse struct {
	FechaUltAct time.Time `json:"FechaUltAct"`
}
