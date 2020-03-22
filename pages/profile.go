package pages

import (
	"github.com/rivo/tview"

	"github.com/bharath-srinivas/giterm/config"
	"github.com/bharath-srinivas/giterm/modules"
)

// ProfilePage returns the profile page.
func ProfilePage(app *tview.Application, config config.Config) *Page {
	user := modules.UserWidget(app, config)
	contribution := modules.ContributionsWidget(app, config)

	layout := tview.NewFlex().
		AddItem(user, 0, 1, false).
		AddItem(contribution, 0, 4, false)

	layout.SetTitle(string('\U0001F642') + " [green::b]Profile").
		SetBorder(true)

	return &Page{
		Name: "Profile",
		Widgets: &Widgets{
			Parent:   layout,
			Children: []tview.Primitive{user, contribution},
		},
	}
}
