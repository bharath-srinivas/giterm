package views

import (
	"github.com/rivo/tview"

	"github.com/bharath-srinivas/giterm/config"
)

type TextWidget struct {
	*Base
	View *tview.TextView
}

func NewTextView(app *tview.Application, config config.Config, bordered bool) *TextWidget {
	widget := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetTextAlign(tview.AlignCenter)
	widget.SetBorder(bordered)
	return &TextWidget{
		Base: NewBase(app, config),
		View: widget,
	}
}

func (t *TextWidget) Redraw(display func()) {
	t.app.QueueUpdateDraw(func() {
		t.View.Clear()
		display()
		t.View.ScrollToBeginning()
	})
}
