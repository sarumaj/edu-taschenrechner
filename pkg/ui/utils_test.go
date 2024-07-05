package ui

import (
	"testing"

	"github.com/sarumaj/edu-taschenrechner/pkg/memory"
	"github.com/sarumaj/edu-taschenrechner/pkg/runes"
)

func TestExampleFor_modifyInput(t *testing.T) {
	setup := func(t testing.TB) (*runes.Input, memory.MemoryCell) {
		t.Helper()
		return runes.NewInput("_"), memory.NewMemoryCell()
	}

	type args struct {
		btnID           string
		requestNewSetup bool
	}

	input, memory := setup(t)

	for _, tt := range []struct {
		name string
		args args
		want string
	}{
		{"test#01", args{"+", false}, "_"},
		{"test#02", args{"-", false}, "-_"},
		{"test#03", args{"6", false}, "-6_"},
		{"test#04", args{"()", false}, "-6×(_"},
		{"test#05", args{"()", false}, "-6×((_"},
		{"test#06", args{"2", false}, "-6×((2_"},
		{"test#07", args{"-", false}, "-6×((2-_"},
		{"test#08", args{"-", false}, "-6×((2+_"},
		{"test#09", args{"-", false}, "-6×((2+-_"},
		{"test#10", args{"7", false}, "-6×((2+-7_"},
		{"test#11", args{"()", false}, "-6×((2+-7)_"},
		{"test#12", args{"0", false}, "-6×((2+-7)×0_"},
		{"test#13", args{".", false}, "-6×((2+-7)×0._"},
		{"test#14", args{"6", false}, "-6×((2+-7)×0.6_"},
		{"test#15", args{"=", false}, "18"},
		{"test#16", args{"+", false}, "ANS+_"},
		{"test#17", args{"9", false}, "ANS+9_"},
		{"test#18", args{"=", false}, "27"},
		{"test#19", args{"=", true}, "_"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.requestNewSetup {
				input, memory = setup(t)
			}

			modifyInput(input, memory, tt.args.btnID)
			got := input.String()

			if got != tt.want {
				t.Errorf(`modifyInput(%T, %T, %q) failed, got: %q, want: %q`, input, memory, tt.args.btnID, got, tt.want)
			}
		})
	}
}
