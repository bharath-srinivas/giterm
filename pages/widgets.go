package pages

import (
	"github.com/rivo/tview"
)

// Widgets represents all the widgets in a page.
type Widgets struct {
	Parent   tview.Primitive
	Children []tview.Primitive
}

// Refreshable provides a method to refresh the widgets.
type Refreshable interface {
	Refresh()
}

// Prev returns the previous widget from the widget list.
func (w *Widgets) Prev() tview.Primitive {
	widgets := w.Children
	widgetLen := len(widgets)
	for i := 0; i < widgetLen; i++ {
		if widgets[i].GetFocusable().HasFocus() {
			prev := (i - 1) % widgetLen
			if prev < 0 {
				return widgets[widgetLen-1]
			}
			return widgets[prev]
		}
	}
	return widgets[0]
}

// Next returns the next widget from the widget list.
func (w *Widgets) Next() tview.Primitive {
	widgets := w.Children
	widgetLen := len(widgets)
	for i := 0; i < widgetLen; i++ {
		if widgets[i].GetFocusable().HasFocus() {
			return widgets[(i+1)%widgetLen]
		}
	}
	return widgets[0]
}

// Refresh refreshes all the widgets if they are refreshable.
func (w *Widgets) Refresh() {
	for _, widget := range w.Children {
		if widget, ok := widget.(Refreshable); ok {
			go widget.Refresh()
		}
	}
}
