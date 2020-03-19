package pages

import (
	"github.com/rivo/tview"
)

type Widgets struct {
	Parent   tview.Primitive
	Children []tview.Primitive
}

type Refreshable interface {
	Refresh()
}

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

func (w *Widgets) Refresh() {
	for _, widget := range w.Children {
		if widget, ok := widget.(Refreshable); ok {
			go widget.Refresh()
		}
	}
}
