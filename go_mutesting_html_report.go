package main

import (
	"embed"
	"encoding/json"
	"flag"
	"html/template"
	"log"
	"os"
	"strings"
)

//go:embed *.tmpl
var reportTmplFS embed.FS

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

type MutatorDetail struct {
	MutatorName string
	Diff        string
	Checksum    string
}

type ReportDetails struct {
	Escaped []MutatorDetail
	Killed  []MutatorDetail
}
type Report struct {
	Stats        Stats
	ReportDetail map[string]ReportDetails
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

func executeTemplate(report Report, templatePath string, outputReportFilePath string) {
	parsedTemplate, err := template.ParseFS(reportTmplFS, templatePath)
	if err != nil {
		panic("Unable to parse template file: " + err.Error())
	}

	template := template.Must(parsedTemplate, err)

	outputReport, err := os.Create(outputReportFilePath)
	if err != nil {
		panic("Unable to create report file: " + err.Error())
	}

	err = template.Execute(outputReport, report)
	if err != nil {
		panic("Error executing template: " + err.Error())
	}
}

func groupByFile(data Data) Report {
	fileMap := make(map[string]ReportDetails)
	escaped := data.Escaped
	killed := data.Killed

	for _, escapedMutantEntry := range escaped {
		out := strings.Split(escapedMutantEntry.ProcessOutput, " ")
		entry := MutatorDetail{MutatorName: escapedMutantEntry.Mutator.MutatorName, Diff: escapedMutantEntry.Diff, Checksum: out[4]}
		escapedEntry := fileMap[escapedMutantEntry.Mutator.OriginalFilePath].Escaped
		updatedEntry := append(escapedEntry, entry)
		fileMap[escapedMutantEntry.Mutator.OriginalFilePath] = ReportDetails{Escaped: updatedEntry}
	}

	for _, killedMutantEntry := range killed {
		out := strings.Split(killedMutantEntry.ProcessOutput, " ")
		entry := MutatorDetail{MutatorName: killedMutantEntry.Mutator.MutatorName, Diff: killedMutantEntry.Diff, Checksum: out[4]}
		killedEntry := fileMap[killedMutantEntry.Mutator.OriginalFilePath].Killed
		updatedEntry := append(killedEntry, entry)
		fileMap[killedMutantEntry.Mutator.OriginalFilePath] = ReportDetails{Escaped: fileMap[killedMutantEntry.Mutator.OriginalFilePath].Escaped, Killed: updatedEntry}
	}

	return Report{Stats: data.Stats, ReportDetail: fileMap}
}

func main() {
	jsonFilePath := flag.String("file", "report.json", "Provide report.json:: -file <PATH_TO_JSON_REPORT>")
	templatePath := flag.String("template", "report.tmpl", "Provide template path:: -template <PATH_TO_TEMPLATE>")
	reportPath := flag.String("out", "report.html", "Provide report output path:: -out <PATH_TO_OUTPUT_HTML_REPORT>")
	flag.Parse()

	data := readJson(*jsonFilePath)
	executeTemplate(groupByFile(data), *templatePath, *reportPath)
}
