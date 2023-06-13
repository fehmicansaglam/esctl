package es

import (
	"fmt"
	"strings"
)

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

type CountResponse struct {
	Count int `json:"count"`
}

func parseTermFilter(filter string) map[string]interface{} {
	parts := strings.SplitN(filter, ":", 2)
	if len(parts) != 2 {
		return nil
	}
	return map[string]interface{}{
		"term": map[string]interface{}{
			parts[0]: parts[1],
		},
	}
}

func parseExistsFilter(filter string) map[string]interface{} {
	return map[string]interface{}{
		"exists": map[string]interface{}{
			"field": filter,
		},
	}
}

func CountDocuments(index string, termFilters []string, existsFilters []string) (int, error) {
	endpoint := index + "/_count"
	query := map[string]interface{}{
		"match_all": map[string]interface{}{},
	}

	if len(termFilters) > 0 || len(existsFilters) > 0 {
		filterQueries := make([]map[string]interface{}, 0, len(termFilters)+len(existsFilters))

		for _, filter := range termFilters {
			filterQuery := parseTermFilter(filter)
			if filterQuery != nil {
				filterQueries = append(filterQueries, filterQuery)
			}
		}

		for _, filter := range existsFilters {
			filterQuery := parseExistsFilter(filter)
			if filterQuery != nil {
				filterQueries = append(filterQueries, filterQuery)
			}
		}

		query = map[string]interface{}{
			"bool": map[string]interface{}{
				"must": filterQueries,
			},
		}
	}

	body := map[string]interface{}{
		"query": query,
	}

	var response CountResponse
	if err := getJSONResponseWithBody(endpoint, &response, body); err != nil {
		return 0, err
	}

	return response.Count, nil
}

func GetDocumentCounts(termFilters []string, existsFilters []string) (map[string]int, error) {
	indices, err := GetIndices("")
	if err != nil {
		return nil, err
	}

	indexCounts := make(map[string]int)
	for _, index := range indices {
		count, err := CountDocuments(index.Index, termFilters, existsFilters)
		if err != nil {
			return nil, err
		}
		indexCounts[index.Index] = count
	}

	return indexCounts, nil
}
