package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/fehmicansaglam/esctl/constants"
	"github.com/fehmicansaglam/esctl/es"
	"github.com/fehmicansaglam/esctl/output"
	"github.com/spf13/cobra"
)

var (
	flagIndex        string
	flagNode         string
	flagShard        int
	flagPrimary      bool
	flagReplica      bool
	flagStarted      bool
	flagRelocating   bool
	flagInitializing bool
	flagUnassigned   bool
	flagActions      []string
	flagSortBy       []string
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get Elasticsearch entities",
	Long: `The 'get' command allows you to retrieve information about Elasticsearch entities. Supported entities include nodes, indices, and shards. This command provides a read-only view of the cluster and does not support data querying.

Usage:
	esctl get [entity]

Available Entities:
	- nodes: List all nodes in the Elasticsearch cluster.
	- indices: List all indices in the Elasticsearch cluster.
	- shards: List detailed information about shards, including their sizes and placement.
	- aliases: List all aliases in the Elasticsearch cluster.
	- tasks: List all tasks in the Elasticsearch cluster.

Options:
	[entity] - Specifies the entity type to retrieve. Supports 'nodes', 'indices', and 'shards'.

Examples:
	esctl get nodes
	Retrieves a list of all nodes in the Elasticsearch cluster.

	esctl get indices
	Retrieves a list of all indices in the Elasticsearch cluster.

	esctl get shards
	Retrieves detailed information about shards in the Elasticsearch cluster.

	esctl get shards --index my_index
	Retrieve shard information for an index.

	esctl get shards --started --relocating
	Retrieve shard information filtered by state.

	esctl get aliases
	Retrieve all aliases.

	esctl get tasks --actions 'index*' --actions '*search*'
	Retrieve tasks filtered by actions using wildcard patterns.

	esctl get tasks
	Retrieve all tasks.

Please note that the 'get' command only provides read-only access and does not support data querying or modification operations.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Please specify an entity.")
			os.Exit(1)
		}

		entity := args[0]

		switch entity {
		case constants.EntityNode, constants.EntityNodes:
			handleNodeLogic()
		case constants.EntityIndex, constants.EntityIndices:
			handleIndexLogic()
		case constants.EntityShard, constants.EntityShards:
			handleShardLogic()
		case constants.EntityAlias, constants.EntityAliases:
			handleAliasLogic()
		case constants.EntityTask, constants.EntityTasks:
			handleTaskLogic()
		default:
			fmt.Printf("Unknown entity: %s\n", entity)
			os.Exit(1)
		}
	},
}

func handleNodeLogic() {
	nodes, err := es.GetNodes()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to retrieve nodes:", err)
		os.Exit(1)
	}

	headers := []string{"NAME", "IP", "NODE-ROLE", "MASTER", "HEAP-MAX", "HEAP-CURRENT", "HEAP-PERCENT", "CPU", "LOAD-1M", "DISK-TOTAL", "DISK-USED", "DISK-AVAILABLE"}
	data := [][]string{}

	for _, node := range nodes {
		row := []string{
			node.Name, node.IP, node.NodeRole, node.Master, node.HeapMax, node.HeapCurrent,
			node.HeapPercent + "%", node.CPU + "%", node.Load1m,
			node.DiskTotal, node.DiskUsed, node.DiskAvail,
		}
		data = append(data, row)
	}

	if len(flagSortBy) > 0 {
		output.PrintTable(headers, data, flagSortBy...)
	} else {
		output.PrintTable(headers, data, "NAME")
	}
}

func handleIndexLogic() {
	indices, err := es.GetIndices(flagIndex)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to retrieve indices:", err)
		os.Exit(1)
	}

	headers := []string{"INDEX", "UUID", "HEALTH", "STATUS", "SHARDS", "REPLICAS", "DOCS-COUNT", "DOCS-DELETED", "CREATION-DATE", "STORE-SIZE", "PRI-STORE-SIZE"}
	data := [][]string{}

	for _, index := range indices {
		row := []string{
			index.Index, index.UUID, index.Health, index.Status, index.Pri, index.Rep,
			index.DocsCount, index.DocsDeleted, index.CreationDate, index.StoreSize, index.PriStoreSize,
		}
		data = append(data, row)
	}

	if len(flagSortBy) > 0 {
		output.PrintTable(headers, data, flagSortBy...)
	} else {
		output.PrintTable(headers, data, "INDEX")
	}
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

func handleShardLogic() {
	shards, err := es.GetShards(flagIndex)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to retrieve shards:", err)
		os.Exit(1)
	}

	headers := []string{"INDEX", "SHARD", "PRI-REP", "STATE", "DOCS", "STORE", "IP", "NODE", "NODE-ID", "UNASSIGNED-REASON", "UNASSIGNED-AT", "SEGMENTS-COUNT"}
	data := [][]string{}

	for _, shard := range shards {
		if includeShardByState(shard) && includeShardByNumber(shard) && includeShardByPriRep(shard) && includeShardByNode(shard) {
			row := []string{
				shard.Index, shard.Shard, humanizePriRep(shard.PriRep), shard.State, shard.Docs, shard.Store, shard.IP, shard.Node, shard.ID, shard.UnassignedReason, shard.UnassignedAt, shard.SegmentsCount,
			}
			data = append(data, row)
		}
	}

	if len(flagSortBy) > 0 {
		output.PrintTable(headers, data, flagSortBy...)
	} else {
		output.PrintTable(headers, data, "INDEX", "SHARD", "PRI-REP")
	}
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

func handleAliasLogic() {
	aliases, err := es.GetAliases(flagIndex)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to retrieve aliases:", err)
		os.Exit(1)
	}

	headers := []string{"ALIAS", "INDEX"}
	data := [][]string{}

	for alias, index := range aliases {
		row := []string{alias, index}
		data = append(data, row)
	}

	output.PrintTable(headers, data, flagSortBy...)
}

func handleTaskLogic() {
	tasksResponse, err := es.GetTasks(flagActions)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to retrieve tasks:", err)
		os.Exit(1)
	}

	headers := []string{"NODE", "ID", "ACTION"}
	data := [][]string{}

	for _, node := range tasksResponse.Nodes {
		for _, task := range node.Tasks {
			row := []string{task.Node, fmt.Sprintf("%d", task.ID), task.Action}
			data = append(data, row)
		}
	}

	if len(flagSortBy) > 0 {
		output.PrintTable(headers, data, flagSortBy...)
	} else {
		output.PrintTable(headers, data, "NODE", "ID")
	}
}

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.Flags().StringVar(&flagIndex, "index", "", "Name of the index")
	getCmd.Flags().StringVar(&flagNode, "node", "", "Filter shards by node")
	getCmd.Flags().IntVar(&flagShard, "shard", -1, "Filter shards by shard number")
	getCmd.Flags().BoolVar(&flagPrimary, "primary", false, "Filter primary shards")
	getCmd.Flags().BoolVar(&flagReplica, "replica", false, "Filter replica shards")
	getCmd.Flags().BoolVar(&flagStarted, "started", false, "Filter shards in STARTED state")
	getCmd.Flags().BoolVar(&flagRelocating, "relocating", false, "Filter shards in RELOCATING state")
	getCmd.Flags().BoolVar(&flagInitializing, "initializing", false, "Filter shards in INITIALIZING state")
	getCmd.Flags().BoolVar(&flagUnassigned, "unassigned", false, "Filter shards in UNASSIGNED state")
	getCmd.Flags().StringSliceVar(&flagActions, "actions", []string{}, "Filter tasks by actions")
	getCmd.Flags().StringSliceVar(&flagSortBy, "sort-by", []string{}, "Columns to sort by (comma-separated)")
}
