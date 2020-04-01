// Package modules implements the different kinds of widgets that are used by the application.
package modules

import (
	"fmt"

	"github.com/google/go-github/v30/github"
	"github.com/rivo/tview"

	"github.com/bharath-srinivas/giterm/config"
	"github.com/bharath-srinivas/giterm/views"
	"github.com/bharath-srinivas/timeago"
)

// Feeds represents the github feeds.
type Feeds struct {
	*views.TextWidget
	*github.ListOptions
	*github.Response
}

// FeedsWidget returns a new instance of feeds widget.
func FeedsWidget(app *tview.Application, config config.Config) *Feeds {
	widget := views.NewTextView(app, config, true)
	widget.SetTextAlign(tview.AlignCenter)
	return &Feeds{
		TextWidget:  widget,
		ListOptions: &github.ListOptions{},
		Response:    &github.Response{},
	}
}

// Refresh refreshes the feeds widget.
func (f *Feeds) Refresh() {
	f.Redraw(func() {
		f.display(f.ListOptions)
	})
}

// SetPageSize sets the page size.
func (f *Feeds) SetPageSize(pageSize int) {
	f.ListOptions.Page = 1
	f.ListOptions.PerPage = pageSize
	go f.Refresh()
}

// First navigates to the first page of the feed.
func (f *Feeds) First() {
	if f.Response != nil {
		f.ListOptions.Page = f.FirstPage
		go f.Refresh()
	}
}

// Last navigates to the last page of the feed.
func (f *Feeds) Last() {
	if f.Response != nil {
		f.ListOptions.Page = f.LastPage
		go f.Refresh()
	}
}

// Prev navigates to the previous page of the feed.
func (f *Feeds) Prev() {
	if f.Response != nil {
		f.ListOptions.Page = f.PrevPage
		go f.Refresh()
	}
}

// Next navigates to the next page of the feed.
func (f *Feeds) Next() {
	if f.Response != nil {
		f.ListOptions.Page = f.NextPage
		go f.Refresh()
	}
}

// display renders the feeds according to the provided pagination and page size options.
func (f *Feeds) display(options *github.ListOptions) {
	events, res, err := f.Client.Activity.ListEventsReceivedByUser(f.Context, f.Username, false, options)
	if err != nil {
		_, _ = fmt.Fprintln(f.TextView, "[::b]an error occurred while retrieving feeds")
		return
	}
	f.Response = res
	for _, event := range events {
		switch *event.Type {
		case "CreateEvent":
			payload, _ := event.ParsePayload()
			if payload.(*github.CreateEvent).GetRefType() == "repository" {
				time := timeago.NoMax(timeago.English).Format(event.GetCreatedAt())
				_, _ = fmt.Fprintf(f.TextView, "[::b]%s [::d]created a repository [::b]%s [gray::d]%s\n\n", event.Actor.GetLogin(), event.Repo.GetName(), time)
			}
		case "PushEvent":
			time := timeago.NoMax(timeago.English).Format(event.GetCreatedAt())
			_, _ = fmt.Fprintf(f.TextView, "[::b]%s [::d]pushed to [::b]%s [gray::d]%s\n\n", event.Actor.GetLogin(), event.Repo.GetName(), time)
		case "ForkEvent":
			payload, _ := event.ParsePayload()
			fork := payload.(*github.ForkEvent).Forkee.GetFullName()
			time := timeago.NoMax(timeago.English).Format(event.GetCreatedAt())
			_, _ = fmt.Fprintf(f.TextView, "[::b]%s [::d]forked [::b]%s [::d]from [::b]%s [gray::d]%s\n\n", event.Actor.GetLogin(), fork, event.Repo.GetName(), time)
		case "WatchEvent":
			time := timeago.NoMax(timeago.English).Format(event.GetCreatedAt())
			_, _ = fmt.Fprintf(f.TextView, "[::b]%s [::d]starred [::b]%s [gray::d]%s\n\n", event.Actor.GetLogin(), event.Repo.GetName(), time)
		}
	}
}
