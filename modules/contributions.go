package modules

import (
	"fmt"

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
	nodes map[string][]*tview.TreeNode
}

// ContributionsWidget returns a new instance of contribution widget.
func ContributionsWidget(app *tview.Application, config config.Config) *Contributions {
	widget := tview.NewTreeView().
		SetTopLevel(1)
	widget.SetTitle(string('\U0001F4C8') + " [green::b]Contribution activity").
		SetBorder(true)
	c := &Contributions{
		app:      app,
		nodes:    map[string][]*tview.TreeNode{},
		Base:     views.NewBase(app, config),
		TreeView: widget,
	}
	c.SetSelectedFunc(func(node *tview.TreeNode) {
		node.SetExpanded(!node.IsExpanded())
	})
	c.render()
	return c
}

// render renders the contribution activities of a user in a tree view.
func (c *Contributions) render() {
	root := tview.NewTreeNode("[::b]Contributions in the last 90 days")
	c.TreeView.
		SetRoot(root).
		SetCurrentNode(root)

	c.createNodes()
	for nodeText, children := range c.nodes {
		childNode := tview.NewTreeNode(nodeText).
			SetSelectable(true).
			SetChildren(children)
		root.AddChild(childNode)
	}
}

// createNodes creates the child nodes based on the month of contribution.
func (c *Contributions) createNodes() {
	contributions, err := c.getContributionData()
	if err != nil {
		root := tview.NewTreeNode("[::b]an error occurred while retrieving contribution data")
		c.TreeView.SetRoot(root).SetCurrentNode(root)
		return
	}

	var children []*tview.TreeNode
	for _, activity := range contributions {
		nodeText := fmt.Sprintf("[::b]%s", activity.GetCreatedAt().Format("January 2006"))
		if _, ok := c.nodes[nodeText]; !ok {
			children = []*tview.TreeNode{}
		}

		switch *activity.Type {
		case "PushEvent":
			if !activity.GetRepo().GetFork() {
				payload, _ := activity.ParsePayload()
				commitCount := len(payload.(*github.PushEvent).Commits)
				commitString := "commit"
				if commitCount > 1 {
					commitString = "commits"
				}
				t := activity.GetCreatedAt().Format("Jan 02")
				s := fmt.Sprintf(" [::b]Created %d %s in %s  [gray::d]%s", commitCount, commitString, activity.GetRepo().GetName(), t)
				child := tview.NewTreeNode(s).SetSelectable(false)
				children = append(children, child)
			}
		case "PullRequestEvent":
			if !activity.GetRepo().GetFork() {
				payload, _ := activity.ParsePayload()
				if payload.(*github.PullRequestEvent).GetAction() == "opened" {
					t := activity.GetCreatedAt().Format("Jan 02")
					s := fmt.Sprintf(" [::b]Opened a pull request in %s  [gray::d]%s", activity.GetRepo().GetName(), t)
					child := tview.NewTreeNode(s).SetSelectable(false)
					children = append(children, child)
				}
			}
		case "PullRequestReviewEvent":
			if !activity.GetRepo().GetFork() {
				t := activity.GetCreatedAt().Format("Jan 02")
				s := fmt.Sprintf(" [::b]Reviewed a pull request in %s  [gray::d]%s", activity.GetRepo().GetName(), t)
				child := tview.NewTreeNode(s).SetSelectable(false)
				children = append(children, child)
			}
		case "IssuesEvent":
			if !activity.GetRepo().GetFork() {
				t := activity.GetCreatedAt().Format("Jan 02")
				s := fmt.Sprintf(" [::b]Created an issue in %s  [gray::d]%s", activity.GetRepo().GetName(), t)
				child := tview.NewTreeNode(s).SetSelectable(false)
				children = append(children, child)
			}
		case "CreateEvent":
			payload, _ := activity.ParsePayload()
			if payload.(*github.CreateEvent).GetRefType() == "repository" {
				t := activity.GetCreatedAt().Format("Jan 02")
				s := fmt.Sprintf(" [::b]Created a repository %s  [gray::d]%s", activity.GetRepo().GetName(), t)
				child := tview.NewTreeNode(s).SetSelectable(false)
				children = append(children, child)
			}
		}
		c.nodes[nodeText] = children
	}
}

// getContributionData retrieves the contribution data of a user for the past 90 days. Refer to
//https://developer.github.com/v3/activity/events/#events for more information.
func (c *Contributions) getContributionData() ([]*github.Event, error) {
	var eventList []*github.Event
	response := &github.Response{}
	response.NextPage = 1
	for response.NextPage > 0 {
		events, res, err := c.Client.Activity.ListEventsPerformedByUser(c.Context, c.Username, false, &github.ListOptions{
			Page:    response.NextPage,
			PerPage: 100,
		})
		if err != nil {
			return nil, err
		}
		response = res
		eventList = append(eventList, events...)
	}
	return eventList, nil
}
