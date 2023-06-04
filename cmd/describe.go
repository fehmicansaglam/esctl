package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/fehmicansaglam/esctl/constants"
	"github.com/fehmicansaglam/esctl/es"
	"github.com/fehmicansaglam/esctl/output"
	"github.com/spf13/cobra"
)

var describeCmd = &cobra.Command{
	Short:     "Print detailed information about an entity",
	Args:      cobra.RangeArgs(1, 2),
	ValidArgs: []string{"cluster", "index"},
	Run: func(cmd *cobra.Command, args []string) {
		entity := args[0]
		switch entity {
		case constants.EntityCluster:
			handleDescribeCluster()
		case constants.EntityIndex:
			if len(args) < 2 {
				fmt.Println("Index name is required.")
				cmd.Help()
				os.Exit(1)
			}
			handleDescribeIndex(args[1])
		default:
			fmt.Printf("Unknown entity: %s\n", entity)
			cmd.Help()
			os.Exit(1)
		}
	},
}

func handleDescribeCluster() {
	cluster, err := es.GetCluster()
	if err != nil {
		fmt.Println("Failed to retrieve cluster information:", err)
		return
	}

	output.PrintYaml(cluster)
}

func handleDescribeIndex(index string) {
	mappings, err := es.GetIndexMappings(index)

	if err != nil {
		fmt.Println("Failed to retrieve cluster information:", err)
		return
	}

	output.PrintJson(mappings)
}

func init() {
	describeCmd.Use = fmt.Sprintf(`describe [%s] [NAME]`, strings.Join(describeCmd.ValidArgs, "|"))
	describeCmd.Long = fmt.Sprintf("Print detailed information about the specified entity.\nAvailable entities: %s.", strings.Join(describeCmd.ValidArgs, ", "))
	rootCmd.AddCommand(describeCmd)
}
