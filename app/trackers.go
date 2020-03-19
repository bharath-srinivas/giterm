package app

import (
	"github.com/gdamore/tcell"
)

func (g *GitApp) inputHandler(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyCtrlQ:
		g.app.Stop()
		return event
	case tcell.KeyTab:
		currentPage, _ := g.pages.GetFrontPage()
		page := g.appPages.Get(currentPage)
		g.app.SetFocus(page.NextWidget())
		return event
	case tcell.KeyBacktab:
		currentPage, _ := g.pages.GetFrontPage()
		page := g.appPages.Get(currentPage)
		g.app.SetFocus(page.PrevWidget())
		return event
	case tcell.KeyCtrlN:
		currentPage, _ := g.pages.GetFrontPage()
		g.pages.SwitchToPage(g.appPages.Next(currentPage))
		return event
	case tcell.KeyCtrlP:
		currentPage, _ := g.pages.GetFrontPage()
		g.pages.SwitchToPage(g.appPages.Prev(currentPage))
		return event
	}
	return event
}
