// Package pages implements all the pages that are displayed in the application.
package pages

import (
	"github.com/rivo/tview"

	"github.com/bharath-srinivas/giterm/config"
)

// Page represents a page in the application.
type Page struct {
	Name string
	*Widgets
}

// MakePage returns a new page according to the given page name.
func MakePage(app *tview.Application, config config.Config, pageName string) *Page {
	var page *Page
	switch pageName {
	case "profile":
		page = ProfilePage(app, config)
	case "repos":
		page = ReposPage(app, config)
	case "feeds":
		page = FeedsPage(app, config)
	case "default":
		page = &Page{
			Name:    "",
			Widgets: &Widgets{Parent: tview.NewBox().SetTitle("Unknown").SetBorder(true)},
		}
	}
	return page
}

// PrevWidget returns the previous widget in a page. If no child widgets are present, it will set the focus on the parent.
func (p *Page) PrevWidget() tview.Primitive {
	if len(p.Children) > 0 {
		return p.Widgets.Prev()
	}
	return p.Widgets.Parent
}

// NextWidget returns the next widget in a page. If no child widgets are present, it will set the focus on the parent.
func (p *Page) NextWidget() tview.Primitive {
	if len(p.Children) > 0 {
		return p.Widgets.Next()
	}
	return p.Widgets.Parent
}

// Refresh refreshes all the widgets in a page.
func (p *Page) Refresh() {
	if len(p.Children) > 0 {
		p.Widgets.Refresh()
	} else if parent, ok := p.Parent.(Refreshable); ok {
		go parent.Refresh()
	}
}

// Pages represents a collection of pages.
type Pages []*Page

// MakePages returns all the pages needed by the application.
func MakePages(app *tview.Application, config config.Config) Pages {
	var p Pages
	pages := []string{"feeds", "profile", "repos"}
	for _, pageName := range pages {
		p = append(p, MakePage(app, config, pageName))
	}
	return p
}

// Get returns the page according to the given page name if it exists.
func (p Pages) Get(name string) *Page {
	for _, page := range p {
		if page.Name == name {
			return page
		}
	}
	return nil
}

// Prev returns the previous page name based on the given page name.
func (p Pages) Prev(currentPage string) string {
	pageCount := len(p)
	currentPageIndex := p.getPageIndex(currentPage)
	prevPageIndex := (currentPageIndex - 1) % pageCount
	if prevPageIndex < 0 {
		return p[pageCount-1].Name
	}
	return p[prevPageIndex].Name
}

// Next returns the next page name based on the given page name.
func (p Pages) Next(currentPage string) string {
	pageCount := len(p)
	currentPageIndex := p.getPageIndex(currentPage)
	nextPageIndex := (currentPageIndex + 1) % pageCount
	return p[nextPageIndex].Name
}

// getPageIndex returns the index of the given page name if it exists else it will return a negative value.
func (p Pages) getPageIndex(name string) int {
	for index, page := range p {
		if page.Name == name {
			return index
		}
	}
	return -1
}
