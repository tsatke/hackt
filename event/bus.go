package event

import "github.com/rs/zerolog"

type Bus struct {
	log zerolog.Logger

	ProjectCreated
	ProjectFileOpen
	ProjectLoaded
}

func NewBus(log zerolog.Logger) *Bus {
	return &Bus{
		log: log,
	}
}
