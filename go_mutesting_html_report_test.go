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
			"escaped": [
			{
      "mutator": {
        "mutatorName": "arithmetic/base",
		"mutatedSourceCode": "package main",
        "originalFilePath": "./go_mutesting_html_report1.go",
        "originalStartLine": 0
      },
      "diff": "--- Original\n+++ New\n@@ -56,7 +56,7 @@\n \tvar data Data\n \terr = json.Unmarshal(jsonData, \u0026data)\n \tif err != nil {\n-\t\tpanic(\"Invalid JSON format: \" + err.Error())\n+\t\tpanic(\"Invalid JSON format: \" - err.Error())\n \t}\n \treturn data\n }\n",
      "processOutput": "PASS \"/var/folders/y_/rcm9fbdn27d9plc8hdd5pttw0000gq/T/go-mutesting-3719870323/./go_mutesting_html_report.go.0\" with checksum 1\n"
    },
	{
      "mutator": {
        "mutatorName": "arithmetic/base",
		"mutatedSourceCode": "package main",
        "originalFilePath": "./go_mutesting_html_report.go",
        "originalStartLine": 0
      },
      "diff": "--- Original\n+++ New\n@@ -56,7 +56,7 @@\n \tvar data Data\n \terr = json.Unmarshal(jsonData, \u0026data)\n \tif err != nil {\n-\t\tpanic(\"Invalid JSON format: \" + err.Error())\n+\t\tpanic(\"Invalid JSON format: \" - err.Error())\n \t}\n \treturn data\n }\n",
      "processOutput": "PASS \"/var/folders/y_/rcm9fbdn27d9plc8hdd5pttw0000gq/T/go-mutesting-3719870323/./go_mutesting_html_report.go.0\" with checksum 2\n"
    },
	{
      "mutator": {
        "mutatorName": "arithmetic/base",
		"mutatedSourceCode": "package main",
        "originalFilePath": "./go_mutesting_html_report1.go",
        "originalStartLine": 0
      },
      "diff": "--- Original\n+++ New\n@@ -56,7 +56,7 @@\n \tvar data Data\n \terr = json.Unmarshal(jsonData, \u0026data)\n \tif err != nil {\n-\t\tpanic(\"Invalid JSON format: \" + err.Error())\n+\t\tpanic(\"Invalid JSON format: \" - err.Error())\n \t}\n \treturn data\n }\n",
      "processOutput": "PASS \"/var/folders/y_/rcm9fbdn27d9plc8hdd5pttw0000gq/T/go-mutesting-3719870323/./go_mutesting_html_report.go.0\" with checksum 3\n"
    },
    {
      "mutator": {
        "mutatorName": "arithmetic/base",
        "originalSourceCode": "package main\n\nim",
        "mutatedSourceCode": "package main\n",
        "originalFilePath": "./go_mutesting_html_report.go",
        "originalStartLine": 0
      },
      "diff": "--- Original\n+++ New\n@@ -64,7 +64,7 @@\n func executeTemplate(data Data, templatePath string, outputReportFilePath string) {\n \tparsedTemplate, err := template.ParseFS(reportTmplFS, templatePath)\n \tif err != nil {\n-\t\tpanic(\"Unable to parse template file: \" + err.Error())\n+\t\tpanic(\"Unable to parse template file: \" - err.Error())\n \t}\n \n \ttemplate := template.Must(parsedTemplate, err)\n",
      "processOutput": "PASS \"/var/folders/y_/rcm9fbdn27d9plc8hdd5pttw0000gq/T/go-mutesting-3719870323/./go_mutesting_html_report.go.1\" with checksum 4\n"
    }],
	"killed": [{
      "mutator": {
        "mutatorName": "arithmetic/base",
		"mutatedSourceCode": "package main",
        "originalFilePath": "./go_mutesting_html_report1.go",
        "originalStartLine": 0
      },
      "diff": "--- Original\n+++ New\n@@ -56,7 +56,7 @@\n \tvar data Data\n \terr = json.Unmarshal(jsonData, \u0026data)\n \tif err != nil {\n-\t\tpanic(\"Invalid JSON format: \" + err.Error())\n+\t\tpanic(\"Invalid JSON format: \" - err.Error())\n \t}\n \treturn data\n }\n",
      "processOutput": "PASS \"/var/folders/y_/rcm9fbdn27d9plc8hdd5pttw0000gq/T/go-mutesting-3719870323/./go_mutesting_html_report.go.0\" with checksum 5\n"
    },
    {
      "mutator": {
        "mutatorName": "arithmetic/some",
        "originalSourceCode": "package main\n\nim",
        "mutatedSourceCode": "package main\n",
        "originalFilePath": "./go_mutesting_html_report.go",
        "originalStartLine": 0
      },
      "diff": "--- Original\n+++ New\n@@ -64,7 +64,7 @@\n func executeTemplate(data Data, templatePath string, outputReportFilePath string) {\n \tparsedTemplate, err := template.ParseFS(reportTmplFS, templatePath)\n \tif err != nil {\n-\t\tpanic(\"Unable to parse template file: \" + err.Error())\n+\t\tpanic(\"Unable to parse template file: \" - err.Error())\n \t}\n \n \ttemplate := template.Must(parsedTemplate, err)\n",
      "processOutput": "PASS \"/var/folders/y_/rcm9fbdn27d9plc8hdd5pttw0000gq/T/go-mutesting-3719870323/./go_mutesting_html_report.go.1\" with checksum 6\n"
    }]
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

		os.Args = []string{"cmd", "-out", reportOutputPath, "-template", "report_test.tmpl"}
		main()

		fileContent, _ := os.ReadFile(reportOutputPath)

		expectedOutput := "<html><body> Total Mutant: 10 | Killed Count: 5 <body><html>"
		assert.Equal(t, expectedOutput, string(fileContent), "Report File content not matched")
	})
}

func TestExecuteTemplate(t *testing.T) {

	t.Run("TestExecuteTemplateUnableToReadTemplate", func(t *testing.T) {
		data := models.Report{}

		defer func() {
			expectedError := "Unable to parse template file: template: pattern matches no files: `testTemplate.html`"
			r := recover()
			assert.NotNil(t, r, "Expected Error, but test passed")
			assert.Equal(t, expectedError, r)
		}()

		executeTemplate(data, "testTemplate.html", "/unknowpath/testOutput.html")
	})

	t.Run("TestExecuteTemplateUnableToCreateReportFile", func(t *testing.T) {
		data := models.Report{}

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

		data := models.Report{}

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
			expectedError := "Error executing template: template: " + tempTemplateFileName + ":1:32: executing \"" + tempTemplateFileName + "\" at <.Stats.totalMutantsCount>: can't evaluate field totalMutantsCount in type models.Stats"
			assert.NotNil(t, r, "Expected Error, but test passed")
			assert.Equal(t, expectedError, r)
		}()

		data := models.Report{Stats: models.Stats{TotalMutantsCount: 10}}

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

		data := models.Report{}

		executeTemplate(data, tempTemplateFileName, reportOutputPath)
		fileContent, _ := os.ReadFile(reportOutputPath)

		expectedOutput := "<html><body> Total Mutant: 0 | Killed Count: 0 <body><html>"
		assert.Equal(t, expectedOutput, string(fileContent), "Report File content not matched")
	})
}
