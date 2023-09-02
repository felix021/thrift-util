package buffer

import "testing"

func TestIsPowerOf2(t *testing.T) {
	type Case struct {
		n        int
		expected bool
	}
	cases := []Case{
		{0, false},
		{1, true},
		{2, true},
		{3, false},
		{4, true},
		{5, false},
	}
	for _, c := range cases {
		if got := IsPowerOf2(c.n); got != c.expected {
			t.Errorf("IsPowerOf2(%d) = %v, != %v", c.n, got, c.expected)
		}
	}
}

func TestGetExponent(t *testing.T) {
	type Case struct {
		n        int
		expected int
	}
	cases := []Case{
		{0, -1},
		{1, 0},
		{2, 1},
		{3, -1},
		{4, 2},
		{8, 3},
	}
	for _, c := range cases {
		if got := GetExponent(c.n); got != c.expected {
			t.Errorf("GetExponent(%d) = %v, != %v", c.n, got, c.expected)
		}
	}
}
