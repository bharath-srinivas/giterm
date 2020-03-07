package main

import (
	"log"
	"os"

	"github.com/rivo/tview"

	"github.com/bharath-srinivas/giterm/modules"
)

func main() {
	gitApp := modules.New(tview.NewApplication())
	gitApp.LoadWidgets()
	gitApp.LoadInputHandler()

	flex := tview.NewFlex()
	flex.SetBorder(false)

	flex.AddItem(gitApp.Widgets["profile"], 0, 1, true)
	flex.AddItem(gitApp.Widgets["repositories"], 0, 5, false)

	gitApp.App.SetRoot(flex, true)
	if err := gitApp.App.Run(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
