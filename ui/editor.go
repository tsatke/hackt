package ui

import (
	"fmt"

	"code.rocketnine.space/tslocum/cview"
	"github.com/davecgh/go-spew/spew"
	"github.com/gdamore/tcell/v2"
	"github.com/rs/zerolog"
	"github.com/spf13/afero"
	"github.com/tsatke/hackt/event"
)

var _ cview.Primitive = (*Editor)(nil)

type Editor struct {
	log    zerolog.Logger
	events *event.Bus

	EditorUI
	cursor CursorPosition

	backingFile afero.File
	content     *EditorContent
	startLine   int
}

type CursorPosition struct {
	line   int
	column int
}

type EditorUI struct {
	cview.Primitive
	contentArea *cview.Box
	layout      cview.Primitive
}

func NewEditorTab(log zerolog.Logger, events *event.Bus, file afero.File) (*Editor, error) {
	layout := cview.NewFlex()
	contentArea := cview.NewBox()
	layout.AddItem(contentArea, 0, 10, false)

	data, err := afero.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	content := NewBuffer(data)

	e := &Editor{
		log:    log,
		events: events,
		EditorUI: EditorUI{
			Primitive:   layout,
			layout:      layout,
			contentArea: contentArea,
		},
		backingFile: file,
		content:     content,
	}

	contentArea.SetInputCapture(e.inputCapture)

	return e, nil
}

func (e *Editor) Close() error {
	return e.backingFile.Close()
}

func (e *Editor) Draw(screen tcell.Screen) {
	e.layout.Draw(screen)

	style := tcell.Style{}.Foreground(tcell.ColorWhite).Background(tcell.ColorBlack)

	lines := e.content.Lines()

	x, y, width, height := e.EditorUI.contentArea.GetInnerRect()

	// draw cursor
	screenCursorX := x + e.cursor.column
	screenCursorY := y + e.cursor.line - e.startLine
	if e.EditorUI.contentArea.InRect(screenCursorX, screenCursorY) {
		screen.ShowCursor(x+e.cursor.column, y+e.cursor.line-e.startLine)
	} else if screenCursorY < y {
		e.startLine--
	} else if screenCursorY > y+height {
		e.startLine++
	} else {
		// FIXME: handle cursor out of bounds right/left of the screen
	}

	for screenLine := y; screenLine < y+height; screenLine++ {
		contentLine := screenLine - y + e.startLine
		if contentLine >= len(lines) {
			// if there are no more lines, don't try to draw any
			break
		}
		lineBytes := lines[contentLine]
		lineRunes := []rune(string(lineBytes))
		for screenColumn := x; screenColumn < x+width; screenColumn++ {

			contentColumn := screenColumn - x

			r := ' '
			if contentColumn < len(lineRunes) {
				r = lineRunes[contentColumn]
			}

			screen.SetCell(screenColumn, screenLine, style, r)

		}
		// render scroll bar on top of the last column
		cview.RenderScrollBar(screen, cview.ScrollBarAuto, x+width-1, screenLine, height, len(lines), e.cursor.line, contentLine-e.startLine, true, tcell.ColorWhite)
	}
}

func (e *Editor) inputCapture(evt *tcell.EventKey) *tcell.EventKey {
	e.log.Info().Msg(spew.Sdump(evt))

	switch evt.Key() {
	case tcell.KeyUp:
		e.cursorUp()
	case tcell.KeyDown:
		e.cursorDown()
	case tcell.KeyLeft:
		e.cursorLeft()
	case tcell.KeyRight:
		e.cursorRight()
	case tcell.KeyDEL:
		e.backspace()
	case tcell.KeyEnter:
		e.insertRune('\n')
		e.cursorRight()
	default:
		r := evt.Rune()
		if r != 0 {
			e.insertRune(r)
			e.cursorRight()
		}
	}

	e.events.UI.UIRedraw.Trigger(event.UIRedrawPayload{
		Components: []cview.Primitive{e},
	})

	return evt
}

func (e *Editor) insertRune(r rune) {
	e.content.InsertAt([]byte(string(r)), e.cursorOffset())
}

func (e *Editor) backspace() {
	if e.cursorOffset() == 0 {
		// don't delete at the start of the file
		return
	}
	e.cursorLeft()
	e.delete()
}

func (e *Editor) delete() {
	e.content.DeleteAt(1, e.cursorOffset())
}

func (e *Editor) cursorOffset() (offset int64) {
	lines := e.content.Lines()
	for i := 0; i < e.cursor.line; i++ {
		offset += int64(len(lines[i]))
		offset++ // add 1 for a linefeed byte
	}

	offset += int64(e.cursor.column)

	return
}

func (e *Editor) cursorUp() {
	lines := e.content.Lines()

	if e.cursor.line > 0 {
		e.cursor.line--
	}
	if e.cursor.column > len(lines[e.cursor.line]) {
		e.cursor.column = len(lines[e.cursor.line])
	}
}

func (e *Editor) cursorDown() {
	lines := e.content.Lines()

	if e.cursor.line < len(lines)-1 {
		e.cursor.line++
	}
	if e.cursor.column > len(lines[e.cursor.line]) {
		e.cursor.column = len(lines[e.cursor.line])
	}
}

func (e *Editor) cursorLeft() {
	lines := e.content.Lines()

	if e.cursor.column > 0 {
		e.cursor.column--
	} else if e.cursor.line > 0 {
		e.cursor.line--
		e.cursor.column = len(lines[e.cursor.line])
	}
}

func (e *Editor) cursorRight() {
	lines := e.content.Lines()

	if e.cursor.column < len(lines[e.cursor.line]) {
		e.cursor.column++
	} else if e.cursor.line < len(lines)-1 {
		e.cursor.line++
		e.cursor.column = 0
	}
}
