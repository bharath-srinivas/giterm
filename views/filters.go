package views

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type Filterer interface {
	Filter(string)
}

func FilterWidget(label string, options []string, filterer Filterer) *tview.DropDown {
	filters := tview.NewDropDown().
		SetLabelColor(tcell.ColorWhite).
		SetFieldTextColor(tcell.ColorWhite).
		SetFieldBackgroundColor(tcell.ColorBlack).
		SetLabel(label).
		SetOptions(options, func(text string, index int) {
			filterer.Filter(text)
		}).SetCurrentOption(0)
	return filters
}
