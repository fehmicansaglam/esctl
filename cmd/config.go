package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Modify and view the configuration",
}

var useContextCmd = &cobra.Command{
	Use:   "use-context",
	Short: "Set the current context",
	Long:  `Set the current context to connect to. This command updates the 'current-context' field in the configuration file.`,
	Args:  cobra.ExactArgs(1),
	Run:   runUseContext,
}

var getContextsCmd = &cobra.Command{
	Use:   "get-contexts",
	Short: "List the contexts defined in the esctl.yml file",
	Run:   runGetContexts,
}

var currentContextCmd = &cobra.Command{
	Use:   "current-context",
	Short: "Display the current context",
	Run:   runCurrentContext,
}

func init() {
	configCmd.AddCommand(useContextCmd)
	configCmd.AddCommand(getContextsCmd)
	configCmd.AddCommand(currentContextCmd)
	rootCmd.AddCommand(configCmd)
}

func runUseContext(cmd *cobra.Command, args []string) {
	contextName := args[0]
	config := parseConfigFile()

	contextExists := false
	for _, context := range config.Contexts {
		if context.Name == contextName {
			contextExists = true
			break
		}
	}

	if !contextExists {
		fmt.Printf("Error: No context found with the name '%s' in the configuration.\n", contextName)
		os.Exit(1)
	}

	viper.Set("current-context", contextName)

	err := viper.WriteConfig()
	if err != nil {
		fmt.Printf("Error writing updated configuration: %s\n", err)
		os.Exit(1)
	}
}

func runGetContexts(cmd *cobra.Command, args []string) {
	config := parseConfigFile()
	for _, context := range config.Contexts {
		contextName := context.Name
		if contextName == config.CurrentContext {
			contextName += "(*)"
		}
		fmt.Printf("- name: %s\n", contextName)
		fmt.Printf("  host: %s\n", context.Host)
		if context.Protocol != "" {
			fmt.Printf("  protocol: %s\n", context.Protocol)
		}
		if context.Port != 0 {
			fmt.Printf("  port: %d\n", context.Port)
		}
		if context.Username != "" {
			fmt.Printf("  username: %s\n", context.Username)
		}
		if context.Password != "" {
			fmt.Printf("  password: %s\n", context.Password)
		}
	}
}

func runCurrentContext(cmd *cobra.Command, args []string) {
	config := parseConfigFile()
	fmt.Println(config.CurrentContext)
}

type Context struct {
	Name     string `mapstructure:"name"`
	Protocol string `mapstructure:"protocol"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type Entity struct {
	Columns []string `mapstructure:"columns"`
}

type Config struct {
	CurrentContext string            `mapstructure:"current-context"`
	Contexts       []Context         `mapstructure:"contexts"`
	Entities       map[string]Entity `mapstructure:"entities"`
}

func parseConfigFile() Config {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error getting user's home directory: %v\n", err)
		os.Exit(1)
	}

	viper.AddConfigPath(filepath.Join(home, ".config"))
	viper.SetConfigName("esctl")
	viper.SetConfigType("yml")

	err = viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Error reading config file: %s\n", err)
		os.Exit(1)
	}

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		fmt.Printf("Error unmarshaling config into struct: %v\n", err)
		os.Exit(1)
	}

	return config
}
