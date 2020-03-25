package cockpit

import "time"

type MonitoringDetails struct {
	Total    int `json:"total"`   // ODP
	Finished int `json:"selesai"` // Selesai ODP
	Ongoing  int `json:"proses"`  // Proses ODP
}

type SurveillanceDetails struct {
	Total    int `json:"total"`   // PDP
	Finished int `json:"selesai"` // Selesai PDP
	Ongoing  int `json:"proses"`  // Proses PDP
}

type InfectedDetails struct {
	Total     int `json:"total"`     // Positif COVID-19
	Recovered int `json:"sembuh"`    // Sembuh
	Died      int `json:"meninggal"` // Meninggal
}

type ProvincialLevelCases struct {
	MonitoringCases   MonitoringDetails   `json:"odp"`     // ODP
	SurveillanceCases SurveillanceDetails `json:"pdp"`     // PDP
	InfectedCases     InfectedDetails     `json:"positif"` // Positif COVID-19
	LastUpdatedAt     time.Time           `json:"last_updated_at"`
	LastFetchedAt     time.Time           `json:"last_fetched_at"`
}

type CityLevelCases struct {
	Name          string
	TotalODP      int       `json:"total_odp"`
	TotalPDP      int       `json:"total_pdp"`
	TotalPositif  int       `json:"total_positif"`
	LastUpdatedAt time.Time `json:"last_updated_at"`
	LastFetchedAt time.Time `json:"last_fetched_at"`
}

type Hospital struct {
	Latitude  float64
	Longitude float64
	Name      string
	Address   string
	Telephone string
}

type ProvinceHospitalReferences struct {
	Hospitals []Hospital
}

type NationalLevelCases struct {
	InfectedCases InfectedDetails `json:"positif"` // Positif COVID-19
}
