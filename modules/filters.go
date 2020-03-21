package modules

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

// Filterer provides a method which can filter the data according to the provided option.
type Filterer interface {
	Filter(string)
}

// FilterWidget returns a new drop down widget with specified label and options. It calls the Filter method of the Filterer
//interface to filter the data whenever an option is selected.
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
