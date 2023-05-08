/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/fehmicansaglam/esctl/es"
	"github.com/spf13/cobra"
)

// shardsCmd represents the shards command
var shardsCmd = &cobra.Command{
	Use:   "shards",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		shards, err := es.GetShards("localhost", 9200, "articles")
		if err != nil {
			panic(err)
		}
		for _, shard := range shards {
			fmt.Println(shard.ID, shard.PriRep, shard.Shard, shard.IP, shard.State)
		}
	},
}

func init() {
	describeCmd.AddCommand(shardsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// shardsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// shardsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
