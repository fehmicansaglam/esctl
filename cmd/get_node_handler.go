package cmd

import (
	"fmt"
	"os"

	"github.com/fehmicansaglam/esctl/es"
	"github.com/fehmicansaglam/esctl/output"
)

var nodeColumns = []output.ColumnDef{
	{Header: "NAME", Type: output.Text},
	{Header: "IP", Type: output.Text},
	{Header: "NODE-ROLE", Type: output.Text},
	{Header: "MASTER", Type: output.Text},
	{Header: "HEAP-MAX", Type: output.DataSize},
	{Header: "HEAP-CURRENT", Type: output.DataSize},
	{Header: "HEAP-PERCENT", Type: output.Percent},
	{Header: "RAM-MAX", Type: output.DataSize},
	{Header: "RAM-CURRENT", Type: output.DataSize},
	{Header: "RAM-PERCENT", Type: output.Percent},
	{Header: "CPU", Type: output.Percent},
	{Header: "LOAD-1M", Type: output.Number},
	{Header: "DISK-TOTAL", Type: output.DataSize},
	{Header: "DISK-USED", Type: output.DataSize},
	{Header: "DISK-AVAILABLE", Type: output.DataSize},
	{Header: "UPTIME", Type: output.Text},
}

func handleNodeLogic(config Config) {
	nodes, err := es.GetNodes()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to retrieve nodes: %v\n", err)
		os.Exit(1)
	}

	columnDefs, err := getColumnDefs(config, "node", nodeColumns)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to get column definitions:", err)
		os.Exit(1)
	}

	data := [][]string{}

	for _, node := range nodes {
		rowData := map[string]string{
			"NAME":           node.Name,
			"IP":             node.IP,
			"NODE-ROLE":      node.NodeRole,
			"MASTER":         node.Master,
			"HEAP-MAX":       node.HeapMax,
			"HEAP-CURRENT":   node.HeapCurrent,
			"HEAP-PERCENT":   node.HeapPercent + "%",
			"RAM-MAX":        node.RAMMax,
			"RAM-CURRENT":    node.RAMCurrent,
			"RAM-PERCENT":    node.RAMPercent + "%",
			"CPU":            node.CPU + "%",
			"LOAD-1M":        node.Load1m,
			"DISK-TOTAL":     node.DiskTotal,
			"DISK-USED":      node.DiskUsed,
			"DISK-AVAILABLE": node.DiskAvail,
			"UPTIME":         node.Uptime,
		}

		row := make([]string, len(columnDefs))
		for i, colDef := range columnDefs {
			row[i] = rowData[colDef.Header]
		}
		data = append(data, row)
	}

	if len(flagSortBy) > 0 {
		output.PrintTable(columnDefs, data, flagSortBy...)
	} else {
		output.PrintTable(columnDefs, data, "NAME")
	}
}
