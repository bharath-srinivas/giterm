package modules

import (
	"fmt"

	"github.com/rivo/tview"

	"github.com/bharath-srinivas/giterm/views"
)

// Pinned represents a pinned github repo or gist.
type Pinned struct {
	*views.PinnedView
	app *tview.Application
}

// PinnedWidget returns an instance of pinned widget.
func PinnedWidget(app *tview.Application) *Pinned {
	widget := views.NewPinnedView()
	widget.SetBorder(true).
		SetTitle("\U0001F4CC [green::b]Pinned")
	p := &Pinned{widget, app}
	go p.Refresh()
	return p
}

// Refresh refreshes the pinned widget.
func (p *Pinned) Refresh() {
	p.app.QueueUpdateDraw(func() {
		p.Clear()
		p.display()
	})
}

// display renders the pinned item info within the custom pinned view.
func (p *Pinned) display() {
	if user == nil {
		return
	}

	var pinnedItems []*views.PinnedItem
	nodes := user.Viewer.PinnedItems.Nodes
	for _, node := range nodes {
		var archived, lang string
		name := node.Repository.Name
		if name == "" {
			name = "[white::b]" + node.Gist.Description
			pinnedItem := views.NewPinnedItem().
				SetContent(&views.PinnedItemContent{
					Name: name,
				})
			pinnedItems = append(pinnedItems, pinnedItem)
			continue
		}

		if node.Repository.Owner.Login != user.Viewer.Login {
			name = "[white::b]" + node.Repository.NameWithOwner
		}

		name = "[white::b]" + name
		description := "[white]" + node.Repository.Description
		if node.Repository.PrimaryLanguage != nil {
			name := node.Repository.PrimaryLanguage.Name
			color := node.Repository.PrimaryLanguage.Color
			lang = fmt.Sprintf("[%s]\u25CF %s", color, name)
		}

		stars := fmt.Sprintf("[white]\u2605 %d", node.Repository.Stargazers.TotalCount)
		if node.Repository.IsArchived {
			archived = "[darkslategray:white:] Archived [:black:]"
		}

		pinnedItem := views.NewPinnedItem().
			SetContent(&views.PinnedItemContent{
				Name:        name,
				Description: description,
				Lang:        lang,
				Stars:       stars,
				Archived:    archived,
			})
		pinnedItems = append(pinnedItems, pinnedItem)
	}
	p.SetItems(pinnedItems)
}
