package models

type Stats struct {
	TotalMutantsCount    int     `json:"totalMutantsCount"`
	KilledCount          int     `json:"killedCount"`
	NotCoveredCount      int     `json:"notCoveredCount"`
	EscapedCount         int     `json:"escapedCount"`
	ErrorCount           int     `json:"errorCount"`
	SkippedCount         int     `json:"skippedCount"`
	TimeOutCount         int     `json:"timeOutCount"`
	Msi                  float64 `json:"msi"`
	MutationCodeCoverage float64 `json:"mutationCodeCoverage"`
	CoveredCodeMsi       float64 `json:"coveredCodeMsi"`
}
