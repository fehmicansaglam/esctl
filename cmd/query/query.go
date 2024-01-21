package query

import (
	"fmt"
	"os"

	"github.com/fehmicansaglam/esctl/cmd/utils"
	"github.com/fehmicansaglam/esctl/es"
	"github.com/fehmicansaglam/esctl/output"
	"github.com/spf13/cobra"
)

var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "Query Elasticsearch",
	Long:  `This command allows you to query Elasticsearch.`,
	Example: utils.TrimAndIndent(`
esctl query articles
esctl query articles --id 61
esctl query articles --term "price:10" --size 1
esctl query articles --sort "price:desc" --from 10 --size 10`),
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		index := args[0]

		response, err := es.SearchDocuments(index, flagId, flagTerm, flagFrom, flagSize, flagNested, flagSort)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to query:", err)
			os.Exit(1)
		}
		output.PrintJson(response["hits"])
	},
}

func Cmd() *cobra.Command {
	return queryCmd
}

func init() {
	queryCmd.Flags().StringArrayVar(&flagId, "id", []string{}, "Document IDs to fetch")
	queryCmd.Flags().StringArrayVarP(&flagTerm, "term", "t", []string{}, "Term filter(s)")
	queryCmd.Flags().StringArrayVar(&flagNested, "nested", []string{}, "Nested path(s)")
	queryCmd.Flags().StringArrayVarP(&flagSort, "sort", "s", []string{}, "Sort definition(s)")
	queryCmd.Flags().IntVar(&flagFrom, "from", 0, "Starting document offset")
	queryCmd.Flags().IntVar(&flagSize, "size", 1, "Number of hits to return")
}
