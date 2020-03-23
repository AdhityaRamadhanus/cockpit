package tasks

import (
	"encoding/json"
	"runtime/debug"

	"github.com/AdhityaRamadhanus/cockpit"
	"github.com/AdhityaRamadhanus/cockpit/pkg/config"
	extractor "github.com/AdhityaRamadhanus/cockpit/pkg/extractors/json"
	"github.com/AdhityaRamadhanus/cockpit/pkg/http"
	"github.com/AdhityaRamadhanus/cockpit/pkg/keygenerator"
	"github.com/sirupsen/logrus"
)

type TaskJabar struct {
	KeyValueService cockpit.KeyValueService
	URL             config.URLConfig
	FeatureToggle   config.FeatureToggleConfig
}

func (t TaskJabar) SaveJabarProvincialLevelData() {
	defer func() {
		if err := recover(); err != nil {
			logrus.WithError(err.(error)).Errorf("Panic error %s", debug.Stack())
		}
	}()

	if !t.FeatureToggle.Jabar {
		logrus.Info("Skipping SaveJabarProvincialLevelData() Task")
		return
	}

	logrus.Info("Starting SaveJabarProvincialLevelData() Task")
	var arr []map[string]interface{}
	if err := http.GetJSON(t.URL.JabarCoreData, &arr); err != nil {
		logTaskError("task-jabar", err)
		return
	}

	// transform to map
	jsonInput := map[string]interface{}{}
	for _, row := range arr {
		date := row["tanggal"].(string)
		delete(row, "tanggal")
		jsonInput[date] = row
	}

	provincialLevelCases, err := extractor.ExtractProvincialLevelCasesJabar(jsonInput)
	if err != nil {
		logTaskError("task-jabar", err)
		return
	}

	jsonData, err := json.Marshal(provincialLevelCases)
	if err != nil {
		logTaskError("task-jabar", err)
		return
	}

	key := keygenerator.ProvincialLevelRedisKey("jawa-barat")
	if err := t.KeyValueService.Set(key, jsonData); err != nil {
		logTaskError("task-jabar", err)
		return
	}
	logrus.Info("Finished SaveJabarProvincialLevelData() Task")
}
