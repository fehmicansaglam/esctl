package es

import (
	"fmt"
	"net/url"
)

type AliasResponse map[string]AliasDetail

type AliasDetail struct {
	Aliases map[string]interface{} `json:"aliases"`
}

func GetAliases(host string, port int, index string) (map[string]string, error) {
	if index == "" {
		index = "_all"
	}

	url := url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("%s:%d", host, port),
		Path:   index + "/_alias",
	}

	var aliasResp AliasResponse
	if err := getJSONResponse(url.String(), &aliasResp); err != nil {
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
