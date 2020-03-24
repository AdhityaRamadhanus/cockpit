package html

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/AdhityaRamadhanus/cockpit"
	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	"go.uber.org/multierr"
)

func ExtractProvincialLevelCasesJatim(htmlText string) (interface{}, error) {
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlText))
	if err != nil {
		return nil, errors.Wrapf(err, "goquery.NewDocumentFromReader() err;")
	}

	selector := "table tbody tr"

	totalODP := 0
	totalPDP := 0
	totalCases := 0

	dataRowStrs := [][]string{}
	// Find the review items
	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		dataRowStrs = append(dataRowStrs, s.Find("td").Map(func(j int, si *goquery.Selection) string {
			return si.Text()
		}))
	})

	if len(dataRowStrs) == 0 {
		return nil, errors.Errorf("Couldn't match element matching selector %q", selector)
	}

	for _, dataRowStr := range dataRowStrs {
		if len(dataRowStr) < 4 {
			return nil, errors.Errorf("Data row contains column less than 4, dataRowStr: %v", dataRowStr)
		}

		odp, errODP := strconv.Atoi(dataRowStr[1])
		pdp, errPDP := strconv.Atoi(dataRowStr[2])
		positif, errPositif := strconv.Atoi(dataRowStr[3])
		// TODO get timestamps
		// lastUpdate, _ := time.Parse("2006-01-02 15:04:05", dataRowStr[4])

		if err := multierr.Combine(errODP, errPDP, errPositif); err != nil {
			return nil, errors.Wrapf(err, "Failed to parse data row columns, dataRowStr: %v", dataRowStr)
		}

		totalODP += odp
		totalPDP += pdp
		totalCases += positif
	}

	provincialLevelCases := cockpit.ProvincialLevelCases{
		MonitoringCases: cockpit.MonitoringDetails{
			Total:    totalODP,
			Finished: -1,
			Ongoing:  -1,
		},
		SurveillanceCases: cockpit.SurveillanceDetails{
			Total:    totalPDP,
			Finished: -1,
			Ongoing:  -1,
		},
		InfectedCases: cockpit.InfectedDetails{
			Total:     totalCases,
			Died:      -1,
			Recovered: -1,
		},
	}

	return provincialLevelCases, nil
}

func ExtractProvincialLevelCasesJateng(htmlText string) (interface{}, error) {
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlText))
	if err != nil {
		return nil, errors.Wrapf(err, "goquery.NewDocumentFromReader() err;")
	}

	selector := "#features.features .card-text"
	dataLabels := []string{"odp", "pdp", "positif", "dirawat", "sembuh", "meninggal"}
	dataValues := []string{}
	dataMap := map[string]int{}

	var errs error
	// Find the review items
	dataValues = doc.Find(selector).Map(func(i int, s *goquery.Selection) string {
		// transform "2.416" to "2416"
		return strings.Replace(s.Text(), ".", "", -1)
	})

	if len(dataValues) == 0 {
		return nil, errors.Errorf("Couldn't match element matching selector %q", selector)
	}

	for i, dataValue := range dataValues {
		dataMap[dataLabels[i]], err = strconv.Atoi(dataValue)
		errs = multierr.Append(errs, err)
	}

	if errs != nil {
		return nil, errors.Wrapf(errs, "Failed to parse data row columns")
	}

	if len(dataMap) != 6 {
		return nil, errors.New("Data row is less than 6")
	}

	provincialLevelCases := cockpit.ProvincialLevelCases{
		MonitoringCases: cockpit.MonitoringDetails{
			Total:    dataMap["odp"],
			Finished: -1,
			Ongoing:  -1,
		},
		SurveillanceCases: cockpit.SurveillanceDetails{
			Total:    dataMap["pdp"],
			Finished: -1,
			Ongoing:  -1,
		},
		InfectedCases: cockpit.InfectedDetails{
			Total:     dataMap["positif"],
			Died:      dataMap["meninggal"],
			Recovered: dataMap["sembuh"],
		},
	}

	return provincialLevelCases, nil
}

func ExtractProvincialLevelCasesBanten(htmlText string) (interface{}, error) {
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlText))
	if err != nil {
		return nil, errors.Wrapf(err, "goquery.NewDocumentFromReader() err;")
	}

	titleMap := map[string]string{
		"ORANG DALAM PEMANTAUAN (ODP)":             "odp",
		"PASIEN DALAM PENGAWASAN (PDP)":            "pdp",
		"KASUS TERKONFIRMASI COVID-19 PROV BANTEN": "positif",
	}
	regexMap := map[string]interface{}{
		"odp": map[string]string{
			"total":   (`(?m)(\d+)\sTOTAL\sODP`),
			"dirawat": (`(?m)(\d+)\sMASIH\sDIPANTAU`),
			"sembuh":  (`(?m)(\d+)\sSEMBUH`),
		},
		"pdp": map[string]string{
			"total":   (`(?m)(\d+)\sTOTAL\sPDP`),
			"dirawat": (`(?m)(\d+)\sMASIH\sDIRAWAT`),
			"sembuh":  (`(?m)(\d+)\sSEMBUH`),
		},
		"positif": map[string]string{
			"total":     (`(?m)(\d+)\sKASUS\sPOSITIF`),
			"dirawat":   (`(?m)(\d+)\sDIRAWAT`),
			"sembuh":    (`(?m)(\d+)\sSEMBUH`),
			"meninggal": (`(?m)(\d+)\sMENINGGAL`),
		},
	}
	selector := ".content-row-no-bg .container .box"
	dataMap := map[string]int{}

	var errs error

	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		boxTitle := s.Find("h4").Text()
		boxText := s.Find("b").Text()

		regexMapTitle := titleMap[boxTitle]
		regexs, ok := regexMap[regexMapTitle]
		if ok {
			for regexName, regex := range regexs.(map[string]string) {
				re := regexp.MustCompile(regex)
				matches := re.FindStringSubmatch(boxText)
				if len(matches) < 2 {
					errs = multierr.Combine(errs, errors.New("Couldn't parse data row"))
				} else {
					dataMap[fmt.Sprintf("%s_%s", regexMapTitle, regexName)], err = strconv.Atoi(matches[1])
					errs = multierr.Combine(errs, err)
				}
			}
		}
	})

	if errs != nil {
		return nil, errors.Wrapf(errs, "Failed to parse data row columns")
	}

	provincialLevelCases := cockpit.ProvincialLevelCases{
		MonitoringCases: cockpit.MonitoringDetails{
			Total:    dataMap["odp_total"],
			Finished: dataMap["odp_sembuh"],
			Ongoing:  dataMap["odp_dirawat"],
		},
		SurveillanceCases: cockpit.SurveillanceDetails{
			Total:    dataMap["pdp_total"],
			Finished: dataMap["pdp_sembuh"],
			Ongoing:  dataMap["pdp_dirawat"],
		},
		InfectedCases: cockpit.InfectedDetails{
			Total:     dataMap["positif_total"],
			Died:      dataMap["positif_meninggal"],
			Recovered: dataMap["positif_sembuh"],
		},
	}

	return provincialLevelCases, nil
}
