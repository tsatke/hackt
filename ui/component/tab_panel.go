package component

import (
	"code.rocketnine.space/tslocum/cview"
	"github.com/gdamore/tcell/v2"
	"github.com/google/uuid"
)

var _ cview.Primitive = (*TabPanel)(nil)

type TabPanel struct {
	*cview.Box

	openTab *Tab
	tabs    []*Tab
}

func (t TabPanel) Tabs() []*Tab { return t.tabs }

func (t TabPanel) Draw(screen tcell.Screen) {
	t.Box.Draw(screen)

	// x, y, width, height := t.Box.GetInnerRect()
}

type Tab struct {
	primitive cview.Primitive

	id       uuid.UUID
	title    string
	dirty    bool
	closable bool
}

func (t *Tab) Content() cview.Primitive {
	return t.primitive
}

func (t *Tab) Title() string         { return t.title }
func (t *Tab) SetTitle(title string) { t.title = title }

func (t *Tab) Dirty() bool         { return t.dirty }
func (t *Tab) SetDirty(dirty bool) { t.dirty = dirty }

func (t *Tab) Closable() bool            { return t.closable }
func (t *Tab) SetClosable(closable bool) { t.closable = closable }
