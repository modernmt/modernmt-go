package modernmt

import (
	"crypto/rsa"
	"fmt"
	"net/http"
)

type ModernMT struct {
	client   *httpClient
	pk       *rsa.PublicKey
	pkTime   int64
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
	Priority        string
	ProjectId       string
	Multiline       *bool
	Timeout         int
	Format          string
	AltTranslations int

	// batch translation
	Metadata       interface{}
	IdempotencyKey string
}

type Translation struct {
	Translation      string
	ContextVector    string
	Characters       int
	BilledCharacters int
	DetectedLanguage string
	AltTranslations  []string
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

	altTranslationsInterface, ok := data["altTranslations"].([]interface{})
	if ok {
		altTranslations := make([]string, len(altTranslationsInterface))
		for i, v := range altTranslationsInterface {
			altTranslations[i] = v.(string)
		}
		translation.AltTranslations = altTranslations
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
	importJob := ImportJob{
		Id:       data["id"].(string),
		Size:     int(data["size"].(float64)),
		Progress: float32(data["progress"].(float64)),
	}

	memory, ok := data["memory"].(float64)
	if ok {
		importJob.Memory = int64(memory)
	}

	return importJob
}

type DetectedLanguage struct {
	BilledCharacters int
	DetectedLanguage string
}

func makeDetectedLanguage(data map[string]interface{}) DetectedLanguage {
	return DetectedLanguage{
		BilledCharacters: int(data["billedCharacters"].(float64)),
		DetectedLanguage: data["detectedLanguage"].(string),
	}
}

type billingPeriod struct {
	Begin           string
	End             string
	Chars           int64
	Plan            string
	PlanDescription string
	PlanForCatTool  bool
	Amount          float32
	Currency        string
	CurrencySymbol  string
}

type User struct {
	Id               int64
	Name             string
	Email            string
	RegistrationDate string
	Country          string
	IsBusiness       int8
	Status           string
	BillingPeriod    billingPeriod
}

func makeUser(data map[string]interface{}) User {
	bp := data["billingPeriod"].(map[string]interface{})
	return User{
		Id:               int64(data["id"].(float64)),
		Name:             data["name"].(string),
		Email:            data["email"].(string),
		RegistrationDate: data["registrationDate"].(string),
		Country:          data["country"].(string),
		IsBusiness:       int8(data["isBusiness"].(float64)),
		Status:           data["status"].(string),
		BillingPeriod: billingPeriod{
			Begin:           bp["begin"].(string),
			End:             bp["end"].(string),
			Chars:           int64(bp["chars"].(float64)),
			Plan:            bp["plan"].(string),
			PlanDescription: bp["planDescription"].(string),
			PlanForCatTool:  bp["planForCatTool"].(bool),
			Amount:          float32(bp["amount"].(float64)),
			Currency:        bp["currency"].(string),
			CurrencySymbol:  bp["currencySymbol"].(string),
		},
	}
}
