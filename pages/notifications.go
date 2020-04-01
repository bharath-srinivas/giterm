package pages

import (
	"github.com/rivo/tview"

	"github.com/bharath-srinivas/giterm/config"
	"github.com/bharath-srinivas/giterm/modules"
)

// NotificationsPage returns the notifications page with pagination options.
func NotificationsPage(app *tview.Application, config config.Config) *Page {
	notifications := modules.NotificationsWidget(app, config)
	pagination := modules.PaginationWidget(notifications)

	footer := tview.NewFlex().
		AddItem(pagination.Prev, 0, 1, false).
		AddItem(pagination.Next, 0, 1, false)

	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(notifications, 0, 15, false).
		AddItem(footer, 0, 1, false)
	layout.SetTitle(string('\U0001F514') + " [green::b]Notifications").
		SetBorder(true)

	return &Page{
		Name: "Notifications",
		Widgets: &Widgets{
			Parent: layout,
			Children: []tview.Primitive{
				notifications,
				pagination.Prev,
				pagination.Next,
			},
		},
	}
}
