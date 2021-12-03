package ui

import (
	"code.rocketnine.space/tslocum/cview"
	"github.com/gdamore/tcell/v2"
	"github.com/rs/zerolog"
	"github.com/tsatke/hackt/event"
)

type MenuBar struct {
	MenuBarUI

	log    zerolog.Logger
	events *event.Bus
}

type MenuBarUI struct {
	*cview.Form
}

func NewMenuBar(log zerolog.Logger, events *event.Bus) *MenuBar {
	form := cview.NewForm()
	form.SetHorizontal(true)
	form.AddButton("open", func() {
		log.Error().
			Msg("open not supported yet")
	})
	form.AddButton("about", func() {
		log.Error().
			Msg("about not supported yet")
	})
	form.AddButton("exit", func() {
		events.MenuBar.MenuBarExit.Trigger(event.MenuBarExitPayload{})
	})

	return &MenuBar{
		log:    log,
		events: events,
		MenuBarUI: MenuBarUI{
			Form: form,
		},
	}
}

func (b *MenuBar) Draw(screen tcell.Screen) {
	b.Form.Draw(screen)

	// x, y, width, height := b.Form.GetRect()
	// for screenLine := y; screenLine < y+height; screenLine++ {
	// 	for screenColumn := x; screenColumn < x+width; screenColumn++ {
	// 		screen.SetCell(screenColumn, screenLine, tcell.StyleDefault, 'x')
	// 	}
	// }
}
