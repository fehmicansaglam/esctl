package es

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/fehmicansaglam/esctl/shared"
)

type EsError struct {
	Error struct {
		Type   string `json:"type"`
		Reason string `json:"reason"`
	} `json:"error"`
	Status int `json:"status"`
}

func getJSONResponse(endpoint string, target interface{}) error {
	baseURL := fmt.Sprintf("%s://%s:%d/%s", shared.ElasticsearchProtocol, shared.ElasticsearchHost, shared.ElasticsearchPort, endpoint)

	req, err := http.NewRequest(http.MethodGet, baseURL, nil)
	if err != nil {
		return err
	}

	if shared.ElasticsearchUsername != "" && shared.ElasticsearchPassword != "" {
		req.SetBasicAuth(shared.ElasticsearchUsername, shared.ElasticsearchPassword)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var esError EsError
		if err := json.NewDecoder(resp.Body).Decode(&esError); err != nil {
			return fmt.Errorf("unexpected http status: %s", resp.Status)
		}
		return errors.New(esError.Error.Reason)
	}

	return json.NewDecoder(resp.Body).Decode(target)
}
