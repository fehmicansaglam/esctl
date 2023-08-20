package es

import (
	"net/url"
	"strings"
)

type TasksResponse struct {
	Nodes map[string]TaskNode `json:"nodes"`
}

type TaskNode struct {
	Name             string            `json:"name"`
	TransportAddress string            `json:"transport_address"`
	Host             string            `json:"host"`
	IP               string            `json:"ip"`
	Roles            []string          `json:"roles"`
	Attributes       map[string]string `json:"attributes"`
	Tasks            map[string]Task   `json:"tasks"`
}

type Task struct {
	Node               string                 `json:"node"`
	ID                 int64                  `json:"id"`
	Type               string                 `json:"type"`
	Action             string                 `json:"action"`
	Description        string                 `json:"description"`
	StartTimeInMillis  int64                  `json:"start_time_in_millis"`
	RunningTimeInNanos int64                  `json:"running_time_in_nanos"`
	Cancellable        bool                   `json:"cancellable"`
	Cancelled          bool                   `json:"cancelled"`
	ParentTaskID       string                 `json:"parent_task_id"`
	Headers            map[string]interface{} `json:"headers"`
}

func GetTasks(actions []string) (TasksResponse, error) {
	baseEndpoint := "_tasks"

	values := url.Values{
		"detailed": {"true"},
	}

	if len(actions) > 0 {
		actionsParam := strings.Join(actions, ",")
		values.Set("actions", actionsParam)
	}

	endpoint := baseEndpoint + "?" + values.Encode()

	var response TasksResponse
	if err := getJSONResponse(endpoint, &response); err != nil {
		return TasksResponse{}, err
	}

	return response, nil
}
