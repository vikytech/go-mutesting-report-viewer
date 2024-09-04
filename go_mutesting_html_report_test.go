package main

import (
	"io"
	"log"
	"os"
	"strings"
	"testing"
)

func setupSuite(tb testing.TB, teardown func(tb testing.TB)) {
	log.Println("Setup suite")
	defer teardown(tb)

}

func teardownSuite(tb testing.TB) {
	log.Println("Teardown suite", os.Args)
	defer func() {
		if r := recover(); r != nil {
			tb.Errorf("Error:: %s", r)
		}
	}()
}

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
		setupSuite(t, teardownSuite)
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

		tempTemplateFile := createTempFile(t, "testTemplate*.html")
		tempTemplateFileName := tempTemplateFile.Name()
		reportOutputPath := "/tmp/testOutput.html"

		writeContent(t, tempTemplateFile, `<html><body> Total Mutant: {{ .Stats.TotalMutantsCount}} | Killed Count: {{ .Stats.KilledCount}} <body><html>`)

		defer func() {
			r := recover()
			if r != nil {
				t.Errorf("\nExpected test to pass, but threw error: %s", r)
			}

			if err := os.Remove(tempTemplateFileName); err != nil {
				t.Errorf("Failed to remove temporary template file: %v", err)
			}

			if _, err := os.Stat(reportOutputPath); err == nil {
				if err := os.Remove(reportOutputPath); err != nil {
					t.Errorf("Failed to remove temporary test output file: %v", err)
				}
			}

			if err := os.Remove(filePath); err != nil {
				t.Errorf("Failed to remove temporary json file: %v", err)
			}
		}()

		os.Args = []string{"cmd", "-file", filePath, "-out", reportOutputPath, "-template", tempTemplateFileName}
		main()

		fileContent, _ := os.ReadFile(reportOutputPath)

		expectedOutput := "<html><body> Total Mutant: 10 | Killed Count: 5 <body><html>"
		if string(fileContent) != expectedOutput {
			t.Errorf("Report File content not matched\nExpected:: %s\nActual::%s", expectedOutput, fileContent)
		}
	})
}

