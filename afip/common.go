package afip

const URLWSAATesting string = "https://wsaahomo.afip.gov.ar/ws/services/LoginCms?WSDL"
const URLWSAAProduction string = "https://wsaa.afip.gov.ar/ws/services/LoginCms?WSDL"

const URLWSFETesting string = "https://wswhomo.afip.gov.ar/wsfev1/service.asmx?WSDL"
const URLWSFEProduction string = "https://servicios1.afip.gov.ar/wsfev1/service.asmx?WSDL"

type Environment int

const (
	TESTING Environment = iota
	PRODUCTION
)
