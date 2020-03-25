package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/AdhityaRamadhanus/cockpit"
	"github.com/AdhityaRamadhanus/cockpit/pkg/http/render"
	"github.com/AdhityaRamadhanus/cockpit/pkg/keygenerator"
	"github.com/gorilla/mux"
)

type CityHandler struct {
	KeyValueService cockpit.KeyValueService
}

func (h CityHandler) RegisterRoutes(router *mux.Router) {
	subRouter := router.PathPrefix("/api").Subrouter()

	subRouter.HandleFunc("/countries/indonesia/provinces/{province}/cities", h.getCitiesData).Methods("GET")
}

func (h CityHandler) getCitiesData(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	province := vars["province"]
	cityLevelCasesHash, err := h.KeyValueService.GetHashAll(keygenerator.CityLevelRedisKey(province))
	if err != nil {
		handleServiceError(res, req, err)
		return
	}

	jsonData := map[string]interface{}{}
	for key, hashStr := range cityLevelCasesHash {
		cityLevelCase := cockpit.CityLevelCases{}
		if err := json.Unmarshal(bytes.NewBufferString(hashStr).Bytes(), &cityLevelCase); err != nil {
			continue
		}

		jsonData[key] = cityLevelCase
	}

	render.JSON(res, http.StatusOK, jsonData)
}
