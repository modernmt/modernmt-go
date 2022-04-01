package modernmt

import "os"

func Create(apiKey string) *ModernMT {
	return CreateWithIdentity(apiKey, "modernmt-go", "1.0.1")
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

	res, err := re.TranslateListAdaptive(source, target, []string{q}, hints, contextVector, options)
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

	res, err := re.GetContextVectors(source, []string{target}, text, hints, limit)
	if err != nil {
		return "", nil
	}

	return res[target].(string), nil
}

func (re *ModernMT) GetContextVectors(source string, targets []string, text string, hints []int64,
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

func (re *ModernMT) GetContextVectorFromFile(source string, target string, file *os.File, hints []int64,
	limit int, compression string) (string, error) {

	res, err := re.GetContextVectorsFromFile(source, []string{target}, file, hints, limit, compression)
	if err != nil {
		return "", err
	}

	return res[target].(string), nil
}

func (re *ModernMT) GetContextVectorsFromFile(source string, targets []string, file *os.File, hints []int64,
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

func (re *ModernMT) GetContextVectorFromFilePath(source string, target string, path string, hints []int64,
	limit int, compression string) (string, error) {

	res, err := re.GetContextVectorsFromFilePath(source, []string{target}, path, hints, limit, compression)
	if err != nil {
		return "", err
	}

	return res[target].(string), nil
}

func (re *ModernMT) GetContextVectorsFromFilePath(source string, targets []string, path string, hints []int64,
	limit int, compression string) (map[string]interface{}, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	return re.GetContextVectorsFromFile(source, targets, file, hints, limit, compression)
}
