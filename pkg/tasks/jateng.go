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

type TaskJateng struct {
	KeyValueService cockpit.KeyValueService
	URL             config.URLConfig
	FeatureToggle   config.FeatureToggleConfig
}

func (t TaskJateng) SaveJatengProvincialLevelData() {
	defer func() {
		if err := recover(); err != nil {
			logrus.WithError(err.(error)).Errorf("Panic error %s", debug.Stack())
		}
	}()

	if !t.FeatureToggle.Jakarta {
		logrus.Info("Skipping SaveJatengProvincialLevelData() Task")
		return
	}

	logrus.Info("Starting SaveJatengProvincialLevelData() Task")

	htmlText, err := http.GetHTML(t.URL.JatengWeb)
	if err != nil {
		logTaskError("task-jateng", err)
		return
	}

	provincialLevelCases, err := extractor.ExtractProvincialLevelCasesJateng(htmlText)
	if err != nil {
		logTaskError("task-jateng", err)
		return
	}

	jsonData, err := json.Marshal(provincialLevelCases)
	if err != nil {
		logTaskError("task-jateng", err)
		return
	}

	key := keygenerator.ProvincialLevelRedisKey("jawa-tengah")
	if err := t.KeyValueService.Set(key, jsonData); err != nil {
		logTaskError("task-jateng", err)
		return
	}
	logrus.Info("Finished SaveJatengProvincialLevelData() Task")
}
