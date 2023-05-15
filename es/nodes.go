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
