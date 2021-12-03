package ui

import (
	"os"
	"path/filepath"
	"regexp"

	"code.rocketnine.space/tslocum/cview"
	"github.com/rs/zerolog"
	"github.com/tsatke/hackt/event"
)

type EditorTabs struct {
	log    zerolog.Logger
	events *event.Bus

	*cview.TabbedPanels
}

var tabNameSanitizingRegexp = regexp.MustCompile(`[^a-zA-Z0-9 ]`)

func NewEditorTabs(log zerolog.Logger, events *event.Bus) *EditorTabs {
	tabs := cview.NewTabbedPanels()
	tabs.SetBorder(true)
	tabs.SetTitle("Editor")

	editor := &EditorTabs{
		log:    log,
		events: events,

		TabbedPanels: tabs,
	}

	events.Project.ProjectFileOpen.Register(editor.processFileOpenRequest)

	return editor
}

func (e *EditorTabs) processFileOpenRequest(evt event.ProjectFileOpenPayload) {
	e.log.Debug().
		Str("path", evt.Path).
		Str("project", evt.Project.Name()).
		Msg("open file in editor")

	file, err := evt.Project.Fs().OpenFile(evt.Path, os.O_RDWR, 0666)
	if err != nil {
		e.log.Error().
			Err(err).
			Str("path", evt.Path).
			Msg("unable to open file")
		return
	}

	tabName := tabNameSanitizingRegexp.ReplaceAllString(evt.Project.Name()+" "+evt.Path, " ")
	tab, err := NewEditorTab(e.log, e.events, file)
	if err != nil {
		e.log.Error().
			Err(err).
			Str("path", evt.Path).
			Msg("unable to create tab from file")
		return
	}

	e.TabbedPanels.AddTab(tabName, filepath.Base(evt.Path), tab)
	e.TabbedPanels.SetCurrentTab(tabName)

	e.events.UI.UIRedraw.Trigger(event.UIRedrawPayload{
		Components: []cview.Primitive{e},
	})
}
