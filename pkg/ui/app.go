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

		// make display
		a.Objects["display"] = NewDisplay("_", a.Window)

		// make buttons
		for _, btnText := range append(runes.Each("1234567890+-×÷=.√π^!e°"), "AC", "()", "<x", "sin", "cos", "tan", "log", "ln") {
			a.Objects[btnText] = NewButton(btnText, a.Objects["display"]).SetAlternateText(map[string]string{
				"√":   "x²",
				"sin": "sin⁻¹",
				"cos": "cos⁻¹",
				"tan": "tan⁻¹",
				"log": "10^",
				"ln":  "e^",
				"°":   "1/°",
			}[btnText])
		}

		// make dropdowns
		for name, relations := range map[string][]string{
			"const": {"π", "e"},
			"func":  {"sin", "cos", "tan", "log", "ln"},
		} {
			a.Objects[name] = NewButtonDropDown(a.getButtons(relations...))
		}

		a.Objects["INV"] = NewButton("INV", a.Objects["display"]).SetOnTapped(func() {
			for i := range a.Objects {
				if o, ok := a.Objects[i].(*Button); ok {
					o.Invert()
				}
			}

			for i := range a.Objects {
				if o, ok := a.Objects[i].(*ButtonDropDown); ok {
					o.Update()
				}
			}
		})

		// lay out the components
		appContent := container.NewBorder(
			a.Objects["display"], nil, nil, nil,
			container.NewGridWithRows(6,
				container.NewGridWithColumns(3,
					container.NewGridWithColumns(2, a.getButtons("√", "^")...),
					container.NewGridWithColumns(2, append(a.getButtons("!"), a.Objects["const"])...),
					a.Objects["func"],
				),
				container.NewGridWithColumns(2,
					container.NewGridWithColumns(3, a.getButtons("INV", "AC", "<x")...),
					container.NewGridWithColumns(2, a.getButtons("()", "÷")...),
				),
				container.NewGridWithColumns(4, a.getButtons(runes.Each("789×")...)...),
				container.NewGridWithColumns(4, a.getButtons(runes.Each("456-")...)...),
				container.NewGridWithColumns(4, a.getButtons(runes.Each("123+")...)...),
				container.NewGridWithColumns(2,
					container.NewGridWithColumns(2,
						a.getButtons("0")[0],
						container.NewGridWithColumns(2, a.getButtons(runes.Each(".°")...)...),
					),
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
