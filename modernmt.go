package modernmt

import (
	"os"
	"strconv"
)

func toSliceOfString(slice []int64) []string {
	res := make([]string, len(slice))
	for i, v := range slice {
		res[i] = strconv.FormatInt(v, 10)
	}
	return res
}

func Create(apiKey string) *ModernMT {
	return CreateWithIdentity(apiKey, "modernmt-go", "1.0.2")
}

func CreateWithIdentity(apiKey string, platform string, platformVersion string) *ModernMT {
	headers := map[string]string{
		"MMT-ApiKey":          apiKey,
		"MMT-Platform":        platform,
		"MMT-PlatformVersion": platformVersion,
	}

	client := CreateHttpClient("https://api.modernmt.com", headers)

	return &ModernMT{
		client: client,
		Memories: memoryServices{
			client: client,
		},
	}
}

func (re *ModernMT) ListSupportedLanguages() ([]string, error) {
	res, err := re.client.send("GET", "/translate/languages", nil, nil)
	if err != nil {
		return nil, err
	}

	var languages []string
	for _, el := range res.([]interface{}) {
		languages = append(languages, el.(string))
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
	}

	res, err := re.client.send("GET", "/translate", data, nil)
	if err != nil {
		return nil, err
	}

	var translations []Translation
	for _, el := range res.([]interface{}) {
		translations = append(translations, makeTranslation(el.(map[string]interface{})))
	}

	return translations, nil
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

	res, err := re.client.send("GET", "/context-vector", data, nil)
	if err != nil {
		return nil, err
	}

	vectors := res.(map[string]interface{})["vectors"]

	return vectors.(map[string]interface{}), nil
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

	res, err := re.client.send("GET", "/context-vector", data, files)
	if err != nil {
		return nil, err
	}

	vectors := res.(map[string]interface{})["vectors"]

	return vectors.(map[string]interface{}), nil
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
