package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/fehmicansaglam/esctl/constants"
	"github.com/fehmicansaglam/esctl/es"
	"github.com/fehmicansaglam/esctl/output"
)

func includeShardByState(shard es.Shard) bool {
	switch {
	case flagStarted && shard.State == constants.ShardStateStarted:
		return true
	case flagRelocating && shard.State == constants.ShardStateRelocating:
		return true
	case flagInitializing && shard.State == constants.ShardStateInitializing:
		return true
	case flagUnassigned && shard.State == constants.ShardStateUnassigned:
		return true
	case !flagStarted && !flagRelocating && !flagInitializing && !flagUnassigned:
		return true
	}
	return false
}

func includeShardByNumber(shard es.Shard) bool {
	shardNumber, err := strconv.Atoi(shard.Shard)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to parse shard number:", err)
		os.Exit(1)
	}
	return flagShard == -1 || flagShard == shardNumber
}

func includeShardByPriRep(shard es.Shard) bool {
	return (flagPrimary && shard.PriRep == constants.ShardPrimary) ||
		(flagReplica && shard.PriRep == constants.ShardReplica) ||
		(!flagPrimary && !flagReplica)
}

func includeShardByNode(shard es.Shard) bool {
	if flagNode == "" {
		return true
	}

	return shard.Node == flagNode
}

func humanizePriRep(priRep string) string {
	switch priRep {
	case constants.ShardPrimary:
		return "primary"
	case constants.ShardReplica:
		return "replica"
	default:
		return priRep
	}
}

var shardColumns = []output.ColumnDef{
	{Header: "INDEX", Type: output.Text},
	{Header: "SHARD", Type: output.Number},
	{Header: "PRI-REP", Type: output.Text},
	{Header: "STATE", Type: output.Text},
	{Header: "DOCS", Type: output.Number},
	{Header: "STORE", Type: output.DataSize},
	{Header: "IP", Type: output.Text},
	{Header: "NODE", Type: output.Text},
	{Header: "NODE-ID", Type: output.Text},
	{Header: "UNASSIGNED-REASON", Type: output.Text},
	{Header: "UNASSIGNED-AT", Type: output.Date},
	{Header: "SEGMENTS-COUNT", Type: output.Number},
}

func handleShardLogic(config Config) {
	shards, err := es.GetShards(flagIndex)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to retrieve shards:", err)
		os.Exit(1)
	}

	columnDefs := getColumnDefs(config, "shard", shardColumns)

	data := [][]string{}

	for _, shard := range shards {
		if includeShardByState(shard) && includeShardByNumber(shard) &&
			includeShardByPriRep(shard) && includeShardByNode(shard) {

			rowData := map[string]string{
				"INDEX":             shard.Index,
				"SHARD":             shard.Shard,
				"PRI-REP":           humanizePriRep(shard.PriRep),
				"STATE":             shard.State,
				"DOCS":              shard.Docs,
				"STORE":             shard.Store,
				"IP":                shard.IP,
				"NODE":              shard.Node,
				"NODE-ID":           shard.ID,
				"UNASSIGNED-REASON": shard.UnassignedReason,
				"UNASSIGNED-AT":     shard.UnassignedAt,
				"SEGMENTS-COUNT":    shard.SegmentsCount,
			}

			row := make([]string, len(columnDefs))
			for i, colDef := range columnDefs {
				row[i] = rowData[colDef.Header]
			}
			data = append(data, row)
		}
	}

	if len(flagSortBy) > 0 {
		output.PrintTable(columnDefs, data, flagSortBy...)
	} else {
		output.PrintTable(columnDefs, data, "INDEX", "SHARD", "PRI-REP")
	}
}
