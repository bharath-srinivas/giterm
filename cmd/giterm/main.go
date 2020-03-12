package main

import (
	"github.com/rivo/tview"

	"github.com/bharath-srinivas/giterm/modules"
)

func main() {
	gitApp := modules.New(tview.NewApplication())
	gitApp.Start()
}
