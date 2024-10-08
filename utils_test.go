package main

import (
	models "gomutestingrhtmlreport/models"
	"io"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	log.SetOutput(io.Discard)
}

func TestReadJSON(t *testing.T) {
	t.Run("TestReadJson_ValidJson", func(t *testing.T) {
		jsonContent := `{
		"stats": {
			"totalMutantsCount": 10,
			"killedCount": 5,
			"msi": 0.50,
			"extraField": "shouldBeIgnored"

		},
		"escaped": [],
		"killed": []
	}`

		file := createTempFile(t, "*.json")
		filePath := file.Name()
		writeContent(t, file, jsonContent)
		defer os.Remove(filePath)

		data := readSingleJsonReportFile(filePath)
		assert.Equal(t, 10, data.Stats.TotalMutantsCount, "TotalMutantsCount")
		assert.Equal(t, 5, data.Stats.KilledCount, "KilledCount")
		assert.Equal(t, 0.50, data.Stats.Msi, "MSI")
	})
	t.Run("TestReadJson_InvalidJson", func(t *testing.T) {
		jsonContent := `{invalid json}`

		file := createTempFile(t, "*.json")
		filePath := file.Name()
		writeContent(t, file, jsonContent)
		defer os.Remove(filePath)

		defer func() {
			expectedError := "Invalid JSON format: invalid character 'i' looking for beginning of object key string"
			r := recover()
			assert.NotNil(t, r, "Expected Error, but test passed")
			assert.Equal(t, expectedError, r)
		}()

		readSingleJsonReportFile(filePath)
	})

	t.Run("TestReadJson_FileNotFound", func(t *testing.T) {
		defer func() {
			expectedError := "Error reading file: open nonexistent.json: no such file or directory"
			r := recover()
			assert.NotNil(t, r, "Expected Error, but test passed")
			assert.Equal(t, expectedError, r)
		}()

		readSingleJsonReportFile("nonexistent.json")
	})

	t.Run("TestReadJson_NoFilePath", func(t *testing.T) {
		defer func() {
			expectedError := "Error reading file: open : no such file or directory"
			r := recover()
			assert.NotNil(t, r, "Expected Error, but test passed")
			assert.Equal(t, expectedError, r)
		}()

		readSingleJsonReportFile("")
	})
}

func TestGroupByFiles(t *testing.T) {
	t.Run("TestGroupFilesByName", func(t *testing.T) {
		data := models.Data{
			Stats: models.Stats{TotalMutantsCount: 10},
			Escaped: []models.Entry{
				{Mutator: models.Mutator{OriginalFilePath: "testFile.go"}, ProcessOutput: "Pass testFile.go with checksum randomchecksumescaped"},
				{Mutator: models.Mutator{OriginalFilePath: "testFile1.go"}, ProcessOutput: "Pass testFile1.go with checksum randomchecksumescaped1"},
			},
			Killed: []models.Entry{
				{Mutator: models.Mutator{OriginalFilePath: "testFile.go"}, ProcessOutput: "Pass testFile.go with checksum randomchecksumkilled"},
				{Mutator: models.Mutator{OriginalFilePath: "testFile1.go"}, ProcessOutput: "Pass testFile.go with checksum randomchecksumkilled1"},
			},
		}
		report := dataToReportMapper(data)
		expectedFileMap := make(map[string]models.ReportDetails)
		escapedFile1 := []models.MutatorDetail{{Checksum: "randomchecksumescaped"}}
		escapedFile2 := []models.MutatorDetail{{Checksum: "randomchecksumescaped1"}}
		killedFile1 := []models.MutatorDetail{{Checksum: "randomchecksumkilled"}}
		killedFile2 := []models.MutatorDetail{{Checksum: "randomchecksumkilled1"}}
		expectedFileMap["testFile.go"] = models.ReportDetails{Escaped: escapedFile1, Killed: killedFile1}
		expectedFileMap["testFile1.go"] = models.ReportDetails{Escaped: escapedFile2, Killed: killedFile2}
		expectedReport := models.Report{
			Stats:        models.Stats{TotalMutantsCount: 10},
			ReportDetail: expectedFileMap,
		}
		assert.Equal(t, expectedReport, report)
	})
}
