package ui

import (
	"fmt"
	"log"
	"math"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/driver/mobile"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/sarumaj/edu-taschenrechner/pkg/cursor"
	"github.com/sarumaj/edu-taschenrechner/pkg/parser"
	"github.com/sarumaj/edu-taschenrechner/pkg/runes"
)

// Make sure the Display widget implements all necessary interfaces.
var _ interface {
	desktop.Cursorable
	desktop.Keyable
	fyne.Tappable
	fyne.DoubleTappable
	fyne.SecondaryTappable
	mobile.Touchable
} = (*Display)(nil)

// Display is a custom label widget that extends the default label with a memory cell.
type Display struct {
	widget.Entry
	parserOpts           []parser.Option
	MaximumContentLength int
}

// Cursor returns the default cursor.
func (*Display) Cursor() desktop.Cursor { return desktop.DefaultCursor }

// CopyToClipboard copies the text of the display widget to the clipboard.
func (display *Display) CopyToClipboard() {
	fyne.CurrentApp().Driver().AllWindows()[0].Clipboard().SetContent(strings.TrimSuffix(display.Text, "_"))
}

// GetParserOptions returns the parser options of the display widget.
func (display *Display) GetParserOptions() []parser.Option { return display.parserOpts }

// GetText returns the text of the display widget.
func (display *Display) GetText() string { return display.Entry.Text }

// GetOnChanged returns a function that sets the cursor to the end of the text.
func (display *Display) GetOnChanged() func(string) {
	return func(text string) {
		column := len([]rune(text))
		if strings.HasSuffix(text, "_") {
			column--
		}

		display.Entry.CursorColumn = column
		display.Entry.FocusGained()
		display.Entry.Refresh()
	}
}

// MeasureDisplayCapacity measures the display capacity interactively.
// It shows a dialog to ask the user if the result is visible in the display.
// If the measurement is completed, it shows an information dialog with the result.
func (display *Display) MeasureDisplayCapacity() {
	config := display.NewMeasurement("Display Capacity Measurement", 1, 1_000)
	dialog.ShowConfirm(config.title, "Do you want to measure the display capacity?", func(b bool) {
		if !b {
			config.Exit(true)
			return
		}

		config.Run()
	}, fyne.CurrentApp().Driver().AllWindows()[0])
}

// NewMeasurement creates a new measurement of the display capacity.
func (display *Display) NewMeasurement(title string, lo, hi int) *measurement {
	return &measurement{
		title:      title,
		display:    display,
		lowerBound: lo,
		upperBound: hi,
		current:    lo + int(math.Floor(float64(hi-lo)/2)),
		result:     display.MaximumContentLength,
	}
}

// SetMaximumContentLength sets the maximum content length of the display widget.
func (display *Display) SetMaximumContentLength(length int) *Display {
	display.MaximumContentLength = length
	return display
}

// SetOnChanged sets the OnChanged function of the display widget.
func (display *Display) SetOnChanged(fn func(string)) *Display {
	display.Entry.OnChanged = fn
	return display
}

// SetParserOptions sets the parser options of the display widget.
func (display *Display) SetParserOptions(options ...parser.Option) *Display {
	display.parserOpts = options
	return display
}

// SetText sets the text of the display widget.
// It moves the cursor to the end of the text and checks the state of the cursor.
// If the cursor is in an invalid state, it shows an error dialog.
func (display *Display) SetText(text string) {
	window := fyne.CurrentApp().Driver().AllWindows()[0]

	if text == "=" && Interactive { // Display waiting dialog when calculating
		dialog.ShowCustomWithoutButtons(
			"Calculating",
			widget.NewLabel(fmt.Sprintf("Evaluating: %q", strings.TrimSuffix(display.Text, "_"))),
			window,
		)
	}

	// define a channel for synchronization
	var sync chan struct{}
	if !Interactive {
		sync = make(chan struct{})
	}

	// run the cursor operation in a separate goroutine
	go func() {
		// create a new cursor
		textCursor := cursor.New(runes.NewSequence(display.Text), time.Minute, display.parserOpts...)

		// perform the operation on the cursor
		result := textCursor.Do(strings.TrimSpace(text)).String()

		// on first exceedance of the maximum content length, show the current value in scientific notation
		if display.MaximumContentLength > 0 && len(result) > display.MaximumContentLength {
			// exploit the capability of the memory cell to display the result in scientific notation
			result = textCursor.EqualsWithFormat('g').String()
		}

		// on second exceedance of the maximum content length, truncate the result
		if display.MaximumContentLength > 0 && len(result) > display.MaximumContentLength {
			result = result[:display.MaximumContentLength-3] + "..."
		}

		// perform the calculation and set the result
		display.Entry.SetText(result)

		// schedule an update of the entry on the main thread
		windowCanvas := window.Canvas()
		windowCanvas.Refresh(&display.Entry)

		// remove all dialog overlays from the canvas
		for _, overlay := range windowCanvas.Overlays().List() {
			log.Printf("%T\n", overlay)
			switch overlay.(type) {
			case dialog.Dialog, *widget.PopUp:
				windowCanvas.Overlays().Remove(overlay)
			}

		}

		// check the state of the cursor and show an error dialog if needed
		if err := textCursor.Check(); err != nil && Interactive {
			dialog.ShowError(err, window)
		}

		if !Interactive {
			sync <- struct{}{}
		}
	}()

	if !Interactive {
		<-sync
	}
}

