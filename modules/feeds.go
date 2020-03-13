package modules

import (
	"fmt"

	"github.com/google/go-github/v29/github"
	"github.com/rivo/tview"

	"github.com/bharath-srinivas/timeago"
)

func (c *Client) FeedsWidget() *Widget {
	widget := tview.NewTextView()
	widget.SetBorder(true)
	widget.SetDynamicColors(true)
	widget.SetScrollable(true)
	widget.SetTextAlign(tview.AlignCenter)

	eventTypes := tview.NewDropDown()
	eventTypes.SetBorder(true)
	eventTypes.SetLabel("Filter by: ")
	eventTypes.SetOptions([]string{"AllEvents", "CreateEvent", "ForkEvent", "StarEvent"}, nil)

	eventTypes.SetSelectedFunc(func(text string, index int) {
		go func() {
			widget.Clear()
			events, _, err := c.client.Activity.ListEventsReceivedByUser(c.ctx, c.username, false, &github.ListOptions{
				PerPage: 100,
			})
			if err != nil {
				_, _ = fmt.Fprintln(widget, "[::b]an error occurred while retrieving events")
				return
			}

			widget.SetTitle(fmt.Sprintf("[green::b]%s", text))
			if text == "StarEvent" {
				text = "WatchEvent"
			}

			for _, event := range events {
				if index > -1 && *event.Type == text {
					data := extractEventData(event)
					_, _ = fmt.Fprintf(widget, data)
				} else if text == "AllEvents" {
					data := extractEventData(event)
					_, _ = fmt.Fprintf(widget, data)
				}
			}
		}()
	})

	eventTypes.SetCurrentOption(0)
	widget.SetChangedFunc(func() {
		widget.ScrollToBeginning()
		c.app.Draw()
	})

	flex := tview.NewFlex()
	flex.SetBorder(true)
	flex.SetTitle(string('\U0001F642') + " [green::b]Feeds")
	flex.SetDirection(tview.FlexRow)

	flex.AddItem(eventTypes, 0, 1, false)
	flex.AddItem(widget, 0, 7, false)

	return &Widget{
		Parent:   flex,
		Children: []tview.Primitive{eventTypes, widget},
	}
}

func extractEventData(event *github.Event) string {
	switch *event.Type {
	case "CreateEvent":
		time := timeago.NoMax(timeago.English).Format(event.GetCreatedAt())
		return fmt.Sprintf("[::b]%s [::d]created a repository [::b]%s [gray::d]%s\n", event.Actor.GetLogin(), event.Repo.GetName(), time)
	case "ForkEvent":
		payload, _ := event.ParsePayload()
		fork := payload.(*github.ForkEvent).Forkee.GetFullName()
		time := timeago.NoMax(timeago.English).Format(event.GetCreatedAt())
		return fmt.Sprintf("[::b]%s [::d]forked [::b]%s [::d]from [::b]%s [gray::d]%s\n", event.Actor.GetLogin(), fork, event.Repo.GetName(), time)
	case "WatchEvent":
		time := timeago.NoMax(timeago.English).Format(event.GetCreatedAt())
		return fmt.Sprintf("[::b]%s [::d]starred [::b]%s [gray::d]%s\n", event.Actor.GetLogin(), event.Repo.GetName(), time)
	default:
		return fmt.Sprint("")
	}
}
