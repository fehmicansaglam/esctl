package es

import "fmt"

type Node struct {
	Name        string `json:"name"`
	IP          string `json:"ip"`
	NodeRole    string `json:"node.role"`
	Master      string `json:"master"`
	HeapMax     string `json:"heap.max"`
	HeapCurrent string `json:"heap.current"`
	HeapPercent string `json:"heap.percent"`
	RAMCurrent  string `json:"ram.current"`
	RAMMax      string `json:"ram.max"`
	RAMPercent  string `json:"ram.percent"`
	CPU         string `json:"cpu"`
	Load1m      string `json:"load_1m"`
	DiskTotal   string `json:"disk.total"`
	DiskUsed    string `json:"disk.used"`
	DiskAvail   string `json:"disk.avail"`
	Uptime      string `json:"uptime"`
}

func GetNodes() ([]Node, error) {
	endpoint := "_cat/nodes?format=json&h=name,ip,node.role,master,heap.max,heap.current,heap.percent,cpu,load_1m,disk.total,disk.used,disk.avail,ram.current,ram.max,ram.percent,uptime"

	var nodes []Node
	if err := getJSONResponse(endpoint, &nodes); err != nil {
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

func GetIndices(index string) ([]Index, error) {
	endpoint := "_cat/indices"

	if index != "" {
		endpoint += fmt.Sprintf("/%s", index)
	}

	endpoint += "?format=json&h=health,status,index,uuid,pri,rep,docs.count,docs.deleted,creation.date.string,store.size,pri.store.size"

	var indices []Index
	if err := getJSONResponse(endpoint, &indices); err != nil {
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

func GetShards(index string) ([]Shard, error) {
	endpoint := "_cat/shards"

	if index != "" {
		endpoint += fmt.Sprintf("/%s", index)
	}

	endpoint += "?format=json&h=index,shard,prirep,state,docs,store,ip,id,node,unassigned.reason,unassigned.at,segments.count"

	var shards []Shard
	err := getJSONResponse(endpoint, &shards)
	if err != nil {
		return nil, err
	}

	return shards, nil
}
