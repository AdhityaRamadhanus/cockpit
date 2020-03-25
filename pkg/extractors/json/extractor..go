package json

import (
	"strconv"
	"time"

	"github.com/AdhityaRamadhanus/cockpit"
	"github.com/pkg/errors"
	"go.uber.org/multierr"
)

func ExtractProvincialLevelCasesJabar(jsonData map[string]interface{}) (interface{}, error) {
	now := time.Now()
	dateStrs := []string{}
	// generate ["2020-21-03", "2020-20-03", "2020-19-03"]
	for i := 0; i <= 3; i++ {
		dateStrs = append(dateStrs, now.AddDate(0, 0, -1*i).Format("02-01-2006"))
	}

	dataMap := map[string]interface{}{}
	var foundDate string

	for _, dateStr := range dateStrs {
		dataRow, ok := jsonData[dateStr]
		if ok && dataRow.(map[string]interface{})["total_odp"] != nil {
			foundDate = dateStr
			dataMap = dataRow.(map[string]interface{})
			break
		}
	}

	// couldn't get any latest date within 3 days
	if len(dataMap) == 0 {
		return nil, errors.New("Couldn't get any match within 3 days")
	}

	totalODP, errTotalODP := strconv.Atoi(dataMap["total_odp"].(string))
	finishedODP, errFinishedODP := strconv.Atoi(dataMap["selesai_pemantauan"].(string))
	ongoingODP, errOngoingODP := strconv.Atoi(dataMap["proses_pemantauan"].(string))

	totalPDP, errTotalPDP := strconv.Atoi(dataMap["total_pdp"].(string))
	finishedPDP, errFinishedPDP := strconv.Atoi(dataMap["selesai_pengawasan"].(string))
	ongoingPDP, errOngoingPDP := strconv.Atoi(dataMap["proses_pengawasan"].(string))

	totalInfected := int(dataMap["total_positif_saat_ini"].(float64))
	recoveredCases := int(dataMap["total_sembuh"].(float64))
	deathCases := int(dataMap["total_meninggal"].(float64))

	errs := multierr.Combine(errTotalODP, errFinishedODP, errOngoingODP, errTotalPDP, errFinishedPDP, errOngoingPDP)
	if errs != nil {
		return nil, errors.Wrapf(errs, "Failed to parse data row columns %v", dataMap)
	}

	lastUpdatedAt, _ := time.Parse("02-01-2006", foundDate)
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
		InfectedCases: cockpit.InfectedDetails{
			Total:     totalInfected,
			Died:      deathCases,
			Recovered: recoveredCases,
		},
		LastFetchedAt: time.Now(),
		LastUpdatedAt: lastUpdatedAt,
	}

	return provincialLevelCases, nil
}
