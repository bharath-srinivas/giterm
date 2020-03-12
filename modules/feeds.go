package modules

import (
	"fmt"

	"github.com/google/go-github/v29/github"
	"github.com/rivo/tview"

	"github.com/bharath-srinivas/timeago"
)

func (g *GitApp) FeedsPage() *Page {
	widget := tview.NewTextView()
	widget.SetBorder(true)
	widget.SetDynamicColors(true)
	widget.SetScrollable(true)
	widget.SetTextAlign(tview.AlignCenter)

	eventTypes := tview.NewDropDown()
	eventTypes.SetBorder(true)
	eventTypes.SetLabel("Filter by: ")
	eventTypes.SetOptions([]string{"CreateEvent", "ForkEvent", "StarEvent"}, nil)

	eventTypes.SetSelectedFunc(func(text string, index int) {
		go func() {
			widget.Clear()
			events, _, err := g.Client.Activity.ListEventsReceivedByUser(g.Context, "", false, &github.ListOptions{
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
					switch *event.Type {
					case "CreateEvent":
						time := timeago.NoMax(timeago.English).Format(event.GetCreatedAt())
						_, _ = fmt.Fprintf(widget, "[::b]%s [::d]created a repository [::b]%s [gray::d]%s\n", event.Actor.GetLogin(), event.Repo.GetName(), time)
					case "ForkEvent":
						payload, _ := event.ParsePayload()
						time := timeago.NoMax(timeago.English).Format(event.GetCreatedAt())
						_, _ = fmt.Fprintf(widget, "[::b]%s [::d]forked [::b]%s [::d]from [::b]%s [gray::d]%s\n", event.Actor.GetLogin(), payload.(*github.ForkEvent).Forkee.GetFullName(), event.Repo.GetName(), time)
					case "WatchEvent":
						time := timeago.NoMax(timeago.English).Format(event.GetCreatedAt())
						_, _ = fmt.Fprintf(widget, "[::b]%s [::d]starred [::b]%s [gray::d]%s\n", event.Actor.GetLogin(), event.Repo.GetName(), time)
					default:
						_, _ = fmt.Fprintf(widget, "")
					}
				}
			}
		}()
	})

	eventTypes.SetCurrentOption(0)
	widget.SetChangedFunc(func() {
		g.App.Draw()
	})

	flex := tview.NewFlex()
	flex.SetBorder(true)
	flex.SetTitle(string('\U0001F642') + " [green::b]Feeds")
	flex.SetDirection(tview.FlexRow)

	flex.AddItem(eventTypes, 0, 1, false)
	flex.AddItem(widget, 0, 7, false)

	return &Page{
		Name:            "Feeds",
		Parent:          flex,
		ChildComponents: []tview.Primitive{eventTypes, widget},
	}
}
