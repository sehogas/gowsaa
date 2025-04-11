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
				resp.Errors[i] = fmt.Sprintf("el parámetro %s es requerido", err.Field())
			case "max":
				resp.Errors[i] = fmt.Sprintf("el parámetro %s debe tener una longitud máxima de %s", err.Field(), err.Param())
			case "url":
				resp.Errors[i] = fmt.Sprintf("el parámetro %s debe ser una URL válida", err.Field())
			case "email":
				resp.Errors[i] = fmt.Sprintf("el parámetro %s debe ser un email válido", err.Field())
			case "alphaspace":
				resp.Errors[i] = fmt.Sprintf("el parámetro %s sólo puede contener caracteres alfabéticos y espacios", err.Field())
			case "datetime":
				resp.Errors[i] = fmt.Sprintf("el parámetro %s debe tener el formato %s", err.Field(), err.Param())
			case "len":
				resp.Errors[i] = fmt.Sprintf("la longitud del parámetro %s debe ser de %s caracteres", err.Field(), err.Param())
			case "startswith":
				resp.Errors[i] = fmt.Sprintf("el parámetro %s debe comenzar con %s", err.Field(), err.Param())
			default:
				resp.Errors[i] = fmt.Sprintf("parámetro %s no cumple los requerimientos de validación; tag %s", err.Field(), err.Tag())
			}
		}

		return &resp
	}

	return nil
}
