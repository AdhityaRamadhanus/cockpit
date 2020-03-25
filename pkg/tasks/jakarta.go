package tasks

import (
	"encoding/json"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/AdhityaRamadhanus/cockpit"
	"github.com/AdhityaRamadhanus/cockpit/pkg/config"
	extractor "github.com/AdhityaRamadhanus/cockpit/pkg/extractors/tableau"
	"github.com/AdhityaRamadhanus/cockpit/pkg/http"
	"github.com/AdhityaRamadhanus/cockpit/pkg/keygenerator"
	"github.com/AdhityaRamadhanus/cockpit/pkg/tableau"
	"github.com/AdhityaRamadhanus/cockpit/pkg/tableau/sheets"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.uber.org/multierr"
)

type TaskJakarta struct {
	KeyValueService cockpit.KeyValueService
	Tableau         config.TableauConfig
	FeatureToggle   config.FeatureToggleConfig
}

func (t TaskJakarta) SaveJakartaProvincialLevelData() {
	defer func() {
		if err := recover(); err != nil {
			logrus.WithError(err.(error)).Errorf("Panic error %s", debug.Stack())
		}
	}()

	if !t.FeatureToggle.Jakarta {
		logrus.Info("Skipping SaveJakartaProvincialLevelData() Task")
		return
	}

	logrus.Info("Starting SaveJakartaProvincialLevelData() Task")
	workbooks := t.Tableau.Workbooks
	dashboards := t.Tableau.Dashboards
	sessionIDs := []string{}

	for i, workBook := range workbooks {
		sessionID, err := tableau.GetSessionID(workBook, dashboards[i])
		logrus.Debugf("Workbook: %s, Dashboard: %s, Session ID: %s\n", workBook, dashboards[i], sessionID)
		if err != nil {
			logTaskError("task-jakarta", err)
			return
		}
		sessionIDs = append(sessionIDs, sessionID)
	}

	jobs := []Job{
		Job{
			WorkBook:  workbooks[0],
			Dashboard: dashboards[0],
			SessionID: sessionIDs[0],
			SheetName: "Selesai ODP",
			SheetID:   sheets.JakartaSheetIDsDashboard1["SelesaiODP"],
			Extractor: extractor.ExtractOneRowSheet,
		},
		Job{
			WorkBook:  workbooks[0],
			Dashboard: dashboards[0],
			SessionID: sessionIDs[0],
			SheetName: "Proses ODP",
			SheetID:   sheets.JakartaSheetIDsDashboard1["ProsesODP"],
			Extractor: extractor.ExtractOneRowSheet,
		},
		Job{
			WorkBook:  workbooks[0],
			Dashboard: dashboards[0],
			SessionID: sessionIDs[0],
			SheetName: "Selesai PDP",
			SheetID:   sheets.JakartaSheetIDsDashboard1["SelesaiPDP"],
			Extractor: extractor.ExtractOneRowSheet,
		},
		Job{
			WorkBook:  workbooks[0],
			Dashboard: dashboards[0],
			SessionID: sessionIDs[0],
			SheetName: "Proses PDP",
			SheetID:   sheets.JakartaSheetIDsDashboard1["ProsesPDP"],
			Extractor: extractor.ExtractOneRowSheet,
		},
		Job{
			WorkBook:  workbooks[1],
			Dashboard: dashboards[1],
			SessionID: sessionIDs[1],
			SheetName: "Dashboard 2",
			SheetID:   "Dashboard%202",
			Extractor: extractor.ExtractJakartaInfectedDetails,
		},
	}

	jobResultChan := make(chan JobResult, len(jobs))
	for _, job := range jobs {
		go func(jobResultChan chan JobResult, job Job) {
			jobResult := JobResult{
				Job: job,
			}

			url := fmt.Sprintf("https://public.tableau.com/vizql/w/%s/v/%s/bootstrapSession/sessions/%s", job.WorkBook, job.Dashboard, job.SessionID)
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
		logTaskError("task-jakarta", errs)
		return
	}

	finishedODP := dataMap["Selesai ODP"].(int)
	ongoingODP := dataMap["Proses ODP"].(int)
	totalODP := finishedODP + ongoingODP

	finishedPDP := dataMap["Selesai PDP"].(int)
	ongoingPDP := dataMap["Proses PDP"].(int)
	totalPDP := finishedPDP + ongoingPDP

	provincialLevelCases := cockpit.ProvincialLevelCases{
		MonitoringCases: cockpit.MonitoringDetails{
			Total:    totalODP,
			Finished: finishedODP,
			Ongoing:  ongoingODP,
		},
		SurveillanceCases: cockpit.SurveillanceDetails{
			Total:    totalPDP,
			Finished: finishedPDP,
			Ongoing:  ongoingPDP,
		},
		InfectedCases: dataMap["Dashboard 2"].(cockpit.InfectedDetails),
		LastFetchedAt: time.Now(),
	}

	// fetch lastUpdate workbook
	now := time.Now().UnixNano()
	url := fmt.Sprintf("https://public.tableau.com/profile/api/single_workbook/%s?no_cache=%v", t.Tableau.Workbooks[0], now)
	lastUpdateResp := struct {
		LastUpdatedMs int64 `json:"lastUpdateDate"`
	}{}
	if err := http.GetJSON(url, &lastUpdateResp); err != nil {
		logTaskError("task-jakarta", errors.Wrapf(err, "Failed to get last update at"))
	}

	provincialLevelCases.LastUpdatedAt = time.Unix(0, lastUpdateResp.LastUpdatedMs*int64(1000000))

	jsonData, err := json.Marshal(provincialLevelCases)
	if err != nil {
		logTaskError("task-jakarta", err)
		return
	}

	key := keygenerator.ProvincialLevelRedisKey("dki-jakarta")
	if err := t.KeyValueService.Set(key, jsonData); err != nil {
		logTaskError("task-jakarta", err)
		return
	}
	logrus.Info("Finished SaveJakartaProvincialLevelData() Task")
}
