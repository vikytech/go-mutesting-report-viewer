package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
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
	jsonFilePath := flag.String("file", "", filePath)
	flag.Parse()

	if *jsonFilePath == "" {
		fmt.Println("Error: No file path provided")
		os.Exit(1)
	}

	jsonData, err := ioutil.ReadFile(*jsonFilePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}

	var data Data
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		fmt.Errorf("Invalid JSON format: " + err.Error())
		os.Exit(1)
	}
	return data
}

func executeTemplate(data Data) {
	tmpl := template.Must(template.ParseFiles("template.html"))
	report, err := os.Create("docs/report.html")
	if err != nil {
		log.Println("Unable to create report file: ", err)
		return
	}

	err = tmpl.Execute(report, data)
	if err != nil {
		fmt.Errorf("Error executing template: " + err.Error())
		return
	}
}

func uploadJson(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("template.html"))
	if r.Method == http.MethodPost {
		file, _, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Unable to get file from form: "+err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()

		jsonData, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "Unable to read file: "+err.Error(), http.StatusInternalServerError)
			return
		}

		var data Data
		err = json.Unmarshal(jsonData, &data)
		if err != nil {
			http.Error(w, "Invalid JSON format: "+err.Error(), http.StatusBadRequest)
			return
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Error executing template: "+err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "docs/index.html")
}

func main() {
	if len(os.Args) > 1 && os.Args[1] != "" {
		data := readJson(os.Args[1])
		executeTemplate(data)
		os.Exit(0)
	}

	http.HandleFunc("/", serveIndex)
	http.HandleFunc("/upload", uploadJson)
	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
