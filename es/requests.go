package es

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

func getJSONResponse(url string, target interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(target)
}

type Node struct {
	Name        string `json:"name"`
	IP          string `json:"ip"`
	NodeRole    string `json:"node.role"`
	Master      string `json:"master"`
	HeapMax     string `json:"heap.max"`
	HeapCurrent string `json:"heap.current"`
	HeapPercent string `json:"heap.percent"`
	CPU         string `json:"cpu"`
	Load1m      string `json:"load_1m"`
	DiskTotal   string `json:"disk.total"`
	DiskUsed    string `json:"disk.used"`
	DiskAvail   string `json:"disk.avail"`
}

func GetNodes(host string, port int) ([]Node, error) {
	url := fmt.Sprintf("http://%s:%d/_cat/nodes?format=json&h=name,ip,node.role,master,heap.max,heap.current,heap.percent,cpu,load_1m,disk.total,disk.used,disk.avail", host, port)

	var nodes []Node
	if err := getJSONResponse(url, &nodes); err != nil {
		return nil, err
	}

	return nodes, nil
}

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

type Shard struct {
	Index            string `json:"index"`
	Shard            string `json:"shard"`
	PriRep           string `json:"prirep"`
	State            string `json:"state"`
	Docs             string `json:"docs"`
	Store            string `json:"store"`
	IP               string `json:"ip"`
	ID               string `json:"id"`
	Node             string `json:"node"`
	UnassignedReason string `json:"unassigned.reason"`
	UnassignedAt     string `json:"unassigned.at"`
	SegmentsCount    string `json:"segments.count"`
}

func GetShards(host string, port int, index string) ([]Shard, error) {
	url := fmt.Sprintf("http://%s:%d/_cat/shards", host, port)

	if index != "" {
		url += fmt.Sprintf("/%s", index)
	}

	url += "?format=json&h=index,shard,prirep,state,docs,store,ip,id,node,unassigned.reason,unassigned.at,segments.count"

	var shards []Shard
	err := getJSONResponse(url, &shards)
	if err != nil {
		return nil, err
	}

	return shards, nil
}

type AliasResponse map[string]AliasDetail

type AliasDetail struct {
	Aliases map[string]interface{} `json:"aliases"`
}

func GetAliases(host string, port int, index string) (map[string]string, error) {
	if index == "" {
		index = "_all"
	}

	url := url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("%s:%d", host, port),
		Path:   index + "/_alias",
	}

	var aliasResp AliasResponse
	if err := getJSONResponse(url.String(), &aliasResp); err != nil {
		return nil, err
	}

	aliases := make(map[string]string)
	for index, detail := range aliasResp {
		for alias := range detail.Aliases {
			aliases[alias] = index
		}
	}

	return aliases, nil
}
