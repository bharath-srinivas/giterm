package modules

import (
	"fmt"
	"strconv"

	"github.com/gdamore/tcell"
	"github.com/google/go-github/v29/github"
	"github.com/rivo/tview"

	"github.com/bharath-srinivas/timeago"
)

func (c *Client) FeedsWidget() *Widget {
	var res *github.Response
	var pageSize int

	feeds := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetTextAlign(tview.AlignCenter)
	feeds.SetBorder(true)
	feeds.SetChangedFunc(func() {
		feeds.ScrollToBeginning()
		c.app.Draw()
	})

	pageSizeOptions := tview.NewDropDown().
		SetLabelColor(tcell.ColorWhite).
		SetFieldTextColor(tcell.ColorWhite).
		SetFieldBackgroundColor(tcell.ColorBlack).
		SetLabel("Items per page: ").
		SetOptions([]string{"25", "50", "75", "100"}, func(text string, index int) {
			pageSize, _ = strconv.Atoi(text)
			res = c.displayEventsData(feeds, 1, pageSize)
		}).SetCurrentOption(0)

	first := createButton(string('\U000000AB')).
		SetSelectedFunc(func() {
			if res != nil {
				res = c.displayEventsData(feeds, res.FirstPage, pageSize)
			}
		})

	last := createButton(string('\U000000BB')).
		SetSelectedFunc(func() {
			if res != nil {
				res = c.displayEventsData(feeds, res.LastPage, pageSize)
			}
		})

	prev := createButton(string('\U000025C4')).
		SetSelectedFunc(func() {
			if res != nil {
				res = c.displayEventsData(feeds, res.PrevPage, pageSize)
			}
		})

	next := createButton(string('\U000025BA')).
		SetSelectedFunc(func() {
			if res != nil {
				res = c.displayEventsData(feeds, res.NextPage, pageSize)
			}
		})

	header := tview.NewFlex().
		AddItem(pageSizeOptions, 0, 1, false)

	footer := tview.NewFlex().
		AddItem(first, 0, 1, false).
		AddItem(prev, 0, 1, false).
		AddItem(next, 0, 1, false).
		AddItem(last, 0, 1, false)

	view := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(header, 0, 1, false).
		AddItem(feeds, 0, 15, false).
		AddItem(footer, 0, 1, false)
	view.SetTitle(string('\U0001F559') + " [green::b]Feeds")
	view.SetBorder(true)

	return &Widget{
		Parent:   view,
		Children: []tview.Primitive{pageSizeOptions, feeds, first, prev, next, last},
	}
}

func (c *Client) displayEventsData(widget *tview.TextView, page, pageSize int) *github.Response {
	events, res, err := c.client.Activity.ListEventsReceivedByUser(c.ctx, c.username, false, &github.ListOptions{
		Page:    page,
		PerPage: pageSize,
	})
	if err != nil {
		_, _ = fmt.Fprintln(widget, "[::b]an error occurred while retrieving events")
		return res
	}
	widget.Clear()
	go func() {
		for _, event := range events {
			switch *event.Type {
			case "CreateEvent":
				time := timeago.NoMax(timeago.English).Format(event.GetCreatedAt())
				_, _ = fmt.Fprintf(widget, "[::b]%s [::d]created a repository [::b]%s [gray::d]%s\n\n", event.Actor.GetLogin(), event.Repo.GetName(), time)
			case "ForkEvent":
				payload, _ := event.ParsePayload()
				fork := payload.(*github.ForkEvent).Forkee.GetFullName()
				time := timeago.NoMax(timeago.English).Format(event.GetCreatedAt())
				_, _ = fmt.Fprintf(widget, "[::b]%s [::d]forked [::b]%s [::d]from [::b]%s [gray::d]%s\n\n", event.Actor.GetLogin(), fork, event.Repo.GetName(), time)
			case "WatchEvent":
				time := timeago.NoMax(timeago.English).Format(event.GetCreatedAt())
				_, _ = fmt.Fprintf(widget, "[::b]%s [::d]starred [::b]%s [gray::d]%s\n\n", event.Actor.GetLogin(), event.Repo.GetName(), time)
			}
		}
	}()
	return res
}

func createButton(label string) *tview.Button {
	button := tview.NewButton(label).
		SetLabelColor(tcell.ColorWhite).
		SetLabelColorActivated(tcell.ColorBlack).
		SetBackgroundColorActivated(tcell.ColorWhite)
	button.SetBackgroundColor(tcell.ColorBlack)
	return button
}
