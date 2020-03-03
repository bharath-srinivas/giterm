package modules

import (
	"encoding/json"
	"fmt"

	"github.com/google/go-github/v29/github"
	"github.com/rivo/tview"
)

type repo struct {
	Name             string `json:"name,omitempty"`
	FullName         string `json:"full_name,omitempty"`
	Description      string `json:"description,omitempty"`
	Homepage         string `json:"homepage,omitempty"`
	DefaultBranch    string `json:"default_branch,omitempty"`
	MasterBranch     string `json:"master_branch,omitempty"`
	HTMLURL          string `json:"html_url,omitempty"`
	CloneURL         string `json:"clone_url,omitempty"`
	GitURL           string `json:"git_url,omitempty"`
	Language         string `json:"language,omitempty"`
	Fork             bool   `json:"fork,omitempty"`
	ForksCount       int    `json:"forks_count,omitempty"`
	OpenIssuesCount  int    `json:"open_issues_count,omitempty"`
	StargazersCount  int    `json:"stargazers_count,omitempty"`
	SubscribersCount int    `json:"subscribers_count,omitempty"`
	WatchersCount    int    `json:"watchers_count,omitempty"`
	Archived         bool   `json:"archived,omitempty"`
	Disabled         bool   `json:"disabled,omitempty"`
}

func (g *GitApp) RepoWidget() tview.Primitive {
	widget := tview.NewTable()
	widget.SetBorder(true)
	widget.SetBorders(true)
	widget.SetTitle(string('\U0001F4D5') + " [green::b]Repositories")

	repositories, err := g.GetRepositories()
	if err != nil {
		widget.SetCellSimple(0, 0, err.Error())
		return widget
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
		//widget.SetCellSimple(row+1, 3, repo.DefaultBranch)
		//widget.SetCellSimple(row+1, 4, repo.MasterBranch)
		//widget.SetCellSimple(row+1, 5, repo.CloneURL)
		widget.SetCellSimple(row+1, 3, "[white::b]"+repo.GitURL)
		//widget.SetCellSimple(row+1, 7, repo.HTMLURL)
		//widget.SetCellSimple(row+1, 8, repo.Language)
		widget.SetCellSimple(row+1, 4, fmt.Sprintf("[white::b]%d", repo.StargazersCount))
		widget.SetCellSimple(row+1, 5, fmt.Sprintf("[white::b]%d", repo.OpenIssuesCount))
		widget.SetCellSimple(row+1, 6, fmt.Sprintf("[white::b]%d", repo.ForksCount))
		//widget.SetCellSimple(row+1, 12, fmt.Sprintf("%d", repo.SubscribersCount))
		widget.SetCellSimple(row+1, 7, fmt.Sprintf("[white::b]%d", repo.WatchersCount))
	}
	return widget
}

func (g *GitApp) GetRepositories() ([]*github.Repository, error) {
	repos, _, err := g.Client.Repositories.List(g.Context, "", nil)
	return repos, err
}

func setTableHeaders(widget *tview.Table) *tview.Table {
	headers := []string{
		"[gray::b]Name",
		"[gray::b]Description",
		"[gray::b]Homepage",
		//"[gray::b]DefaultBranch",
		//"[gray::b]MasterBranch",
		//"[gray::b]CloneURL",
		//"[gray::b]HTMLURL",
		"[gray::b]URL",
		//"[gray::b]Language",
		"[gray::b]Stargazers",
		"[gray::b]Open Issues Count",
		"[gray::b]Forks",
		//"[gray::b]SubscribersCount",
		"[gray::b]Watchers"}

	for i, field := range headers {
		widget.SetCellSimple(0, i, field)
	}
	return widget
}
