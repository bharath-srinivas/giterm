package modules

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"github.com/shurcooL/githubv4"

	"github.com/bharath-srinivas/giterm/config"
	"github.com/bharath-srinivas/giterm/views"
)

// ContributionsCollection represents the github contributions collection.
type ContributionsCollection struct {
	HasAnyContributions                                bool
	HasActivityInThePast                               bool
	TotalCommitContributions                           int
	TotalRepositoriesWithContributedCommits            int
	TotalRepositoryContributions                       int
	TotalPullRequestContributions                      int
	TotalRepositoriesWithContributedPullRequests       int
	TotalPullRequestReviewContributions                int
	TotalRepositoriesWithContributedPullRequestReviews int
	TotalIssueContributions                            int
	TotalRepositoriesWithContributedIssues             int

	JoinedGitHubContribution struct {
		OccurredAt *time.Time
	}

	FirstIssueContribution struct {
		CreatedIssueContribution struct {
			Issue struct {
				Title string
				State string
			}
			OccurredAt *time.Time
		} `graphql:"... on CreatedIssueContribution"`
	}

	FirstPullRequestContribution struct {
		CreatedPullRequestContribution struct {
			PullRequest struct {
				Title string
				State string
			}
			OccurredAt *time.Time
		} `graphql:"... on CreatedPullRequestContribution"`
	}

	FirstRepositoryContribution struct {
		CreatedRepositoryContribution struct {
			Repository struct {
				Name            string
				PrimaryLanguage struct {
					Name  string
					Color string
				}
			}
			OccurredAt *time.Time
		} `graphql:"... on CreatedRepositoryContribution"`
	}

	CommitContributionsByRepository []struct {
		Repository struct {
			NameWithOwner string
		}
		Contributions struct {
			TotalCount int
		}
	}

	RepositoryContributions struct {
		TotalCount int
		Nodes      []struct {
			OccurredAt *time.Time
			Repository struct {
				NameWithOwner   string
				PrimaryLanguage struct {
					Name  string
					Color string
				}
			}
		}
	} `graphql:"repositoryContributions(first: 25)"`

	PullRequestContributionsByRepository []struct {
		Contributions struct {
			TotalCount int
			Nodes      []struct {
				OccurredAt  *time.Time
				PullRequest struct {
					Title    string
					State    string
					Comments struct {
						TotalCount int
					}
				}
			}
		} `graphql:"contributions(first: 25)"`
		Repository struct {
			NameWithOwner string
		}
	}

	PullRequestReviewContributionsByRepository []struct {
		Contributions struct {
			TotalCount int
			Nodes      []struct {
				OccurredAt  *time.Time
				PullRequest struct {
					Title string
				}
				PullRequestReview struct {
					State string
				}
			}
		} `graphql:"contributions(first: 25)"`
		Repository struct {
			NameWithOwner string
		}
	}

	IssueContributionsByRepository []struct {
		Contributions struct {
			Nodes []struct {
				Issue struct {
					Title    string
					State    string
					Comments struct {
						TotalCount int
					}
				}
				OccurredAt *time.Time
			}
			TotalCount int
		} `graphql:"contributions (first: 25)"`
		Repository struct {
			NameWithOwner string
		}
	}
}

// contributionQuery represents a graphql query
type contributionQuery struct {
	Viewer struct {
		Login                   string
		ContributionsCollection `graphql:"contributionsCollection(from: $from, to: $to)"`
	}
}

// contributions holds the github contributions collection of a user.
var contributions contributionQuery

// Contributions represents the github contributions of a user.
type Contributions struct {
	*views.Base
	*tview.TreeView

	app   *tview.Application
	nodes map[string]ContributionsCollection
	keys  []string
}

