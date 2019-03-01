package storage

import (
	"os"
	"io/ioutil"
	"encoding/json"
)

type WikiThread struct {
	Thread []string `json:"thread"`
}

var wikiThread WikiThread

func LoadWikiThreads() error {
	wikiJSON, err := os.OpenFile("storage/wiki_threads.json", os.O_RDONLY, 0755)
	if err != nil {
		return err
	}

	data, _ := ioutil.ReadAll(wikiJSON)

	json.Unmarshal(data, &wikiThread)

	wikiJSON.Close()

	return nil
}

func GetWikiThread() *WikiThread {
	return &wikiThread
}