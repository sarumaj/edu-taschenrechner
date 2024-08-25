/*
Package runes provides a set of functions to work with runes.
It features functions to split a string into a slice of strings where each contains just one rune,
check if a rune is within a given range, check if a rune is a digit, letter, or word character,
and more.

Sequences can be used to work with runes in a more programmer-friendly way.

Example:

	Each("Hello") // returns []string{"H", "e", "l", "l", "o"}
	IsDigit('1') // returns true
	IsLetter('a') // returns true
	IsWord('_') // returns true
	Sequence("Hello").BeginsWith("He") // returns true
	Sequence("Hello").Backspace() // removes the last rune
	Sequence("Hello").Delete() // removes the first rune
	Sequence("Hello").Clear() // removes all runes
	Sequence("Hello").Contains("He") // returns true
	Sequence("Hello").String() // returns "Hello"
	Sequence("Hello").Shift().String() // returns "Hello" without the first rune
	Sequence("Hello").Back().String() // returns "o"
	HowManyOpen(Sequence("Hello (world)")) // returns 1
	IsDotted(Sequence("Hello 123.45")) // returns true
	IsAnyOf('a', "abc") // returns true
*/
package runes

// Each splits a string into a slice of strings where each contains just one rune.
func Each(str string) (out []string) {
	for _, c := range str {
		out = append(out, string(c))
	}

	return out
}

// InRange checks if a rune is within a given range.
func InRange(r rune, min, max rune) bool {
	if min > max {
		min, max = max, min
	}

	return r >= min && r <= max
}

// IsAnyOf reports if given rune is within provided sequence.
func IsAnyOf(r rune, o string) bool {
	for _, c := range o {
		if r == c {
			return true
		}
	}

	return false
}

// IsDigit determines if a rune is a digit.
func IsDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

// IsLetter determines if a rune is a letter.
func IsLetter(r rune) bool {
	return InRange(r, 'a', 'z') || InRange(r, 'A', 'Z')
}

// IsWord determines if a rune is a word character.
func IsWord(r rune) bool {
	return IsDigit(r) || IsLetter(r) || r == '_'
}

// IsValid checks if the rune is in valid unicode point range 0 < x < 2_147_483_647.
func IsValid(r rune) bool {
	return r >= 0 && r < int32(1<<31-1)
}

// HowManyOpen returns the number of not closed brackets in the Sequence runes.
func HowManyOpen(text *Sequence) int {
	back, cnt := text.Back(), 0
	for _, b := back(); b >= 0; _, b = back() {
		switch b {
		case '(':
			cnt++
		case ')':
			cnt--
		}
	}

	return cnt
}

// IsDotted verifies the last consecutive sequence of digits contained in the Sequence runes
// upon the existence of a decimal floating point.
func IsDotted(text *Sequence) bool {
	back := text.Back()
	for _, b := back(); b > 0; _, b = back() {
		if b == '.' {
			return true
		}

		if !IsDigit(b) {
			break
		}
	}

	return false
}
