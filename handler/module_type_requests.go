package handler

import "github.com/dell/csm-deployment/model"

type moduleResponse struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	Version    string `json:"version"`
	Standalone bool   `json:"stand_alone"`
} //@name ModuleResponse

func newModuleResponse(t *model.ModuleType) *moduleResponse {
	r := moduleResponse{}
	r.ID = t.ID
	r.Name = t.Name
	r.Version = t.Version
	r.Standalone = t.Standalone
	return &r
}
