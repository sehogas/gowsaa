package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/sehogas/gowsaa/afip"
	"github.com/sehogas/gowsaa/internal/dto"
	"github.com/sehogas/gowsaa/internal/util"
	"github.com/sehogas/gowsaa/internal/util/validator"
)

// DummyHandler godoc
//
//	@Summary		Muestra el estado del servicio
//	@Description	Visualizar el estado del servicio web, del servicio de autenticación y de la base de datos de ARCA
//	@Tags			Dummy
//	@Produce		json
//	@Success		200	{object}	dto.DummyResponse
//	@Failure		500	{object}	dto.ErrorResponse
//	@Router			/dummy [get]
func DummyHandler(w http.ResponseWriter, r *http.Request) {
	appServer, authServer, dbServer, err := Wscoem.Dummy()
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}
	util.HttpResponseJSON(w, http.StatusOK, &dto.DummyResponse{
		AppServer:  appServer,
		AuthServer: authServer,
		DbServer:   dbServer,
	}, nil)
}

// RegistrarCaratulaHandler godoc
//
//	@Summary		Registrar Carátula
//	@Description	Registra una nueva Carátula
//	@Tags			RegistrarCaratula
//	@Accept			json
//	@Produce		json
//	@Param			request	body		afip.CaratulaParams	true	"RegistrarCaratulaRequest"
//	@Success		200		{object}	dto.MessageResponse
//	@Failure		400		{object}	dto.ErrorResponse
//	@Failure		500		{object}	dto.ErrorResponse
//	@Router			/registrar-caratula [post]
func RegistrarCaratulaHandler(w http.ResponseWriter, r *http.Request) {
	var post afip.CaratulaParams
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: "error leyendo parámetros de la solicitud"}, err)
		return
	}

	if err := validate.Struct(post); err != nil {
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: strings.Join(validator.ToErrResponse(err).Errors, ", ")}, err)
		return
	}

	identificador, err := Wscoem.RegistrarCaratula(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}

	util.HttpResponseJSON(w, http.StatusOK, &dto.MessageResponse{Message: identificador}, nil)
}

// RectificarCaratulaHandler godoc
//
//	@Summary		Rectificar Carátula
//	@Description	Rectificar una Carátula sin COEMs
//	@Tags			RectificarCaratula
//	@Accept			json
//	@Produce		json
//	@Param			request	body		afip.RectificarCaratulaParams	true	"RectificarCaratulaRequest"
//	@Success		200		{object}	dto.MessageResponse
//	@Failure		400		{object}	dto.ErrorResponse
//	@Failure		500		{object}	dto.ErrorResponse
//	@Router			/rectificar-caratula [put]
func RectificarCaratulaHandler(w http.ResponseWriter, r *http.Request) {
	var post afip.RectificarCaratulaParams
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: "error leyendo parámetros de la solicitud"}, err)
		return
	}

	if err := validate.Struct(post); err != nil {
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: strings.Join(validator.ToErrResponse(err).Errors, ", ")}, err)
		return
	}

	identificador, err := Wscoem.RectificarCaratula(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}

	util.HttpResponseJSON(w, http.StatusOK, &dto.MessageResponse{Message: identificador}, nil)
}

// AnularCaratulaHandler godoc
//
//	@Summary		Anular Carátula
//	@Description	Anular Carátula sin COEMs
//	@Tags			AnularCaratula
//	@Accept			json
//	@Produce		json
//	@Param			request	body		afip.IdentificadorCaraturaParams	true	"AnularCaratulaRequest"
//	@Success		200		{object}	dto.MessageResponse
//	@Failure		400		{object}	dto.ErrorResponse
//	@Failure		500		{object}	dto.ErrorResponse
//	@Router			/anular-caratula [delete]
func AnularCaratulaHandler(w http.ResponseWriter, r *http.Request) {
	var post afip.IdentificadorCaraturaParams
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: "error leyendo parámetros de la solicitud"}, err)
		return
	}

	if err := validate.Struct(post); err != nil {
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: strings.Join(validator.ToErrResponse(err).Errors, ", ")}, err)
		return
	}

	identificador, err := Wscoem.AnularCaratula(&afip.IdentificadorCaraturaParams{
		IdentificadorCaratula: post.IdentificadorCaratula,
	})
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}

	util.HttpResponseJSON(w, http.StatusOK, &dto.MessageResponse{Message: identificador}, nil)
}

