package event

import "github.com/rs/zerolog"

type Bus struct {
	log zerolog.Logger

	UI struct {
		UIRedraw
	}
	Project struct {
		ProjectCreated
		ProjectFileOpen
		ProjectLoaded
	}
	MenuBar struct {
		MenuBarExit
	}
}

func NewBus(log zerolog.Logger) *Bus {
	return &Bus{
		log: log,
	}
}
