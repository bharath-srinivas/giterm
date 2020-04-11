package modules

import (
	"fmt"

	"github.com/google/go-github/v30/github"
	"github.com/olekukonko/tablewriter"
	"github.com/rivo/tview"

	"github.com/bharath-srinivas/giterm/config"
	"github.com/bharath-srinivas/giterm/views"
	"github.com/bharath-srinivas/timeago"
)

// Notifications represents the github notifications.
type Notifications struct {
	*views.TextWidget
	*github.NotificationListOptions
	*github.Response
}

// NotificationsWidget returns a new instance of notifications widget.
func NotificationsWidget(app *tview.Application, config config.Config) *Notifications {
	widget := views.NewTextView(app, config, true)
	widget.SetWrap(false).
		SetTitleAlign(tview.AlignCenter)
	n := &Notifications{
		widget,
		&github.NotificationListOptions{
			All: true,
			ListOptions: github.ListOptions{
				Page:    1,
				PerPage: 25,
			},
		},
		&github.Response{},
	}
	go n.Refresh()
	return n
}

// Refresh refreshes the notifications widget.
func (n *Notifications) Refresh() {
	n.Redraw(func() {
		n.display(n.NotificationListOptions)
	})
}

// First navigates to the first page of the notifications.
func (n *Notifications) First() {
	if n.Response != nil {
		n.Page = n.FirstPage
		go n.Refresh()
	}
}

// Last navigates to the last page of the notifications.
func (n *Notifications) Last() {
	if n.Response != nil {
		n.Page = n.LastPage
		go n.Refresh()
	}
}

// Prev navigates to the previous page of the notifications.
func (n *Notifications) Prev() {
	if n.Response != nil {
		n.Page = n.PrevPage
		go n.Refresh()
	}
}

// Next navigates to the next page of the notifications.
func (n *Notifications) Next() {
	if n.Response != nil {
		n.Page = n.NextPage
		go n.Refresh()
	}
}

// display renders the notifications according to the provided options.
func (n *Notifications) display(options *github.NotificationListOptions) {
	notifications, res, err := n.Activity.ListNotifications(n.Context, options)
	if err != nil {
		_, _ = fmt.Fprint(n.TextView, "[::b]an error occurred while retrieving notifications")
		return
	}
	n.Response = res

	table := tablewriter.NewWriter(n.TextView)
	table.SetBorder(false)
	table.SetAutoWrapText(false)
	table.SetColumnSeparator(" ")
	table.SetRowLine(true)
	table.SetRowSeparator("_")
	table.SetCenterSeparator("")
	table.SetColMinWidth(1, 165)
	for _, notification := range notifications {
		reason := fmt.Sprintf("[:green:d] %s [:black:]", notification.GetReason())
		notificationType := notification.GetSubject().GetType()
		repoName := "[::b]" + notification.GetRepository().GetFullName()
		subject := tview.Escape(notification.GetSubject().GetTitle())
		time := timeago.NoMax(timeago.English).Format(notification.GetUpdatedAt())
		table.Append([]string{
			"\n" + notificationType,
			"\n" + repoName + "\n" + subject,
			"\n" + reason,
			"\n" + "[gray::d]" + time})
	}
	table.Render()
}
