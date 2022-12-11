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

	"gopkg.in/yaml.v2"
	"maxwolffe.com/recipeUploader/v2/data"
)

type ButlerService struct {
	ButlerConfig data.ButlerConfig
	APIBase      string
}

func NewButlerService() *ButlerService {
	butService := ButlerService{}

	// TODO make configuration source configurable
	absPath, _ := filepath.Abs("secrets.yaml")
	fmt.Println("Absolute Path: " + absPath)
	fileContent, _ := os.ReadFile(absPath)

	butlerConfig := &data.ButlerConfig{}
	yaml.Unmarshal(fileContent, butlerConfig)

	butService.ButlerConfig = *butlerConfig
	butService.APIBase = "https://app.butlerlabs.ai/api/queues/" + butlerConfig.QueueID

	fmt.Println("Created ButlerService with APIKey: " + string(butService.ButlerConfig.ApiKey))
	fmt.Println("Created ButlerService with QueueID: " + string(butService.ButlerConfig.ApiKey))
	fmt.Println("Created ButlerService with BaseAPI: " + string(butService.APIBase))

	return &butService
}

func (butService *ButlerService) ProcessSingleImage(filePath string) error {
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
		return err
	}
	defer res.Body.Close()

	cnt, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	uploadResponse := data.UploadResponse{}
	json.Unmarshal(cnt, &uploadResponse)
	fmt.Println("Printing upload response")
	fmt.Print(uploadResponse)

	for {
		result, err := butService.GetExtractionResults(uploadResponse.UploadID)
		if err != nil {
			log.Fatal(err)
		}

		if result.Ready {
			fmt.Println(result.Response)
			return nil
		}
		time.Sleep(60 * time.Second)
	}
}

func (butService *ButlerService) ProcessRecipesInDir(dir string) error {
	// TODO - process HEIC files instead of having to export to JPEG.
	// Open the image file

	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		fmt.Println(file.Name(), file.IsDir())
		if !file.IsDir() {
			// TODO Add response to a list
			// TODO return a structured object from this filename
			butService.ProcessSingleImage(dir + "/" + file.Name())
		}
	}

	// TODO collate responses into csv

	// TODO - short term - write the processed files to a csv
	// TODO - longterm write the processed files to Paprika directly
	return nil
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
