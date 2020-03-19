package modules

import (
	"fmt"

	"github.com/google/go-github/v29/github"
	"github.com/rivo/tview"

	"github.com/bharath-srinivas/giterm/config"
	"github.com/bharath-srinivas/giterm/views"
	"github.com/bharath-srinivas/timeago"
)

type Feeds struct {
	*views.TextWidget
}

func FeedsWidget(app *tview.Application, config config.Config) *Feeds {
	widget := views.NewTextView(app, config, true)
	return &Feeds{
		TextWidget: widget,
	}
}

func (f *Feeds) Refresh() {
	f.Redraw(func() {
		f.display(f.ListOptions)
	})
}

func (f *Feeds) display(options *github.ListOptions) {
	events, res, err := f.Client.Activity.ListEventsReceivedByUser(f.Context, f.Username, false, options)
	if err != nil {
		_, _ = fmt.Fprintln(f.View, "[::b]an error occurred while retrieving feeds")
		return
	}
	f.Response = res
	for _, event := range events {
		switch *event.Type {
		case "CreateEvent":
			time := timeago.NoMax(timeago.English).Format(event.GetCreatedAt())
			_, _ = fmt.Fprintf(f.View, "[::b]%s [::d]created a repository [::b]%s [gray::d]%s\n\n", event.Actor.GetLogin(), event.Repo.GetName(), time)
		case "ForkEvent":
			payload, _ := event.ParsePayload()
			fork := payload.(*github.ForkEvent).Forkee.GetFullName()
			time := timeago.NoMax(timeago.English).Format(event.GetCreatedAt())
			_, _ = fmt.Fprintf(f.View, "[::b]%s [::d]forked [::b]%s [::d]from [::b]%s [gray::d]%s\n\n", event.Actor.GetLogin(), fork, event.Repo.GetName(), time)
		case "WatchEvent":
			time := timeago.NoMax(timeago.English).Format(event.GetCreatedAt())
			_, _ = fmt.Fprintf(f.View, "[::b]%s [::d]starred [::b]%s [gray::d]%s\n\n", event.Actor.GetLogin(), event.Repo.GetName(), time)
		}
	}
}

func (f *Feeds) SetPageSize(pageSize int) {
	f.ListOptions.Page = 1
	f.ListOptions.PerPage = pageSize
	go f.Refresh()
}

func (f *Feeds) First() {
	if f.Response != nil {
		f.ListOptions.Page = f.FirstPage
		go f.Refresh()
	}
}

func (f *Feeds) Last() {
	if f.Response != nil {
		f.ListOptions.Page = f.LastPage
		go f.Refresh()
	}
}

func (f *Feeds) Prev() {
	if f.Response != nil {
		f.ListOptions.Page = f.PrevPage
		go f.Refresh()
	}
}

func (f *Feeds) Next() {
	if f.Response != nil {
		f.ListOptions.Page = f.NextPage
		go f.Refresh()
	}
}
