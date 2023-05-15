package es

import "fmt"

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
