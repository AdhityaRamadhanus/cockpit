package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/AdhityaRamadhanus/cockpit"
	"github.com/AdhityaRamadhanus/cockpit/pkg/http/render"
	"github.com/AdhityaRamadhanus/cockpit/pkg/keygenerator"
	"github.com/gorilla/mux"
)

type ProvincialHandler struct {
	KeyValueService cockpit.KeyValueService
}

func (h ProvincialHandler) RegisterRoutes(router *mux.Router) {
	subRouter := router.PathPrefix("/api").Subrouter()

	subRouter.HandleFunc("/countries/indonesia/provinces/{province}", h.getProvincialData).Methods("GET")
}

func (h ProvincialHandler) getProvincialData(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	province := vars["province"]
	provincialLvlCasesBytes, err := h.KeyValueService.Get(keygenerator.ProvincialLevelRedisKey(province))
	if err != nil {
		handleServiceError(res, req, err)
		return
	}

	provincialLevelCases := cockpit.ProvincialLevelCases{}
	if err := json.Unmarshal(provincialLvlCasesBytes, &provincialLevelCases); err != nil {
		handleServiceError(res, req, err)
		return
	}

	render.JSON(res, http.StatusOK, provincialLevelCases)
}
