package pages

import (
	"github.com/rivo/tview"

	"github.com/bharath-srinivas/giterm/config"
)

type Page struct {
	Name string
	*Widgets
}

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

func (p *Page) PrevWidget() tview.Primitive {
	if len(p.Children) > 0 {
		return p.Widgets.Prev()
	}
	return p.Widgets.Parent
}

func (p *Page) NextWidget() tview.Primitive {
	if len(p.Children) > 0 {
		return p.Widgets.Next()
	}
	return p.Widgets.Parent
}

func (p *Page) Refresh() {
	if len(p.Children) > 0 {
		p.Widgets.Refresh()
	} else if parent, ok := p.Parent.(Refreshable); ok {
		go parent.Refresh()
	}
}

type Pages []*Page

func MakePages(app *tview.Application, config config.Config) Pages {
	var p Pages
	pages := []string{"feeds", "profile", "repos"}
	for _, pageName := range pages {
		p = append(p, MakePage(app, config, pageName))
	}
	return p
}

func (p Pages) Get(name string) *Page {
	for _, page := range p {
		if page.Name == name {
			return page
		}
	}
	return nil
}

func (p Pages) Prev(currentPage string) string {
	pageCount := len(p)
	currentPageIndex := p.getPageIndex(currentPage)
	prevPageIndex := (currentPageIndex - 1) % pageCount
	if prevPageIndex < 0 {
		return p[pageCount-1].Name
	}
	return p[prevPageIndex].Name
}

func (p Pages) Next(currentPage string) string {
	pageCount := len(p)
	currentPageIndex := p.getPageIndex(currentPage)
	nextPageIndex := (currentPageIndex + 1) % pageCount
	return p[nextPageIndex].Name
}

func (p Pages) getPageIndex(name string) int {
	for index, page := range p {
		if page.Name == name {
			return index
		}
	}
	return -1
}
