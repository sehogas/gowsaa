package dto

type DummyResponse struct {
	AppServer  string `json:"app_server"`
	AuthServer string `json:"auth_server"`
	DbServer   string `json:"db_server"`
}
