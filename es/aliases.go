package es

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
