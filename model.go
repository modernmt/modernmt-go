package modernmt

import (
	"fmt"
	"net/http"
)

type ModernMT struct {
	client   *httpClient
	Memories memoryServices
}

type memoryServices struct {
	client *httpClient
}

type httpClient struct {
	baseUrl string
	headers map[string]string
	client  *http.Client
}

type APIError struct {
	Status  int
	Type    string
	Message string
}

func (re APIError) Error() string {
	return fmt.Sprintf("%s: %s", re.Type, re.Message)
}

type TranslateOptions struct {
	Priority  string
	ProjectId string
	Multiline *bool
	Timeout   int
}

type Translation struct {
	Translation      string
	ContextVector    string
	Characters       int
	BilledCharacters int
	DetectedLanguage string
}

func makeTranslation(data map[string]interface{}) Translation {
	translation := Translation{
		Translation:      data["translation"].(string),
		Characters:       int(data["characters"].(float64)),
		BilledCharacters: int(data["billedCharacters"].(float64)),
	}

	contextVector, ok := data["contextVector"].(string)
	if ok {
		translation.ContextVector = contextVector
	}

	detectedLanguage, ok := data["detectedLanguage"].(string)
	if ok {
		translation.DetectedLanguage = detectedLanguage
	}

	return translation
}

type Memory struct {
	Id           int64
	Name         string
	Description  string
	CreationDate string
}

func makeMemory(data map[string]interface{}) Memory {
	memory := Memory{
		Id:           int64(data["id"].(float64)),
		Name:         data["name"].(string),
		CreationDate: data["creationDate"].(string),
	}

	description, ok := data["description"].(string)
	if ok {
		memory.Description = description
	}

	return memory
}

type ImportJob struct {
	Id       string
	Memory   int64
	Size     int
	Progress float32
}

func makeImportJob(data map[string]interface{}) ImportJob {
	return ImportJob{
		Id:       data["id"].(string),
		Memory:   int64(data["memory"].(float64)),
		Size:     int(data["size"].(float64)),
		Progress: float32(data["progress"].(float64)),
	}
}
