package tasks

import "github.com/AdhityaRamadhanus/cockpit"

type Job struct {
	WorkBook  string
	Dashboard string
	SheetName string
	SheetID   string
	SessionID string
	Extractor cockpit.Extractor
}

type JobResult struct {
	Job    Job
	Result interface{}
	Error  error
}