// ContributionsWidget returns a new instance of contribution widget.
func ContributionsWidget(app *tview.Application, config config.Config) *Contributions {
	widget := tview.NewTreeView().
		SetTopLevel(1)
	widget.SetTitle(string('\U0001F4C8') + " [green::b]Contribution activity").
		SetBorder(true)
	c := &Contributions{
		app:      app,
		Base:     views.NewBase(app, config),
		TreeView: widget,
	}
	c.SetSelectedFunc(func(node *tview.TreeNode) {
		node.SetExpanded(!node.IsExpanded())
	})
	c.createRootNode("[::b]Contributions in the last 90 days")
	go c.Refresh()
	return c
}

// Refresh refreshes the contributions widget.
func (c *Contributions) Refresh() {
	c.app.QueueUpdateDraw(c.display)
}

// display renders the contribution activities of a user in a tree view.
func (c *Contributions) display() {
	c.keys = nil
	c.nodes = map[string]ContributionsCollection{}
	if err := c.getContributionData(); err != nil {
		c.createRootNode("[::b]an error occurred while retrieving your contributions")
		return
	}

	if len(c.nodes) == 0 {
		c.createRootNode("[::b]Nothing to display.")
		return
	}

	root := c.TreeView.GetRoot().ClearChildren()
	for _, key := range c.keys {
		childNode := tview.NewTreeNode("[::b]" + key).
			SetSelectable(true)

		if !c.nodes[key].HasAnyContributions {
			text := fmt.Sprintf(" [::d]" + contributions.Viewer.Login + " had no activity during this period.")
			node := tview.NewTreeNode(text).SetSelectable(false)
			childNode.AddChild(node)
			root.AddChild(childNode)
			continue
		}

		if commitNode := c.getCommitNode(key); commitNode != nil {
			childNode.AddChild(commitNode)
		}
		if createNode := c.getRepoNode(key); createNode != nil {
			childNode.AddChild(createNode)
		}
		if pullRequestNode := c.getPullRequestNode(key); pullRequestNode != nil {
			childNode.AddChild(pullRequestNode)
		}
		if pullRequestReviewNode := c.getPullRequestReviewNode(key); pullRequestReviewNode != nil {
			childNode.AddChild(pullRequestReviewNode)
		}
		if issueNode := c.getIssueNode(key); issueNode != nil {
			childNode.AddChild(issueNode)
		}
		root.AddChild(childNode)
	}
}

// createRootNode creates the root tree node for the tree view with the provided text.
func (c *Contributions) createRootNode(text string) {
	root := tview.NewTreeNode(text)
	c.TreeView.
		SetRoot(root).
		SetCurrentNode(root).
		SetGraphics(false)
}

// getCommitNode returns the tree node with commits data for a given key, which is the month.
func (c *Contributions) getCommitNode(key string) *tview.TreeNode {
	totalCommits := c.nodes[key].TotalCommitContributions
	if totalCommits == 0 {
		return nil
	}

	totalRepos := c.nodes[key].TotalRepositoriesWithContributedCommits
	nodes := c.nodes[key].CommitContributionsByRepository
	var childNodes []*tview.TreeNode
	for _, node := range nodes {
		repo := node.Repository.NameWithOwner
		commitCount := node.Contributions.TotalCount
		commitText := pluralize("commit", commitCount)
		text := fmt.Sprintf(" [white]%s  [gray::d]%d %s", repo, commitCount, commitText)
		child := tview.NewTreeNode(text).SetSelectable(false)
		childNodes = append(childNodes, child)
	}

	repoCount := len(nodes)
	if totalRepos > repoCount {
		repoText := pluralize("repository", repoCount)
		text := fmt.Sprintf("[gray::d]%d %s not shown", totalRepos-repoCount, repoText)
		child := tview.NewTreeNode(text).SetSelectable(false)
		childNodes = append(childNodes, child)
	}

	commitText := pluralize("commit", totalCommits)
	repoText := pluralize("repository", totalRepos)
	text := fmt.Sprintf(" [::b]Created %d %s in %d %s", totalCommits, commitText, totalRepos, repoText)
	node := tview.NewTreeNode(text).SetSelectable(true)
	node.SetChildren(childNodes)
	return node
}

