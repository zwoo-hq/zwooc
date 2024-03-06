package config

import (
	"reflect"
	"testing"
)

func TestNewContext(t *testing.T) {
	tests := []struct {
		name string
		opts LoadOptions
		want loadingContext
	}{
		{
			name: "Test with empty LoadOptions",
			opts: LoadOptions{},
			want: loadingContext{
				skipHooks:    false,
				excludedKeys: []string{},
				extraArgs:    []string{},
				callStack:    []string{},
			},
		},
		{
			name: "Test with non-empty LoadOptions",
			opts: LoadOptions{
				SkipHooks: true,
				Exclude:   []string{"key1", "key2"},
				ExtraArgs: []string{"arg1", "arg2"},
			},
			want: loadingContext{
				skipHooks:    true,
				excludedKeys: []string{"key1", "key2"},
				extraArgs:    []string{"arg1", "arg2"},
				callStack:    []string{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewContext(tt.opts); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewContext() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoadingContext_getArgs(t *testing.T) {
	t.Run("should return args on first call", func(t *testing.T) {
		ctx := loadingContext{
			callStack: []string{},
			extraArgs: []string{"arg1", "arg2"},
		}
		if got := ctx.getArgs(); !reflect.DeepEqual(got, []string{"arg1", "arg2"}) {
			t.Errorf("getArgs() = %v, want %v", got, []string{"arg1", "arg2"})
		}
	})

	t.Run("should not return args on subsequent calls", func(t *testing.T) {
		ctx := loadingContext{
			callStack: []string{"caller1"},
			extraArgs: []string{"arg1", "arg2"},
		}
		if got := ctx.getArgs(); !reflect.DeepEqual(got, []string{}) {
			t.Errorf("getArgs() = %v, want %v", got, []string{})
		}
	})
}

func TestLoadingContext_withCaller(t *testing.T) {
	t.Run("should add caller to callStack", func(t *testing.T) {
		ctx := loadingContext{
			callStack: []string{},
		}
		ctx = ctx.withCaller("caller1")
		if !reflect.DeepEqual(ctx.callStack, []string{"caller1"}) {
			t.Errorf("withCaller() = %v, want %v", ctx.callStack, []string{"caller1"})
		}
	})
}

func TestLoadingContext_hasCaller(t *testing.T) {
	t.Run("should return true if caller is in callStack", func(t *testing.T) {
		ctx := loadingContext{
			callStack: []string{"a", "caller1", "b"},
		}
		if !ctx.hasCaller("caller1") {
			t.Errorf("hasCaller() = false, want true")
		}
	})

	t.Run("should return false if caller is not in callStack", func(t *testing.T) {
		ctx := loadingContext{
			callStack: []string{"caller1"},
		}
		if ctx.hasCaller("caller2") {
			t.Errorf("hasCaller() = true, want false")
		}
	})
}

func TestLoadingContext_excludes(t *testing.T) {
	t.Run("should return true if target is in excludedKeys", func(t *testing.T) {
		ctx := loadingContext{
			excludedKeys: []string{"key1"},
		}
		if !ctx.excludes("key1") {
			t.Errorf("excludes() = false, want true")
		}
	})

	t.Run("should return false if target is not in excludedKeys", func(t *testing.T) {
		ctx := loadingContext{
			excludedKeys: []string{"key1"},
		}
		if ctx.excludes("key2") {
			t.Errorf("excludes() = true, want false")
		}
	})
}
