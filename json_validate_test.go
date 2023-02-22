package generics

import (
	"testing"
)

type RuneStack []rune

// Push adds e on top of the stack. The time complexity is O(1) amortized.
func (s *RuneStack) Push(e interface{}) { *s = append(*s, e.(rune)) }

// Pop removes and returns the last added rune element from this stack. The time complexity is O(1) amortized.
func (s *RuneStack) Pop() (e interface{}) {
	if s.Len() == 0 {
		return nil
	}

	e = (*s)[s.Len()-1]

	if cap(*s) > 64 && float64(s.Len()/cap(*s)) < 0.75 { // Free memory when the length of a slice shrunk enough.
		*s = append([]rune(nil), (*s)[:s.Len()-1]...)
		return e
	}

	*s = (*s)[:s.Len()-1]
	return e
}

func (s *RuneStack) Len() int { return len(*s) }

// IsWellFormed returns true is given string of brackets is well-formed. The time complexity is O(n), where n is the number of characters in s. The space complexity is O(n).
func IsWellFormed(input string) bool {
	stack := new(RuneStack)
	for _, c := range input {
		switch p, ok := map[rune]rune{'"': '"', '[': ']', '{': '}', ',': ','}[c]; {
		case ok:
			if p != ',' {
				stack.Push(p)
			}
		case stack.Len() == 0 || c != stack.Pop().(rune):
			return false
		}
	}

	if stack.Len() != 0 {
		return false
	}

	return true
}

func TestIsWellFormed(t *testing.T) {
	for _, test := range []struct {
		in   string
		want bool
	}{
		{"[]", true},
		{"{}", true},
		{"[{},[]]", true},
		{"[],{}", true},
		{"[[],{}]", true},
		{"}", false},
		{"[", false},
		{"{", false},
		{"]", false},
		{"[[],{", false},
		{"0}", false},
		{"[a", false},
	} {
		if got := IsWellFormed(test.in); got != test.want {
			t.Errorf("IsWellFormed(%q) = %t; want %t", test.in, got, test.want)
		}
	}
}
