package main

import (
	"errors"
	"io"
	"log"
	"os"
	"os/exec"
	"testing"
	"text/template"
)

var osCreate = os.Create
var osReadFile = os.ReadFile
var templateParseFiles = template.ParseFiles
var execCommand = exec.Command

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
			"killedCount": 5
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
}

// Test readJson with invalid JSON input
func TestReadJson_InvalidJson(t *testing.T) {
	jsonContent := `{invalid json}`

	filePath := createTempJSONFile(t, jsonContent)
	defer os.Remove(filePath)

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic due to invalid JSON")
		}
	}()

	readJson(filePath)
}

// Test readJson with missing file
func TestReadJson_FileNotFound(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic due to missing file")
		}
	}()

	readJson("nonexistent.json")
}

// Test readJson with no file path provided
func TestReadJson_NoFilePath(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic due to no file path provided")
		}
	}()

	readJson("")
}

func TestExecuteTemplate_TemplateParseError(t *testing.T) {
	data := Data{}

	// Mock template.ParseFiles to return an error
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

// Test main with valid arguments
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

	os.Args = []string{"cmd", filePath}
	main()
}
