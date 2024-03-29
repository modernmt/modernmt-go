package modernmt

import (
	"os"
	"strconv"
)

func (re *memoryServices) List() ([]Memory, error) {
	res, err := re.client.send("GET", "/memories", nil, nil, nil)
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
	res, err := re.client.send("GET", path, nil, nil, nil)
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

	res, err := re.client.send("POST", "/memories", data, nil, nil)
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

	res, err := re.client.send("POST", "/memories", data, nil, nil)
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
	res, err := re.client.send("PUT", path, data, nil, nil)
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
	res, err := re.client.send("DELETE", path, nil, nil, nil)
	if err != nil {
		return Memory{}, err
	}

	return makeMemory(res.(map[string]interface{})), nil
}

func (re *memoryServices) Add(id int64, source string, target string, sentence string, translation string,
	tuid string) (ImportJob, error) {
	_id := strconv.FormatInt(id, 10)
	return re.AddWithSessionByKey(_id, source, target, sentence, translation, tuid, "")
}

func (re *memoryServices) AddByKey(id string, source string, target string, sentence string, translation string,
	tuid string) (ImportJob, error) {
	return re.AddWithSessionByKey(id, source, target, sentence, translation, tuid, "")
}

func (re *memoryServices) AddWithSession(id int64, source string, target string, sentence string, translation string,
	tuid string, session string) (ImportJob, error) {
	_id := strconv.FormatInt(id, 10)
	return re.AddWithSessionByKey(_id, source, target, sentence, translation, tuid, session)
}

func (re *memoryServices) AddWithSessionByKey(id string, source string, target string,
	sentence string, translation string, tuid string, session string) (ImportJob, error) {

	data := map[string]interface{}{
		"source":      source,
		"target":      target,
		"sentence":    sentence,
		"translation": translation,
	}

	if tuid != "" {
		data["tuid"] = tuid
	}

	if session != "" {
		data["session"] = session
	}

	path := "/memories/" + id + "/content"
	res, err := re.client.send("POST", path, data, nil, nil)
	if err != nil {
		return ImportJob{}, err
	}

	return makeImportJob(res.(map[string]interface{})), nil
}

func (re *memoryServices) Replace(id int64, tuid string, source string, target string, sentence string,
	translation string) (ImportJob, error) {
	_id := strconv.FormatInt(id, 10)
	return re.ReplaceWithSessionByKey(_id, tuid, source, target, sentence, translation, "")
}

func (re *memoryServices) ReplaceByKey(id string, tuid string, source string, target string, sentence string,
	translation string) (ImportJob, error) {
	return re.ReplaceWithSessionByKey(id, tuid, source, target, sentence, translation, "")
}

func (re *memoryServices) ReplaceWithSession(id int64, tuid string, source string, target string, sentence string,
	translation string, session string) (ImportJob, error) {
	_id := strconv.FormatInt(id, 10)
	return re.ReplaceWithSessionByKey(_id, tuid, source, target, sentence, translation, session)
}

func (re *memoryServices) ReplaceWithSessionByKey(id string, tuid string, source string, target string, sentence string,
	translation string, session string) (ImportJob, error) {

	data := map[string]interface{}{
		"tuid":        tuid,
		"source":      source,
		"target":      target,
		"sentence":    sentence,
		"translation": translation,
	}

	if session != "" {
		data["session"] = session
	}

	path := "/memories/" + id + "/content"
	res, err := re.client.send("PUT", path, data, nil, nil)
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
	res, err := re.client.send("POST", path, data, files, nil)
	if err != nil {
		return ImportJob{}, err
	}

	return makeImportJob(res.(map[string]interface{})), nil
}

func (re *memoryServices) AddToGlossary(id int64, terms []GlossaryTerm, _type string, tuid string) (ImportJob, error) {
	_id := strconv.FormatInt(id, 10)
	return re.AddToGlossaryByKey(_id, terms, _type, tuid)
}

func (re *memoryServices) AddToGlossaryByKey(id string, terms []GlossaryTerm, _type string,
	tuid string) (ImportJob, error) {

	data := map[string]interface{}{
		"terms": terms,
		"type":  _type,
	}

	if tuid != "" {
		data["tuid"] = tuid
	}

	path := "/memories/" + id + "/glossary"
	res, err := re.client.send("POST", path, data, nil, nil)
	if err != nil {
		return ImportJob{}, err
	}

	return makeImportJob(res.(map[string]interface{})), nil
}

func (re *memoryServices) ReplaceInGlossary(id int64, terms []GlossaryTerm, _type string,
	tuid string) (ImportJob, error) {

	_id := strconv.FormatInt(id, 10)
	return re.ReplaceInGlossaryByKey(_id, terms, _type, tuid)
}

func (re *memoryServices) ReplaceInGlossaryByKey(id string, terms []GlossaryTerm, _type string,
	tuid string) (ImportJob, error) {

	data := map[string]interface{}{
		"terms": terms,
		"type":  _type,
	}

	if tuid != "" {
		data["tuid"] = tuid
	}

	path := "/memories/" + id + "/glossary"
	res, err := re.client.send("PUT", path, data, nil, nil)
	if err != nil {
		return ImportJob{}, err
	}

	return makeImportJob(res.(map[string]interface{})), nil
}

func (re *memoryServices) ImportGlossaryPath(id int64, path string, _type string,
	compression string) (ImportJob, error) {

	_id := strconv.FormatInt(id, 10)
	return re.ImportGlossaryPathByKey(_id, path, _type, compression)
}

func (re *memoryServices) ImportGlossaryPathByKey(id string, path string, _type string,
	compression string) (ImportJob, error) {

	file, err := os.Open(path)
	if err != nil {
		return ImportJob{}, err
	}

	return re.ImportGlossaryByKey(id, file, _type, compression)
}

func (re *memoryServices) ImportGlossary(id int64, csv *os.File, _type string, compression string) (ImportJob, error) {
	_id := strconv.FormatInt(id, 10)
	return re.ImportGlossaryByKey(_id, csv, _type, compression)
}

func (re *memoryServices) ImportGlossaryByKey(id string, csv *os.File, _type string,
	compression string) (ImportJob, error) {

	data := map[string]interface{}{
		"type": _type,
	}

	if compression != "" {
		data["compression"] = compression
	}

	files := map[string]*os.File{
		"csv": csv,
	}

	path := "/memories/" + id + "/glossary"
	res, err := re.client.send("POST", path, data, files, nil)
	if err != nil {
		return ImportJob{}, err
	}

	return makeImportJob(res.(map[string]interface{})), nil
}

func (re *memoryServices) GetImportStatus(uuid string) (ImportJob, error) {
	res, err := re.client.send("GET", "/import-jobs/"+uuid, nil, nil, nil)
	if err != nil {
		return ImportJob{}, err
	}

	return makeImportJob(res.(map[string]interface{})), nil
}
