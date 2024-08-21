package ui

import (
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"github.com/sarumaj/edu-taschenrechner/pkg/runes"
)

// App is a custom interface for bringing up the application window.
type App struct {
	sync.Once
	fyne.App
	fyne.Window
	Objects map[string]fyne.CanvasObject
}

// Build renders the application window and sets up all widgets.
func (a *App) Build() {
	a.Do(func() {
		a.Objects = make(map[string]fyne.CanvasObject)

		// make label
		display := NewDisplay("_")
		a.Objects["display"] = display

		// make buttons
		for _, btnText := range append(runes.Each("1234567890+-×÷=."), "AC", "()", "<x") {
			a.Objects[btnText] = NewButton(btnText, display)
		}

		// lay out the components
		appContent := container.NewBorder(
			display, nil, nil, nil,
			container.NewGridWithRows(5,
				container.NewGridWithColumns(4, a.getButtons("AC", "<x", "()", "÷")...),
				container.NewGridWithColumns(4, a.getButtons(runes.Each("789×")...)...),
				container.NewGridWithColumns(4, a.getButtons(runes.Each("456-")...)...),
				container.NewGridWithColumns(4, a.getButtons(runes.Each("123+")...)...),
				container.NewGridWithColumns(2,
					container.NewGridWithColumns(2, a.getButtons(runes.Each("0.")...)...),
					a.Objects["="],
				),
			),
		)

		// fill in the application window
		a.Window.SetContent(container.NewCenter(appContent))
	})
}

// getButtons provides a selector to select buttons..
func (a *App) getButtons(in ...string) (out []fyne.CanvasObject) {
	for _, c := range in {
		switch a.Objects[c].(type) {
		case *Button:
			out = append(out, a.Objects[c])

		}
	}

	return
}

// Create new application window.
// Call ShowAndRun to display the window.
func NewApp(title string) *App {
	a := app.New()
	w := a.NewWindow(title)
	w.Resize(fyne.NewSize(600, 400))
	a.Settings().SetTheme(NewDoubleSizeTheme(theme.DefaultTheme()))

	return &App{App: a, Window: w}
}
