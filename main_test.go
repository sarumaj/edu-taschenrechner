package main

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/cucumber/godog"
	"github.com/sarumaj/edu-taschenrechner/pkg/ui"
)

// InitializeScenario is being used to set up the feature steps.
func InitializeScenario(sc *godog.ScenarioContext) {
	// used to store and retrieve values in context.Context
	type appInstanceKey string

	// used to verify test step input
	validButtonListRegex := regexp.MustCompile(`^(?:[0-9\+-×÷=\.]|AC|\(\)|<x)(?:\s(?:[0-9\+-×÷=\.]|AC|\(\))|<x)*$`)

	// static errors to be reused
	taschenrechnerNotReadyErr := errors.New("start Taschenrechner first")
	displayNotFoundErr := errors.New("display not found")

	// setup app instance and build window content
	sc.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		app := ui.NewApp("Test Taschenrechner")
		app.Build()

		return context.WithValue(ctx, appInstanceKey(""), app), nil
	})

	// when step describing user interaction: input through a sequence of button clicks
	sc.When(
		`^I press following buttons: "([^"]*)"$`,
		func(ctx context.Context, button_list string) (context.Context, error) {
			// check if our button_list is valid (buttons are supposed to be delimited with one white space)
			if !validButtonListRegex.MatchString(button_list) {
				return ctx, fmt.Errorf("the sequence %q does not match %q", button_list, validButtonListRegex)
			}

			// get app instances
			app, ok := ctx.Value(appInstanceKey("")).(*ui.App)
			if !ok {
				return ctx, taschenrechnerNotReadyErr
			}

			// retrieve buttons and mock up clicks
			for _, btnId := range strings.Split(button_list, " ") {
				object, ok := app.Objects[btnId]
				if !ok {
					// uppsala, something went wrong
					return ctx, fmt.Errorf("object %q not found", btnId)
				}

				// assert it is really a button with a defined click action
				btn, ok := object.(*ui.Button)
				if ok && btn.OnTapped != nil {
					btn.OnTapped() // click
					continue
				}

				// no valid button found
				return ctx, fmt.Errorf("invalid button %q", btnId)
			}

			return ctx, nil
		},
	)

	// then step describing the desired state of the display
	sc.Then(
		`^I get following result: "([^"]*)"$`,
		func(ctx context.Context, result string) (context.Context, error) {
			// retrieve app instance
			app, ok := ctx.Value(appInstanceKey("")).(*ui.App)
			if !ok {
				return ctx, taschenrechnerNotReadyErr
			}

			// get display used as display
			display, ok := app.Objects["display"].(*ui.Display)
			if !ok {
				return ctx, displayNotFoundErr
			}

			// compare display state with the desired one
			if result != display.Text {
				return ctx, fmt.Errorf("invalid result, got: %q, expected: %q", display.Text, result)
			}

			return ctx, nil
		},
	)

}

func TestExampleFor_BDT(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario, // step definitions & setup of test app instance
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{filepath.Join("go-test", "features")},
			TestingT: t, // Testing instance that will run sub-tests.
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
