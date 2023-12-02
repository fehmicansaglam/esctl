package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/fehmicansaglam/esctl/cmd/config"
	"github.com/fehmicansaglam/esctl/cmd/count"
	"github.com/fehmicansaglam/esctl/cmd/describe"
	"github.com/fehmicansaglam/esctl/cmd/get"
	"github.com/fehmicansaglam/esctl/cmd/query"
	"github.com/fehmicansaglam/esctl/constants"
	"github.com/fehmicansaglam/esctl/shared"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "esctl",
	Short: "esctl is CLI for Elasticsearch",
	Long:  `esctl is a read-only CLI for Elasticsearch that allows users to manage and monitor their Elasticsearch clusters.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initialize)

	initProtocolFlag()
	initHostFlag()
	initPortFlag()
	initUsernameFlag()
	initPasswordFlag()

	rootCmd.PersistentFlags().BoolVar(&shared.Debug, "debug", false, "Enable debug mode")

	rootCmd.AddCommand(config.Cmd())
	rootCmd.AddCommand(count.Cmd())
	rootCmd.AddCommand(describe.Cmd())
	rootCmd.AddCommand(get.Cmd())
	rootCmd.AddCommand(query.Cmd())
}

func initialize() {
	if shared.ElasticsearchHost == "" {
		conf := config.ParseConfigFile()
		readContextFromConfig(conf)
	}
}

func readContextFromConfig(conf config.Config) {
	if len(conf.Contexts) == 0 {
		fmt.Println("Error: No contexts defined in the configuration.")
		os.Exit(1)
	}
	if conf.CurrentContext == "" {
		conf.CurrentContext = conf.Contexts[0].Name
	}

	clusterFound := false
	for _, cluster := range conf.Contexts {
		if cluster.Name == conf.CurrentContext {
			shared.ElasticsearchProtocol = cluster.Protocol
			if shared.ElasticsearchProtocol == "" {
				shared.ElasticsearchProtocol = constants.DefaultElasticsearchProtocol
			}
			shared.ElasticsearchPort = cluster.Port
			if shared.ElasticsearchPort == 0 {
				shared.ElasticsearchPort = constants.DefaultElasticsearchPort
			}
			shared.ElasticsearchUsername = cluster.Username
			shared.ElasticsearchPassword = cluster.Password
			shared.ElasticsearchHost = cluster.Host
			if shared.ElasticsearchHost == "" {
				fmt.Println("Error: 'host' field is not specified in the configuration for the current cluster.")
				os.Exit(1)
			}
			clusterFound = true
			break
		}
	}

	if !clusterFound {
		fmt.Printf("Error: No cluster found with the name '%s' in the configuration.\n", conf.CurrentContext)
		os.Exit(1)
	}
}

func initProtocolFlag() {
	defaultProtocol := constants.DefaultElasticsearchProtocol
	defaultProtocolEnv := os.Getenv(constants.ElasticsearchProtocolEnvVar)
	if defaultProtocolEnv != "" {
		defaultProtocol = defaultProtocolEnv
	}
	rootCmd.PersistentFlags().StringVar(&shared.ElasticsearchProtocol, "protocol", defaultProtocol, "Elasticsearch protocol")
}

func initHostFlag() {
	defaultHost := os.Getenv(constants.ElasticsearchHostEnvVar)
	rootCmd.PersistentFlags().StringVar(&shared.ElasticsearchHost, "host", defaultHost, "Elasticsearch host")
}

func initPortFlag() {
	defaultPort := constants.DefaultElasticsearchPort
	defaultPortStr := os.Getenv(constants.ElasticsearchPortEnvVar)
	if defaultPortStr != "" {
		parsedPort, err := strconv.Atoi(defaultPortStr)
		if err != nil || parsedPort <= 0 {
			fmt.Printf("Invalid value for %s environment variable: %s\n", constants.ElasticsearchPortEnvVar, defaultPortStr)
			os.Exit(1)
		}
		defaultPort = parsedPort
	}
	rootCmd.PersistentFlags().IntVar(&shared.ElasticsearchPort, "port", defaultPort, "Elasticsearch port")
}

func initUsernameFlag() {
	defaultUsername := os.Getenv(constants.ElasticsearchUsernameEnvVar)
	rootCmd.PersistentFlags().StringVar(&shared.ElasticsearchUsername, "username", defaultUsername, "Elasticsearch username")
}

func initPasswordFlag() {
	defaultPassword := os.Getenv(constants.ElasticsearchPasswordEnvVar)
	rootCmd.PersistentFlags().StringVar(&shared.ElasticsearchPassword, "password", defaultPassword, "Elasticsearch password")
}
