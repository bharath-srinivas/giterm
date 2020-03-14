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

	feeds := tview.NewTextView()
	feeds.SetBorder(true)
	feeds.SetDynamicColors(true)
	feeds.SetScrollable(true)
	feeds.SetTextAlign(tview.AlignCenter)
	feeds.SetChangedFunc(func() {
		feeds.ScrollToBeginning()
		c.app.Draw()
	})

	pageSizeOptions := tview.NewDropDown()
	pageSizeOptions.SetLabelColor(tcell.ColorWhite)
	pageSizeOptions.SetFieldTextColor(tcell.ColorWhite)
	pageSizeOptions.SetFieldBackgroundColor(tcell.ColorBlack)
	pageSizeOptions.SetLabel("Items per page: ")
	pageSizeOptions.SetOptions([]string{"25", "50", "75", "100"}, func(text string, index int) {
		pageSize, _ = strconv.Atoi(text)
		res = c.displayEventsData(feeds, 1, pageSize)
	})
	pageSizeOptions.SetCurrentOption(0)

	first := tview.NewButton(string('\U000000AB')).SetLabelColor(tcell.ColorWhite)
	first.SetLabelColorActivated(tcell.ColorBlack)
	first.SetBackgroundColorActivated(tcell.ColorWhite)
	first.SetBackgroundColor(tcell.ColorBlack)
	first.SetSelectedFunc(func() {
		if res != nil {
			res = c.displayEventsData(feeds, res.FirstPage, pageSize)
		}
	})

	last := tview.NewButton(string('\U000000BB'))
	last.SetLabelColor(tcell.ColorWhite)
	last.SetLabelColorActivated(tcell.ColorBlack)
	last.SetBackgroundColorActivated(tcell.ColorWhite)
	last.SetBackgroundColor(tcell.ColorBlack)
	last.SetSelectedFunc(func() {
		if res != nil {
			res = c.displayEventsData(feeds, res.LastPage, pageSize)
		}
	})

	prev := tview.NewButton(string('\U000025C4'))
	prev.SetLabelColor(tcell.ColorWhite)
	prev.SetLabelColorActivated(tcell.ColorBlack)
	prev.SetBackgroundColorActivated(tcell.ColorWhite)
	prev.SetBackgroundColor(tcell.ColorBlack)
	prev.SetSelectedFunc(func() {
		if res != nil {
			res = c.displayEventsData(feeds, res.PrevPage, pageSize)
		}
	})

	next := tview.NewButton(string('\U000025BA'))
	next.SetLabelColor(tcell.ColorWhite)
	next.SetLabelColorActivated(tcell.ColorBlack)
	next.SetBackgroundColorActivated(tcell.ColorWhite)
	next.SetBackgroundColor(tcell.ColorBlack)
	next.SetSelectedFunc(func() {
		if res != nil {
			res = c.displayEventsData(feeds, res.NextPage, pageSize)
		}
	})

	header := tview.NewFlex()
	header.AddItem(pageSizeOptions, 0, 1, false)

	footer := tview.NewFlex()
	footer.AddItem(first, 0, 1, false)
	footer.AddItem(prev, 0, 1, false)
	footer.AddItem(next, 0, 1, false)
	footer.AddItem(last, 0, 1, false)

	view := tview.NewFlex()
	view.SetBorder(true)
	view.SetTitle(string('\U0001F559') + " [green::b]Feeds")
	view.SetDirection(tview.FlexRow)

	view.AddItem(header, 0, 1, false)
	view.AddItem(feeds, 0, 15, false)
	view.AddItem(footer, 0, 1, false)

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
