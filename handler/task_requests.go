package handler

import "github.com/dell/csm-deployment/model"

type taskResponse struct {
	ID            uint                         `json:"id"`
	Status        string                       `json:"status"`
	ApplicationID uint                         `json:"application_id"`
	Logs          string                       `json:"logs"`
	Links         map[string]map[string]string `json:"_links"`
} //@name TaskResponse

func newTaskResponse(t *model.Task) *taskResponse {
	r := taskResponse{}
	r.ID = t.ID
	r.Status = t.Status
	r.ApplicationID = t.ApplicationID
	r.Logs = string(t.Logs)
	return &r
}

func newTaskResponseWithLinks(t *model.Task, links map[string]map[string]string) *taskResponse {
	r := newTaskResponse(t)
	r.Links = links
	return r
}
