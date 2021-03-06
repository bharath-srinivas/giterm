package pages

import (
	"github.com/rivo/tview"

	"github.com/bharath-srinivas/giterm/config"
	"github.com/bharath-srinivas/giterm/modules"
)

// ReposPage returns the repository page with filters, pagination and page size options.
func ReposPage(app *tview.Application, config config.Config) *Page {
	repos := modules.RepoWidget(app, config)
	pageSizes := modules.PageSizeWidget(repos)
	pagination := modules.PaginationWidget(repos)
	filters := modules.FilterWidget("Type: ", []string{"All", "Public", "Private", "Sources", "Forks"}, repos)

	header := tview.NewFlex().
		AddItem(pageSizes, 0, 1, false).
		AddItem(filters, 0, 1, false)

	footer := tview.NewFlex().
		AddItem(pagination.Prev, 0, 1, false).
		AddItem(pagination.Next, 0, 1, false)

	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(header, 0, 1, false).
		AddItem(repos, 0, 15, false).
		AddItem(footer, 0, 1, false)
	layout.SetTitle("\U0001F4D5 [green::b]Repositories").
		SetBorder(true)

	return &Page{
		Name: "Repos",
		Widgets: &Widgets{
			Parent: layout,
			Children: []tview.Primitive{
				pageSizes,
				filters,
				repos,
				pagination.Prev,
				pagination.Next,
			},
		},
	}
}
