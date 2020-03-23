package tasks

import (
	"encoding/json"
	"runtime/debug"

	"github.com/AdhityaRamadhanus/cockpit"
	"github.com/AdhityaRamadhanus/cockpit/pkg/config"
	extractor "github.com/AdhityaRamadhanus/cockpit/pkg/extractors/html"
	"github.com/AdhityaRamadhanus/cockpit/pkg/http"
	"github.com/AdhityaRamadhanus/cockpit/pkg/keygenerator"
	"github.com/sirupsen/logrus"
)

type TaskJatim struct {
	KeyValueService cockpit.KeyValueService
	URL             config.URLConfig
	FeatureToggle   config.FeatureToggleConfig
}

func (t TaskJatim) SaveJatimProvincialLevelData() {
	defer func() {
		if err := recover(); err != nil {
			logrus.WithError(err.(error)).Errorf("Panic error %s", debug.Stack())
		}
	}()

	if !t.FeatureToggle.Jatim {
		logrus.Info("Skipping SaveJatimProvincialLevelData() Task")
		return
	}

	logrus.Info("Starting SaveJatimProvincialLevelData() Task")
	htmlText, err := http.GetHTML(t.URL.JatimDraxiData)
	if err != nil {
		logTaskError("task-jatim", err)
		return
	}

	provincialLevelCases, err := extractor.ExtractProvincialLevelCasesJatim(htmlText)
	if err != nil {
		logTaskError("task-jatim", err)
		return
	}

	jsonData, err := json.Marshal(provincialLevelCases)
	if err != nil {
		logTaskError("task-jatim", err)
		return
	}

	key := keygenerator.ProvincialLevelRedisKey("jawa-timur")
	if err := t.KeyValueService.Set(key, jsonData); err != nil {
		logTaskError("task-jatim", err)
		return
	}
	logrus.Info("Finished SaveJatimProvincialLevelData() Task")
}
