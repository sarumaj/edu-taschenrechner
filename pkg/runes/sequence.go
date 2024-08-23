package runes

// Sequence is used to ease processing user interactions
type Sequence []rune

// Append writes content at the end.
func (i *Sequence) Append(r string) {
	for _, c := range r {
		*i = append(*i, c)
	}
}

// Back provides a closure to iterate over runes in reserved order.
func (i Sequence) Back() func() (int, rune) {
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
func (i *Sequence) Backspace() {
	if len(*i) > 0 {
		*i = (*i)[: len(*i)-1 : len(*i)-1]
	}
}

// BeginsWith checks whether runes begin with given sequence.
func (i Sequence) BeginsWith(r string) bool {
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
func (i *Sequence) Clear() {
	for len(*i) > 0 {
		*i = (*i)[: len(*i)-1 : len(*i)-1]
	}
}

// Delete removes the first rune in the runes.
func (i *Sequence) Delete() {
	if len(*i) > 0 {
		*i = (*i)[1:len(*i):len(*i)]
	}
}

// Contains checks whether the runes contain given sequence.
func (i Sequence) Contains(r string) bool {
	return i.Index(r) >= 0
}

// EndsWith checks whether the runes end with given sequence.
func (i Sequence) EndsWith(r string) bool {
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
func (i Sequence) Equals(r string) bool {
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

// First returns the foremost rune from the runes.
// If runes are empty, it returns -1.
func (i Sequence) First() rune {
	if len(i) > 0 {
		return i[0]
	}

	return -1
}

// Forward returns a closure to iterate over runes from the beginning to the end.
func (i Sequence) Forward() func() (int, rune) {
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
// If the runes do not contain the sequence, it returns -1.
func (i Sequence) Index(r string) int {
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

// Last returns the last rune from the runes.
// If runes are empty, it returns -1.
func (i Sequence) Last() rune {
	if len(i) > 0 {
		return i[len(i)-1]
	}

	return -1
}

// LastIndex returns the last index at which given sequence occurs.
// If the runes do not contain the sequence, it returns -1.
func (i Sequence) LastIndex(r string) int {
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

// Prepend adds sequence at the beginning of the sequence.
func (i *Sequence) Prepend(r string) {
	p := make([]rune, len([]rune(r)))
	_ = copy(p, []rune(r))
	*i = append(p, (*i)...)
}

// Shift moves one rune from the end of the runes and returns a shortened sequence of runes.
func (i Sequence) Shift() *Sequence {
	var v Sequence
	if len(i) > 0 {
		v = make(Sequence, len(i)-1)
		_ = copy(v, i[:len(i)-1:len(i)-1])
	}
	return &v
}

// String converts runes into a string.
func (i Sequence) String() string {
	return string([]rune(i))
}

// NewSequence creates new runes from a string.
func NewSequence(str string) *Sequence {
	i := Sequence([]rune(str))
	return &i
}
