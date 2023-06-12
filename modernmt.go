package modernmt

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"strconv"
	"time"
)

func toSliceOfString(slice []int64) []string {
	res := make([]string, len(slice))
	for i, v := range slice {
		res[i] = strconv.FormatInt(v, 10)
	}
	return res
}

func Create(apiKey string) *ModernMT {
	return CreateWithClientId(apiKey, 0)
}

func CreateWithIdentity(apiKey string, platform string, platformVersion string) *ModernMT {
	return CreateWithIdentityAndClientId(apiKey, platform, platformVersion, 0)
}

func CreateWithClientId(apiKey string, apiClient int64) *ModernMT {
	return CreateWithIdentityAndClientId(apiKey, "modernmt-go", "1.2.0", apiClient)
}

func CreateWithIdentityAndClientId(apiKey string, platform string, platformVersion string, apiClient int64) *ModernMT {
	headers := map[string]string{
		"MMT-ApiKey":          apiKey,
		"MMT-Platform":        platform,
		"MMT-PlatformVersion": platformVersion,
	}

	if apiClient != 0 {
		headers["MMT-ApiClient"] = strconv.FormatInt(apiClient, 10)
	}

	client := createHttpClient("https://api.modernmt.com", headers)

	return &ModernMT{
		client: client,
		pk:     nil,
		pkTime: 0,
		Memories: memoryServices{
			client: client,
		},
	}
}

func (re *ModernMT) ListSupportedLanguages() ([]string, error) {
	res, err := re.client.send("GET", "/translate/languages", nil, nil, nil)
	if err != nil {
		return nil, err
	}

	var languages []string
	for _, el := range res.([]interface{}) {
		languages = append(languages, el.(string))
	}

	return languages, nil
}

func (re *ModernMT) DetectLanguage(q string, format string) (DetectedLanguage, error) {
	res, err := re.DetectLanguages([]string{q}, format)
	if err != nil {
		return DetectedLanguage{}, err
	}

	return res[0], nil
}

func (re *ModernMT) DetectLanguages(q []string, format string) ([]DetectedLanguage, error) {
	data := map[string]interface{}{
		"q": q,
	}

	if format != "" {
		data["format"] = format
	}
	res, err := re.client.send("GET", "/translate/detect", data, nil, nil)
	if err != nil {
		return nil, err
	}

	var languages []DetectedLanguage
	for _, el := range res.([]interface{}) {
		languages = append(languages, makeDetectedLanguage(el.(map[string]interface{})))
	}

	return languages, nil
}

func (re *ModernMT) Translate(source string, target string, q string, options *TranslateOptions) (Translation, error) {
	return re.TranslateAdaptive(source, target, q, nil, "", options)
}

func (re *ModernMT) TranslateAdaptive(source string, target string, q string, hints []int64, contextVector string,
	options *TranslateOptions) (Translation, error) {
	_hints := toSliceOfString(hints)
	return re.TranslateAdaptiveWithKeys(source, target, q, _hints, contextVector, options)
}

func (re *ModernMT) TranslateAdaptiveWithKeys(source string, target string, q string, hints []string,
	contextVector string, options *TranslateOptions) (Translation, error) {

	res, err := re.TranslateListAdaptiveWithKeys(source, target, []string{q}, hints, contextVector, options)
	if err != nil {
		return Translation{}, err
	}

	return res[0], nil
}

func (re *ModernMT) TranslateList(source string, target string, q []string,
	options *TranslateOptions) ([]Translation, error) {

	return re.TranslateListAdaptive(source, target, q, nil, "", options)
}

func (re *ModernMT) TranslateListAdaptive(source string, target string, q []string, hints []int64,
	contextVector string, options *TranslateOptions) ([]Translation, error) {
	_hints := toSliceOfString(hints)
	return re.TranslateListAdaptiveWithKeys(source, target, q, _hints, contextVector, options)
}

func (re *ModernMT) TranslateListAdaptiveWithKeys(source string, target string, q []string, hints []string,
	contextVector string, options *TranslateOptions) ([]Translation, error) {

	data := map[string]interface{}{
		"source": source,
		"target": target,
		"q":      q,
	}

	if contextVector != "" {
		data["context_vector"] = contextVector
	}

	if hints != nil {
		data["hints"] = hints
	}

	if options != nil {
		if options.Priority != "" {
			data["priority"] = options.Priority
		}
		if options.ProjectId != "" {
			data["project_id"] = options.ProjectId
		}
		if options.Multiline != nil {
			data["multiline"] = *options.Multiline
		}
		if options.Timeout != 0 {
			data["timeout"] = options.Timeout
		}
		if options.Format != "" {
			data["format"] = options.Format
		}
		if options.AltTranslations != 0 {
			data["alt_translations"] = options.AltTranslations
		}
	}

	res, err := re.client.send("GET", "/translate", data, nil, nil)
	if err != nil {
		return nil, err
	}

	var translations []Translation
	for _, el := range res.([]interface{}) {
		translations = append(translations, makeTranslation(el.(map[string]interface{})))
	}

	return translations, nil
}

