package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/sehogas/gowsaa/afip"
	"github.com/sehogas/gowsaa/internal/dto"
	"github.com/sehogas/gowsaa/internal/util"
	"github.com/sehogas/gowsaa/internal/util/validator"
)

//validate := validator.New(validator.WithRequiredStructEnabled())

func DummyHandler(w http.ResponseWriter, r *http.Request) {
	appServer, authServer, dbServer, err := Wscoem.Dummy()
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.GenericResponse{
			Status:  false,
			Message: err.Error(),
		}, err)
		return
	}
	util.HttpResponseJSON(w, http.StatusOK, &dto.GenericResponse{
		Status: true,
		Data: &dto.DummyResponse{
			AppServer:  appServer,
			AuthServer: authServer,
			DbServer:   dbServer,
		},
	}, nil)
}

func RegistrarCaratulaHandler(w http.ResponseWriter, r *http.Request) {
	var post afip.CaratulaParams
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.GenericResponse{
			Status:  false,
			Message: "error decodificando requerimiento"},
			err)
		return
	}

	if err := validate.Struct(post); err != nil {
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.GenericResponse{
			Status:  false,
			Message: strings.Join(validator.ToErrResponse(err).Errors, ", ")}, err)
		return
	}

	identificador, err := Wscoem.RegistrarCaratula(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.GenericResponse{
			Status:  false,
			Message: err.Error(),
		}, err)
		return
	}

	util.HttpResponseJSON(w, http.StatusOK, &dto.GenericResponse{
		Status: true,
		Data:   identificador,
	}, nil)
}

func RectificarCaratulaHandler(w http.ResponseWriter, r *http.Request) {
	var post afip.RectificarCaratulaParams
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.GenericResponse{
			Status:  false,
			Message: "error decodificando requerimiento"},
			err)
		return
	}

	if err := validate.Struct(post); err != nil {
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.GenericResponse{
			Status:  false,
			Message: strings.Join(validator.ToErrResponse(err).Errors, ", ")}, err)
		return
	}

	identificador, err := Wscoem.RectificarCaratula(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.GenericResponse{
			Status:  false,
			Message: err.Error(),
		}, err)
		return
	}

	util.HttpResponseJSON(w, http.StatusOK, &dto.GenericResponse{
		Status: true,
		Data:   identificador,
	}, nil)
}

func AnularCaratulaHandler(w http.ResponseWriter, r *http.Request) {
	identificadorCaratula := r.URL.Query().Get("identificador")
	if len(identificadorCaratula) != 16 {
		err := errors.New("parámetro Identificador faltante o inválido")
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.GenericResponse{
			Status:  false,
			Message: err.Error(),
		}, err)
		return
	}

	identificador, err := Wscoem.AnularCaratula(identificadorCaratula)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.GenericResponse{
			Status:  false,
			Message: err.Error(),
		}, err)
		return
	}

	util.HttpResponseJSON(w, http.StatusOK, &dto.GenericResponse{
		Status: true,
		Data:   identificador,
	}, nil)
}

/*
func ObtenerConsultaNoAbordoHandler(w http.ResponseWriter, r *http.Request) {
	identificadorCaratula := r.URL.Query().Get("identificador")
	if len(identificadorCaratula) != 16 {
		err := errors.New("parámetro Identificador faltante o inválido")
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.GenericResponse{
			Status:  false,
			Message: err.Error(),
		}, err)
		return
	}

	data, err := Wscoemcons.ObtenerConsultaNoAbordo(identificadorCaratula)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.GenericResponse{
			Status:  false,
			Message: err.Error(),
		}, err)
		return
	}

	util.HttpResponseJSON(w, http.StatusOK, &dto.GenericResponse{
		Status: true,
		Data:   data,
	}, nil)
}

func ObtenerConsultaSolicitudesHandler(w http.ResponseWriter, r *http.Request) {
	identificadorCaratula := r.URL.Query().Get("identificador")
	if len(identificadorCaratula) != 16 {
		err := errors.New("parámetro Identificador faltante o inválido")
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.GenericResponse{
			Status:  false,
			Message: err.Error(),
		}, err)
		return
	}

	data, err := Wscoemcons.ObtenerConsultaSolicitudes(identificadorCaratula)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.GenericResponse{
			Status:  false,
			Message: err.Error(),
		}, err)
		return
	}

	util.HttpResponseJSON(w, http.StatusOK, &dto.GenericResponse{
		Status: true,
		Data:   data,
	}, nil)
}
*/
