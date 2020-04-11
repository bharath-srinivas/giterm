package pages

import (
	"github.com/rivo/tview"

	"github.com/bharath-srinivas/giterm/config"
	"github.com/bharath-srinivas/giterm/modules"
)

// FeedsPage returns the feeds page with pagination and page size options.
func FeedsPage(app *tview.Application, config config.Config) *Page {
	feeds := modules.FeedsWidget(app, config)
	return &Page{
		Name:    "Feeds",
		Widgets: &Widgets{Parent: feeds},
	}
}
