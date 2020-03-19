package pages

import (
	"github.com/rivo/tview"

	"github.com/bharath-srinivas/giterm/config"
	"github.com/bharath-srinivas/giterm/modules"
	"github.com/bharath-srinivas/giterm/views"
)

func FeedsPage(app *tview.Application, config config.Config) *Page {
	feeds := modules.FeedsWidget(app, config)
	pageSizes := views.PageSizeWidget(feeds)
	pagination := views.PaginationWidget(feeds)

	header := tview.NewFlex().
		AddItem(pageSizes, 0, 1, false)

	footer := tview.NewFlex().
		AddItem(pagination.First, 0, 1, false).
		AddItem(pagination.Prev, 0, 1, false).
		AddItem(pagination.Next, 0, 1, false).
		AddItem(pagination.Last, 0, 1, false)

	view := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(header, 0, 1, false).
		AddItem(feeds.View, 0, 15, false).
		AddItem(footer, 0, 1, false)
	view.SetTitle(string('\U0001F559') + " [green::b]Feeds")
	view.SetBorder(true)

	return &Page{
		Name: "Feeds",
		Widgets: &Widgets{
			Parent: view,
			Children: []tview.Primitive{
				pageSizes,
				feeds.View,
				pagination.First,
				pagination.Prev,
				pagination.Next,
				pagination.Last,
			},
		},
	}
}