func (re *ModernMT) BatchTranslate(webhook string, source string, target string, q string, options *TranslateOptions) (bool, error) {
	return re.BatchTranslateAdaptive(webhook, source, target, q, nil, "", options)
}

func (re *ModernMT) BatchTranslateAdaptive(webhook string, source string, target string, q string, hints []int64, contextVector string,
	options *TranslateOptions) (bool, error) {
	_hints := toSliceOfString(hints)
	return re.BatchTranslateAdaptiveWithKeys(webhook, source, target, q, _hints, contextVector, options)
}

func (re *ModernMT) BatchTranslateAdaptiveWithKeys(webhook string, source string, target string, q string, hints []string,
	contextVector string, options *TranslateOptions) (bool, error) {
	return re.BatchTranslateListAdaptiveWithKeys(webhook, source, target, []string{q}, hints, contextVector, options)
}

func (re *ModernMT) BatchTranslateList(webhook string, source string, target string, q []string,
	options *TranslateOptions) (bool, error) {

	return re.BatchTranslateListAdaptive(webhook, source, target, q, nil, "", options)
}

func (re *ModernMT) BatchTranslateListAdaptive(webhook string, source string, target string, q []string, hints []int64,
	contextVector string, options *TranslateOptions) (bool, error) {
	_hints := toSliceOfString(hints)
	return re.BatchTranslateListAdaptiveWithKeys(webhook, source, target, q, _hints, contextVector, options)
}

func (re *ModernMT) BatchTranslateListAdaptiveWithKeys(webhook string, source string, target string, q []string, hints []string,
	contextVector string, options *TranslateOptions) (bool, error) {

	data := map[string]interface{}{
		"webhook": webhook,
		"source":  source,
		"target":  target,
		"q":       q,
	}

	if contextVector != "" {
		data["context_vector"] = contextVector
	}

	if hints != nil {
		data["hints"] = hints
	}

	headers := map[string]string{}

	if options != nil {
		if options.ProjectId != "" {
			data["project_id"] = options.ProjectId
		}
		if options.Multiline != nil {
			data["multiline"] = *options.Multiline
		}
		if options.Format != "" {
			data["format"] = options.Format
		}
		if options.AltTranslations != 0 {
			data["alt_translations"] = options.AltTranslations
		}
		if options.Metadata != nil {
			data["metadata"] = options.Metadata
		}
		if options.IdempotencyKey != "" {
			headers["x-idempotency-key"] = options.IdempotencyKey
		}
	}

	res, err := re.client.send("POST", "/translate/batch", data, nil, headers)
	if err != nil {
		return false, err
	}

	return res.(map[string]interface{})["enqueued"].(bool), nil
}

func (re *ModernMT) GetContextVector(source string, target string, text string, hints []int64,
	limit int) (string, error) {
	_hints := toSliceOfString(hints)
	return re.GetContextVectorByKeys(source, target, text, _hints, limit)
}

func (re *ModernMT) GetContextVectors(source string, targets []string, text string, hints []int64,
	limit int) (map[string]interface{}, error) {
	_hints := toSliceOfString(hints)
	return re.GetContextVectorsByKeys(source, targets, text, _hints, limit)
}

func (re *ModernMT) GetContextVectorByKeys(source string, target string, text string, hints []string,
	limit int) (string, error) {

	res, err := re.GetContextVectorsByKeys(source, []string{target}, text, hints, limit)
	if err != nil {
		return "", nil
	}

	return res[target].(string), nil
}

func (re *ModernMT) GetContextVectorsByKeys(source string, targets []string, text string, hints []string,
	limit int) (map[string]interface{}, error) {

	data := map[string]interface{}{
		"source":  source,
		"targets": targets,
		"text":    text,
	}

	if hints != nil {
		data["hints"] = hints
	}

	if limit != 0 {
		data["limit"] = limit
	}

	res, err := re.client.send("GET", "/context-vector", data, nil, nil)
	if err != nil {
		return nil, err
	}

	vectors := res.(map[string]interface{})["vectors"]

	return vectors.(map[string]interface{}), nil
}

func (re *ModernMT) GetContextVectorFromFile(source string, target string, file *os.File, hints []int64,
	limit int, compression string) (string, error) {
	_hints := toSliceOfString(hints)
	return re.GetContextVectorFromFileByKeys(source, target, file, _hints, limit, compression)
}

func (re *ModernMT) GetContextVectorsFromFile(source string, targets []string, file *os.File, hints []int64,
	limit int, compression string) (map[string]interface{}, error) {
	_hints := toSliceOfString(hints)
	return re.GetContextVectorsFromFileByKeys(source, targets, file, _hints, limit, compression)
}

func (re *ModernMT) GetContextVectorFromFilePath(source string, target string, path string, hints []int64,
	limit int, compression string) (string, error) {
	_hints := toSliceOfString(hints)
	return re.GetContextVectorFromFilePathByKeys(source, target, path, _hints, limit, compression)
}