// getRepoNode returns the tree node with created repositories data for a given key, which is the month.
func (c *Contributions) getRepoNode(key string) *tview.TreeNode {
	totalRepoCount := c.nodes[key].TotalRepositoryContributions
	if totalRepoCount == 0 {
		return nil
	}

	nodes := c.nodes[key].RepositoryContributions.Nodes
	var childNodes []*tview.TreeNode
	var emptyLang struct {
		Name  string
		Color string
	}
	for _, node := range nodes {
		repo := node.Repository.NameWithOwner

		var lang string
		if node.Repository.PrimaryLanguage != (emptyLang) {
			lang = fmt.Sprintf("[%s]● %s", node.Repository.PrimaryLanguage.Color, node.Repository.PrimaryLanguage.Name)
		}
		createdAt := node.OccurredAt.Format("Jan 02")
		text := fmt.Sprintf(" [white]%s \t%s  [gray::d]%s", repo, lang, createdAt)
		child := tview.NewTreeNode(text).SetSelectable(false)
		childNodes = append(childNodes, child)
	}

	repoCount := len(nodes)
	if totalRepoCount > repoCount {
		repoText := pluralize("repository", repoCount)
		text := fmt.Sprintf("[gray::d]%d %s not shown", totalRepoCount-repoCount, repoText)
		child := tview.NewTreeNode(text).SetSelectable(false)
		childNodes = append(childNodes, child)
	}

	repoText := pluralize("repository", totalRepoCount)
	text := fmt.Sprintf(" [::b]Created %d %s", totalRepoCount, repoText)
	node := tview.NewTreeNode(text).SetSelectable(true)
	node.SetChildren(childNodes)
	return node
}

// getPullRequestNode returns the tree node with pull requests data for a given key, which is the month.
func (c *Contributions) getPullRequestNode(key string) *tview.TreeNode {
	totalPullRequests := c.nodes[key].TotalPullRequestContributions
	if totalPullRequests == 0 {
		return nil
	}

	totalRepoCount := c.nodes[key].TotalRepositoriesWithContributedPullRequests
	nodes := c.nodes[key].PullRequestContributionsByRepository
	var childNodes []*tview.TreeNode
	for _, node := range nodes {
		pullRequestCount := node.Contributions.TotalCount
		repo := node.Repository.NameWithOwner
		text := " [white::]" + repo
		child := tview.NewTreeNode(text).
			SetSelectable(true).
			SetExpanded(false).
			SetColor(tcell.ColorGray)
		var openCount, mergedCount, closedCount int
		for _, node := range node.Contributions.Nodes {
			title := tview.Escape(node.PullRequest.Title)
			createdAt := node.OccurredAt.Format("Jan 02")

			var subText string
			switch node.PullRequest.State {
			case "OPEN":
				openCount += 1
				subText = fmt.Sprintf(" [green::d]%s  [gray::d]%s", title, createdAt)
			case "MERGED":
				mergedCount += 1
				subText = fmt.Sprintf(" [rebeccapurple::d]%s  [gray::d]%s", title, createdAt)
			case "CLOSED":
				closedCount += 1
				subText = fmt.Sprintf(" [indianred::d]%s  [gray::d]%s", title, createdAt)
			}

			subChild := tview.NewTreeNode(subText).SetSelectable(false)
			child.AddChild(subChild)
		}

		if openCount > 0 {
			text += fmt.Sprintf("  [white:green:] %d [:black:] [gray::d]%s", openCount, "open")
		}
		if mergedCount > 0 {
			text += fmt.Sprintf("  [white:rebeccapurple:] %d [:black:] [gray::d]%s", mergedCount, "merged")
		}
		if closedCount > 0 {
			text += fmt.Sprintf("  [white:indianred:] %d [:black:] [gray::d]%s", closedCount, "closed")
		}

		if pullRequestCount > 25 {
			subChild := tview.NewTreeNode("[gray::d]Some pull requests not shown.").
				SetSelectable(false)
			child.AddChild(subChild)
		}
		child.SetText(text)
		childNodes = append(childNodes, child)
	}

	repoCount := len(nodes)
	if totalRepoCount > repoCount {
		repoText := pluralize("repository", repoCount)
		text := fmt.Sprintf("[gray::d]%d %s not shown", totalRepoCount-repoCount, repoText)
		child := tview.NewTreeNode(text).SetSelectable(false)
		childNodes = append(childNodes, child)
	}

	prText := pluralize("pull request", totalPullRequests)
	repoText := pluralize("repository", totalRepoCount)
	text := fmt.Sprintf(" [::b]Opened %d %s in %d %s", totalPullRequests, prText, totalRepoCount, repoText)
	node := tview.NewTreeNode(text).SetSelectable(true)
	node.SetChildren(childNodes)
	return node
}

