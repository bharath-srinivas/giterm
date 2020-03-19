package pages

import (
	"github.com/rivo/tview"

	"github.com/bharath-srinivas/giterm/config"
	"github.com/bharath-srinivas/giterm/modules"
)

func ProfilePage(app *tview.Application, config config.Config) *Page {
	profile := modules.ProfileWidget(app, config)
	return &Page{
		Name:    "Profile",
		Widgets: &Widgets{Parent: profile},
	}
}
