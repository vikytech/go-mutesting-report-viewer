package main

import (
	"errors"
	"io"
	"log"
	"os"
	"testing"
	"text/template"
)

var osCreate = os.Create
var templateParseFiles = template.ParseFiles

func init() {
	log.SetOutput(io.Discard)
}

func createTempJSONFile(t *testing.T, content string) string {
	tempFile, err := os.CreateTemp("", "*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer tempFile.Close()

	if _, err := tempFile.WriteString(content); err != nil {
		t.Fatal(err)
	}

	return tempFile.Name()
}

func TestReadJson_ValidJson(t *testing.T) {
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

	filePath := createTempJSONFile(t, jsonContent)
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
}

func TestReadJson_InvalidJson(t *testing.T) {
	jsonContent := `{invalid json}`

	filePath := createTempJSONFile(t, jsonContent)
	defer os.Remove(filePath)

	defer func() {
		if r := recover(); r != nil {
			if r != "Invalid JSON format: invalid character 'i' looking for beginning of object key string" {
				t.Errorf("\nExpected:: Invalid JSON format: invalid character 'i' looking for beginning of object key string, Got:: %s", r)
			}
		}
	}()

	readJson(filePath)
}

func TestReadJson_FileNotFound(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			if r != "Error reading file: open nonexistent.json: no such file or directory" {
				t.Errorf("\nExpected:: Error reading file: open nonexistent.json: no such file or directory, Got:: %s", r)
			}
		}
	}()

	readJson("nonexistent.json")
}

func TestReadJson_NoFilePath(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			if r != "Error reading file: open : no such file or directory" {
				t.Errorf("\nExpected:: Error reading file: open : no such file or directory, Got:: %s", r)
			}
		}
	}()

	readJson("")
}

func TestExecuteTemplate_TemplateParseError(t *testing.T) {
	data := Data{}

	templateParseFiles = func(filenames ...string) (*template.Template, error) {
		return nil, errors.New("mocked parse error")
	}

	executeTemplate(data)
}

func TestExecuteTemplate_CreateFileError(t *testing.T) {
	data := Data{}

	osCreate = func(name string) (*os.File, error) {
		return nil, errors.New("mocked create error")
	}

	executeTemplate(data)
}

func TestMain_ValidFile(t *testing.T) {
	jsonContent := `{
		"stats": {
			"totalMutantsCount": 10,
			"killedCount": 5
		},
		"escaped": [],
		"killed": []
	}`
	filePath := createTempJSONFile(t, jsonContent)
	defer os.Remove(filePath)

	os.Args = []string{"cmd", "-file", filePath}
	main()
}

func TestMain_EmptyFilePath(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			if r != "Error: No file path provided.\nUsage: go run go_mutesting_html_report.go -file <PATH_TO_REPORT>" {
				t.Errorf("Error: No file path provided.\nUsage: go run go_mutesting_html_report.go -file <PATH_TO_REPORT>, Got:: %s", r)
			}
		}
	}()

	os.Args = []string{"-file", ""}
	main()
}
