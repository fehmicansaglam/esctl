/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/fehmicansaglam/esctl/es"
	"github.com/spf13/cobra"
)

// indicesCmd represents the indices command
var indicesCmd = &cobra.Command{
	Use:   "indices",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		indices, err := es.GetIndices("localhost", 9200)
		if err != nil {
			panic(err)
		}
		for _, index := range indices {
			fmt.Println("Index\tHealth\tShards\tReplicas")
			fmt.Printf("%-9s(%s):\t%s\tshards,\t%s\treplicas. %s docs, %s\n", index.Index, index.Health, index.Shards, index.Replicas, index.DocsCount, index.StoreSize)
		}
	},
}

func init() {
	describeCmd.AddCommand(indicesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// indicesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// indicesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
