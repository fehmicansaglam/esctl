package es

import "fmt"

type IndexMappings struct {
	Mappings interface{} `json:"mappings"`
}

type IndexSettings struct {
	Settings interface{} `json:"settings"`
}

type MappingsResponse map[string]IndexMappings
type SettingsResponse map[string]IndexSettings

type IndexDetails struct {
	Settings interface{} `json:"settings,omitempty"`
	Mappings interface{} `json:"mappings,omitempty"`
}

type IndexDetailsResponse map[string]IndexDetails

func GetIndexDetails(index string, shouldGetMappings, shouldGetSettings bool) (IndexDetailsResponse, error) {
	var mappingsResponse MappingsResponse
	var settingsResponse SettingsResponse

	if shouldGetMappings {
		mappingsEndpoint := fmt.Sprintf("%s/_mappings", index)
		if err := getJSONResponse(mappingsEndpoint, &mappingsResponse); err != nil {
			return nil, fmt.Errorf("failed to get index mappings: %w", err)
		}
	}

	if shouldGetSettings {
		settingsEndpoint := fmt.Sprintf("%s/_settings", index)
		if err := getJSONResponse(settingsEndpoint, &settingsResponse); err != nil {
			return nil, fmt.Errorf("failed to get index settings: %w", err)
		}
	}

	// Merge settings and mappings.
	merged := make(IndexDetailsResponse)
	for index, indexMappings := range mappingsResponse {
		details := IndexDetails{
			Mappings: indexMappings.Mappings,
		}

		if indexSettings, ok := settingsResponse[index]; ok {
			details.Settings = indexSettings.Settings
		}

		merged[index] = details
	}

	// Process settings for indices that don't exist in mappings.
	for index, indexSettings := range settingsResponse {
		if _, exists := merged[index]; !exists {
			merged[index] = IndexDetails{
				Settings: indexSettings.Settings, // Only settings are added.
			}
		}
	}

	return merged, nil
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
