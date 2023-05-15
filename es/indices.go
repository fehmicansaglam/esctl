package es

import "fmt"

type Index struct {
	Health       string `json:"health"`
	Status       string `json:"status"`
	Index        string `json:"index"`
	UUID         string `json:"uuid"`
	Pri          string `json:"pri"`
	Rep          string `json:"rep"`
	DocsCount    string `json:"docs.count"`
	DocsDeleted  string `json:"docs.deleted"`
	CreationDate string `json:"creation.date.string"`
	StoreSize    string `json:"store.size"`
	PriStoreSize string `json:"pri.store.size"`
}

func GetIndices(host string, port int) ([]Index, error) {
	url := fmt.Sprintf("http://%s:%d/_cat/indices?format=json&h=health,status,index,uuid,pri,rep,docs.count,docs.deleted,creation.date.string,store.size,pri.store.size", host, port)

	var indices []Index
	if err := getJSONResponse(url, &indices); err != nil {
		return nil, err
	}

	return indices, nil
}