// SolicitarCambioBuqueHandler godoc
//
//	@Summary		Solicitar cambio de Buque
//	@Description	Solicitar cambio de Buque para Carátulas con COEMs
//	@Tags			SolicitarCambioBuque
//	@Accept			json
//	@Produce		json
//	@Param			request	body		afip.CambioBuqueParams	true	"SolicitarCambioBuqueRequest"
//	@Success		200		{object}	dto.MessageResponse
//	@Failure		400		{object}	dto.ErrorResponse
//	@Failure		500		{object}	dto.ErrorResponse
//	@Router			/solicitar-cambio-buque [put]
func SolicitarCambioBuqueHandler(w http.ResponseWriter, r *http.Request) {
	var post afip.CambioBuqueParams
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: "error leyendo parámetros de la solicitud"}, err)
		return
	}

	if err := validate.Struct(post); err != nil {
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: strings.Join(validator.ToErrResponse(err).Errors, ", ")}, err)
		return
	}

	identificador, err := Wscoem.SolicitarCambioBuque(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}

	util.HttpResponseJSON(w, http.StatusOK, &dto.MessageResponse{Message: identificador}, nil)
}

// SolicitarCambioFechasHandler godoc
//
//	@Summary		Solicitar cambio de Fechas
//	@Description	Solicitar cambio de Fechas para Carátulas con COEMs
//	@Tags			SolicitarCambioFechas
//	@Accept			json
//	@Produce		json
//	@Param			request	body		afip.CambioFechasParams	true	"CambioFechasParamsRequest"
//	@Success		200		{object}	dto.MessageResponse
//	@Failure		400		{object}	dto.ErrorResponse
//	@Failure		500		{object}	dto.ErrorResponse
//	@Router			/solicitar-cambio-fechas [put]
func SolicitarCambioFechasHandler(w http.ResponseWriter, r *http.Request) {
	var post afip.CambioFechasParams
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: "error leyendo parámetros de la solicitud"}, err)
		return
	}

	if err := validate.Struct(post); err != nil {
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: strings.Join(validator.ToErrResponse(err).Errors, ", ")}, err)
		return
	}

	identificador, err := Wscoem.SolicitarCambioFechas(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}

	util.HttpResponseJSON(w, http.StatusOK, &dto.MessageResponse{Message: identificador}, nil)
}

// SolicitarCambioLOTHandler godoc
//
//	@Summary		Solicitar cambio de LOT
//	@Description	Solicitar cambio de Lugar Operativo para Carátulas con COEMs
//	@Tags			SolicitarCambioLOT
//	@Accept			json
//	@Produce		json
//	@Param			request	body		afip.CambioLOTParams	true	"CambioLOTParamsRequest"
//	@Success		200		{object}	dto.MessageResponse
//	@Failure		400		{object}	dto.ErrorResponse
//	@Failure		500		{object}	dto.ErrorResponse
//	@Router			/solicitar-cambio-lot [put]
func SolicitarCambioLOTHandler(w http.ResponseWriter, r *http.Request) {
	var post afip.CambioLOTParams
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: "error leyendo parámetros de la solicitud"}, err)
		return
	}

	if err := validate.Struct(post); err != nil {
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: strings.Join(validator.ToErrResponse(err).Errors, ", ")}, err)
		return
	}

	identificador, err := Wscoem.SolicitarCambioLOT(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}

	util.HttpResponseJSON(w, http.StatusOK, &dto.MessageResponse{Message: identificador}, nil)
}

// RegistrarCOEMHandler godoc
//
//	@Summary		Registrar COEM
//	@Description	Registrar COEM en Carátula
//	@Tags			RegistrarCOEM
//	@Accept			json
//	@Produce		json
//	@Param			request	body		afip.RegistrarCOEMParams	true	"RegistrarCOEMRequest"
//	@Success		200		{object}	dto.MessageResponse
//	@Failure		400		{object}	dto.ErrorResponse
//	@Failure		500		{object}	dto.ErrorResponse
//	@Router			/registrar-coem [post]
func RegistrarCOEMHandler(w http.ResponseWriter, r *http.Request) {
	var post afip.RegistrarCOEMParams
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: "error leyendo parámetros de la solicitud"}, err)
		return
	}

	if err := validate.Struct(post); err != nil {
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: strings.Join(validator.ToErrResponse(err).Errors, ", ")}, err)
		return
	}

	identificador, err := Wscoem.RegistrarCOEM(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}

	util.HttpResponseJSON(w, http.StatusOK, &dto.MessageResponse{Message: identificador}, nil)
}

