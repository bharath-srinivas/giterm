package modules

import (
	"errors"
	"fmt"
	"sync"

	"github.com/google/go-github/v29/github"
	"github.com/rivo/tview"

	"github.com/bharath-srinivas/giterm/config"
	"github.com/bharath-srinivas/giterm/views"
)

// Contributions represents the github contributions of a user.
type Contributions struct {
	*views.Base
	*tview.TreeView

	app   *tview.Application
	wg    sync.WaitGroup
	nodes map[string]contributionActivities
}

type stats struct {
	statMap map[string]int
	mu      sync.Mutex
}

type creations struct {
	creationMap map[string]string
	mu          sync.Mutex
}

// contributionActivities represents various github contribution activities.
type contributionActivities struct {
	commits            *stats
	creates            *creations
	pullRequests       *stats
	pullRequestReviews *stats
	issues             *stats
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
		wg:       sync.WaitGroup{},
		nodes:    map[string]contributionActivities{},
		Base:     views.NewBase(app, config),
		TreeView: widget,
	}
	c.SetSelectedFunc(func(node *tview.TreeNode) {
		node.SetExpanded(!node.IsExpanded())
	})
	c.display()
	return c
}

// display renders the contribution activities of a user in a tree view.
func (c *Contributions) display() {
	if err := c.parseContributions(); err != nil {
		c.createRootNode(err.Error())
		return
	}

	if c.nodes == nil {
		c.createRootNode("[::b]" + c.Username + " had no activity during this period.")
		return
	}

	c.createRootNode("[::b]Contributions in the last 90 days")
	root := c.TreeView.GetRoot()
	for node := range c.nodes {
		childNode := tview.NewTreeNode("[::b]" + node).
			SetSelectable(true)
		if commitNode := c.getCommitNode(node); commitNode != nil {
			childNode.AddChild(commitNode)
		}
		if createNode := c.getCreateNode(node); createNode != nil {
			childNode.AddChild(createNode)
		}
		if pullRequestNode := c.getPullRequestNode(node); pullRequestNode != nil {
			childNode.AddChild(pullRequestNode)
		}
		if pullRequestReviewNode := c.getPullRequestReviewNode(node); pullRequestReviewNode != nil {
			childNode.AddChild(pullRequestReviewNode)
		}
		if issueNode := c.getIssueNode(node); issueNode != nil {
			childNode.AddChild(issueNode)
		}
		root.AddChild(childNode)
	}
}

// parseContributions parses the contribution data based on the month of contribution and populates the nodes.
func (c *Contributions) parseContributions() error {
	contributions, err := c.getContributionData()
	if err != nil {
		return errors.New("[::b]an error occurred while retrieving contribution data")
	}

	var contribActivities contributionActivities
	for _, activity := range contributions {
		node := activity.GetCreatedAt().Format("January 2006")
		if _, ok := c.nodes[node]; !ok {
			contribActivities = contributionActivities{
				commits:            &stats{make(map[string]int, 0), sync.Mutex{}},
				creates:            &creations{make(map[string]string, 0), sync.Mutex{}},
				pullRequests:       &stats{make(map[string]int, 0), sync.Mutex{}},
				pullRequestReviews: &stats{make(map[string]int, 0), sync.Mutex{}},
				issues:             &stats{make(map[string]int, 0), sync.Mutex{}},
			}
		}
		switch *activity.Type {
		case "PushEvent":
			c.wg.Add(1)
			go c.parseCommits(activity, contribActivities.commits)
		case "PullRequestEvent":
			c.wg.Add(1)
			go c.parsePullRequests(activity, contribActivities.pullRequests)
		case "PullRequestReviewEvent":
			c.wg.Add(1)
			go c.parsePullRequestReviews(activity, contribActivities.pullRequestReviews)
		case "IssuesEvent":
			c.wg.Add(1)
			go c.parseIssues(activity, contribActivities.issues)
		case "CreateEvent":
			c.wg.Add(1)
			go c.parseCreates(activity, contribActivities.creates)
		case "ForkEvent":
			c.wg.Add(1)
			go c.parseForks(activity, contribActivities.creates)
		}
		c.nodes[node] = contribActivities
	}
	c.wg.Wait()
	return nil
}

