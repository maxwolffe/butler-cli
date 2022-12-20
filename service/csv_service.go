package service

import (
	"bufio"
	"encoding/csv"
	"os"

	"github.com/maxwolffe/butler-cli/v2/data"
)

func getValuesFromFormField(document data.Document) []string {
	formValues := make([]string, 0)
	for _, formField := range document.FormFields {
		formValues = append(formValues, formField.Value)
	}
	return formValues
}

func GenerateCsv(documents []data.Document, filepath string) {
	// Given a list of documents and a file path. Create a CSV and save it with the contents of the documents in the csv.

	f, _ := os.Create(filepath)
	defer f.Close()
	w := bufio.NewWriter(f)
	csvWriter := csv.NewWriter(w)
	for _, document := range documents {
		csvFields := make([]string, 0)
		csvFields = append(csvFields, document.Filename)
		csvFields = append(csvFields, getValuesFromFormField(document)...)
		csvWriter.Write(csvFields)
	}
	csvWriter.Flush()
	w.Flush()
}
