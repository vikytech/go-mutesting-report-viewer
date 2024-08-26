package main

import (
	"encoding/json"
	"flag"
	"html/template"
	"log"
	"os"
)

type Mutator struct {
	MutatorName        string `json:"mutatorName"`
	OriginalSourceCode string `json:"originalSourceCode"`
	MutatedSourceCode  string `json:"mutatedSourceCode"`
	OriginalFilePath   string `json:"originalFilePath"`
	OriginalStartLine  int    `json:"originalStartLine"`
}

type Entry struct {
	Mutator       Mutator `json:"mutator"`
	Diff          string  `json:"diff"`
	ProcessOutput string  `json:"processOutput"`
}

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

type Data struct {
	Stats     Stats       `json:"stats"`
	Escaped   []Entry     `json:"escaped"`
	Timeouted interface{} `json:"timeouted"`
	Killed    []Entry     `json:"killed"`
	Errored   interface{} `json:"errored"`
}

func readJson(filePath string) Data {
	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		log.Panicf("Error reading file: %s", err.Error())
	}

	var data Data
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		panic("Invalid JSON format: " + err.Error())
	}
	return data
}

func executeTemplate(data Data) {
	tmpl := template.Must(template.ParseFiles("template.html"))
	outputReportFilePath := "report.html"
	report, err := os.Create(outputReportFilePath)
	if err != nil {
		panic("Unable to create report file: " + err.Error())
	}

	err = tmpl.Execute(report, data)
	if err != nil {
		panic("Error executing template: " + err.Error())
	}
}

func main() {
	log.SetFlags(0)
	if len(os.Args) > 1 && os.Args[1] == "" {
		panic("Error: No file path provided.\nUsage: go run go_mutesting_html_report.go -file <PATH_TO_REPORT>")
	}
	jsonFilePath := flag.String("file", "", "Provide report.json:: -file <PATH_TO_REPORT>")
	flag.Parse()

	if *jsonFilePath == "" {
		panic("Error: No file path provided")
	}
	data := readJson(*jsonFilePath)
	executeTemplate(data)
}
