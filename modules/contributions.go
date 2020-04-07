package modules

import (
	"fmt"
	"time"

	"github.com/rivo/tview"
	"github.com/shurcooL/githubv4"

	"github.com/bharath-srinivas/giterm/config"
	"github.com/bharath-srinivas/giterm/views"
)

// ContributionsCollections represents the github contributions collection.
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
	} `graphql:"repositoryContributions(first: 100)"`

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
		} `graphql:"contributions(first: 100)"`
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
		} `graphql:"contributions(first: 100)"`
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
		} `graphql:"contributions (first: 100)"`
		Repository struct {
			NameWithOwner string
		}
	}
}

// contributionQuery represents a graphql query
type contributionQuery struct {
	Viewer struct {
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
		SetTopLevel(1).
		SetAlign(true)
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
	c.app.QueueUpdateDraw(func() {
		c.display()
	})
}

// display renders the contribution activities of a user in a tree view.
func (c *Contributions) display() {
	c.keys = nil
	c.nodes = map[string]ContributionsCollection{}
	if err := c.getContributionData(); err != nil {
		c.createRootNode(err.Error())
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
			text := fmt.Sprintf("[::d]" + c.Username + " had no activity during this period.")
			node := tview.NewTreeNode(text).SetSelectable(true)
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
		SetCurrentNode(root)
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

	var childNodes []*tview.TreeNode
	for _, node := range c.nodes[key].RepositoryContributions.Nodes {
		createdAt := node.OccurredAt.Format("Jan 02")
		text := fmt.Sprintf(" [white]%s  [gray::d]%s", node.Repository.NameWithOwner, createdAt)
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
	var childNodes []*tview.TreeNode
	for _, node := range c.nodes[key].PullRequestContributionsByRepository {
		repo := node.Repository.NameWithOwner
		prCount := node.Contributions.TotalCount
		prText := pluralize("pull request", prCount)
		text := fmt.Sprintf(" [white]%s  [gray::d]%d %s", repo, prCount, prText)
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
	var childNodes []*tview.TreeNode
	for _, node := range c.nodes[key].PullRequestReviewContributionsByRepository {
		reviewCount := node.Contributions.TotalCount
		repo := node.Repository.NameWithOwner
		reviewText := pluralize("pull request", reviewCount)
		text := fmt.Sprintf(" [white]%s  [gray::d]%d %s", repo, reviewCount, reviewText)
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
	var childNodes []*tview.TreeNode
	for _, node := range c.nodes[key].IssueContributionsByRepository {
		repo := node.Repository.NameWithOwner
		issueCount := node.Contributions.TotalCount
		issueText := pluralize("issue", issueCount)
		text := fmt.Sprintf(" [white]%s  [gray::d]%d %s", repo, issueCount, issueText)
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
		firstOfMonth, lastOfMonth := getPeriod(i)
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

// getPeriod returns the first date and last date of a month based on the provided offset value. The offset can be a negative
//value to get the first and last date of previous month.
func getPeriod(offset time.Month) (time.Time, time.Time) {
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