func (re *ModernMT) GetContextVectorsFromFilePath(source string, targets []string, path string, hints []int64,
	limit int, compression string) (map[string]interface{}, error) {
	_hints := toSliceOfString(hints)
	return re.GetContextVectorsFromFilePathByKeys(source, targets, path, _hints, limit, compression)
}

func (re *ModernMT) GetContextVectorFromFilePathByKeys(source string, target string, path string, hints []string,
	limit int, compression string) (string, error) {

	res, err := re.GetContextVectorsFromFilePathByKeys(source, []string{target}, path, hints, limit, compression)
	if err != nil {
		return "", err
	}

	return res[target].(string), nil
}

func (re *ModernMT) GetContextVectorsFromFilePathByKeys(source string, targets []string, path string, hints []string,
	limit int, compression string) (map[string]interface{}, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	return re.GetContextVectorsFromFileByKeys(source, targets, file, hints, limit, compression)
}

func (re *ModernMT) GetContextVectorFromFileByKeys(source string, target string, file *os.File, hints []string,
	limit int, compression string) (string, error) {

	res, err := re.GetContextVectorsFromFileByKeys(source, []string{target}, file, hints, limit, compression)
	if err != nil {
		return "", err
	}

	return res[target].(string), nil
}

func (re *ModernMT) GetContextVectorsFromFileByKeys(source string, targets []string, file *os.File, hints []string,
	limit int, compression string) (map[string]interface{}, error) {

	files := map[string]*os.File{
		"content": file,
	}

	data := map[string]interface{}{
		"source":  source,
		"targets": targets,
	}

	if hints != nil {
		data["hints"] = hints
	}

	if limit != 0 {
		data["limit"] = limit
	}

	if compression != "" {
		data["compression"] = compression
	}

	res, err := re.client.send("GET", "/context-vector", data, files, nil)
	if err != nil {
		return nil, err
	}

	vectors := res.(map[string]interface{})["vectors"]

	return vectors.(map[string]interface{}), nil
}

func (re *ModernMT) HandleTranslateCallback(body []byte, signature string) (Translation, error) {
	res, err := re.HandleTranslateListCallback(body, signature)
	if err != nil {
		return Translation{}, err
	}

	return res[0], nil
}

func (re *ModernMT) HandleTranslateCallbackWithMetadata(body []byte, signature string, metadata interface{}) (Translation, error) {
	res, err := re.HandleTranslateListCallbackWithMetadata(body, signature, metadata)
	if err != nil {
		return Translation{}, err
	}

	return res[0], nil
}

func (re *ModernMT) HandleTranslateListCallback(body []byte, signature string) ([]Translation, error) {
	return re.HandleTranslateListCallbackWithMetadata(body, signature, nil)
}

func (re *ModernMT) HandleTranslateListCallbackWithMetadata(body []byte, signature string, metadata interface{}) ([]Translation, error) {
	err := re.verifyCallbackSignature(signature)
	if err != nil {
		return nil, err
	}

	var jBody map[string]interface{}
	err = json.Unmarshal(body, &jBody)
	if err != nil {
		return nil, err
	}

	if jMetadata, ok := jBody["metadata"]; ok {
		if metadata != nil {
			// not a fan of it, but it seems the easiest way to do it
			bytes, err := json.Marshal(jMetadata)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(bytes, metadata)
			if err != nil {
				return nil, err
			}
		}
	}

	result := jBody["result"].(map[string]interface{})

	status := int(result["status"].(float64))
	if status >= 300 || status < 200 {
		e := result["error"].(map[string]interface{})
		return nil, APIError{
			Status:  status,
			Type:    e["type"].(string),
			Message: e["message"].(string),
		}
	} else {
		var translations []Translation
		for _, el := range result["data"].([]interface{}) {
			translations = append(translations, makeTranslation(el.(map[string]interface{})))
		}

		return translations, nil
	}
}

func (re *ModernMT) verifyCallbackSignature(signature string) error {
	token, err := jwt.Parse(signature, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		pk, err := re.getPublicKey()
		if err != nil {
			return nil, err
		}

		return pk, nil
	})

	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return nil
	} else {
		return err
	}
}

func (re *ModernMT) getPublicKey() (*rsa.PublicKey, error) {
	if re.pk == nil || re.pkTime+3600 < time.Now().Unix() {
		res, err := re.retrievePublicKey()
		if err == nil {
			re.pk = res
			re.pkTime = time.Now().Unix()
		} else if re.pk == nil { //  if previous version ok pk is available, ignore API exception
			return nil, err
		}
	}

	return re.pk, nil
}

func (re *ModernMT) retrievePublicKey() (*rsa.PublicKey, error) {
	res, err := re.client.send("GET", "/translate/batch/key", nil, nil, nil)
	if err != nil {
		return nil, err
	}

	encoded := res.(map[string]interface{})["publicKey"].(string)
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}

	return jwt.ParseRSAPublicKeyFromPEM(decoded)
}

func (re *ModernMT) Me() (User, error) {
	res, err := re.client.send("GET", "/users/me", nil, nil, nil)
	if err != nil {
		return User{}, err
	}

	return makeUser(res.(map[string]interface{})), nil
}