// getPullRequestReviewNode returns the tree node with pull request reviews data for a given key, which is the month.
func (c *Contributions) getPullRequestReviewNode(key string) *tview.TreeNode {
	totalPullRequestReviews := c.nodes[key].TotalPullRequestReviewContributions
	if totalPullRequestReviews == 0 {
		return nil
	}

	totalRepoCount := c.nodes[key].TotalRepositoriesWithContributedPullRequestReviews
	nodes := c.nodes[key].PullRequestReviewContributionsByRepository
	var childNodes []*tview.TreeNode
	for _, node := range nodes {
		reviewCount := node.Contributions.TotalCount
		repo := node.Repository.NameWithOwner
		reviewText := pluralize("pull request", reviewCount)
		text := fmt.Sprintf(" [white]%s  [gray::d]%d %s", repo, reviewCount, reviewText)
		child := tview.NewTreeNode(text).
			SetSelectable(true).
			SetExpanded(false).
			SetColor(tcell.ColorGray)

		var pending, commented, approved, changeRequested, dismissed int
		for _, node := range node.Contributions.Nodes {
			title := tview.Escape(node.PullRequest.Title)
			createdAt := node.OccurredAt.Format("Jan 02")

			var subText string
			switch node.PullRequestReview.State {
			case "PENDING":
				pending += 1
				subText = fmt.Sprintf(" [gray::d]%s  [gray::d]%s", title, createdAt)
			case "COMMENTED":
				commented += 1
				subText = fmt.Sprintf(" [darkmagenta::d]%s  [gray::d]%s", title, createdAt)
			case "APPROVED":
				approved += 1
				subText = fmt.Sprintf(" [green::d]%s  [gray::d]%s", title, createdAt)
			case "CHANGES_REQUESTED":
				changeRequested += 1
				subText = fmt.Sprintf(" [yellow::d]%s  [gray::d]%s", title, createdAt)
			case "DISMISSED":
				dismissed += 1
				subText = fmt.Sprintf(" [indianred::d]%s  [gray::d]%s", title, createdAt)
			}

			subChild := tview.NewTreeNode(subText).SetSelectable(false)
			child.AddChild(subChild)
		}

		if reviewCount > 25 {
			subChild := tview.NewTreeNode("[gray::d]Some pull request reviews not shown.").
				SetSelectable(false)
			child.AddChild(subChild)
		}
		childNodes = append(childNodes, child)
	}

	repoCount := len(nodes)
	if totalRepoCount > repoCount {
		repoText := pluralize("repository", repoCount)
		text := fmt.Sprintf("[gray::d]%d %s not shown", totalRepoCount-repoCount, repoText)
		child := tview.NewTreeNode(text).SetSelectable(false)
		childNodes = append(childNodes, child)
	}

	reviewText := pluralize("pull request", totalPullRequestReviews)
	repoText := pluralize("repository", totalRepoCount)
	text := fmt.Sprintf(" [::b]Reviewed %d %s in %d %s", totalPullRequestReviews, reviewText, totalRepoCount, repoText)
	node := tview.NewTreeNode(text).SetSelectable(true)
	node.SetChildren(childNodes)
	return node
}

