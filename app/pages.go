package app

import (
	"github.com/bharath-srinivas/giterm/modules"
)

type page struct {
	name   string
	widget *modules.Widget
}

func (g *GitApp) LoadPages() {
	g.loadAppPages()
	for _, page := range g.appPages {
		g.pages.AddPage(page.name, page.widget.Parent, true, false)
	}
	g.pages.SwitchToPage(g.appPages[0].name)
}

func (g *GitApp) loadAppPages() {
	pages := []string{"profile", "repos", "feeds"}
	for _, moduleName := range pages {
		g.appPages = append(g.appPages, &page{
			name:   moduleName,
			widget: modules.MakeWidget(g.app, g.config, moduleName),
		})
	}
}
