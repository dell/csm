package handler

import "github.com/dell/csm-deployment/model"

type taskResponse struct {
	Task struct {
		ID            uint                         `json:"id"`
		Status        string                       `json:"status"`
		ApplicationID uint                         `json:"application_id"`
		Logs          string                       `json:"logs"`
		Links         map[string]map[string]string `json:"_links"`
	} `json:"task"`
}

func newTaskResponse(t *model.Task) *taskResponse {
	r := taskResponse{}
	r.Task.ID = t.ID
	r.Task.Status = t.Status
	r.Task.ApplicationID = t.ApplicationID
	r.Task.Logs = string(t.Logs)
	return &r
}

func newTaskResponseWithLinks(t *model.Task, links map[string]map[string]string) *taskResponse {
	r := newTaskResponse(t)
	r.Task.Links = links
	return r
}
