package ui

import (
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/sarumaj/edu-taschenrechner/pkg/memory"
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

		label := widget.NewLabel("_")
		label.Alignment = fyne.TextAlignCenter
		label.Truncation = fyne.TextTruncateEllipsis
		a.Objects["label"] = label

		memory := memory.NewMemoryCell()
		a.Objects["memory"] = memory

		// make buttons
		for _, btnText := range append(runes.Each("1234567890+-×÷=."), "AC", "()", "<x") {
			a.Objects[btnText] = widget.NewButton(btnText, doBtnClick(label, memory, btnText))
		}

		// build app content
		btnSelector := getButtons(a.Objects)
		appContent := container.NewGridWithColumns(1,
			label,
			container.NewGridWithColumns(4, btnSelector("AC", "<x", "()", "÷")...),
			container.NewGridWithColumns(4, btnSelector(runes.Each("789×")...)...),
			container.NewGridWithColumns(4, btnSelector(runes.Each("456-")...)...),
			container.NewGridWithColumns(4, btnSelector(runes.Each("123+")...)...),
			container.NewGridWithColumns(2,
				container.NewGridWithColumns(2, btnSelector(runes.Each("0.")...)...),
				a.Objects["="],
			),
		)

		// fill in the application window
		a.Window.SetContent(appContent)
	})
}

// Create new application window.
// Call ShowAndRun to display the window.
func NewApp(title string) *App {
	a := app.New()
	w := a.NewWindow(title)
	w.Resize(fyne.NewSize(600, 400))

	return &App{App: a, Window: w}
}
