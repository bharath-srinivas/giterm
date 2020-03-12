package main

import (
	"log"
	"os"

	"github.com/rivo/tview"

	"github.com/bharath-srinivas/giterm/modules"
)

func main() {
	gitApp := modules.New(tview.NewApplication())
	gitApp.LoadInputHandler()

	if err := gitApp.App.Run(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
