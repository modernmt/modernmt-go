package modernmt

import (
	"os"
	"strconv"
)

func (re *memoryServices) List() ([]Memory, error) {
	res, err := re.client.send("GET", "/memories", nil, nil)
	if err != nil {
		return nil, err
	}

	var memories []Memory
	for _, el := range res.([]interface{}) {
		memories = append(memories, makeMemory(el.(map[string]interface{})))
	}

	return memories, nil
}

func (re *memoryServices) Get(id int64) (Memory, error) {
	_id := strconv.FormatInt(id, 10)
	return re.GetByKey(_id)
}

func (re *memoryServices) GetByKey(id string) (Memory, error) {
	path := "/memories/" + id
	res, err := re.client.send("GET", path, nil, nil)
	if err != nil {
		return Memory{}, err
	}

	return makeMemory(res.(map[string]interface{})), nil
}

func (re *memoryServices) Create(name string, description string) (Memory, error) {
	data := map[string]interface{}{
		"name": name,
	}

	if description != "" {
		data["description"] = description
	}

	res, err := re.client.send("POST", "/memories", data, nil)
	if err != nil {
		return Memory{}, err
	}

	return makeMemory(res.(map[string]interface{})), nil
}

func (re *memoryServices) Connect(name string, description string, externalId string) (Memory, error) {
	data := map[string]interface{}{
		"name": name,
	}

	if description != "" {
		data["description"] = description
	}

	if externalId != "" {
		data["external_id"] = externalId
	}

	res, err := re.client.send("POST", "/memories", data, nil)
	if err != nil {
		return Memory{}, err
	}

	return makeMemory(res.(map[string]interface{})), nil
}

func (re *memoryServices) Edit(id int64, name string, description string) (Memory, error) {
	_id := strconv.FormatInt(id, 10)
	return re.EditByKey(_id, name, description)
}

func (re *memoryServices) EditByKey(id string, name string, description string) (Memory, error) {
	data := map[string]interface{}{}

	if name != "" {
		data["name"] = name
	}

	if description != "" {
		data["description"] = description
	}

	path := "/memories/" + id
	res, err := re.client.send("PUT", path, data, nil)
	if err != nil {
		return Memory{}, err
	}

	return makeMemory(res.(map[string]interface{})), nil
}

func (re *memoryServices) Delete(id int64) (Memory, error) {
	_id := strconv.FormatInt(id, 10)
	return re.DeleteByKey(_id)
}

func (re *memoryServices) DeleteByKey(id string) (Memory, error) {
	path := "/memories/" + id
	res, err := re.client.send("DELETE", path, nil, nil)
	if err != nil {
		return Memory{}, err
	}

	return makeMemory(res.(map[string]interface{})), nil
}

func (re *memoryServices) Add(id int64, source string, target string, sentence string, translation string,
	tuid string) (ImportJob, error) {
	_id := strconv.FormatInt(id, 10)
	return re.AddByKey(_id, source, target, sentence, translation, tuid)
}

func (re *memoryServices) AddByKey(id string, source string, target string, sentence string, translation string,
	tuid string) (ImportJob, error) {

	data := map[string]interface{}{
		"source":      source,
		"target":      target,
		"sentence":    sentence,
		"translation": translation,
	}

	if tuid != "" {
		data["tuid"] = tuid
	}

	path := "/memories/" + id + "/content"
	res, err := re.client.send("POST", path, data, nil)
	if err != nil {
		return ImportJob{}, err
	}

	return makeImportJob(res.(map[string]interface{})), nil
}

func (re *memoryServices) Replace(id int64, tuid string, source string, target string, sentence string,
	translation string) (ImportJob, error) {
	_id := strconv.FormatInt(id, 10)
	return re.ReplaceByKey(_id, tuid, source, target, sentence, translation)
}

func (re *memoryServices) ReplaceByKey(id string, tuid string, source string, target string, sentence string,
	translation string) (ImportJob, error) {

	data := map[string]interface{}{
		"tuid":        tuid,
		"source":      source,
		"target":      target,
		"sentence":    sentence,
		"translation": translation,
	}

	path := "/memories/" + id + "/content"
	res, err := re.client.send("PUT", path, data, nil)
	if err != nil {
		return ImportJob{}, err
	}

	return makeImportJob(res.(map[string]interface{})), nil
}

func (re *memoryServices) ImportTmxPath(id int64, path string, compression string) (ImportJob, error) {
	_id := strconv.FormatInt(id, 10)
	return re.ImportTmxPathByKey(_id, path, compression)
}

func (re *memoryServices) ImportTmxPathByKey(id string, path string, compression string) (ImportJob, error) {
	file, err := os.Open(path)
	if err != nil {
		return ImportJob{}, err
	}

	return re.ImportTmxByKey(id, file, compression)
}

func (re *memoryServices) ImportTmx(id int64, tmx *os.File, compression string) (ImportJob, error) {
	_id := strconv.FormatInt(id, 10)
	return re.ImportTmxByKey(_id, tmx, compression)
}

func (re *memoryServices) ImportTmxByKey(id string, tmx *os.File, compression string) (ImportJob, error) {
	data := map[string]interface{}{}

	if compression != "" {
		data["compression"] = compression
	}

	files := map[string]*os.File{
		"tmx": tmx,
	}

	path := "/memories/" + id + "/content"
	res, err := re.client.send("POST", path, data, files)
	if err != nil {
		return ImportJob{}, err
	}

	return makeImportJob(res.(map[string]interface{})), nil
}

func (re *memoryServices) GetImportStatus(uuid string) (ImportJob, error) {
	res, err := re.client.send("GET", "/import-jobs/"+uuid, nil, nil)
	if err != nil {
		return ImportJob{}, err
	}

	return makeImportJob(res.(map[string]interface{})), nil
}
