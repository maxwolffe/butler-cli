package data

type UploadResponse struct {
	UploadID  string
	Documents []Document
}

type Document struct {
	Filename       string      `json:"fileName,omitempty"`
	DocumentId     string      `json:"documentId,omitempty"`
	DocumentStatus string      `json:"documentStatus,omitempty"`
	FormFields     []FormField `json:"formFields,omitempty"`
}

type FormField struct {
	FieldName       string
	Value           string
	ConfidenceScore string
}

type ExtractionResult struct {
	Response ExtractionResponse
	Ready    bool
}

type ExtractionResponse struct {
	Items      []Document
	TotalCount int
}
