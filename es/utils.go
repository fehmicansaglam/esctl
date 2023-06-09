package es

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

func httpRequest(method, endpoint string, body, target interface{}, expectedStatusCode int) error {
	baseURL := fmt.Sprintf("%s://%s:%d/%s", shared.ElasticsearchProtocol, shared.ElasticsearchHost, shared.ElasticsearchPort, endpoint)

	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return err
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequest(method, baseURL, bodyReader)
	if err != nil {
		return err
	}

	if shared.ElasticsearchUsername != "" && shared.ElasticsearchPassword != "" {
		req.SetBasicAuth(shared.ElasticsearchUsername, shared.ElasticsearchPassword)
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != expectedStatusCode {
		var esError EsError
		if err := json.NewDecoder(resp.Body).Decode(&esError); err != nil {
			return fmt.Errorf("unexpected http status: %s", resp.Status)
		}
		return errors.New(esError.Error.Reason)
	}

	return json.NewDecoder(resp.Body).Decode(target)
}

func getJSONResponse(endpoint string, target interface{}) error {
	return httpRequest(http.MethodGet, endpoint, nil, target, http.StatusOK)
}

func getJSONResponseWithBody(endpoint string, target interface{}, body interface{}) error {
	return httpRequest(http.MethodPost, endpoint, body, target, http.StatusOK)
}
