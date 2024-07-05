package main

import "github.com/sarumaj/edu-taschenrechner/pkg/ui"

func main() {
	// create new application window
	app := ui.NewApp("Taschenrechner")
	// build and render window components
	app.Build()
	// execute
	app.ShowAndRun()
}
