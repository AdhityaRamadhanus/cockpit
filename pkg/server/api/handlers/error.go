package handlers

import (
	"fmt"
	"net/http"

	"github.com/AdhityaRamadhanus/cockpit"
	"github.com/AdhityaRamadhanus/cockpit/pkg/http/render"
	"github.com/sirupsen/logrus"
)

var (
	ServiceErrorsHTTPMapping = map[error]struct {
		HTTPCode int
		ErrCode  string
	}{
		cockpit.ErrKeyNotFound: {
			HTTPCode: http.StatusNotFound,
			ErrCode:  "ErrKeyNotFOund",
		},
	}
)

func handleServiceError(res http.ResponseWriter, req *http.Request, err error) {
	errorMapping, isErrorMapped := ServiceErrorsHTTPMapping[err]
	if isErrorMapped {
		render.JSON(res, errorMapping.HTTPCode, map[string]interface{}{
			"status": errorMapping.HTTPCode,
			"error": map[string]interface{}{
				"code":    errorMapping.ErrCode,
				"message": err.Error(),
			},
		})
		return
	}

	logrus.WithFields(logrus.Fields{
		"method":       req.Method,
		"route":        req.URL.Path,
		"status_code":  http.StatusInternalServerError,
		"stack_trace":  fmt.Sprintf("%+v", err),
		"client":       req.Header.Get("X-API-ClientID"),
		"x-request-id": req.Header.Get("X-Request-ID"),
	}).WithError(err).Error("Error Handler")

	render.JSON(res, http.StatusInternalServerError, map[string]interface{}{
		"status": http.StatusInternalServerError,
		"error": map[string]interface{}{
			"code":    "ErrInternalServer",
			"message": "covid api unable to process the request",
		},
	})
	return
}
