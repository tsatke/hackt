package ui

import (
	"time"

	"code.rocketnine.space/tslocum/cview"
	"github.com/rs/zerolog"
	"github.com/tsatke/hackt/event"
)

type ApplicationView struct {
	log    zerolog.Logger
	events *event.Bus

	doneCh     chan error
	ui         *cview.Application
	explorer   *Explorer
	editorArea *Editor
}

func NewApplicationView(log zerolog.Logger, events *event.Bus) *ApplicationView {
	ui := cview.NewApplication()

	layout := cview.NewFlex()
	ui.SetRoot(layout, true)

	explorer := NewExplorer(log, events)
	layout.AddItem(explorer, 0, 2, false)

	editorArea := NewEditor(log, events)
	layout.AddItem(editorArea, 0, 8, true)

	go func() {
		// FIXME: ugly workaround for the TabbedPanels not repainting as new ones are added
		for {
			ui.Draw(editorArea)
			time.Sleep(34 * time.Millisecond) // about 30 FPS
		}
	}()

	return &ApplicationView{
		log:    log,
		events: events,

		ui:         ui,
		doneCh:     make(chan error, 1),
		explorer:   explorer,
		editorArea: editorArea,
	}
}

func (view *ApplicationView) Run() {
	go func() {
		defer view.ui.HandlePanic()

		// enable mouse support
		view.ui.EnableMouse(true)
		view.ui.SetDoubleClickInterval(cview.StandardDoubleClick)

		view.doneCh <- view.ui.Run()
	}()
}

func (view *ApplicationView) Done() <-chan error {
	return view.doneCh
}
