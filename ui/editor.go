package ui

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"

	"code.rocketnine.space/tslocum/cview"
	"github.com/gdamore/tcell/v2"
	"github.com/rs/zerolog"
	"github.com/spf13/afero"
	"github.com/tsatke/hackt/event"
)

type Editor struct {
	log    zerolog.Logger
	events *event.Bus

	tabs *cview.TabbedPanels
}

var tabNameSanitizingRegexp = regexp.MustCompile(`[^a-zA-Z0-9 ]`)

func NewEditor(log zerolog.Logger, events *event.Bus) *Editor {
	tabs := cview.NewTabbedPanels()
	tabs.SetBorder(true)
	tabs.SetTitle("Editor")

	editor := &Editor{
		log:    log,
		events: events,

		tabs: tabs,
	}

	events.ProjectFileOpen.Register(editor.processFileOpenRequest)

	return editor
}

func (e *Editor) processFileOpenRequest(event event.ProjectFileOpenPayload) {
	e.log.Debug().
		Str("path", event.Path).
		Str("project", event.Project.Name()).
		Msg("open file in editor")

	file, err := event.Project.Fs().OpenFile(event.Path, os.O_RDWR, 0666)
	if err != nil {
		e.log.Error().
			Err(err).
			Str("path", event.Path).
			Msg("unable to open file")
		return
	}

	tabName := tabNameSanitizingRegexp.ReplaceAllString(event.Project.Name()+" "+event.Path, " ")
	tab, err := NewEditorTab(e.log, e.events, file)
	if err != nil {
		e.log.Error().
			Err(err).
			Str("path", event.Path).
			Msg("unable to create tab from file")
		return
	}

	e.tabs.AddTab(tabName, filepath.Base(event.Path), tab)
}

// implement cview.Primitive, but delegate everything to the actual underlying TabbedPanels

func (e Editor) Draw(screen tcell.Screen) {
	e.tabs.Draw(screen)
}

func (e Editor) GetRect() (int, int, int, int) {
	return e.tabs.GetRect()
}

func (e Editor) SetRect(x, y, width, height int) {
	e.tabs.SetRect(x, y, width, height)
}

func (e Editor) GetVisible() bool {
	return e.tabs.GetVisible()
}

func (e Editor) SetVisible(v bool) {
	e.tabs.SetVisible(v)
}

func (e Editor) InputHandler() func(event *tcell.EventKey, setFocus func(p cview.Primitive)) {
	return e.tabs.InputHandler()
}

func (e Editor) Focus(delegate func(p cview.Primitive)) {
	e.tabs.Focus(delegate)
}

func (e Editor) Blur() {
	e.tabs.Blur()
}

func (e Editor) GetFocusable() cview.Focusable {
	return e.tabs.GetFocusable()
}

func (e Editor) MouseHandler() func(action cview.MouseAction, event *tcell.EventMouse, setFocus func(p cview.Primitive)) (consumed bool, capture cview.Primitive) {
	return e.tabs.MouseHandler()
}

type EditorTab struct {
	field       *cview.TextView
	backingFile afero.File
}

func NewEditorTab(log zerolog.Logger, events *event.Bus, file afero.File) (*EditorTab, error) {
	field := cview.NewTextView()
	_, err := io.Copy(field, file)
	if err != nil {
		return nil, fmt.Errorf("read: %w", err)
	}

	return &EditorTab{
		field: field,
	}, nil
}

func (tab *EditorTab) Close() error {
	return tab.backingFile.Close()
}

// implement cview.Primitive, but delegate everything to the actual underlying InputField

func (tab EditorTab) Draw(screen tcell.Screen) {
	tab.field.Draw(screen)
}

func (tab EditorTab) GetRect() (int, int, int, int) {
	return tab.field.GetRect()
}

func (tab EditorTab) SetRect(x, y, width, height int) {
	tab.field.SetRect(x, y, width, height)
}

func (tab EditorTab) GetVisible() bool {
	return tab.field.GetVisible()
}

func (tab EditorTab) SetVisible(v bool) {
	tab.field.SetVisible(v)
}

func (tab EditorTab) InputHandler() func(event *tcell.EventKey, setFocus func(p cview.Primitive)) {
	return tab.field.InputHandler()
}

func (tab EditorTab) Focus(delegate func(p cview.Primitive)) {
	tab.field.Focus(delegate)
}

func (tab EditorTab) Blur() {
	tab.field.Blur()
}

func (tab EditorTab) GetFocusable() cview.Focusable {
	return tab.field.GetFocusable()
}

func (tab EditorTab) MouseHandler() func(action cview.MouseAction, event *tcell.EventMouse, setFocus func(p cview.Primitive)) (consumed bool, capture cview.Primitive) {
	return tab.field.MouseHandler()
}
