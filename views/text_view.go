package views

import (
	"github.com/rivo/tview"

	"github.com/bharath-srinivas/giterm/config"
)

type TextWidget struct {
	*Base
	*tview.TextView
}

func NewTextView(app *tview.Application, config config.Config, bordered bool) *TextWidget {
	widget := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetTextAlign(tview.AlignCenter)
	widget.SetBorder(bordered)
	return &TextWidget{
		Base:     NewBase(app, config),
		TextView: widget,
	}
}

func (t *TextWidget) Redraw(display func()) {
	t.app.QueueUpdateDraw(func() {
		t.Clear()
		display()
		t.ScrollToBeginning()
	})
}
