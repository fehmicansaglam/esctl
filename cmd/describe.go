package cmd

import (
	"fmt"
	"os"

	"github.com/fehmicansaglam/esctl/constants"
	"github.com/fehmicansaglam/esctl/es"
	"github.com/fehmicansaglam/esctl/output"
	"github.com/spf13/cobra"
)

var describeCmd = &cobra.Command{
	Use:       "describe ENTITY",
	Short:     "Print detailed information about an entity",
	Long:      "Print detailed information about the specified entity.",
	Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	ValidArgs: []string{"cluster"},
	Run: func(cmd *cobra.Command, args []string) {
		entity := args[0]
		switch entity {
		case constants.EntityCluster:
			handleDescribeCluster()
		default:
			fmt.Printf("Unknown entity: %s\n", entity)
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

func init() {
	rootCmd.AddCommand(describeCmd)
}
