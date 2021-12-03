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
	menuBar    *MenuBar
	explorer   *Explorer
	editorTabs *EditorTabs
}

func NewApplicationView(log zerolog.Logger, events *event.Bus) *ApplicationView {
	ui := cview.NewApplication()

	vlayout := cview.NewFlex()
	vlayout.SetDirection(cview.FlexRow)
	ui.SetRoot(vlayout, true)

	menuBar := NewMenuBar(log, events)
	vlayout.AddItem(menuBar, 3, 0, false)

	hlayout := cview.NewFlex()
	vlayout.AddItem(hlayout, 0, 1, true)

	explorer := NewExplorer(log, events)
	hlayout.AddItem(explorer, 0, 2, false)

	editorTabs := NewEditorTabs(log, events)
	hlayout.AddItem(editorTabs, 0, 8, true)

	go func() {
		// FIXME: ugly workaround for the TabbedPanels not repainting as new ones are added
		for {
			ui.Draw(editorTabs)
			time.Sleep(34 * time.Millisecond) // about 30 FPS
		}
	}()

	return &ApplicationView{
		log:    log,
		events: events,

		ui:         ui,
		doneCh:     make(chan error, 1),
		explorer:   explorer,
		editorTabs: editorTabs,
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

func (view *ApplicationView) Stop() {
	view.ui.Stop()
	view.doneCh <- nil // FIXME: calling stop for some reason doesn't seem to make Run return...
}

func (view *ApplicationView) Done() <-chan error {
	return view.doneCh
}
