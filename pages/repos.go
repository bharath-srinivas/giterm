package pages

import (
	"github.com/rivo/tview"

	"github.com/bharath-srinivas/giterm/config"
	"github.com/bharath-srinivas/giterm/modules"
)

func ReposPage(app *tview.Application, config config.Config) *Page {
	repos := modules.RepoWidget(app, config)
	return &Page{
		Name:    "Repos",
		Widgets: &Widgets{Parent: repos},
	}
}
