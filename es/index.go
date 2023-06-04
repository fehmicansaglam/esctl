package es

import "fmt"

type MappingResponse map[string]interface{}

func GetIndexMappings(index string) (MappingResponse, error) {
	endpoint := fmt.Sprintf("%s/_mappings", index)

	var mappingResponse MappingResponse
	if err := getJSONResponse(endpoint, &mappingResponse); err != nil {
		return nil, err
	}

	return mappingResponse, nil
}

type AliasResponse map[string]AliasDetail

type AliasDetail struct {
	Aliases map[string]interface{} `json:"aliases"`
}

func GetAliases(index string) (map[string]string, error) {
	if index == "" {
		index = "_all"
	}

	var aliasResp AliasResponse
	if err := getJSONResponse(index+"/_alias", &aliasResp); err != nil {
		return nil, err
	}

	aliases := make(map[string]string)
	for index, detail := range aliasResp {
		for alias := range detail.Aliases {
			aliases[alias] = index
		}
	}

	return aliases, nil
}
