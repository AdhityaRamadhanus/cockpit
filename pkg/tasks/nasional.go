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
	"github.com/AdhityaRamadhanus/cockpit/pkg/tableau/sheets"
	"github.com/sirupsen/logrus"
	"go.uber.org/multierr"
)

type TaskIndonesia struct {
	KeyValueService cockpit.KeyValueService
	Tableau         config.TableauConfig
	FeatureToggle   config.FeatureToggleConfig
}

func (t TaskIndonesia) SaveIndonesiaNationalLevelData() {
	defer func() {
		if err := recover(); err != nil {
			logrus.WithError(err.(error)).Errorf("Panic error %s", debug.Stack())
		}
	}()

	workBook := t.Tableau.Workbooks[0]
	dashboard := t.Tableau.Dashboards[0]
	sessionID, err := tableau.GetSessionID(workBook, dashboard)

	if !t.FeatureToggle.Indonesia {
		logrus.Info("Skipping SaveIndonesiaNationalLevelData() Task")
		return
	}

	logrus.Info("Starting SaveIndonesiaNationalLevelData() Task")

	logrus.Debugf("Workbook: %s, Dashboard: %s, Session ID: %s\n", workBook, dashboard, sessionID)
	if err != nil {
		logTaskError("task-indonesia", err)
		return
	}

	jobs := []Job{
		Job{
			SheetName: "Positif Nasional",
			SheetID:   sheets.JakartaSheetIDsDashboard1["PositifNasional"],
			Extractor: extractor.ExtractOneRowSheet,
		},
		Job{
			SheetName: "Meninggal Nasional",
			SheetID:   sheets.JakartaSheetIDsDashboard1["MeninggalNasional"],
			Extractor: extractor.ExtractOneRowSheet,
		},
		Job{
			SheetName: "Sembuh Nasional",
			SheetID:   sheets.JakartaSheetIDsDashboard1["SembuhNasional"],
			Extractor: extractor.ExtractOneRowSheet,
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
		logTaskError("task-indonesia", errs)
		return
	}

	nationalLevelCases := cockpit.NationalLevelCases{
		InfectedCases: cockpit.InfectedDetails{
			Total:     dataMap["Positif Nasional"].(int),
			Recovered: dataMap["Sembuh Nasional"].(int),
			Died:      dataMap["Meninggal Nasional"].(int),
		},
	}

	jsonData, err := json.Marshal(nationalLevelCases)
	if err != nil {
		logTaskError("task-indonesia", err)
		return
	}

	key := keygenerator.NationalLevelRedisKey()
	if err := t.KeyValueService.Set(key, jsonData); err != nil {
		logTaskError("task-indonesia", err)
		return
	}
	logrus.Info("Finished SaveIndonesiaNationalLevelData() Task")
}
