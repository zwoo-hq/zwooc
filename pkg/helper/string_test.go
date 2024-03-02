package helper

import "testing"

func TestRepeat(t *testing.T) {
	tests := []struct {
		name   string
		s      string
		count  int
		result string
	}{
		{"should repeat string 0 times", "foo", 0, ""},
		{"should repeat string 1 time", "foo", 1, "foo"},
		{"should repeat string 3 times", "foo", 3, "foofoofoo"},
		{"should repeat string 5 times", "foo", 5, "foofoofoofoofoo"},
		{"should not fail with negativ number", "foo", -2, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Repeat(tt.s, tt.count); got != tt.result {
				t.Errorf("Repeat(%s, %d) = %s, want %s", tt.s, tt.count, got, tt.result)
			}
		})
	}
}
