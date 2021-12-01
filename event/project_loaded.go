package event

import "github.com/tsatke/hackt/workspace"

type ProjectLoaded struct {
	handlers []ProjectLoadedHandler
}

type ProjectLoadedHandler func(ProjectLoadedPayload)

type ProjectLoadedPayload struct {
	Project *workspace.Project
}

func (evt *ProjectLoaded) Register(handler ProjectLoadedHandler) {
	evt.handlers = append(evt.handlers, handler)
}

func (evt *ProjectLoaded) Trigger(payload ProjectLoadedPayload) {
	for _, handler := range evt.handlers {
		go handler(payload)
	}
}
