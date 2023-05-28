package output

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

func PrintYaml(data interface{}) {
	yamlData, err := yaml.Marshal(data)
	if err != nil {
		fmt.Println("Failed to marshal data to YAML:", err)
		os.Exit(1)
	}

	fmt.Println(string(yamlData))
}
