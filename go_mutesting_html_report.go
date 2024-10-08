package main

import (
	"embed"
	"flag"
	models "gomutestingrhtmlreport/models"
	"html/template"
	"os"
	"path/filepath"
)

//go:embed *.tmpl
var reportTmplFS embed.FS

func executeTemplate(data models.Report, templatePath string, outputReportFilePath string) {
	parsedTemplate, err := template.ParseFS(reportTmplFS, templatePath)
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

var reports = []models.Data{}

func readAll(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if !info.IsDir() && filepath.Ext(info.Name()) == ".json" {
		reports = append(reports, readSingleJsonReportFile(path))
	}
	return nil
}

func mergeReports() models.Report {
	directory := "reports"
	err := filepath.Walk(directory, readAll)

	if err != nil {
		panic("Error reading directory:" + err.Error())
	}

	mergedReports := models.Report{}
	mergedReports.GlobalStats = models.Stats{}
	fileMap := make(map[string]models.ReportDetails)

	for _, report := range reports {
		filePath := report.Killed[0].Mutator.OriginalFilePath
		fileMap[filePath] = dataToReportMapper(report).ReportDetail[filePath]
		mergedReports.ReportDetail = fileMap
	}
	return mergedReports
}

func main() {
	templatePath := flag.String("template", "report.tmpl", "Provide template path:: -template <PATH_TO_TEMPLATE>")
	reportPath := flag.String("out", "reports/report.html", "Provide report output path:: -out <PATH_TO_OUTPUT_HTML_REPORT>")
	flag.Parse()

	executeTemplate(mergeReports(), *templatePath, *reportPath)
}