// RectificarCOEMHandler godoc
//
//	@Summary		Rectificar COEM
//	@Description	Rectificar COEM
//	@Tags			RectificarCOEM
//	@Accept			json
//	@Produce		json
//	@Param			request	body		afip.RectificarCOEMParams	true	"RectificarCOEMRequest"
//	@Success		200		{object}	dto.MessageResponse
//	@Failure		400		{object}	dto.ErrorResponse
//	@Failure		500		{object}	dto.ErrorResponse
//	@Router			/rectificar-coem [put]
func RectificarCOEMHandler(w http.ResponseWriter, r *http.Request) {
	var post afip.RectificarCOEMParams
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: "error leyendo parámetros de la solicitud"}, err)
		return
	}

	if err := validate.Struct(post); err != nil {
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: strings.Join(validator.ToErrResponse(err).Errors, ", ")}, err)
		return
	}

	identificador, err := Wscoem.RectificarCOEM(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}

	util.HttpResponseJSON(w, http.StatusOK, &dto.MessageResponse{Message: identificador}, nil)
}

// CerrarCOEMHandler godoc
//
//	@Summary		Cerrar COEM
//	@Description	Cerrar COEM
//	@Tags			CerrarCOEM
//	@Accept			json
//	@Produce		json
//	@Param			request	body		afip.IdentificadorCOEMParams	true	"CerrarCOEMRequest"
//	@Success		200		{object}	dto.MessageResponse
//	@Failure		400		{object}	dto.ErrorResponse
//	@Failure		500		{object}	dto.ErrorResponse
//	@Router			/cerrar-coem [post]
func CerrarCOEMHandler(w http.ResponseWriter, r *http.Request) {
	var post afip.IdentificadorCOEMParams
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: "error leyendo parámetros de la solicitud"}, err)
		return
	}

	if err := validate.Struct(post); err != nil {
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: strings.Join(validator.ToErrResponse(err).Errors, ", ")}, err)
		return
	}

	identificador, err := Wscoem.CerrarCOEM(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}

	util.HttpResponseJSON(w, http.StatusOK, &dto.MessageResponse{Message: identificador}, nil)
}

// AnularCOEMHandler godoc
//
//	@Summary		Anular COEM
//	@Description	Anular COEM
//	@Tags			AnularCOEM
//	@Accept			json
//	@Produce		json
//	@Param			request	body		afip.IdentificadorCOEMParams	true	"AnularCOEMRequest"
//	@Success		200		{object}	dto.MessageResponse
//	@Failure		400		{object}	dto.ErrorResponse
//	@Failure		500		{object}	dto.ErrorResponse
//	@Router			/anular-coem [delete]
func AnularCOEMHandler(w http.ResponseWriter, r *http.Request) {
	var post afip.IdentificadorCOEMParams
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: "error leyendo parámetros de la solicitud"}, err)
		return
	}

	if err := validate.Struct(post); err != nil {
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: strings.Join(validator.ToErrResponse(err).Errors, ", ")}, err)
		return
	}

	identificador, err := Wscoem.AnularCOEM(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}

	util.HttpResponseJSON(w, http.StatusOK, &dto.MessageResponse{Message: identificador}, nil)
}

// SolicitarAnulacionCOEMHandler godoc
//
//	@Summary		Solicitar Anulación COEM
//	@Description	Solicitar anulación COEM
//	@Tags			SolicitarAnulacionCOEM
//	@Accept			json
//	@Produce		json
//	@Param			request	body		afip.IdentificadorCOEMParams	true	"SolicitarAnulacionCOEMRequest"
//	@Success		200		{object}	dto.MessageResponse
//	@Failure		400		{object}	dto.ErrorResponse
//	@Failure		500		{object}	dto.ErrorResponse
//	@Router			/solicitar-anulacion-coem [post]
func SolicitarAnulacionCOEMHandler(w http.ResponseWriter, r *http.Request) {
	var post afip.IdentificadorCOEMParams
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: "error leyendo parámetros de la solicitud"}, err)
		return
	}

	if err := validate.Struct(post); err != nil {
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: strings.Join(validator.ToErrResponse(err).Errors, ", ")}, err)
		return
	}

	identificador, err := Wscoem.SolicitarAnulacionCOEM(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}

	util.HttpResponseJSON(w, http.StatusOK, &dto.MessageResponse{Message: identificador}, nil)
}

