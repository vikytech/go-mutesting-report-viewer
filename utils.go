package main

import (
	"encoding/json"
	models "gomutestingrhtmlreport/models"
	"log"
	"os"
	"strings"
)

func readSingleJsonReportFile(filePath string) models.Data {
	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		log.Panicf("Error reading file: %s", err.Error())
	}

	var data models.Data
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		panic("Invalid JSON format: " + err.Error())
	}
	return data
}

func dataToReportMapper(data models.Data) models.Report {
	fileMap := make(map[string]models.ReportDetails)
	escaped := data.Escaped
	killed := data.Killed

	for _, escapedMutantEntry := range escaped {
		checksum := strings.Split(escapedMutantEntry.ProcessOutput, " ")[4]
		entry := models.MutatorDetail{MutatorName: escapedMutantEntry.Mutator.MutatorName, Diff: escapedMutantEntry.Diff, Checksum: checksum}
		escapedEntry := fileMap[escapedMutantEntry.Mutator.OriginalFilePath].Escaped
		updatedEntry := append(escapedEntry, entry)
		fileMap[escapedMutantEntry.Mutator.OriginalFilePath] = models.ReportDetails{Escaped: updatedEntry}
	}

	for _, killedMutantEntry := range killed {
		checksum := strings.Split(killedMutantEntry.ProcessOutput, " ")[4]
		entry := models.MutatorDetail{MutatorName: killedMutantEntry.Mutator.MutatorName, Diff: killedMutantEntry.Diff, Checksum: checksum}
		killedEntry := fileMap[killedMutantEntry.Mutator.OriginalFilePath].Killed
		updatedEntry := append(killedEntry, entry)
		fileMap[killedMutantEntry.Mutator.OriginalFilePath] = models.ReportDetails{Escaped: fileMap[killedMutantEntry.Mutator.OriginalFilePath].Escaped, Killed: updatedEntry}
	}

	return models.Report{GlobalStats: data.Stats, ReportDetail: fileMap}
}