// Overwrite methods to prevent the display from being editable
func (*Display) DoubleTapped(*fyne.PointEvent)    {}
func (*Display) KeyDown(*fyne.KeyEvent)           {}
func (*Display) KeyUp(*fyne.KeyEvent)             {}
func (*Display) MouseDown(*desktop.MouseEvent)    {}
func (*Display) MouseUp(*desktop.MouseEvent)      {}
func (*Display) Tapped(*fyne.PointEvent)          {}
func (*Display) TappedSecondary(*fyne.PointEvent) {}
func (*Display) TouchCancel(*mobile.TouchEvent)   {}
func (*Display) TouchDown(*mobile.TouchEvent)     {}
func (*Display) TouchUp(*mobile.TouchEvent)       {}
func (*Display) TypedKey(*fyne.KeyEvent)          {}
func (*Display) TypedRune(rune)                   {}
func (*Display) TypedShortcut(fyne.Shortcut)      {}

// measurement is a struct that represents a measurement configuration of the display capacity.
type measurement struct {
	title                           string
	display                         *Display
	exit                            bool
	lowerBound, upperBound, current int
	result                          int
}

// backward decreases the upper bound of the measurement.
func (m *measurement) backward() {
	m.upperBound = m.current - 1 // Decrease the upper bound
	m.current = m.lowerBound + int(math.Floor(float64(m.upperBound-m.lowerBound)/2))
}

// forward increases the lower bound of the measurement.
func (m *measurement) forward() {
	m.lowerBound = m.current + 1 // Increase the lower bound
	m.current = m.lowerBound + int(math.Floor(float64(m.upperBound-m.lowerBound)/2))
}

// Exit sets the exit flag of the measurement and adjusts the display capacity.
func (m *measurement) Exit(exit bool) {
	m.exit = exit
	if m.exit {
		m.display.MaximumContentLength = m.result
	}
}

// Run runs the measurement of the display capacity interactively.
// It shows a dialog to ask the user if the result is visible in the display.
// If the measurement is completed, it shows an information dialog with the result.
func (m *measurement) Run() {
	window := fyne.CurrentApp().Driver().AllWindows()[0]
	if m.exit || m.lowerBound > m.upperBound { // Measurement aborted or completed
		m.display.MaximumContentLength = m.result
		m.display.Entry.SetText("_") // Clear the display
		dialog.ShowInformation(
			m.title,
			fmt.Sprintf("Measured maximum displayable sequence length of %d.", m.result),
			window,
		)
		return
	}

	// Create a new cursor
	textCursor := cursor.New(runes.NewSequence("_"), time.Minute, m.display.parserOpts...)

	// Submit the current value to the display
	for _, ch := range fmt.Sprintf("%d", m.current) {
		textCursor.Do(string(ch))
	}

	// Calculate the factorial of the current value and submit the result to the display
	m.display.Entry.SetText(textCursor.Factorial().Equals().String())
	m.result = len(m.display.Entry.Text) // Get the length of the result

	// Create a new dialog to ask the user if the result is visible in the display
	ask := dialog.NewCustomWithoutButtons(
		m.title,
		&widget.Label{
			Text:      fmt.Sprintf("Can you see the result in the display?\n(Current content length: %d)", m.result),
			Wrapping:  fyne.TextWrapWord,
			Alignment: fyne.TextAlignCenter,
		},
		window,
	)

	// Define the buttons of the dialog
	ask.SetButtons([]fyne.CanvasObject{
		widget.NewButton("Yes", func() {
			m.forward()
			ask.Hide()
			m.Run()
		}),
		widget.NewButton("No", func() {
			m.backward()
			ask.Hide()
			m.Run()
		}),
		widget.NewButton("Abort", func() {
			m.exit = true
			ask.Hide()
			m.Run()
		}),
	})

	ask.Show()
}

// NewDisplay creates a new label widget with the given text and memory cell.
func NewDisplay(text string, options ...parser.Option) *Display {
	display := &Display{
		Entry: widget.Entry{
			Wrapping:  fyne.TextWrapOff,
			Text:      text,
			TextStyle: fyne.TextStyle{Monospace: true},
			Scroll:    container.ScrollHorizontalOnly,
		},
		MaximumContentLength: 500,
		parserOpts:           options,
	}

	display.Entry.ActionItem = NewToolbarItem(theme.SettingsIcon()).SetOnTapped(display.MeasureDisplayCapacity)
	display.SetOnChanged(display.GetOnChanged())
	display.ExtendBaseWidget(display)
	display.Entry.FocusGained()
	display.Entry.Refresh()

	return display
}
