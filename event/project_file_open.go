package event

import "github.com/tsatke/hackt/workspace"

type ProjectFileOpen struct {
	handlers []ProjectFileOpenHandler
}

type ProjectFileOpenHandler func(ProjectFileOpenPayload)

type ProjectFileOpenPayload struct {
	Project *workspace.Project
	Path    string
}

func (evt *ProjectFileOpen) Register(handler ProjectFileOpenHandler) {
	evt.handlers = append(evt.handlers, handler)
}

func (evt *ProjectFileOpen) Trigger(payload ProjectFileOpenPayload) {
	for _, handler := range evt.handlers {
		go handler(payload)
	}
}
