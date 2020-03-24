package tasks

import (
	"encoding/json"
	"fmt"
	"runtime/debug"

	"github.com/AdhityaRamadhanus/cockpit"
	"github.com/AdhityaRamadhanus/cockpit/pkg/config"
	extractor "github.com/AdhityaRamadhanus/cockpit/pkg/extractors/tableau"
	"github.com/AdhityaRamadhanus/cockpit/pkg/keygenerator"
	"github.com/AdhityaRamadhanus/cockpit/pkg/tableau"
	"github.com/sirupsen/logrus"
	"go.uber.org/multierr"
)

type TaskJogja struct {
	KeyValueService cockpit.KeyValueService
	Tableau         config.TableauConfig
	FeatureToggle   config.FeatureToggleConfig
}

func (t TaskJogja) SaveJogjaProvincialLevelData() {
	defer func() {
		if err := recover(); err != nil {
			logrus.WithError(err.(error)).Errorf("Panic error %s", debug.Stack())
		}
	}()

	workBook := "Covid-19DIY"
	dashboard := "JumlahStatusbyTable"
	sessionID, err := tableau.GetSessionID(workBook, dashboard)

	if !t.FeatureToggle.Jogja {
		logrus.Info("Skipping SaveJogjaProvincialLevelData() Task")
		return
	}

	logrus.Info("Starting SaveJogjaProvincialLevelData() Task")

	logrus.Debugf("Workbook: %s, Dashboard: %s, Session ID: %s\n", workBook, dashboard, sessionID)
	if err != nil {
		logTaskError("task-jogja", err)
		return
	}

	jobs := []Job{
		Job{
			SheetName: "Kasus Jogja",
			SheetID:   "Jumlah%20Status%20by%20Table",
			Extractor: extractor.ExtractJogjaProvincialLevelCases,
		},
	}

	jobResultChan := make(chan JobResult, len(jobs))
	for _, job := range jobs {
		go func(jobResultChan chan JobResult, job Job) {
			jobResult := JobResult{
				Job: job,
			}

			url := fmt.Sprintf("https://public.tableau.com/vizql/w/%s/v/%s/bootstrapSession/sessions/%s", workBook, dashboard, sessionID)
			sheet, err := tableau.GetSheet(url, job.SheetID)
			if err != nil {
				jobResult.Error = err
			} else {
				res, err := job.Extractor(sheet)
				if err != nil {
					jobResult.Error = err
				} else {
					jobResult.Result = res
				}
			}

			jobResultChan <- jobResult

		}(jobResultChan, job)
	}

	dataMap := map[string]interface{}{}
	var errs error
	for range jobs {
		jobResult := <-jobResultChan
		if jobResult.Error != nil {
			errs = multierr.Append(errs, jobResult.Error)
			continue
		}
		dataMap[jobResult.Job.SheetName] = jobResult.Result
	}

	if errs != nil {
		logTaskError("task-jogja", errs)
		return
	}

	provincialLevelCases := dataMap["Kasus Jogja"].(cockpit.ProvincialLevelCases)

	jsonData, err := json.Marshal(provincialLevelCases)
	if err != nil {
		logTaskError("task-jogja", err)
		return
	}

	key := keygenerator.ProvincialLevelRedisKey("di-jogjakarta")
	if err := t.KeyValueService.Set(key, jsonData); err != nil {
		logTaskError("task-jogja", err)
		return
	}
	logrus.Info("Finished SaveJogjaProvincialLevelData() Task")
}