// getIssueNode returns the tree node with issues data for a given key, which is the month.
func (c *Contributions) getIssueNode(key string) *tview.TreeNode {
	totalIssues := c.nodes[key].TotalIssueContributions
	if totalIssues == 0 {
		return nil
	}

	totalRepoCount := c.nodes[key].TotalRepositoriesWithContributedIssues
	nodes := c.nodes[key].IssueContributionsByRepository
	var childNodes []*tview.TreeNode
	for _, node := range nodes {
		issueCount := node.Contributions.TotalCount
		text := " [white]" + node.Repository.NameWithOwner
		child := tview.NewTreeNode(text).
			SetSelectable(true).
			SetExpanded(false).
			SetColor(tcell.ColorGray)

		var openCount, closedCount int
		for _, node := range node.Contributions.Nodes {
			title := tview.Escape(node.Issue.Title)
			createdAt := node.OccurredAt.Format("Jan 02")

			var subText string
			switch node.Issue.State {
			case "OPEN":
				openCount += 1
				subText = fmt.Sprintf(" [green::d]%s  [gray::d]%s", title, createdAt)
			case "CLOSED":
				closedCount += 1
				subText = fmt.Sprintf(" [indianred::d]%s  [gray::d]%s", title, createdAt)
			}

			subChild := tview.NewTreeNode(subText).SetSelectable(false)
			child.AddChild(subChild)
		}

		if openCount > 0 {
			text += fmt.Sprintf("  [white:green:] %d [:black:] [gray::d]%s", openCount, "open")
		}
		if closedCount > 0 {
			text += fmt.Sprintf("  [white:indianred:] %d [:black:] [gray::d]%s", closedCount, "closed")
		}

		if issueCount > 25 {
			subChild := tview.NewTreeNode("[gray::d]Some issues not shown.").
				SetSelectable(false)
			child.AddChild(subChild)
		}
		child.SetText(text)
		childNodes = append(childNodes, child)
	}

	repoCount := len(nodes)
	if totalRepoCount > repoCount {
		repoText := pluralize("repository", repoCount)
		text := fmt.Sprintf("[gray::d]%d %s not shown", totalRepoCount-repoCount, repoText)
		child := tview.NewTreeNode(text).SetSelectable(false)
		childNodes = append(childNodes, child)
	}

	issueText := pluralize("issue", totalIssues)
	repoText := pluralize("repository", totalRepoCount)
	text := fmt.Sprintf(" [::b]Opened %d %s in %d %s", totalIssues, issueText, totalRepoCount, repoText)
	node := tview.NewTreeNode(text).SetSelectable(true)
	node.SetChildren(childNodes)
	return node
}

// getContributionData retrieves the contribution data of a user for the past 3 months.
func (c *Contributions) getContributionData() error {
	for i := time.Month(0); i > -3; i-- {
		firstOfMonth, lastOfMonth := getTimePeriod(i)
		variables := map[string]interface{}{
			"from": githubv4.DateTime{Time: firstOfMonth},
			"to":   githubv4.DateTime{Time: lastOfMonth},
		}

		if err := c.GqlClient.Query(c.Context, &contributions, variables); err != nil {
			return err
		}

		key := firstOfMonth.Format("January 2006")
		c.nodes[key] = contributions.Viewer.ContributionsCollection
		c.keys = append(c.keys, key)

		if !contributions.Viewer.HasActivityInThePast {
			break
		}
	}
	return nil
}

// getTimePeriod returns the first date and last date of a month based on the provided offset value. The offset can be a negative
//value to get the first and last date of previous month.
func getTimePeriod(offset time.Month) (time.Time, time.Time) {
	now := time.Now()
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()
	firstOfMonth := time.Date(currentYear, currentMonth+offset, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1).Add(86399 * time.Second)
	return firstOfMonth, lastOfMonth
}

// pluralize is a helper function which will return pluralized text if the count is greater than 1. Otherwise returns
//the text as is.
func pluralize(text string, count int) string {
	if count > 1 {
		switch text {
		case "repository":
			return "repositories"
		default:
			return text + "s"
		}
	}
	return text
}
