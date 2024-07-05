package runes

// Each splits a string into a slice of strings where each contains just one rune.
func Each(str string) (out []string) {
	for _, c := range str {
		out = append(out, string(c))
	}

	return out
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
	return r <= '9' && r >= '0'
}

// IsValid checks if the rune is in valid unicode point range 0 < x < 2_147_483_647.
func IsValid(r rune) bool {
	return r >= 0 && r < int32(1<<31-1)
}

// HowManyOpen returns the number of not closed brackets in the input runes.
func HowManyOpen(text *Input) int {
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

// IsDotted verifies the last consecutive sequence of digits contained in the input runes
// upon the existence of a decimal floating point.
func IsDotted(text *Input) bool {
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
