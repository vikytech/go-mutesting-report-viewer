package main

import (
	"io"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	log.SetOutput(io.Discard)
}

func writeContent(t *testing.T, tempFile *os.File, content string) {
	file, _ := os.OpenFile(tempFile.Name(), os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if _, err := file.WriteString(content); err != nil {
		t.Fatal(err)
	}
	defer file.Close()
}

func createTempFile(t *testing.T, fileNamePattern string) *os.File {
	tempFile, err := os.CreateTemp("/tmp", fileNamePattern)
	if err != nil {
		t.Fatal(err)
	}
	defer tempFile.Close()

	return tempFile
}

func TestMain(t *testing.T) {
	t.Run("TestMain_ValidFile", func(t *testing.T) {
		jsonContent := `{
			"stats": {
				"totalMutantsCount": 10,
				"killedCount": 5
			},
			"escaped": [],
			"killed": []
		}`
		file := createTempFile(t, "*.json")
		filePath := file.Name()
		writeContent(t, file, jsonContent)

		reportOutputPath := "/tmp/testOutput.html"

		defer func() {
			r := recover()
			assert.Nil(t, r)

			if _, err := os.Stat(reportOutputPath); err == nil {
				err := os.Remove(reportOutputPath)
				assert.Nil(t, err, "Failed to remove temporary test output file")
			}

			err := os.Remove(filePath)
			assert.Nil(t, err, "Failed to remove temporary json file")
		}()

		os.Args = []string{"cmd", "-file", filePath, "-out", reportOutputPath, "-template", "report_test.tmpl"}
		main()

		fileContent, _ := os.ReadFile(reportOutputPath)

		expectedOutput := "<html><body> Total Mutant: 10 | Killed Count: 5 <body><html>"
		assert.Equal(t, expectedOutput, string(fileContent), "Report File content not matched")
	})
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

		data := readJson(filePath)
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

		readJson(filePath)
	})

	t.Run("TestReadJson_FileNotFound", func(t *testing.T) {
		defer func() {
			expectedError := "Error reading file: open nonexistent.json: no such file or directory"
			r := recover()
			assert.NotNil(t, r, "Expected Error, but test passed")
			assert.Equal(t, expectedError, r)
		}()

		readJson("nonexistent.json")
	})

	t.Run("TestReadJson_NoFilePath", func(t *testing.T) {
		defer func() {
			expectedError := "Error reading file: open : no such file or directory"
			r := recover()
			assert.NotNil(t, r, "Expected Error, but test passed")
			assert.Equal(t, expectedError, r)
		}()

		readJson("")
	})

}

func TestExecuteTemplate(t *testing.T) {

	t.Run("TestExecuteTemplateUnableToReadTemplate", func(t *testing.T) {
		data := Data{}

		defer func() {
			expectedError := "Unable to parse template file: template: pattern matches no files: `testTemplate.html`"
			r := recover()
			assert.NotNil(t, r, "Expected Error, but test passed")
			assert.Equal(t, expectedError, r)
		}()

		executeTemplate(data, "testTemplate.html", "/unknowpath/testOutput.html")
	})

	t.Run("TestExecuteTemplateUnableToCreateReportFile", func(t *testing.T) {
		data := Data{}

		defer func() {
			expectedError := "Unable to create report file: open /unknowpath/testOutput.html: no such file or directory"
			r := recover()
			assert.NotNil(t, r, "Expected Error, but test passed")
			assert.Equal(t, expectedError, r)
		}()

		executeTemplate(data, "report_test.tmpl", "/unknowpath/testOutput.html")
	})

	t.Run("TestUnableToParseTemplate", func(t *testing.T) {
		tempTemplateFileName := "report_test_parse_error.tmpl"
		reportOutputPath := "/tmp/testOutput.html"

		defer func() {
			if _, err := os.Stat(reportOutputPath); err == nil {
				if err := os.Remove(reportOutputPath); err != nil {
					t.Errorf("Failed to remove temporary test output file: %v", err)
				}
			}

			r := recover()
			expectedError := "Unable to parse template file: template: " + tempTemplateFileName + ":1: unexpected \"}\" in define clause"
			assert.NotNil(t, r, "Expected Error, but test passed")
			assert.Equal(t, expectedError, r)
		}()

		data := Data{}

		executeTemplate(data, tempTemplateFileName, reportOutputPath)
	})

	t.Run("TestUnableToExecuteTemplate", func(t *testing.T) {
		tempTemplateFileName := "report_test_execute_error.tmpl"
		reportOutputPath := "/tmp/testOutput.html"

		defer func() {
			if _, err := os.Stat(reportOutputPath); err == nil {
				if err := os.Remove(reportOutputPath); err != nil {
					t.Errorf("Failed to remove temporary test output file: %v", err)
				}
			}
			r := recover()
			expectedError := "Error executing template: template: " + tempTemplateFileName + ":1:32: executing \"" + tempTemplateFileName + "\" at <.Stats.totalMutantsCount>: can't evaluate field totalMutantsCount in type main.Stats"
			assert.NotNil(t, r, "Expected Error, but test passed")
			assert.Equal(t, expectedError, r)
		}()

		data := Data{Stats: Stats{TotalMutantsCount: 10}}

		executeTemplate(data, tempTemplateFileName, reportOutputPath)
	})

	t.Run("TestShouldBeAbleToExecuteTemplate", func(t *testing.T) {
		tempTemplateFileName := "report_test.tmpl"
		reportOutputPath := "/tmp/testOutput.html"

		defer func() {
			if _, err := os.Stat(reportOutputPath); err == nil {
				if err := os.Remove(reportOutputPath); err != nil {
					t.Errorf("Failed to remove temporary test output file: %v", err)
				}
			}
			r := recover()
			assert.Nil(t, r, "Expected test to pass, but threw err")
		}()

		data := Data{}

		executeTemplate(data, tempTemplateFileName, reportOutputPath)
		fileContent, _ := os.ReadFile(reportOutputPath)

		expectedOutput := "<html><body> Total Mutant: 0 | Killed Count: 0 <body><html>"
		assert.Equal(t, expectedOutput, string(fileContent), "Report File content not matched")
	})
}
