package modules

import (
	"context"
	"log"
	"os"

	"github.com/gdamore/tcell"
	"github.com/google/go-github/v29/github"
	"github.com/rivo/tview"
	"golang.org/x/oauth2"

	"github.com/bharath-srinivas/giterm/config"
)

type Page struct {
	Name            string
	Parent          tview.Primitive
	ChildComponents []tview.Primitive
}

type GitApp struct {
	App         *tview.Application
	Context     context.Context
	Client      *github.Client
	GitAppPages []*Page

	pages  *tview.Pages
	config config.Config
}

func New(app *tview.Application) *GitApp {
	config.Init()
	cfg := config.GetConfig()

	ctx := context.Background()
	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: cfg.Token},
	)
	httpClient := oauth2.NewClient(ctx, tokenSource)
	client := github.NewClient(httpClient)

	gitApp := &GitApp{
		App:     app,
		Client:  client,
		Context: ctx,
		pages:   tview.NewPages(),
		config:  cfg,
	}
	gitApp.GitAppPages = append(gitApp.GitAppPages, gitApp.ProfilePage(), gitApp.RepoPage())
	for _, gitAppPage := range gitApp.GitAppPages {
		gitApp.pages.AddPage(gitAppPage.Name, gitAppPage.Parent, true, false)
	}
	gitApp.pages.SwitchToPage(gitApp.GitAppPages[0].Name)
	gitApp.App.SetInputCapture(gitApp.inputHandler)
	gitApp.App.SetRoot(gitApp.pages, true).SetFocus(gitApp.pages)
	return gitApp
}

func (g *GitApp) Start() {
	if err := g.App.Run(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func (g *GitApp) GetPage(name string) *Page {
	for _, page := range g.GitAppPages {
		if page.Name == name {
			return page
		}
	}
	return nil
}

func (g *GitApp) GetPageIndex(name string) int {
	for i := 0; i < len(g.GitAppPages); i++ {
		if g.GitAppPages[i].Name == name {
			return i
		}
	}
	return -1
}

func (g *GitApp) inputHandler(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyCtrlQ:
		g.App.Stop()
		return event
	case tcell.KeyCtrlN:
		pageCount := g.pages.GetPageCount()
		currentPage, _ := g.pages.GetFrontPage()
		currentPageIndex := g.GetPageIndex(currentPage)
		nextPageIndex := (currentPageIndex + 1) % pageCount
		if currentPageIndex > -1 {
			g.pages.SwitchToPage(g.GitAppPages[nextPageIndex].Name)
		}
		return event
	case tcell.KeyCtrlP:
		pageCount := g.pages.GetPageCount()
		currentPage, _ := g.pages.GetFrontPage()
		currentPageIndex := g.GetPageIndex(currentPage)
		prevPageIndex := (currentPageIndex - 1) % pageCount
		if prevPageIndex < 0 {
			g.pages.SwitchToPage(g.GitAppPages[pageCount-1].Name)
			return event
		}
		g.pages.SwitchToPage(g.GitAppPages[prevPageIndex].Name)
		return event
	}
	return event
}
