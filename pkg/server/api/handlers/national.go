package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/AdhityaRamadhanus/cockpit"
	"github.com/AdhityaRamadhanus/cockpit/pkg/http/render"
	"github.com/AdhityaRamadhanus/cockpit/pkg/keygenerator"
	"github.com/gorilla/mux"
)

type NationalHandler struct {
	KeyValueService cockpit.KeyValueService
}

func (h NationalHandler) RegisterRoutes(router *mux.Router) {
	subRouter := router.PathPrefix("/api").Subrouter()

	subRouter.HandleFunc("/countries/indonesia", h.getNationalData).Methods("GET")

}

func (h NationalHandler) getNationalData(res http.ResponseWriter, req *http.Request) {
	nationalLvlCasesBytes, err := h.KeyValueService.Get(keygenerator.NationalLevelRedisKey())
	if err != nil {
		handleServiceError(res, req, err)
		return
	}

	nationalLvlCases := cockpit.NationalLevelCases{}
	if err := json.Unmarshal(nationalLvlCasesBytes, &nationalLvlCases); err != nil {
		handleServiceError(res, req, err)
		return
	}

	render.JSON(res, http.StatusOK, nationalLvlCases)
}
