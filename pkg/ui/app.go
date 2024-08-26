//go:build !headless

/*
Package ui provides a user interface for the calculator.
The calculator can be used to perform basic arithmetic operations.

Example:

	app := ui.NewApp("Calculator")
	app.Build()
	app.ShowAndRun()
*/
package ui

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/sarumaj/edu-taschenrechner/pkg/calc"
	"github.com/sarumaj/edu-taschenrechner/pkg/memory"
	"github.com/sarumaj/edu-taschenrechner/pkg/parser"
	"github.com/sarumaj/edu-taschenrechner/pkg/runes"
)

const (
	appID        = "com.github.sarumaj.edu-taschenrechner"
	githubLink   = "https://github.com/sarumaj/edu-taschenrechner"
	linkedinLink = "https://www.linkedin.com/in/dawid-ciepiela"
)

// Interactive is a flag to enable/disable interactive mode.
// Used for tests to disable user interaction.
var Interactive bool = true

// App is a custom interface for bringing up the application window.
type App struct {
	sync.Once
	fyne.App
	fyne.Window
	fyne.Theme
	memory.MemoryCell
	changeListener chan fyne.Settings
	objects        ObjectStorage
}

// Build renders the application window and sets up all widgets.
func (a *App) Build() {
	a.Do(func() {
		// define options for parser
		options := []parser.Option{
			parser.WithVar("ANS", a.MemoryCell.Get),
			parser.WithConst("PI", big.NewFloat(math.Pi)),
			parser.WithConst("E", big.NewFloat(math.E)),
			parser.WithFunc("save", func(arg *big.Float) (*big.Float, error) {
				if err := a.MemoryCell.Set(arg); err != nil {
					return nil, err
				}
				return a.MemoryCell.Get(), nil
			}),
			parser.WithFunc("sin", math.Sin),
			parser.WithFunc("cos", math.Cos),
			parser.WithFunc("tan", math.Tan),
			parser.WithFunc("arcsin", func(f float64) (float64, error) {
				if f < -1 || f > 1 {
					return 0, fmt.Errorf("arcsin(%g) is undefined", f)
				}
				return math.Asin(f), nil
			}),
			parser.WithFunc("arccos", func(f float64) (float64, error) {
				if f < -1 || f > 1 {
					return 0, fmt.Errorf("arccos(%g) is undefined", f)
				}
				return math.Acos(f), nil
			}),
			parser.WithFunc("arctan", math.Atan),
			parser.WithFunc("log", func(f float64) (float64, error) {
				if f <= 0 {
					return 0, fmt.Errorf("log(%g) is undefined", f)
				}
				return math.Log10(f), nil
			}),
			parser.WithFunc("ln", func(f float64) (float64, error) {
				if f <= 0 {
					return 0, fmt.Errorf("ln(%g) is undefined", f)
				}
				return math.Log(f), nil
			}),
			parser.WithFunc("gdc", func(x, y *big.Float) (*big.Float, error) {
				return calc.GreatestCommonDivisor(context.Background(), x, y)
			}),
			parser.WithFunc("lcm", func(x, y *big.Float) (*big.Float, error) {
				return calc.LeastCommonMultiple(context.Background(), x, y)
			}),
			parser.WithReplacements("×", "*", "÷", "/", "π", "PI", "e", "E"),
		}

		// make display using options
		a.objects["display"] = NewDisplay("_", options...)

		// make buttons (some with alternate text)
		for _, btnText := range append(runes.Each("1234567890+-×÷=.π!e°√"),
			"xⁿ", "AC", "()", "↩",
			"sin", "cos", "tan", "log", "ln", "gdc") {

			a.objects[btnText] = NewButton(btnText, a.objects.SelectDisplay("display")).
				SetAlternateText(map[string]string{
					"√":   "x²",
					"sin": "sin⁻¹",
					"cos": "cos⁻¹",
					"tan": "tan⁻¹",
					"log": "10ⁿ",
					"ln":  "eⁿ",
					"°":   "1/°",
					".":   ",",
					"gdc": "lcm",
				}[btnText])
		}

		// make dropdowns
		for name, relations := range map[string][]string{
			"const": {"π", "e"},
			"func":  {"sin", "cos", "tan", "log", "ln", "gdc"},
		} {
			a.objects[name] = NewButtonDropDown(a.objects.SelectButtons(relations...))
		}

		// make toolbars
		var actions []widget.ToolbarItem
		for link, resources := range map[string][]fyne.Resource{
			githubLink:   {resourceGithubPng, resourceGithubWhitePng},
			linkedinLink: {resourceLinkedinPng, nil},
		} {
			actions = append(actions, NewIcon(link, resources[0], resources[1]).Update())
		}
		a.objects["toolbar"] = NewDisplayToolbar(a.objects.SelectDisplay("display"), actions...)

		// make invert button
		a.objects["INV"] = NewButton("INV", a.objects.SelectDisplay("display")).SetOnTapped(func() {
			for i := range a.objects {
				switch o := a.objects[i].(type) {
				case *Button:
					o.Invert()

				case *ButtonDropDown:
					// Important: call with defer, since the update must occur after all buttons have been inverted
					defer o.Update()

				}
			}
		})

		// lay out the components
		var appContent fyne.CanvasObject = container.NewGridWithRows(7,
			container.NewVBox(
				container.NewBorder(
					nil, nil, nil, a.objects["toolbar"],
					a.objects["display"],
				),
				widget.NewSeparator(),
			),
			container.NewGridWithColumns(3,
				container.NewGridWithColumns(2, a.objects.SelectCanvasObjects("√", "xⁿ")...),
				container.NewGridWithColumns(2, append(a.objects.SelectCanvasObjects("!"), a.objects["const"])...),
				a.objects["func"],
			),
			container.NewGridWithColumns(2,
				container.NewGridWithColumns(3, a.objects.SelectCanvasObjects("INV", "AC", "↩")...),
				container.NewGridWithColumns(2, a.objects.SelectCanvasObjects("()", "÷")...),
			),
			container.NewGridWithColumns(4, a.objects.SelectCanvasObjects(runes.Each("789×")...)...),
			container.NewGridWithColumns(4, a.objects.SelectCanvasObjects(runes.Each("456-")...)...),
			container.NewGridWithColumns(4, a.objects.SelectCanvasObjects(runes.Each("123+")...)...),
			container.NewGridWithColumns(2,
				container.NewGridWithColumns(2,
					a.objects["0"],
					container.NewGridWithColumns(2, a.objects.SelectCanvasObjects(runes.Each(".°")...)...),
				),
				a.objects["="],
			),
		)

		// add title and subtitle if running in a browser
		if a.Driver().Device().IsBrowser() {
			title := canvas.NewText("Taschenrechner", nil)
			title.Alignment = fyne.TextAlignCenter
			title.TextStyle = fyne.TextStyle{Bold: true}
			title.TextSize = 52
			a.objects["title"] = title

			subTitle := canvas.NewText("A simple calculator app", nil)
			subTitle.Alignment = fyne.TextAlignCenter
			subTitle.TextSize = 32
			a.objects["subTitle"] = subTitle

			appContent = container.NewBorder(container.NewVBox(title, subTitle, widget.NewSeparator()), nil, nil, nil, appContent)
		}

		// fill in the application window
		a.Window.SetContent(container.NewCenter(appContent))
	})
}

// ShowAndRun displays the application window and starts the event loop.
func (a *App) ShowAndRun() {
	go func() {
		for range a.changeListener {
			for _, o := range a.objects {
				switch o := o.(type) {
				case *Icon:
					o.Update()
				}
			}
		}
	}()

	a.Window.ShowAndRun()
}

// Objects returns the objects of the application.
func (a *App) Objects() map[string]fyne.CanvasObject {
	out := make(map[string]fyne.CanvasObject)
	for k, v := range a.objects {
		out[k] = v
	}
	return out
}

// Create new application window.
// Call ShowAndRun to display the window.
func NewApp(title string) *App {
	a := app.NewWithID(appID)
	a.SetIcon(resourceAppIco)

	w := a.NewWindow(title)
	w.Resize(fyne.NewSize(600, 400))

	i := &App{
		App:            a,
		changeListener: make(chan fyne.Settings),
		MemoryCell:     memory.NewMemoryCell(),
		Window:         w,
		objects:        make(ObjectStorage),
	}

	a.Settings().AddChangeListener(i.changeListener)
	a.Settings().SetTheme(NewDoubleSizeTheme(theme.DefaultTheme()))

	return i
}
