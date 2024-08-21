package runes

// Input is used to ease processing user interactions
type Input []rune

// Append writes content at the end.
func (i *Input) Append(r string) {
	for _, c := range r {
		*i = append(*i, c)
	}
}

// Back provides a closure to iterate over input runes in reserved order.
func (i Input) Back() func() (int, rune) {
	current := len(i) - 1

	return func() (int, rune) {

		if current >= 0 {
			defer func() { current-- }()
			return current, i[current]
		}

		return -1, -1
	}
}

// Backspace removes last rune.
func (i *Input) Backspace() {
	if len(*i) > 0 {
		*i = (*i)[: len(*i)-1 : len(*i)-1]
	}
}

// BeginsWith checks whether input runes begin with given sequence.
func (i Input) BeginsWith(r string) bool {
	if len(i) < len([]rune(r)) {
		return false
	}

	for j, c := range r {
		if i[j] != c {
			return false
		}
	}

	return true
}

// Clear removes all runes one by one.
func (i *Input) Clear() {
	for len(*i) > 0 {
		*i = (*i)[: len(*i)-1 : len(*i)-1]
	}
}

// Delete removes the first rune in the input runes.
func (i *Input) Delete() {
	if len(*i) > 0 {
		*i = (*i)[1:len(*i):len(*i)]
	}
}

// Contains checks whether the input runes contain given sequence.
func (i Input) Contains(r string) bool {
	return i.Index(r) >= 0
}

// EndsWith checks whether the input runes end with given sequence.
func (i Input) EndsWith(r string) bool {
	if len(i) < len([]rune(r)) {
		return false
	}

	offset := len(i) - len([]rune(r))
	for j, c := range r {
		if i[offset+j] != c {
			return false
		}
	}

	return true
}

// Equals checks
func (i Input) Equals(r string) bool {
	if len(i) != len([]rune(r)) {
		return false
	}

	for j, c := range r {
		if c != i[j] {
			return false
		}
	}

	return true
}

// First returns the foremost rune from the input runes.
// If input runes are empty, it returns -1.
func (i Input) First() rune {
	if len(i) > 0 {
		return i[0]
	}

	return -1
}

// Forward returns a closure to iterate over input runes from the beginning to the end.
func (i Input) Forward() func() (int, rune) {
	current := 0

	return func() (int, rune) {

		if current < len(i) {
			defer func() { current++ }()
			return current, i[current]
		}

		return -1, -1
	}
}

// Index returns the first index at which given sequence occurs.
// If the input runes do not contain the sequence, it returns -1.
func (i Input) Index(r string) int {
	if len(i) < len([]rune(r)) {
		return -1
	}

	for k, d := range i {
		for j, c := range r {

			if j == 0 && c == d {
				if (i[k:]).BeginsWith(r) {
					return k
				}

			}

		}
	}

	return -1
}

// Last returns the last rune from the input runes.
// If input runes are empty, it returns -1.
func (i Input) Last() rune {
	if len(i) > 0 {
		return i[len(i)-1]
	}

	return -1
}

// LastIndex returns the last index at which given sequence occurs.
// If the input runes do not contain the sequence, it returns -1.
func (i Input) LastIndex(r string) int {
	if len(i) < len([]rune(r)) {
		return -1
	}

	for k := len(i) - 1; k > 0; k++ {
		d := i[k]

		for j, c := range r {
			if j == 0 && c == d {
				if (i[k:]).BeginsWith(r) {
					return k
				}

			}

		}
	}

	return -1
}

// Prepend adds sequence at the beginning of the input sequence.
func (i *Input) Prepend(r string) {
	p := make([]rune, len([]rune(r)))
	_ = copy(p, []rune(r))
	*i = append(p, (*i)...)
}

// Shift moves one rune from the end of the input runes and returns a shortened sequence of runes.
func (i Input) Shift() *Input {
	var v Input
	if len(i) > 0 {
		v = make(Input, len(i)-1)
		_ = copy(v, i[:len(i)-1:len(i)-1])
	}
	return &v
}

// String converts input runes into a string.
func (i Input) String() string {
	return string([]rune(i))
}

// NewInput creates new input runes from a string.
func NewInput(str string) *Input {
	i := Input([]rune(str))
	return &i
}
