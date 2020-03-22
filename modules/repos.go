package modules

import (
	"fmt"
	"strings"

	"github.com/google/go-github/v29/github"
	"github.com/rivo/tview"

	"github.com/bharath-srinivas/giterm/config"
	"github.com/bharath-srinivas/giterm/views"
)

// Repos represents the github repositories.
type Repos struct {
	app *tview.Application
	*views.Base
	*tview.Table
	*github.RepositoryListOptions
	*github.Response
}

// RepoWidget returns a new instance of repo widget.
// TODO: convert this widget as text widget
func RepoWidget(app *tview.Application, config config.Config) *Repos {
	widget := tview.NewTable().
		SetBorders(true)
	widget.SetBorder(true)
	r := &Repos{
		app:                   app,
		Base:                  views.NewBase(app, config),
		Table:                 widget,
		RepositoryListOptions: &github.RepositoryListOptions{},
		Response:              &github.Response{},
	}
	return r
}

// Refresh refreshes the repository list.
func (r *Repos) Refresh() {
	r.app.QueueUpdateDraw(func() {
		r.Clear()
		r.display(r.RepositoryListOptions)
		r.ScrollToBeginning()
	})
}

// Filter filters the repository list according to the given repository type.
func (r *Repos) Filter(repoType string) {
	r.RepositoryListOptions.Type = strings.ToLower(repoType)
	go r.Refresh()
}

// SetPageSize sets the page size.
func (r *Repos) SetPageSize(pageSize int) {
	r.Page = 1
	r.PerPage = pageSize
	go r.Refresh()
}

// First navigates to the first page of the repository list.
func (r *Repos) First() {
	if r.Response != nil {
		r.Page = r.FirstPage
		go r.Refresh()
	}
}

// Last navigates to the last page of the repository list.
func (r *Repos) Last() {
	if r.Response != nil {
		r.Page = r.LastPage
		go r.Refresh()
	}
}

// Prev navigates to the previous page of the repository list.
func (r *Repos) Prev() {
	if r.Response != nil {
		r.Page = r.PrevPage
		go r.Refresh()
	}
}

// Next navigates to the next page of the repository list.
func (r *Repos) Next() {
	if r.Response != nil {
		r.Page = r.NextPage
		go r.Refresh()
	}
}

// display renders the repository list according to the provided filter, pagination and page size options.
func (r *Repos) display(options *github.RepositoryListOptions) {
	repositories, res, err := r.Client.Repositories.List(r.Context, "", options)
	if err != nil {
		r.Table.SetCellSimple(1, 0, "[::b]an error occurred while retrieving repositories")
		return
	}
	r.Response = res
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

// setTableHeaders sets the table headers.
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
