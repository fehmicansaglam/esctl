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
