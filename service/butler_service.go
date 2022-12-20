package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
	"time"

	"github.com/maxwolffe/butler-cli/v2/data"
	"gopkg.in/yaml.v2"
)

type ButlerService struct {
	ButlerConfig data.ButlerConfig
	APIBase      string
}

func NewButlerService() *ButlerService {
	butService := ButlerService{}

	// TODO make configuration source configurable
	// TODO make this default not depend on WHERE the binary is run - make it always relative to the module source
	absPath, _ := filepath.Abs("secrets.yaml")
	fmt.Println("Absolute Path: " + absPath)
	fileContent, _ := os.ReadFile(absPath)

	butlerConfig := &data.ButlerConfig{}
	yaml.Unmarshal(fileContent, butlerConfig)

	butService.ButlerConfig = *butlerConfig
	butService.APIBase = "https://app.butlerlabs.ai/api/queues/" + butlerConfig.QueueID

	// TODO output these as debug logs if extra verbosity requested
	// TODO move to a proper logging framework.
	// fmt.Println("Created ButlerService with APIKey: " + string(butService.ButlerConfig.ApiKey))
	// fmt.Println("Created ButlerService with QueueID: " + string(butService.ButlerConfig.ApiKey))
	// fmt.Println("Created ButlerService with BaseAPI: " + string(butService.APIBase))

	return &butService
}

func (butService *ButlerService) ProcessSingleImage(filePath string) ([]data.Document, error) {
	imgFile, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer imgFile.Close()

	// Create a new buffer to hold the image data
	buf := new(bytes.Buffer)

	// Create a new multipart writer
	w := multipart.NewWriter(buf)

	// Create a new form field for the image
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
			"files", filePath))
	h.Set("Content-Type", "image/jpeg")
	field, err := w.CreatePart(h)
	if err != nil {
		panic(err)
	}

	// Copy the image data to the form field
	if _, err = io.Copy(field, imgFile); err != nil {
		panic(err)
	}

	// Close the multipart writer to finalize the form data
	w.Close()

	client := &http.Client{}

	request, err := http.NewRequest("POST", butService.APIBase+"/uploads", buf)
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Add("Authorization", "Bearer "+butService.ButlerConfig.ApiKey)
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", w.FormDataContentType())

	res, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	cnt, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	uploadResponse := data.UploadResponse{}
	json.Unmarshal(cnt, &uploadResponse)

	for {
		result, err := butService.GetExtractionResults(uploadResponse.UploadID)
		if err != nil {
			log.Fatal(err)
		}

		if result.Ready {
			fmt.Println(result.Response)
			return result.Response.Items, nil
		}
		time.Sleep(60 * time.Second)
	}
}

func (butService *ButlerService) ProcessRecipesInDir(dir string) ([]data.Document, error) {
	// TODO - process HEIC files instead of having to export to JPEG.
	// TODO - would be cool to confirm all the files which were going to be uploaded first.

	// Open the image file
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	completeDocuments := make([]data.Document, 0)

	for _, file := range files {
		fmt.Println(file.Name(), file.IsDir())
		if !file.IsDir() {
			// TODO Add response to a list
			// TODO return a structured object from this filename
			documents, _ := butService.ProcessSingleImage(dir + "/" + file.Name())
			completeDocuments = append(completeDocuments, documents...)
		}
	}

	// TODO - short term - write the processed files to a csv

	return completeDocuments, nil
}

func (butService *ButlerService) GetExtractionResults(uploadId string) (*data.ExtractionResult, error) {
	fmt.Println("Attempting to extract response: " + uploadId)

	client := &http.Client{}

	request, err := http.NewRequest("GET", butService.APIBase+"/extraction_results", nil)
	if err != nil {
		log.Fatal(err)
	}

	request.Header.Add("Authorization", "Bearer "+butService.ButlerConfig.ApiKey)
	request.Header.Add("Accept", "application/json")
	q := request.URL.Query()
	q.Add("uploadId", uploadId)
	request.URL.RawQuery = q.Encode()

	res, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	cnt, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	extractionResponse := data.ExtractionResponse{}
	json.Unmarshal(cnt, &extractionResponse)
	fmt.Println("Extraction response:")
	fmt.Print(extractionResponse)

	fullyDone := true
	for _, doc := range extractionResponse.Items {
		if doc.DocumentStatus == "InProgress" {
			fullyDone = false
		}
	}

	result := data.ExtractionResult{
		Response: extractionResponse,
		Ready:    fullyDone,
	}

	return &result, nil
}
