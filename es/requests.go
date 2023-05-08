package es

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func getJSONResponse(url string, target interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(target)
}

type Index struct {
	Health    string `json:"health"`
	Status    string `json:"status"`
	Index     string `json:"index"`
	UUID      string `json:"uuid"`
	Shards    string `json:"pri"`
	Replicas  string `json:"rep"`
	DocsCount string `json:"docs.count"`
	StoreSize string `json:"pri.store.size"`
}

func GetIndices(host string, port int) ([]Index, error) {
	url := fmt.Sprintf("http://%s:%d/_cat/indices?format=json", host, port)

	var indices []Index
	if err := getJSONResponse(url, &indices); err != nil {
		return nil, err
	}

	return indices, nil
}

type Shard struct {
	Index        string `json:"index"`
	Shard        string `json:"shard"`
	PriRep       string `json:"prirep"`
	State        string `json:"state"`
	Docs         string `json:"docs"`
	Store        string `json:"store"`
	IP           string `json:"ip"`
	ID           string `json:"id"`
	Node         string `json:"node"`
	Unassigned   string `json:"unassigned.reason"`
	UnassignedAt string `json:"unassigned.at"`
	Segments     string `json:"segments.count"`
}

func GetShards(host string, port int, index string) ([]Shard, error) {
	url := fmt.Sprintf("http://%s:%d/_cat/shards/%s?format=json&h=index,shard,prirep,state,docs,store,ip,id,node,unassigned.reason,unassigned.at,segments.count", host, port, index)

	var shards []Shard
	err := getJSONResponse(url, &shards)
	if err != nil {
		return nil, err
	}

	return shards, nil
}
