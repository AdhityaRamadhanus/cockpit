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

type TaskBanten struct {
	KeyValueService cockpit.KeyValueService
	URL             config.URLConfig
	FeatureToggle   config.FeatureToggleConfig
}

func (t TaskBanten) SaveBantenProvincialLevelData() {
	defer func() {
		if err := recover(); err != nil {
			logrus.WithError(err.(error)).Errorf("Panic error %s", debug.Stack())
		}
	}()

	if !t.FeatureToggle.Banten {
		logrus.Info("Skipping SaveBantenProvincialLevelData() Task")
		return
	}

	logrus.Info("Starting SaveBantenProvincialLevelData() Task")

	htmlText, err := http.GetHTML(t.URL.BantenWeb)
	if err != nil {
		logTaskError("task-banten", err)
		return
	}

	provincialLevelCases, err := extractor.ExtractProvincialLevelCasesBanten(htmlText)
	if err != nil {
		logTaskError("task-banten", err)
		return
	}

	jsonData, err := json.Marshal(provincialLevelCases)
	if err != nil {
		logTaskError("task-banten", err)
		return
	}

	key := keygenerator.ProvincialLevelRedisKey("banten")
	if err := t.KeyValueService.Set(key, jsonData); err != nil {
		logTaskError("task-banten", err)
		return
	}
	logrus.Info("Finished SaveBantenProvincialLevelData() Task")
}
