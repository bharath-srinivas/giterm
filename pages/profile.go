package pages

import (
	"github.com/rivo/tview"

	"github.com/bharath-srinivas/giterm/config"
	"github.com/bharath-srinivas/giterm/modules"
)

// ProfilePage returns the profile page.
func ProfilePage(app *tview.Application, config config.Config) *Page {
	user := modules.UserWidget(app, config)
	userStats := modules.UserStatsWidget(app, config)
	avatar := modules.AvatarWidget(app, config)
	contribution := modules.ContributionsWidget(app, config)

	userFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(avatar, 0, 2, false).
		AddItem(user, 0, 1, false).
		AddItem(userStats, 0, 1, false)

	layout := tview.NewFlex().
		AddItem(userFlex, 0, 1, false).
		AddItem(contribution, 0, 3, false)

	layout.SetTitle(string('\U0001F642') + " [green::b]Profile").
		SetBorder(true)

	return &Page{
		Name: "Profile",
		Widgets: &Widgets{
			Parent:   layout,
			Children: []tview.Primitive{avatar, user, userStats, contribution},
		},
	}
}
