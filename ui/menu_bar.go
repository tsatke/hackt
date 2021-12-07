package ui

import (
	"github.com/rs/zerolog"
	"github.com/tsatke/hackt/event"
	"github.com/tsatke/hackt/ui/component"
)

type MenuBar struct {
	*component.MenuBar

	log    zerolog.Logger
	events *event.Bus
}

func NewMenuBar(log zerolog.Logger, events *event.Bus) *MenuBar {
	menuBar := component.NewMenuBar()
	menuBar.AddButton("open", func() {
		log.Error().
			Msg("open not supported yet")
	})
	menuBar.AddButton("about", func() {
		log.Error().
			Msg("about not supported yet")
	})
	menuBar.AddButton("exit", func() {
		events.MenuBar.MenuBarExit.Trigger(event.MenuBarExitPayload{})
	})

	return &MenuBar{
		MenuBar: menuBar,

		log:    log,
		events: events,
	}
}
