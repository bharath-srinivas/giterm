package app

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

func (g *GitApp) inputHandler(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyCtrlQ:
		g.app.Stop()
		return event
	case tcell.KeyTab:
		currentPage, _ := g.pages.GetFrontPage()
		page := g.getPage(currentPage)
		if page != nil {
			if len(page.widget.Children) > 0 {
				g.app.SetFocus(nextWidget(page.widget.Children))
			} else {
				g.app.SetFocus(page.widget.Parent)
			}
		}
		return event
	case tcell.KeyCtrlN:
		pageCount := g.pages.GetPageCount()
		currentPage, _ := g.pages.GetFrontPage()
		currentPageIndex := g.getPageIndex(currentPage)
		nextPageIndex := (currentPageIndex + 1) % pageCount
		if currentPageIndex > -1 {
			g.pages.SwitchToPage(g.appPages[nextPageIndex].name)
		}
		return event
	case tcell.KeyCtrlP:
		pageCount := g.pages.GetPageCount()
		currentPage, _ := g.pages.GetFrontPage()
		currentPageIndex := g.getPageIndex(currentPage)
		prevPageIndex := (currentPageIndex - 1) % pageCount
		if prevPageIndex < 0 {
			g.pages.SwitchToPage(g.appPages[pageCount-1].name)
			return event
		}
		g.pages.SwitchToPage(g.appPages[prevPageIndex].name)
		return event
	}
	return event
}

func (g *GitApp) getPage(name string) *page {
	for _, page := range g.appPages {
		if page.name == name {
			return page
		}
	}
	return nil
}

func (g *GitApp) getPageIndex(name string) int {
	for i := 0; i < len(g.appPages); i++ {
		if g.appPages[i].name == name {
			return i
		}
	}
	return -1
}

func nextWidget(widgets []tview.Primitive) tview.Primitive {
	widgetLen := len(widgets)
	for i := 0; i < widgetLen; i++ {
		if widgets[i].GetFocusable().HasFocus() {
			return widgets[(i+1)%widgetLen]
		}
	}
	return widgets[0]
}
