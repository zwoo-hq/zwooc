package config

import (
	"testing"
)

func TestIsReserved(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{"$default should be true", KeyDefault, true},
		{"$adapter should be true", KeyAdapter, true},
		{"$compound should be true", KeyCompound, true},
		{"$fragment should be true", KeyFragment, true},
		{"$post should be true", KeyPost, true},
		{"$pre should be true", KeyPre, true},
		{"deafult should be false", "default", false},
		{"foo should be false", "foo", false},
		{"x$default should be false", "x$default", false},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if got := IsReservedKey(tt.value); got != tt.want {
				t.Errorf("IsReservedKey(%s) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}
