package es

import (
	"fmt"
	"strings"
)

type JsonResponse map[string]interface{}

func extractFieldAndValue(term string) (string, string, error) {
	parts := strings.SplitN(term, ":", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid term format: %s", term)
	}
	return parts[0], parts[1], nil
}

func SearchDocuments(index string, ids []string, terms []string, size int, nestedPaths []string) (JsonResponse, error) {
	endpoint := fmt.Sprintf("%s/_search", index)

	filters := make([]map[string]interface{}, 0)

	for _, term := range terms {
		field, value, err := extractFieldAndValue(term)
		if err != nil {
			return nil, err
		}
		nestedPath, isNestedPath := getNestedPath(field, nestedPaths)
		if isNestedPath {
			termFilter := map[string]interface{}{
				"nested": map[string]interface{}{
					"path": nestedPath,
					"query": map[string]interface{}{
						"term": map[string]interface{}{
							field: value,
						},
					},
				},
			}
			filters = append(filters, termFilter)
		} else {
			termFilter := map[string]interface{}{
				"term": map[string]interface{}{
					field: value,
				},
			}
			filters = append(filters, termFilter)
		}
	}

	if len(ids) > 0 {
		idsFilter := map[string]interface{}{
			"ids": map[string]interface{}{
				"values": ids,
			},
		}
		filters = append(filters, idsFilter)
	}

	query := map[string]interface{}{
		"bool": map[string]interface{}{
			"filter": filters,
		},
	}

	requestBody := map[string]interface{}{
		"size":  max(size, len(ids)),
		"query": query,
	}

	var response JsonResponse
	err := getJSONResponseWithBody(endpoint, &response, requestBody)
	if err != nil {
		return nil, err
	}

	return response, nil
}
