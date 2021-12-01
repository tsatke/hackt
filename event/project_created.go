package event

import "github.com/tsatke/hackt/workspace"

type ProjectCreated struct {
	handlers []ProjectCreatedHandler
}

type ProjectCreatedHandler func(ProjectCreatedPayload)

type ProjectCreatedPayload struct {
	Project *workspace.Project
}

func (evt *ProjectCreated) Register(handler ProjectCreatedHandler) {
	evt.handlers = append(evt.handlers, handler)
}

func (evt *ProjectCreated) Trigger(payload ProjectCreatedPayload) {
	for _, handler := range evt.handlers {
		go handler(payload)
	}
}