// parseCommits parses the commits data from the provided event and updates the provided commit map with parsed information.
func (c *Contributions) parseCommits(event *github.Event, stats *stats) {
	defer c.wg.Done()
	repoId := event.GetRepo().GetID()
	repo, _ := c.getRepoById(repoId)
	payload, _ := event.ParsePayload()
	ref := "refs/heads/" + repo.GetDefaultBranch()
	if !repo.GetFork() && payload.(*github.PushEvent).GetRef() == ref {
		commitCount := payload.(*github.PushEvent).GetSize()
		repoName := event.GetRepo().GetName()
		stats.mu.Lock()
		stats.statMap[repoName] += commitCount
		stats.mu.Unlock()
	}
}

// parsePullRequests parses the pull requests data from the provided event and updates the provided  pull request map
//with parsed information.
func (c *Contributions) parsePullRequests(event *github.Event, stats *stats) {
	defer c.wg.Done()
	repoId := event.GetRepo().GetID()
	repo, _ := c.getRepoById(repoId)
	if !repo.GetFork() {
		payload, _ := event.ParsePayload()
		if payload.(*github.PullRequestEvent).GetAction() == "opened" {
			repoName := event.GetRepo().GetName()
			stats.mu.Lock()
			stats.statMap[repoName] += 1
			stats.mu.Unlock()
		}
	}
}

// parsePullRequestReviews parses the pull request reviews data from the provided event and updates the provided
//pull request review map with parsed information.
func (c *Contributions) parsePullRequestReviews(event *github.Event, stats *stats) {
	defer c.wg.Done()
	repoId := event.GetRepo().GetID()
	repo, _ := c.getRepoById(repoId)
	if !repo.GetFork() {
		repoName := event.GetRepo().GetName()
		stats.mu.Lock()
		stats.statMap[repoName] += 1
		stats.mu.Unlock()
	}
}

// parseIssues parses the issues data from the provided event and updates the provided issues map with parsed information.
func (c *Contributions) parseIssues(event *github.Event, stats *stats) {
	defer c.wg.Done()
	repoId := event.GetRepo().GetID()
	repo, _ := c.getRepoById(repoId)
	payload, _ := event.ParsePayload()
	action := payload.(*github.IssuesEvent).GetAction()
	if !repo.GetFork() && action == "opened" {
		repoName := event.GetRepo().GetName()
		stats.mu.Lock()
		stats.statMap[repoName] += 1
		stats.mu.Unlock()
	}
}

// parseCreates parses the created repositories data from the provided event and updates both the provided create map
//and commit map with parsed information.
func (c *Contributions) parseCreates(event *github.Event, creations *creations) {
	defer c.wg.Done()
	payload, _ := event.ParsePayload()
	if payload.(*github.CreateEvent).GetRefType() == "repository" {
		repoName := event.GetRepo().GetName()
		createdAt := event.GetCreatedAt().Format("Jan 02")
		creations.mu.Lock()
		creations.creationMap[repoName] = createdAt
		creations.mu.Unlock()
	}
}

// parseForks parses the forks data from the provided event and updates the provided fork map with parsed information.
func (c *Contributions) parseForks(event *github.Event, creations *creations) {
	defer c.wg.Done()
	payload, _ := event.ParsePayload()
	repoName := payload.(*github.ForkEvent).GetForkee().GetFullName()
	createdAt := event.GetCreatedAt().Format("Jan 02")
	creations.mu.Lock()
	creations.creationMap[repoName] = createdAt
	creations.mu.Unlock()
}

// getContributionData retrieves the contribution data of a user for the past 90 days. Refer to
//https://developer.github.com/v3/activity/events/#events for more information.
func (c *Contributions) getContributionData() ([]*github.Event, error) {
	var eventList []*github.Event
	response := &github.Response{}
	response.NextPage = 1
	for response.NextPage > 0 {
		options := &github.ListOptions{
			Page:    response.NextPage,
			PerPage: 100,
		}
		events, res, err := c.Client.Activity.ListEventsPerformedByUser(c.Context, c.Username, false, options)
		if err != nil {
			return nil, err
		}
		response = res
		eventList = append(eventList, events...)
	}
	return eventList, nil
}

