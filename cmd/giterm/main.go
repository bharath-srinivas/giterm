package main

import (
	"log"
	"os"

	"github.com/bharath-srinivas/giterm/modules"
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()
	gitApp := modules.New(app)

	grid := tview.NewGrid()
	grid.SetRows()
	grid.SetColumns(35)
	grid.SetBorder(false)

	grid.AddItem(gitApp.ProfileWidget(), 0, 0, 1, 1, 0, 0, false)
	grid.AddItem(gitApp.RepoWidget(), 0, 1, 1, 1, 0, 0, true)

	app.SetRoot(grid, true)
	if err := app.Run(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
