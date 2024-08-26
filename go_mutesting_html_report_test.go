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

func setupSuite(tb testing.TB, teardown func(tb testing.TB)) {
	osCreate = os.Create
	templateParseFiles = template.ParseFiles
	log.Println("Setup suite")
	defer teardown(tb)

}

func teardownSuite(tb testing.TB) {
	log.Println("Teardown suite", os.Args)
	os.Args = []string{}

	defer func() {
		if r := recover(); r != nil {
			tb.Errorf("Error:: %s", r)
		}
	}()
}

func init() {
	log.SetOutput(io.Discard)
}

func createTempJSONFile(t *testing.T, content string) string {
	setupSuite(t, teardownSuite)
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
		filePath := createTempJSONFile(t, jsonContent)
		defer os.Remove(filePath)

		os.Args = []string{"cmd", "-file", filePath}
		main()
	})

	t.Run("TestMain_EmptyFilePath", func(t *testing.T) {
		setupSuite(t, teardownSuite)
		defer func() {
			if r := recover(); r != nil {
				if r != "Error: No file path provided.\nUsage: go run go_mutesting_html_report.go -file <PATH_TO_REPORT>" {
					t.Errorf("Expected:: Error: No file path provided.\nUsage: go run go_mutesting_html_report.go -file <PATH_TO_REPORT>, Got:: %s", r)
				}
			}
		}()

		os.Args = []string{"-file", ""}
		main()
	})

	t.Run("TestMain_NoFilePath", func(t *testing.T) {
		setupSuite(t, teardownSuite)
		defer func() {
			if r := recover(); r != nil {
				if r != "Error: No file path provided" {
					t.Errorf("Expected:: Error: No file path provided, Got:: %s", r)
				}
			}
		}()

		os.Args = []string{"cmd", "-file", ""}
		// main()
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
	})
	t.Run("TestReadJson_InvalidJson", func(t *testing.T) {
		setupSuite(t, teardownSuite)
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
	})

	t.Run("TestReadJson_FileNotFound", func(t *testing.T) {
		setupSuite(t, teardownSuite)
		defer func() {
			if r := recover(); r != nil {
				if r != "Error reading file: open nonexistent.json: no such file or directory" {
					t.Errorf("\nExpected:: Error reading file: open nonexistent.json: no such file or directory, Got:: %s", r)
				}
			}
		}()

		readJson("nonexistent.json")
	})

	t.Run("TestReadJson_NoFilePath", func(t *testing.T) {
		setupSuite(t, teardownSuite)
		defer func() {
			if r := recover(); r != nil {
				if r != "Error reading file: open : no such file or directory" {
					t.Errorf("\nExpected:: Error reading file: open : no such file or directory, Got:: %s", r)
				}
			}
		}()

		readJson("")
	})

}

func TestExecuteTemplate(t *testing.T) {
	t.Run("TestExecuteTemplate_TemplateParseError", func(t *testing.T) {
		setupSuite(t, teardownSuite)
		data := Data{}

		templateParseFiles = func(filenames ...string) (*template.Template, error) {
			return nil, errors.New("mocked parse error")
		}

		executeTemplate(data)
	})

	t.Run("TestExecuteTemplate_CreateFileError", func(t *testing.T) {
		setupSuite(t, teardownSuite)
		data := Data{}

		osCreate = func(name string) (*os.File, error) {
			return nil, errors.New("mocked create error")
		}

		executeTemplate(data)
	})
}
