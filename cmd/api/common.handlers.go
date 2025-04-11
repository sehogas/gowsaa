package main

import (
	"net/http"

	"github.com/sehogas/gowsaa/internal/dto"
	"github.com/sehogas/gowsaa/internal/util"
)

// InfoHandler godoc
//
//	@Summary		Muesta información de la API
//	@Description	Muesta información de la API
//	@Tags			API
//	@Produce		json
//	@Success		200	{object}	dto.DummyResponse
//	@Router			/info [get]
func InfoHandler(w http.ResponseWriter, r *http.Request) {
	util.HttpResponseJSON(w, http.StatusOK, &dto.InfoResponse{
		Version: Version,
	}, nil)
}
