package config

import (
	"testing"

	"github.com/zwoo-hq/zwooc/pkg/model"
)

func TestIsReserved(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{"$default should be true", model.KeyDefault, true},
		{"$adapter should be true", model.KeyAdapter, true},
		{"$compound should be true", model.KeyCompound, true},
		{"$fragment should be true", model.KeyFragment, true},
		{"$post should be true", model.KeyPost, true},
		{"$pre should be true", model.KeyPre, true},
		{"default should be false", "default", false},
		{"foo should be false", "foo", false},
		{"x$default should be false", "x$default", false},
		{"$schema should be true", "$schema", true},
		{"$dir should be true", "$dir", true},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if got := IsReservedKey(tt.value); got != tt.want {
				t.Errorf("IsReservedKey(%s) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestIsValidRunMode(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{"run should be true", "run", true},
		{"build should be true", "build", true},
		{"watch should be false", "watch", true},
		{"xxx should be false", "xxx", false},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if got := IsValidRunMode(tt.value); got != tt.want {
				t.Errorf("IsValidRunMode(%s) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}
