package main

import (
	"errors"
	"net/http"

	"github.com/sehogas/gowsaa/internal/dto"
	"github.com/sehogas/gowsaa/internal/util"
)

//validate := validator.New(validator.WithRequiredStructEnabled())

func DummyHandler(w http.ResponseWriter, r *http.Request) {
	appServer, authServer, dbServer, err := Wscoemcons.Dummy()
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

func ObtenerConsultaEstadosCOEMHandler(w http.ResponseWriter, r *http.Request) {
	identificadorCaratula := r.URL.Query().Get("identificador")
	if len(identificadorCaratula) != 16 {
		err := errors.New("parámetro Identificador faltante o inválido")
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.GenericResponse{
			Status:  false,
			Message: err.Error(),
		}, err)
		return
	}

	data, err := Wscoemcons.ObtenerConsultaEstadosCOEM(identificadorCaratula)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.GenericResponse{
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
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.GenericResponse{
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
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.GenericResponse{
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
