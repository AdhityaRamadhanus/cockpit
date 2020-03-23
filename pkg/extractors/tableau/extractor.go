package tableau

import (
	"encoding/json"
	"regexp"

	"github.com/AdhityaRamadhanus/cockpit"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
)

func transformRawSheetToJSONStr(rawSheet string) (map[string]string, error) {
	/* raw sheet is in form of "477574;{jsonstr}4545;{jsonstr}" so we need to remove the garbage string */
	var re = regexp.MustCompile(`(?m)[\d]+;\{`)
	matches := re.Split(rawSheet, -1)

	if len(matches) < 2 {
		return nil, errors.New("Error in cleaning raw sheet")
	}

	primaryJSONStr := "{" + matches[1]
	secondaryJSONStr := "{" + matches[2]

	// check as json
	if err := json.Unmarshal([]byte(primaryJSONStr), &map[string]interface{}{}); err != nil {
		return nil, errors.Wrap(err, "Invalid JSON")
	}

	if err := json.Unmarshal([]byte(secondaryJSONStr), &map[string]interface{}{}); err != nil {
		return nil, errors.Wrap(err, "Invalid JSON")
	}

	return map[string]string{
		"primary":   primaryJSONStr,
		"secondary": secondaryJSONStr,
	}, nil
}

// the most common type of sheet
func ExtractOneRowSheet(rawSheet string) (interface{}, error) {
	jsonStr, err := transformRawSheetToJSONStr(rawSheet)
	if err != nil {
		return -1, err
	}

	jsonPath := "secondaryInfo.presModelMap.dataDictionary.presModelHolder.genDataDictionaryPresModel.dataSegments.0.dataColumns.0.dataValues.0"
	result := gjson.Get(jsonStr["secondary"], jsonPath)

	return int(result.Int()), nil
}

func ExtractJakartaInfectedDetails(rawSheet string) (interface{}, error) {
	jsonStr, err := transformRawSheetToJSONStr(rawSheet)
	if err != nil {
		return -1, err
	}

	positifPath := "worldUpdate.applicationPresModel.workbookPresModel.dashboardPresModel.zones.75.zoneCommon.name"
	// dirawatPath := "worldUpdate.applicationPresModel.workbookPresModel.dashboardPresModel.zones.81.zoneCommon.name"
	sembuhPath := "worldUpdate.applicationPresModel.workbookPresModel.dashboardPresModel.zones.82.zoneCommon.name"
	meninggalPath := "worldUpdate.applicationPresModel.workbookPresModel.dashboardPresModel.zones.83.zoneCommon.name"
	// rawatSendiriPath := "worldUpdate.applicationPresModel.workbookPresModel.dashboardPresModel.zones.84.zoneCommon.name"

	positif := gjson.Get(jsonStr["primary"], positifPath)
	// dirawat := gjson.Get(jsonStr["primary"], dirawatPath)
	sembuh := gjson.Get(jsonStr["primary"], sembuhPath)
	meninggal := gjson.Get(jsonStr["primary"], meninggalPath)
	// rawatSendiri := gjson.Get(jsonStr["primary"], rawatSendiriPath)

	infectedDetails := cockpit.InfectedDetails{
		Total:     int(positif.Int()),
		Recovered: int(sembuh.Int()),
		Died:      int(meninggal.Int()),
	}

	return infectedDetails, nil
}

// func ExtractMappingPerKotaSheet(rawSheet string) (interface{}, error) {
// 	jsonStr, err := transformRawSheetToJSONStr(rawSheet)
// 	if err != nil {
// 		return -1, err
// 	}

// 	jsonPath := "secondaryInfo.presModelMap.dataDictionary.presModelHolder.genDataDictionaryPresModel.dataSegments.0.dataColumns"
// 	result := gjson.Get(jsonStr["secondary"], jsonPath)

// 	dataColumns := result.Array()

// 	dataLabels := []string{}
// 	dataValues := []int{}

// 	for _, dataColumn := range dataColumns {
// 		dataType := dataColumn.Map()["dataType"].String()
// 		values := dataColumn.Map()["dataValues"].Array()
// 		if dataType == "cstring" {
// 			for _, value := range values {
// 				dataLabels = append(dataLabels, value.String())
// 			}
// 		}

// 		if dataType == "integer" {
// 			for _, value := range values {
// 				dataValues = append(dataValues, int(value.Int()))
// 			}
// 		}
// 	}

// 	dataMap := map[string]int{}

// 	for i, label := range dataLabels {
// 		dataMap[label] = dataValues[i]
// 	}
// 	return dataMap, nil
// }
