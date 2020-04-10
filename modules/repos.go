package modules

import (
	"fmt"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/rivo/tview"
	"github.com/shurcooL/githubv4"

	"github.com/bharath-srinivas/giterm/config"
	"github.com/bharath-srinivas/giterm/views"
)

// repoQuery represents a graphql query.
type repoQuery struct {
	Viewer struct {
		Repositories struct {
			TotalCount int
			Nodes      []repository
			PageInfo   struct {
				StartCursor     string
				EndCursor       string
				HasPreviousPage bool
				HasNextPage     bool
			}
		} `graphql:"repositories (first: $first, last: $last, before: $before, after:$after, privacy: $privacy, isFork: $isFork, orderBy: $orderBy)"`
	}
}

// repository represents a github repository.
type repository struct {
	Name             string
	Description      *string
	IsArchived       bool
	IsPrivate        bool
	RepositoryTopics struct {
		Nodes []struct {
			Topic struct {
				Name string
			}
		}
	} `graphql:"repositoryTopics (first: 25)"`
	PrimaryLanguage *struct {
		Name  string
		Color string
	}
	Stargazers struct {
		TotalCount int
	}
	IsMirror bool
	IsFork   bool
	Parent   struct {
		NameWithOwner string
	}
	ForkCount   int
	LicenseInfo *struct {
		Name string
	}
	PushedAt *time.Time
}

// repositories holds the github repositories info.
var repositories repoQuery

// Repos represents the github repositories of the authenticated user.
type Repos struct {
	*views.TextWidget

	first   *githubv4.Int
	last    *githubv4.Int
	privacy *githubv4.RepositoryPrivacy
	isFork  *githubv4.Boolean
	before  *githubv4.String
	after   *githubv4.String
}

// RepoWidget returns a new instance of repo widget.
func RepoWidget(app *tview.Application, config config.Config) *Repos {
	widget := views.NewTextView(app, config, true)
	r := &Repos{TextWidget: widget}
	return r
}

// Refresh refreshes the repository list.
func (r *Repos) Refresh() {
	r.Redraw(r.display)
}

// Filter filters the repository list according to the given repository type.
func (r *Repos) Filter(repoType string) {
	r.before = nil
	r.after = nil
	switch strings.ToLower(repoType) {
	case "public":
		r.isFork = nil
		privacy := githubv4.RepositoryPrivacyPublic
		r.privacy = &privacy
	case "private":
		r.isFork = nil
		privacy := githubv4.RepositoryPrivacyPrivate
		r.privacy = &privacy
	case "sources":
		r.privacy = nil
		isFork := githubv4.Boolean(false)
		r.isFork = &isFork
	case "forks":
		r.privacy = nil
		isFork := githubv4.Boolean(true)
		r.isFork = &isFork
	default:
		r.privacy = nil
		r.isFork = nil
	}
	go r.Refresh()
}

// SetPageSize sets the page size.
func (r *Repos) SetPageSize(pageSize int) {
	r.before = nil
	r.after = nil
	r.last = nil
	r.first = githubv4.NewInt(githubv4.Int(pageSize))
	go r.Refresh()
}

// First navigates to the first page of the repository list.
func (r *Repos) First() {}

// Last navigates to the last page of the repository list.
func (r *Repos) Last() {}

// Prev navigates to the previous page of the repository list.
func (r *Repos) Prev() {
	pageInfo := repositories.Viewer.Repositories.PageInfo
	r.before = githubv4.NewString(githubv4.String(pageInfo.StartCursor))
	r.after = githubv4.NewString(githubv4.String(pageInfo.EndCursor))
	if pageInfo.HasPreviousPage {
		if r.first != nil {
			r.last = r.first
			r.first = nil
		}
		r.after = nil
		go r.Refresh()
	}
}