func TestReadJSON(t *testing.T) {
	t.Run("TestReadJson_ValidJson", func(t *testing.T) {
		setupSuite(t, teardownSuite)
		jsonContent := `{
		"stats": {
			"totalMutantsCount": 10,
			"killedCount": 5,
			"msi": 50.0,
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
		if data.Stats.TotalMutantsCount != 10 {
			t.Errorf("Expected 10 TotalMutantsCount, got %d", data.Stats.TotalMutantsCount)
		}
		if data.Stats.KilledCount != 5 {
			t.Errorf("Expected 5 KilledCount, got %d", data.Stats.KilledCount)
		}
		if data.Stats.Msi != 50.0 {
			t.Errorf("Expected 50.0 MSI, got %f", data.Stats.Msi)
		}
	})
	t.Run("TestReadJson_InvalidJson", func(t *testing.T) {
		setupSuite(t, teardownSuite)
		jsonContent := `{invalid json}`

		file := createTempFile(t, "*.json")
		filePath := file.Name()
		writeContent(t, file, jsonContent)
		defer os.Remove(filePath)

		defer func() {
			expectedError := "Invalid JSON format: invalid character 'i' looking for beginning of object key string"
			r := recover()
			if r == nil {
				t.Error("\nExpected Error, but test passed")
			} else if r != nil {
				if r != expectedError {
					t.Errorf("\nExpected:: %s, Got:: %s", expectedError, r)
				}
			}
		}()

		readJson(filePath)
	})

	t.Run("TestReadJson_FileNotFound", func(t *testing.T) {
		setupSuite(t, teardownSuite)
		defer func() {
			expectedError := "Error reading file: open nonexistent.json: no such file or directory"
			r := recover()
			if r == nil {
				t.Error("\nExpected Error, but test passed")
			} else if r != nil {
				if r != expectedError {
					t.Errorf("\nExpected:: %s, Got:: %s", expectedError, r)
				}
			}
		}()

		readJson("nonexistent.json")
	})

	t.Run("TestReadJson_NoFilePath", func(t *testing.T) {
		setupSuite(t, teardownSuite)
		defer func() {
			expectedError := "Error reading file: open : no such file or directory"
			r := recover()
			if r == nil {
				t.Error("\nExpected Error, but test passed")
			} else if r != nil {
				if r != expectedError {
					t.Errorf("\nExpected:: %s, Got:: %s", expectedError, r)
				}
			}
		}()

		readJson("")
	})

}

func TestExecuteTemplate(t *testing.T) {

	t.Run("TestExecuteTemplateUnableToReadTemplate", func(t *testing.T) {
		setupSuite(t, teardownSuite)
		data := Data{}

		defer func() {
			expectedError := "Unable to parse template file: open testTemplate.html: no such file or directory"
			r := recover()
			if r == nil {
				t.Error("\nExpected Error, but test passed")
			} else if r != nil {
				if r != expectedError {
					t.Errorf("\nExpected:: %s, Got:: %s", expectedError, r)
				}
			}
		}()

		executeTemplate(data, "testTemplate.html", "/unknowpath/testOutput.html")
	})

	t.Run("TestExecuteTemplateUnableToCreateReportFile", func(t *testing.T) {
		setupSuite(t, teardownSuite)
		tempTemplateFile := createTempFile(t, "testTemplate.html")

		defer func() {
			if err := os.Remove(tempTemplateFile.Name()); err != nil {
				t.Errorf("Failed to remove temporary template file: %v", err)
			}
		}()

		data := Data{}

		defer func() {
			expectedError := "Unable to create report file: open /unknowpath/testOutput.html: no such file or directory"
			r := recover()
			if r == nil {
				t.Error("\nExpected Error, but test passed")
			} else if r != nil {
				if r != expectedError {
					t.Errorf("\nExpected:: %s, Got:: %s", expectedError, r)
				}
			}
		}()

		executeTemplate(data, tempTemplateFile.Name(), "/unknowpath/testOutput.html")
	})

	t.Run("TestUnableToParseTemplate", func(t *testing.T) {
		setupSuite(t, teardownSuite)
		tempTemplateFile := createTempFile(t, "testTemplate*.html")
		tempTemplateFileName := tempTemplateFile.Name()
		reportOutputPath := "/tmp/testOutput.html"

		writeContent(t, tempTemplateFile, `{{define "doesnotexist"}the missing piece{{end}}`)

		defer func() {
			if err := os.Remove(tempTemplateFile.Name()); err != nil {
				t.Errorf("Failed to remove temporary template file: %v", err)
			}
			if _, err := os.Stat(reportOutputPath); err == nil {
				if err := os.Remove(reportOutputPath); err != nil {
					t.Errorf("Failed to remove temporary test output file: %v", err)
				}
			}
		}()

		data := Data{}

		expectedError := "Unable to parse template file: template: " + strings.Split(tempTemplateFileName, "/")[2] + ":1: unexpected \"}\" in define clause"

		defer func() {
			r := recover()
			if r == nil {
				t.Error("\nExpected Error, but test passed")
			} else if r != expectedError {
				t.Errorf("\nExpected:: %s, Got:: %s", expectedError, r)
			}
		}()

		executeTemplate(data, tempTemplateFileName, reportOutputPath)
	})

	t.Run("TestUnableToExecuteTemplate", func(t *testing.T) {
		setupSuite(t, teardownSuite)
		tempTemplateFile := createTempFile(t, "testTemplate*.html")
		tempTemplateFileName := tempTemplateFile.Name()
		reportOutputPath := "/tmp/testOutput.html"

		writeContent(t, tempTemplateFile, `<html><body> something {{ .Stats.totalMutantsCount}} <body><html>`)

		defer func() {
			if err := os.Remove(tempTemplateFile.Name()); err != nil {
				t.Errorf("Failed to remove temporary template file: %v", err)
			}
			if _, err := os.Stat(reportOutputPath); err == nil {
				if err := os.Remove(reportOutputPath); err != nil {
					t.Errorf("Failed to remove temporary test output file: %v", err)
				}
			}
		}()

		data := Data{Stats: Stats{TotalMutantsCount: 10}}

		testTemplateFileName := strings.Split(tempTemplateFileName, "/")[2]
		expectedError := "Error executing template: template: " + testTemplateFileName + ":1:32: executing \"" + testTemplateFileName + "\" at <.Stats.totalMutantsCount>: can't evaluate field totalMutantsCount in type main.Stats"

		defer func() {
			r := recover()
			if r == nil {
				t.Error("\nExpected Error, but test passed")
			} else if r != expectedError {
				t.Errorf("\nExpected:: %s, Got:: %s", expectedError, r)
			}
		}()

		executeTemplate(data, tempTemplateFileName, reportOutputPath)
	})

	t.Run("TestShouldBeAbleToExecuteTemplate", func(t *testing.T) {
		setupSuite(t, teardownSuite)
		tempTemplateFile := createTempFile(t, "testTemplate*.html")
		tempTemplateFileName := tempTemplateFile.Name()
		reportOutputPath := "/tmp/testOutput.html"

		writeContent(t, tempTemplateFile, `<html><body> something {{define "value"}}more something{{end}}<body><html>`)

		defer func() {
			if err := os.Remove(tempTemplateFile.Name()); err != nil {
				t.Errorf("Failed to remove temporary template file: %v", err)
			}
			if _, err := os.Stat(reportOutputPath); err == nil {
				if err := os.Remove(reportOutputPath); err != nil {
					t.Errorf("Failed to remove temporary test output file: %v", err)
				}
			}
		}()

		data := Data{}

		defer func() {
			r := recover()
			if r != nil {
				t.Errorf("\nExpected test to pass, but threw error: %s", r)
			}
		}()

		executeTemplate(data, tempTemplateFileName, reportOutputPath)
		fileContent, _ := os.ReadFile(reportOutputPath)

		expectedOutput := "<html><body> something <body><html>"
		if string(fileContent) != expectedOutput {
			t.Errorf("Report File content not matched\nExpected:: %s\nActual::%s", expectedOutput, fileContent)
		}
	})
}
