package ui

import (
	"fmt"

	"code.rocketnine.space/tslocum/cview"
	"github.com/gdamore/tcell/v2"
	"github.com/rs/zerolog"
	"github.com/spf13/afero"
	"github.com/tsatke/hackt/event"
)

var _ cview.Primitive = (*Editor)(nil)

type Editor struct {
	log zerolog.Logger

	EditorUI
	cursor CursorPosition

	backingFile afero.File
	content     *EditorContent
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
		log: log,
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
	for screenLine := y; screenLine < y+height; screenLine++ {
		contentLine := screenLine - y
		if contentLine >= len(lines) {
			break
		}
		lineBytes := lines[contentLine]
		lineRunes := []rune(string(lineBytes))
		for screenColumn := x; screenColumn < x+width; screenColumn++ {
			contentColumn := screenColumn - x

			// check if we should draw the cursor here
			if e.cursor.line == contentLine && e.cursor.column == contentColumn {
				screen.ShowCursor(screenColumn, screenLine)
			}

			r := ' '
			if contentColumn < len(lineRunes) {
				r = lineRunes[contentColumn]
			}

			screen.SetCell(screenColumn, screenLine, style, r)
		}
	}
}

func (e *Editor) inputCapture(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyUp:
		e.cursorUp()
	case tcell.KeyDown:
		e.cursorDown()
	case tcell.KeyLeft:
		e.cursorLeft()
	case tcell.KeyRight:
		e.cursorRight()
	default:
		r := event.Rune()
		if r != 0 {
			e.insertRune(r)
			e.cursorRight()
		}
	}

	return event
}

func (e *Editor) insertRune(r rune) {
	_, _ = e.content.InsertAt([]byte(string(r)), e.cursorOffset())
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
