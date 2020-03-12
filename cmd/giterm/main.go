package main

import (
	"github.com/rivo/tview"

	"github.com/bharath-srinivas/giterm/app"
)

func main() {
	gitApp := app.New(tview.NewApplication())
	gitApp.Start()
}
