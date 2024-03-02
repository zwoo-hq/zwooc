package helper

import "testing"

func TestBuildName(t *testing.T) {
	tests := []struct {
		name   string
		parts  []string
		result string
	}{
		{"should return empty string for empty parts", []string{}, ""},
		{"should return single part", []string{"foo"}, "foo"},
		{"should return multiple parts", []string{"foo", "bar", "baz"}, "foo/bar/baz"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BuildName(tt.parts...); got != tt.result {
				t.Errorf("BuildName(%v) = %s, want %s", tt.parts, got, tt.result)
			}
		})
	}
}
