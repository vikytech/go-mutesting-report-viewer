package models

type Report struct {
	GlobalStats  Stats
	ReportDetail map[string]ReportDetails
}
