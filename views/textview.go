package views

import (
	"github.com/rivo/tview"

	"github.com/bharath-srinivas/giterm/config"
)

// TextWidget represents a text widget used by the application.
type TextWidget struct {
	*Base
	*tview.TextView
}

// NewTextView returns a new text widget with the provided options.
func NewTextView(app *tview.Application, config config.Config, bordered bool) *TextWidget {
	widget := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true)
	widget.SetBorder(bordered)
	return &TextWidget{
		NewBase(app, config),
		widget,
	}
}

// Redraw refreshes the text widget.
func (t *TextWidget) Redraw(display func()) {
	t.app.QueueUpdateDraw(func() {
		t.Clear()
		display()
		t.ScrollToBeginning()
	})
}
