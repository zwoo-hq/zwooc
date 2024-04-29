package helper

import (
	"testing"
)

func TestFindBy(t *testing.T) {
	tests := []struct {
		name      string
		slice     []int
		predicate func(int) bool
		result    int
		found     bool
	}{
		{"should find element", []int{1, 2, 3}, func(i int) bool { return i == 2 }, 2, true},
		{"should not find element", []int{1, 2, 3}, func(i int) bool { return i == 4 }, 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, found := FindBy(tt.slice, tt.predicate)
			if found != tt.found || (got != nil && *got != tt.result) {
				t.Errorf("FindBy() = %v, %v, want %v, %v", got, found, tt.result, tt.found)
			}
		})
	}
}

func slicesAreEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func TestSlicesAreEqual(t *testing.T) {
	tests := []struct {
		name   string
		a      []int
		b      []int
		result bool
	}{
		{"should return true for empty slices", []int{}, []int{}, true},
		{"should return true for equal slices", []int{1, 2, 3}, []int{1, 2, 3}, true},
		{"should return false for different slices", []int{1, 2, 3}, []int{1, 2, 4}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := slicesAreEqual(tt.a, tt.b); got != tt.result {
				t.Errorf("slicesAreEqual() = %v, want %v", got, tt.result)
			}
		})
	}
}

func TestConcat(t *testing.T) {
	tests := []struct {
		name   string
		slices [][]int
		result []int
	}{
		{"should concat empty slices", [][]int{}, []int{}},
		{"should concat two slices", [][]int{{1, 2}, {3, 4}}, []int{1, 2, 3, 4}},
		{"should concat three slices", [][]int{{1, 2}, {3, 4}, {5, 6}}, []int{1, 2, 3, 4, 5, 6}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Concat(tt.slices...); !slicesAreEqual(got, tt.result) {
				t.Errorf("Concat() = %v, want %v", got, tt.result)
			}
		})
	}
}

func TestMapTo(t *testing.T) {
	tests := []struct {
		name   string
		slice  []int
		mapper func(int) int
		result []int
	}{
		{"should map slices", []int{1, 3, 5}, func(in int) int { return in + 1 }, []int{2, 4, 6}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MapTo(tt.slice, tt.mapper); !slicesAreEqual(got, tt.result) {
				t.Errorf("MapTo() = %v, want %v", got, tt.result)
			}
		})
	}
}
