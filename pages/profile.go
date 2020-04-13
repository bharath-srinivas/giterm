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
	pinned := modules.PinnedWidget(app)
	contribution := modules.ContributionsWidget(app, config)

	leftPane := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(avatar, 0, 2, false).
		AddItem(user, 0, 1, false).
		AddItem(userStats, 0, 1, false)

	rightPane := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(pinned, 0, 1, false).
		AddItem(contribution, 0, 3, false)

	layout := tview.NewFlex().
		AddItem(leftPane, 0, 1, false).
		AddItem(rightPane, 0, 3, false)

	layout.SetTitle("\U0001F642 [green::b]Profile").
		SetBorder(true)

	return &Page{
		Name: "Profile",
		Widgets: &Widgets{
			Parent:   layout,
			Children: []tview.Primitive{avatar, user, userStats, pinned, contribution},
		},
	}
}
