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

func executeTemplate(data Data, templatePath string, outputReportFilePath string) {
	parsedTemplate, err := template.ParseFiles(templatePath)

	if err != nil {
		panic("Unable to parse template file: " + err.Error())
	}

	template := template.Must(parsedTemplate, err)

	report, err := os.Create(outputReportFilePath)
	if err != nil {
		panic("Unable to create report file: " + err.Error())
	}

	err = template.Execute(report, data)
	if err != nil {
		panic("Error executing template: " + err.Error())
	}
}

func main() {
	jsonFilePath := flag.String("file", "report.json", "Provide report.json:: -file <PATH_TO_JSON_REPORT>")
	templatePath := flag.String("template", "report.tmpl", "Provide template path:: -template <PATH_TO_TEMPLATE>")
	reportPath := flag.String("out", "report.html", "Provide report output path:: -out <PATH_TO_OUTPUT_HTML_REPORT>")
	flag.Parse()

	data := readJson(*jsonFilePath)
	executeTemplate(data, *templatePath, *reportPath)
}
