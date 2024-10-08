package models

type ReportDetails struct {
	FileStats Stats
	Escaped   []MutatorDetail
	Killed    []MutatorDetail
}
