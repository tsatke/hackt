package component

import (
	"code.rocketnine.space/tslocum/cview"
	"github.com/gdamore/tcell/v2"
)

var _ cview.Primitive = (*MenuBar)(nil)

type MenuBar struct {
	*cview.Form
}

func NewMenuBar() *MenuBar {
	form := cview.NewForm()
	form.SetHorizontal(true)

	return &MenuBar{
		Form: form,
	}
}

func (b *MenuBar) AddButton(label string, selected func()) {
	b.Form.AddButton(label, selected)
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
