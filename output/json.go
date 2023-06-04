package output

import (
	"encoding/json"
	"fmt"
)

func PrintJson(data interface{}) {
	prettyJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println("Failed to generate pretty JSON:", err)
		return
	}

	fmt.Println(string(prettyJSON))
}
