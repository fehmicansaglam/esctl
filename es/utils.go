package es

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type EsError struct {
	Error struct {
		Type   string `json:"type"`
		Reason string `json:"reason"`
	} `json:"error"`
	Status int `json:"status"`
}

func getJSONResponse(url string, target interface{}) error {
	resp, err := http.Get(url)
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
