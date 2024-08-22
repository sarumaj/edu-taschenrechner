package runes

import "testing"

func TestExampleFor_Sequence(t *testing.T) {
	input := NewSequence("")

	t.Run("Append", func(t *testing.T) {
		args := "world!"
		want := "world!"

		input.Append(args)

		if got := input.String(); got != want {
			t.Errorf(`(*Input).Append(%q) failed, got: %q, want: %q`, args, got, want)
		}
	})

	t.Run("Prepend", func(t *testing.T) {
		args := "Hello, "
		want := "Hello, world!"

		input.Prepend(args)

		if got := input.String(); got != want {
			t.Errorf(`(*Input).Prepend(%q) failed, got: %q, want: %q`, args, got, want)
		}
	})

	t.Run("Backspace", func(t *testing.T) {
		want := "Hello, world"

		input.Backspace()

		if got := input.String(); got != want {
			t.Errorf(`(*Input).Backspace() failed, got: %q, want: %q`, got, want)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		// cspell: disable-next-line
		want := "ello, world"

		input.Delete()

		if got := input.String(); got != want {
			t.Errorf(`(*Input).Delete() failed, got: %q, want: %q`, got, want)
		}
	})

	t.Run("Shift", func(t *testing.T) {
		// cspell: disable-next-line
		want := "ello, worl"

		got := input.Shift().String()

		if got != want {
			t.Errorf(`(*Input).Shift() failed, got: %q, want: %q`, got, want)
		}

		if got == input.String() {
			t.Errorf(`(*Input).Shift() failed, the input has been unexpectedly affected`)
		}
	})
}
