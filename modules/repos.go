package modules

import (
	"fmt"

	"github.com/rivo/tview"

	"github.com/bharath-srinivas/giterm/config"
	"github.com/bharath-srinivas/giterm/views"
)

type Repos struct {
	*views.Base
	*tview.Table
}

func RepoWidget(app *tview.Application, config config.Config) *Repos {
	widget := tview.NewTable().
		SetBorders(true)
	widget.SetTitle(string('\U0001F4D5') + " [green::b]Repositories").
		SetBorder(true)

	r := &Repos{
		views.NewBase(app, config),
		widget,
	}
	r.display()
	return r
}

func (r *Repos) display() {
	repositories, _, err := r.Client.Repositories.List(r.Context, "", nil)
	if err != nil {
		r.Table.SetCellSimple(1, 0, "[::b]an error occurred while retrieving repositories")
		return
	}

	r.setTableHeaders()
	for row, repo := range repositories {
		r.Table.SetCellSimple(row+1, 0, "[white::b]"+repo.GetName()).
			SetCellSimple(row+1, 1, "[white::b]"+repo.GetDescription()).
			SetCellSimple(row+1, 2, "[white::b]"+repo.GetHomepage()).
			SetCellSimple(row+1, 3, "[white::b]"+repo.GetGitURL()).
			SetCellSimple(row+1, 4, fmt.Sprintf("[white::b]%d", repo.GetStargazersCount())).
			SetCellSimple(row+1, 5, fmt.Sprintf("[white::b]%d", repo.GetOpenIssuesCount())).
			SetCellSimple(row+1, 6, fmt.Sprintf("[white::b]%d", repo.GetForksCount()))
	}
}

func (r *Repos) setTableHeaders() {
	headers := []string{
		"[gray::b]Name",
		"[gray::b]Description",
		"[gray::b]Homepage",
		"[gray::b]URL",
		"[gray::b]Stargazers",
		"[gray::b]Open Issues Count",
		"[gray::b]Forks"}

	for i, field := range headers {
		r.Table.SetCellSimple(0, i, field)
	}
}