// Next navigates to the next page of the repository list.
func (r *Repos) Next() {
	pageInfo := repositories.Viewer.Repositories.PageInfo
	r.before = githubv4.NewString(githubv4.String(pageInfo.StartCursor))
	r.after = githubv4.NewString(githubv4.String(pageInfo.EndCursor))
	if pageInfo.HasNextPage {
		if r.last != nil {
			r.first = r.last
			r.last = nil
		}
		r.before = nil
		go r.Refresh()
	}
}

// display renders the repository list according to the provided filter, pagination and page size options.
func (r *Repos) display() {
	variables := map[string]interface{}{
		"first":   r.first,
		"last":    r.last,
		"before":  r.before,
		"after":   r.after,
		"privacy": r.privacy,
		"isFork":  r.isFork,
		"orderBy": githubv4.RepositoryOrder{
			Field:     "PUSHED_AT",
			Direction: "DESC",
		},
	}

	if err := r.GqlClient.Query(r.Context, &repositories, variables); err != nil {
		_, _ = fmt.Fprintln(r.TextView, "[::b]an error occurred while retrieving your repositories")
		return
	}

	if repositories.Viewer.Repositories.TotalCount == 0 {
		_, _ = fmt.Fprintln(r.TextView, "[::b]Nothing to display")
		return
	}

	_, _, width, _ := r.GetInnerRect()
	writer := tabwriter.NewWriter(r.TextView, 0, 4, 2, '\t', 0)
	for _, repo := range repositories.Viewer.Repositories.Nodes {
		repoInfo := getRepoInfo(repo) + strings.Repeat("_", width)
		_, _ = fmt.Fprintln(writer, repoInfo)
	}
	_ = writer.Flush()
}

// getRepoInfo returns the tab indented string representation of the repository information.
func getRepoInfo(repo repository) string {
	var forked, mirrored, description, topicTags, lang, stars, forks, license, updatedAt string
	repoInfo := "\n[::b]" + repo.Name
	if repo.IsArchived {
		repoInfo += " [darkslategray:white:d] Archived [:black:-] "
	}
	if repo.IsPrivate {
		repoInfo += " [darkslategray:white:d] Private [:black:-] "
	}

	if repo.IsFork {
		forked = "[gray::d]Forked from " + repo.Parent.NameWithOwner + "\n"
	} else if repo.IsMirror {
		repoInfo += " [darkslategray:white:d] Mirror [:black:-]"
		mirrored = "[gray::d]Mirrored from " + repo.Parent.NameWithOwner + "\n"
	}

	if len(repo.RepositoryTopics.Nodes) > 0 {
		var topics []string
		for _, topic := range repo.RepositoryTopics.Nodes {
			topicName := "[white:green:d] " + topic.Topic.Name + " [:black:-]"
			topics = append(topics, topicName)
		}
		topicTags = strings.Join(topics, " ") + "\n\n"
	}

	if repo.Description != nil {
		description = "[white::d]" + *repo.Description + "\n\n"
	}
	if repo.PrimaryLanguage != nil {
		lang = fmt.Sprintf("[%s]%s %s\t", repo.PrimaryLanguage.Color, string('\u25CF'), repo.PrimaryLanguage.Name)
	}
	if repo.Stargazers.TotalCount > 0 {
		stars = fmt.Sprintf("[white]%s %d\t", string('\u2605'), repo.Stargazers.TotalCount)
	}
	if repo.ForkCount > 0 {
		forks = fmt.Sprintf("[white]%s %d\t", string('\u2442'), repo.ForkCount)
	}
	if repo.LicenseInfo != nil {
		license = fmt.Sprintf("[white]%s %s\t", string('\u2696'), repo.LicenseInfo.Name)
	}
	if repo.PushedAt != nil {
		updatedAt = "[gray::d]Updated on " + repo.PushedAt.Format("Jan 02, 2006")
	}
	repoInfo += fmt.Sprint("\n", mirrored, forked, "\n", description, topicTags, lang, stars, forks, license, updatedAt, "\n")
	return repoInfo
}
