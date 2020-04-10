// Package modules implements the different kinds of widgets that are used by the application.
package modules

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/mmcdole/gofeed/atom"
	"github.com/rivo/tview"

	"github.com/bharath-srinivas/giterm/config"
	"github.com/bharath-srinivas/giterm/views"
	"github.com/bharath-srinivas/timeago"
)

// Feeds represents the github feeds.
type Feeds struct {
	*views.TextWidget
	feedsUrl string
}

// FeedsWidget returns a new instance of feeds widget.
func FeedsWidget(app *tview.Application, config config.Config) *Feeds {
	widget := views.NewTextView(app, config, true)
	widget.SetTextAlign(tview.AlignCenter).
		SetTitle(string('\U0001F559') + " [green::b]Feeds").
		SetBorder(true)

	f := &Feeds{
		TextWidget: widget,
		feedsUrl:   config.FeedsUrl,
	}
	go f.Refresh()
	return f
}

// Refresh refreshes the feeds widget.
func (f *Feeds) Refresh() {
	f.Redraw(f.display)
}

// display renders the feeds in a text view.
func (f *Feeds) display() {
	page := 1
	parser := atom.Parser{}
	for {
		f.feedsUrl = fmt.Sprintf("%s&page=%d", f.feedsUrl, page)
		resp, err := http.Get(f.feedsUrl)
		if err != nil {
			_, _ = fmt.Fprintln(f.TextView, "[::b]an error occurred while retrieving feeds")
			return
		}

		feeds, err := parser.Parse(resp.Body)
		if err != nil {
			_ = resp.Body.Close()
			_, _ = fmt.Fprintln(f.TextView, "[::b]an error occurred while retrieving feeds")
			return
		}

		_ = resp.Body.Close()
		if len(feeds.Entries) == 0 {
			break
		}

		for _, entry := range feeds.Entries {
			time := "[gray::d]" + timeago.NoMax(timeago.English).Format(*entry.PublishedParsed)
			title := strings.Split(entry.Title, " ")
			user, repo := title[0], title[len(title)-1]
			title[0] = "[::b]" + user + "[::d]"
			title[len(title)-1] = "[::b]" + repo + "[::d]"
			feed := strings.Join(title, " ")
			_, _ = fmt.Fprintln(f.TextView, feed+"\t"+time+"\n")
		}
		page += 1
	}
}
