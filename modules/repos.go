package modules

import (
	"encoding/json"
	"fmt"

	"github.com/rivo/tview"
)

type repo struct {
	Name            string `json:"name,omitempty"`
	Description     string `json:"description,omitempty"`
	Homepage        string `json:"homepage,omitempty"`
	GitURL          string `json:"git_url,omitempty"`
	ForksCount      int    `json:"forks_count,omitempty"`
	OpenIssuesCount int    `json:"open_issues_count,omitempty"`
	StargazersCount int    `json:"stargazers_count,omitempty"`
}

func (g *GitApp) RepoPage() *Page {
	widget := tview.NewTable()
	widget.SetBorder(true)
	widget.SetBorders(true)
	widget.SetTitle(string('\U0001F4D5') + " [green::b]Repositories")

	repositories, _, err := g.Client.Repositories.List(g.Context, "", nil)
	if err != nil {
		widget.SetCellSimple(1, 0, err.Error())
		return &Page{}
	}

	var m []map[string]interface{}
	repoJson, _ := json.Marshal(repositories)
	_ = json.Unmarshal(repoJson, &m)

	var repos []repo
	_ = json.Unmarshal(repoJson, &repos)

	widget = setTableHeaders(widget)
	for row, repo := range repos {
		widget.SetCellSimple(row+1, 0, "[white::b]"+repo.Name)
		widget.SetCellSimple(row+1, 1, "[white::b]"+repo.Description)
		widget.SetCellSimple(row+1, 2, "[white::b]"+repo.Homepage)
		widget.SetCellSimple(row+1, 3, "[white::b]"+repo.GitURL)
		widget.SetCellSimple(row+1, 4, fmt.Sprintf("[white::b]%d", repo.StargazersCount))
		widget.SetCellSimple(row+1, 5, fmt.Sprintf("[white::b]%d", repo.OpenIssuesCount))
		widget.SetCellSimple(row+1, 6, fmt.Sprintf("[white::b]%d", repo.ForksCount))
	}
	return &Page{
		Name:            "Repositories",
		Parent:          widget,
		ChildComponents: nil,
	}
}

func setTableHeaders(widget *tview.Table) *tview.Table {
	headers := []string{
		"[gray::b]Name",
		"[gray::b]Description",
		"[gray::b]Homepage",
		"[gray::b]URL",
		"[gray::b]Stargazers",
		"[gray::b]Open Issues Count",
		"[gray::b]Forks"}

	for i, field := range headers {
		widget.SetCellSimple(0, i, field)
	}
	return widget
}
