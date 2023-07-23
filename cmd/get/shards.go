package get

import (
	"fmt"
	"os"
	"strconv"

	"github.com/fehmicansaglam/esctl/cmd/config"
	"github.com/fehmicansaglam/esctl/cmd/utils"
	"github.com/fehmicansaglam/esctl/constants"
	"github.com/fehmicansaglam/esctl/es"
	"github.com/fehmicansaglam/esctl/output"
	"github.com/spf13/cobra"
)

var getShardsCmd = &cobra.Command{
	Use:   "shards",
	Short: "Get shards in Elasticsearch cluster",
	Long: utils.Trim(`
The 'shards' command provides detailed information about each shard in the Elasticsearch cluster.

This includes:
  - Shard number
  - State of the shard (e.g., whether it's started, relocating, initializing, or unassigned)
  - Whether the shard is a primary or a replica
  - Size of the shard
  - Node on which the shard is located

Filters can be applied to only show shards in certain states, with a specific number, located on a particular node, or designated as primary or replica.`),
	Example: utils.TrimAndIndent(`
# Retrieve detailed information about shards in the Elasticsearch cluster.
esctl get shards

# Retrieve shard information for an index.
esctl get shards --index my_index

# Retrieve shard information filtered by state.
esctl get shards --started --relocating`),
	Run: func(cmd *cobra.Command, args []string) {
		conf := config.ParseConfigFile()
		handleShardLogic(conf)
	},
}

func init() {
	getShardsCmd.Flags().StringVarP(&flagIndex, "index", "i", "", "Name of the index")
	getShardsCmd.Flags().StringVar(&flagNode, "node", "", "Filter shards by node")
	getShardsCmd.Flags().IntVar(&flagShard, "shard", -1, "Filter shards by shard number")
	getShardsCmd.Flags().BoolVar(&flagPrimary, "primary", false, "Filter primary shards")
	getShardsCmd.Flags().BoolVar(&flagReplica, "replica", false, "Filter replica shards")
	getShardsCmd.Flags().BoolVar(&flagStarted, "started", false, "Filter shards in STARTED state")
	getShardsCmd.Flags().BoolVar(&flagRelocating, "relocating", false, "Filter shards in RELOCATING state")
	getShardsCmd.Flags().BoolVar(&flagInitializing, "initializing", false, "Filter shards in INITIALIZING state")
	getShardsCmd.Flags().BoolVar(&flagUnassigned, "unassigned", false, "Filter shards in UNASSIGNED state")
}

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

func handleShardLogic(conf config.Config) {
	shards, err := es.GetShards(flagIndex)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to retrieve shards:", err)
		os.Exit(1)
	}

	columnDefs, err := getColumnDefs(conf, "shard", shardColumns)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to get column definitions:", err)
		os.Exit(1)
	}

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
