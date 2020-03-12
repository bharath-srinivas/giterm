package modules

import (
	"github.com/rivo/tview"

	"github.com/bharath-srinivas/giterm/config"
)

type Widget struct {
	Parent   tview.Primitive
	Children []tview.Primitive
}

func MakeWidget(app *tview.Application, config config.Config, moduleName string) *Widget {
	client := NewClient(app, config)
	var widget *Widget
	switch moduleName {
	case "profile":
		widget = client.ProfileWidget()
	case "repos":
		widget = client.RepoWidget()
	case "feeds":
		widget = client.FeedsWidget()
	case "default":
		widget = &Widget{Parent: tview.NewBox().SetTitle("Unknown").SetBorder(true)}
	}
	return widget
}