// getRepoById returns the github repository based on the provided ID.
func (c *Contributions) getRepoById(id int64) (*github.Repository, error) {
	repo, _, err := c.Client.Repositories.GetByID(c.Context, id)
	return repo, err
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
	totalRepoCount := len(c.nodes[key].commits.statMap)
	if totalRepoCount < 1 {
		return nil
	}

	createdRepos := c.nodes[key].creates.creationMap
	var totalCommits int
	var childNodes []*tview.TreeNode
	for repo, commitCount := range c.nodes[key].commits.statMap {
		// special condition to keep track of initial commit in a repository
		if _, ok := createdRepos[repo]; ok {
			commitCount += 1
		}
		commitText := pluralize("commit", commitCount)
		text := fmt.Sprintf(" [white]%s  [gray::d]%d %s", repo, commitCount, commitText)
		child := tview.NewTreeNode(text).SetSelectable(false)
		childNodes = append(childNodes, child)
		totalCommits += commitCount
	}

	commitText := pluralize("commit", totalCommits)
	repoText := pluralize("repository", totalRepoCount)
	text := fmt.Sprintf(" [::b]Created %d %s in %d %s", totalCommits, commitText, totalRepoCount, repoText)
	node := tview.NewTreeNode(text).SetSelectable(true)
	node.SetChildren(childNodes)
	return node
}

// getCreateNode returns the tree node with created repositories data for a given key, which is the month.
func (c *Contributions) getCreateNode(key string) *tview.TreeNode {
	totalRepoCount := len(c.nodes[key].creates.creationMap)
	if totalRepoCount < 1 {
		return nil
	}

	var childNodes []*tview.TreeNode
	for repo, createdAt := range c.nodes[key].creates.creationMap {
		text := fmt.Sprintf(" [white]%s  [gray::d]%s", repo, createdAt)
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
	totalRepoCount := len(c.nodes[key].pullRequests.statMap)
	if totalRepoCount < 1 {
		return nil
	}

	var totalPrCount int
	var childNodes []*tview.TreeNode
	for repo, prCount := range c.nodes[key].pullRequests.statMap {
		prText := pluralize("pull request", prCount)
		text := fmt.Sprintf(" [white]%s  [gray::d]%d %s", repo, prCount, prText)
		child := tview.NewTreeNode(text).SetSelectable(false)
		childNodes = append(childNodes, child)
		totalPrCount += prCount
	}

	prText := pluralize("pull request", totalPrCount)
	repoText := pluralize("repository", totalRepoCount)
	text := fmt.Sprintf(" [::b]Opened %d %s in %d %s", totalPrCount, prText, totalRepoCount, repoText)
	node := tview.NewTreeNode(text).SetSelectable(true)
	node.SetChildren(childNodes)
	return node
}

// getPullRequestReviewNode returns the tree node with pull request reviews data for a given key, which is the month.
func (c *Contributions) getPullRequestReviewNode(key string) *tview.TreeNode {
	totalRepoCount := len(c.nodes[key].pullRequestReviews.statMap)
	if totalRepoCount < 1 {
		return nil
	}

	var totalReviewCount int
	var childNodes []*tview.TreeNode
	for repo, reviewCount := range c.nodes[key].pullRequestReviews.statMap {
		reviewText := pluralize("pull request", reviewCount)
		text := fmt.Sprintf(" [white]%s  [gray::d]%d %s", repo, reviewCount, reviewText)
		child := tview.NewTreeNode(text).SetSelectable(false)
		childNodes = append(childNodes, child)
		totalReviewCount += reviewCount
	}

	reviewText := pluralize("pull request", totalReviewCount)
	repoText := pluralize("repository", totalRepoCount)
	text := fmt.Sprintf(" [::b]Reviewed %d %s in %d %s", totalReviewCount, reviewText, totalRepoCount, repoText)
	node := tview.NewTreeNode(text).SetSelectable(true)
	node.SetChildren(childNodes)
	return node
}

// getIssueNode returns the tree node with issues data for a given key, which is the month.
func (c *Contributions) getIssueNode(key string) *tview.TreeNode {
	totalRepoCount := len(c.nodes[key].issues.statMap)
	if totalRepoCount < 1 {
		return nil
	}

	var totalIssueCount int
	var childNodes []*tview.TreeNode
	for repo, issueCount := range c.nodes[key].issues.statMap {
		issueText := pluralize("issue", issueCount)
		text := fmt.Sprintf(" [white]%s  [gray::d]%d %s", repo, issueCount, issueText)
		child := tview.NewTreeNode(text).SetSelectable(false)
		childNodes = append(childNodes, child)
		totalIssueCount += issueCount
	}

	issueText := pluralize("issue", totalIssueCount)
	repoText := pluralize("repository", totalRepoCount)
	text := fmt.Sprintf(" [::b]Opened %d %s in %d %s", totalIssueCount, issueText, totalRepoCount, repoText)
	node := tview.NewTreeNode(text).SetSelectable(true)
	node.SetChildren(childNodes)
	return node
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
