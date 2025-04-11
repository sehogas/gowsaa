package dto

type GenericResponse struct {
	Status     bool        `json:"status" default:"false"`
	StatusCode int         `json:"statusCode,omitempty"`
	Message    string      `json:"message,omitempty"`
	Data       interface{} `json:"data,omitempty"`
	TotalRows  int64       `json:"total_rows,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type SendEmailRequest struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

type InfoResponse struct {
	Version string `json:"version"`
}
