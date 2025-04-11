package dto

type DummyResponse struct {
	AppServer  string `json:"AppServer"`
	AuthServer string `json:"AuthServer"`
	DbServer   string `json:"DbServer"`
}
