package modernmt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func CreateHttpClient(baseUrl string, headers map[string]string) *httpClient {
	return &httpClient{
		baseUrl: baseUrl,
		headers: headers,
		client:  &http.Client{},
	}
}

func (re *httpClient) _createMultipartRequest(path string, data map[string]interface{},
	files map[string]*os.File) (*http.Request, error) {

	var body bytes.Buffer
	w := multipart.NewWriter(&body)

	for param, file := range files {
		fw, err := w.CreateFormFile(param, file.Name())
		if err != nil {
			return nil, err
		}

		_, err = io.Copy(fw, file)
		if err != nil {
			return nil, err
		}

		err = file.Close()
		if err != nil {
			return nil, err
		}
	}

	for key, val := range data {
		var s string
		switch val.(type) {
		case []string:
			s = strings.Join(val.([]string), ",")
		case []int64:
			s = strings.Trim(strings.Join(strings.Fields(fmt.Sprint(val)), ","), "[]")
		case int:
			s = strconv.Itoa(val.(int))
		default:
			s = val.(string)
		}

		err := w.WriteField(key, s)
		if err != nil {
			return nil, err
		}
	}

	err := w.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", re.baseUrl+path, &body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", w.FormDataContentType())

	return req, nil
}

func (re *httpClient) _createJsonRequest(path string, data map[string]interface{}) (*http.Request, error) {
	var body bytes.Buffer

	if data != nil {
		jsonBytes, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}

		body = *bytes.NewBuffer(jsonBytes)
	}

	req, err := http.NewRequest("POST", re.baseUrl+path, &body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	return req, nil
}

func (re *httpClient) _createRequest(path string, data map[string]interface{},
	files map[string]*os.File) (*http.Request, error) {

	if files != nil {
		return re._createMultipartRequest(path, data, files)
	}

	return re._createJsonRequest(path, data)
}

func (re *httpClient) send(method string, path string, data map[string]interface{}, files map[string]*os.File) (interface{}, error) {

	req, err := re._createRequest(path, data, files)
	if err != nil {
		return nil, err
	}

	req.Header.Add("X-HTTP-Method-Override", method)
	for key, val := range re.headers {
		req.Header.Add(key, val)
	}

	res, err := re.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	status := int(result["status"].(float64))
	if status >= 300 || status < 200 {
		e := result["error"].(map[string]interface{})
		return nil, APIError{
			Status:  status,
			Type:    e["type"].(string),
			Message: e["message"].(string),
		}
	}

	return result["data"], nil
}