// SolicitarNoABordoHandler godoc
//
//	@Summary		Solicitar No Abordo
//	@Description	Solicitar No Abordo
//	@Tags			SolicitarNoABordo
//	@Accept			json
//	@Produce		json
//	@Param			request	body		afip.IdentificadorCOEMParams	true	"SolicitarNoABordoRequest"
//	@Success		200		{object}	dto.MessageResponse
//	@Failure		400		{object}	dto.ErrorResponse
//	@Failure		500		{object}	dto.ErrorResponse
//	@Router			/solicitar-no-abordo [post]
func SolicitarNoABordoHandler(w http.ResponseWriter, r *http.Request) {
	var post afip.SolicitarNoABordoParams
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: "error leyendo parámetros de la solicitud"}, err)
		return
	}

	if err := validate.Struct(post); err != nil {
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: strings.Join(validator.ToErrResponse(err).Errors, ", ")}, err)
		return
	}

	identificador, err := Wscoem.SolicitarNoABordo(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}

	util.HttpResponseJSON(w, http.StatusOK, &dto.MessageResponse{Message: identificador}, nil)
}

// SolicitarCierreCargaContoBultoHandler godoc
//
//	@Summary		Solicitar Cierre de Carga Contenedores y/o Bultos
//	@Description	Solicitar Cierre de Carga Contenedores y/o Bultos
//	@Tags			SolicitarCierreCargaContoBulto
//	@Accept			json
//	@Produce		json
//	@Param			request	body		afip.SolicitarCierreCargaContoBultoParams	true	"SolicitarCierreCargaContoBultoRequest"
//	@Success		200		{object}	dto.MessageResponse
//	@Failure		400		{object}	dto.ErrorResponse
//	@Failure		500		{object}	dto.ErrorResponse
//	@Router			/solicitar-cierre-carga-conto-bulto [post]
func SolicitarCierreCargaContoBultoHandler(w http.ResponseWriter, r *http.Request) {
	var post afip.SolicitarCierreCargaContoBultoParams
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: "error leyendo parámetros de la solicitud"}, err)
		return
	}

	if err := validate.Struct(post); err != nil {
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: strings.Join(validator.ToErrResponse(err).Errors, ", ")}, err)
		return
	}

	identificador, err := Wscoem.SolicitarCierreCargaContoBulto(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}

	util.HttpResponseJSON(w, http.StatusOK, &dto.MessageResponse{Message: identificador}, nil)
}

// SolicitarCierreCargaGranelHandler godoc
//
//	@Summary		Solicitar Cierre de Carga Granel
//	@Description	Solicitar Cierre de Carga Granel
//	@Tags			SolicitarCierreCargaGranel
//	@Accept			json
//	@Produce		json
//	@Param			request	body		afip.SolicitarCierreCargaGranelParams	true	"SolicitarCierreCargaGranel"
//	@Success		200		{object}	dto.MessageResponse
//	@Failure		400		{object}	dto.ErrorResponse
//	@Failure		500		{object}	dto.ErrorResponse
//	@Router			/solicitar-cierre-carga-granel [post]
func SolicitarCierreCargaGranelHandler(w http.ResponseWriter, r *http.Request) {
	var post afip.SolicitarCierreCargaGranelParams
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: "error leyendo parámetros de la solicitud"}, err)
		return
	}

	if err := validate.Struct(post); err != nil {
		util.HttpResponseJSON(w, http.StatusBadRequest, &dto.ErrorResponse{Error: strings.Join(validator.ToErrResponse(err).Errors, ", ")}, err)
		return
	}

	identificador, err := Wscoem.SolicitarCierreCargaGranel(&post)
	if err != nil {
		util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.ErrorResponse{Error: err.Error()}, err)
		return
	}

	util.HttpResponseJSON(w, http.StatusOK, &dto.MessageResponse{Message: identificador}, nil)
}
