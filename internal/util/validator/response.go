package validator

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type ErrResponse struct {
	Errors []string `json:"errors"`
}

func ToErrResponse(err error) *ErrResponse {
	if fieldErrors, ok := err.(validator.ValidationErrors); ok {
		resp := ErrResponse{
			Errors: make([]string, len(fieldErrors)),
		}

		for i, err := range fieldErrors {
			switch err.Tag() {
			case "required":
				resp.Errors[i] = fmt.Sprintf("%s es un campo requerido", err.Field())
			case "max":
				resp.Errors[i] = fmt.Sprintf("%s debe tener una longitud máxima de %s", err.Field(), err.Param())
			case "url":
				resp.Errors[i] = fmt.Sprintf("%s debe ser una URL válida", err.Field())
			case "email":
				resp.Errors[i] = fmt.Sprintf("%s debe ser un email válido", err.Field())
			case "alphaspace":
				resp.Errors[i] = fmt.Sprintf("%s sólo puede contener caracteres alfabéticos y espacios", err.Field())
			case "datetime":
				if err.Param() == "2006-01-02" {
					resp.Errors[i] = fmt.Sprintf("%s debe ser una fecha válida", err.Field())
				} else {
					resp.Errors[i] = fmt.Sprintf("%s debe tener formato %s", err.Field(), err.Param())
				}
			default:
				resp.Errors[i] = fmt.Sprintf("Algo anduvo mal con %s; %s", err.Field(), err.Tag())
			}
		}

		return &resp
	}

	return nil
}
