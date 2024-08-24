package cursor

import (
	"testing"

	"github.com/sarumaj/edu-taschenrechner/pkg/memory"
	"github.com/sarumaj/edu-taschenrechner/pkg/runes"
)

func TestExampleFor_Cursor(t *testing.T) {
	get := func() Cursor {
		return New(runes.NewSequence(""), memory.NewMemoryCell())
	}

	for _, tt := range []struct {
		name string
		args Cursor
		want string
	}{
		{"test#01", get().Clear().Minus().Three().Brackets().Minus().Nine().Zero().Two().Two(), "-3×(-9022_"},
		{"test#02", get().Brackets().Divide().Delete().Eight().Divide().Six().Eight().Minus().Minus().Five(), "ANS×8÷68+5_"},
		{"test#03", get().Nine().Eight().Eight().Three().DecimalPoint().Six().Two().Clear().Four().Brackets(), "4×(_"},
		{"test#04", get().Seven().Two().Plus().DecimalPoint().Plus().Divide().Equals().Minus().Delete().Four(), "ANS×724_"},
		{"test#05", get().Four().Brackets().Four().Clear().DecimalPoint().Six().Brackets().Plus().Five().Times(), "6×(5×_"},
		{"test#06", get().Clear().One().Zero().Eight().Five().Six().DecimalPoint(), "10856._"},
		{"test#07", get().Two().Brackets().Zero().Nine().Zero().Nine().Seven().Divide().Minus(), "ANS×2×(09097÷-_"},
		{"test#08", get().Four().Delete().One().One().Divide().Two().Nine().Times().Equals(), "ANS×11÷29_"},
		{"test#09", get().Equals().Delete().Zero().One().Two().Clear().Equals().Times().Five(), "5_"},
		{"test#10", get().Three().Brackets().Equals().Nine().Two().Four().Zero().One().Times(), "ANS×392401×_"},
		{"test#11", get().Six().Two().Equals().Seven().Divide().Four().Five().One().One().Brackets(), "7÷4511×(_"},
		{"test#12", get().Brackets().Equals().Nine().Nine().Divide().Clear().Six().Delete().Clear().Two(), "2_"},
		{"test#13", get().Clear().One().Seven().One().Equals().Times().Times().Zero().Zero(), "ANS×00_"},
		{"test#14", get().Equals().Equals().Brackets().Brackets().Six().Plus().Equals(), "((6_"},
		{"test#15", get().Times().Delete().Equals().Six().Brackets().Four().Six().Four().Clear().Five(), "5_"},
		{"test#16", get().Delete().Brackets().Seven().Delete().Zero().Seven().Minus().Clear(), "_"},
		{"test#17", get().Zero().Two().Seven().Five().DecimalPoint().Three().Brackets().Seven().Times(), "ANS×0275.3×(7×_"},
		{"test#18", get().One().Minus().Zero().Delete().Minus().One().Equals().Brackets().Two().Three(), "(23_"},
		{"test#19", get().Six().Plus().Five().Eight().Minus().Clear().Clear().Seven().Five(), "75_"},
		{"test#20", get().Nine().DecimalPoint().Six().Three().DecimalPoint().DecimalPoint().Plus().Three().Seven(), "ANS×9.63+37_"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.String(); got != tt.want {
				t.Errorf("Cursor.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExampleFor_Do(t *testing.T) {
	setup := func(t testing.TB) (*runes.Sequence, memory.MemoryCell) {
		t.Helper()
		return runes.NewSequence("_"), memory.NewMemoryCell()
	}

	type args struct {
		operator        string
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
		{"test#09", args{"-", false}, "-6×((2-_"},
		{"test#10", args{"7", false}, "-6×((2-7_"},
		{"test#11", args{"()", false}, "-6×((2-7)_"},
		{"test#12", args{"0", false}, "-6×((2-7)×0_"},
		{"test#13", args{".", false}, "-6×((2-7)×0._"},
		{"test#14", args{"6", false}, "-6×((2-7)×0.6_"},
		{"test#15", args{"=", false}, "-6×((2-7)×0.6)_"},
		{"test#16", args{"+", false}, "-6×((2-7)×0.6)+_"},
		{"test#17", args{"9", false}, "-6×((2-7)×0.6)+9_"},
		{"test#18", args{"=", false}, "27"},
		{"test#19", args{"=", true}, "_"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.requestNewSetup {
				input, memory = setup(t)
			}

			Do(tt.args.operator, input, memory)
			got := input.String()

			if got != tt.want {
				t.Errorf(`Evaluate(%q, %T, %T) failed, got: %q, want: %q`, tt.args.operator, input, memory, got, tt.want)
			}
		})
	}
}
