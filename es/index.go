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
	Count        int                    `json:"count"`
	Aggregations map[string]interface{} `json:"aggregations,omitempty"`
}

type GroupCount map[string]int
type IndexGroupCount map[string]GroupCount

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

func countDocumentsOfIndex(index string, termFilters []string, existsFilters []string) (int, error) {
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

func groupDocumentsOfIndex(index string, termFilters []string, existsFilters []string, groupBy string) (GroupCount, error) {
	endpoint := index + "/_search"
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
		"aggs": map[string]interface{}{
			"group_by": map[string]interface{}{
				"terms": map[string]interface{}{
					"field": groupBy,
				},
			},
		},
	}

	var response CountResponse
	if err := getJSONResponseWithBody(endpoint, &response, body); err != nil {
		return nil, err
	}

	groupCount := make(GroupCount)
	buckets, ok := response.Aggregations["group_by"].(map[string]interface{})
	if ok {
		terms, ok := buckets["buckets"].([]interface{})
		if ok {
			for _, term := range terms {
				termData, ok := term.(map[string]interface{})
				if ok {
					key := fmt.Sprint(termData["key"])
					count := int(termData["doc_count"].(float64))
					groupCount[key] = count
				}
			}
		}
	}

	return groupCount, nil
}

func getCountForIndex(index string, termFilters []string, existsFilters []string, groupBy string) (GroupCount, error) {
	if groupBy == "" {
		count, err := countDocumentsOfIndex(index, termFilters, existsFilters)
		if err != nil {
			return nil, err
		}
		return map[string]int{"": count}, nil
	}

	groupCount, err := groupDocumentsOfIndex(index, termFilters, existsFilters, groupBy)
	if err != nil {
		return nil, err
	}
	return groupCount, nil
}

func CountDocuments(index string, termFilters []string, existsFilters []string, groupBy string) (map[string]GroupCount, error) {
	indexCounts := make(map[string]GroupCount)

	if index == "" {
		indices, err := GetIndices("")
		if err != nil {
			return nil, err
		}

		for _, index := range indices {
			groupCount, err := getCountForIndex(index.Index, termFilters, existsFilters, groupBy)
			if err != nil {
				return nil, err
			}
			indexCounts[index.Index] = groupCount
		}
	} else {
		groupCount, err := getCountForIndex(index, termFilters, existsFilters, groupBy)
		if err != nil {
			return nil, err
		}
		indexCounts[index] = groupCount
	}

	return indexCounts, nil
}
