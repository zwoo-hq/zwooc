package config

import (
	"testing"
)

func TestNormalizeFragmentKey(t *testing.T) {
	tests := []struct {
		name        string
		fullKey     string
		wantKey     string
		wantMode    string
		wantProfile string
	}{
		{
			name:        "Test with one part",
			fullKey:     "key",
			wantKey:     "key",
			wantMode:    "",
			wantProfile: "",
		},
		{
			name:        "Test with two parts",
			fullKey:     "key:mode",
			wantKey:     "key",
			wantMode:    "mode",
			wantProfile: "",
		},
		{
			name:        "Test with three parts",
			fullKey:     "key:mode:profile",
			wantKey:     "key",
			wantMode:    "mode",
			wantProfile: "profile",
		},
		{
			name:        "Test with more than three parts",
			fullKey:     "key1:key2:mode:profile",
			wantKey:     "key1:key2",
			wantMode:    "mode",
			wantProfile: "profile",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotKey, gotMode, gotProfile := normalizeFragmentKey(tt.fullKey)
			if gotKey != tt.wantKey {
				t.Errorf("normalizeFragmentKey() gotKey = %v, want %v", gotKey, tt.wantKey)
			}
			if gotMode != tt.wantMode {
				t.Errorf("normalizeFragmentKey() gotMode = %v, want %v", gotMode, tt.wantMode)
			}
			if gotProfile != tt.wantProfile {
				t.Errorf("normalizeFragmentKey() gotProfile = %v, want %v", gotProfile, tt.wantProfile)
			}
		})
	}
}

func TestCombineFragmentKey(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		mode    string
		profile string
		want    string
	}{
		{
			name:    "Test with empty strings",
			key:     "",
			mode:    "",
			profile: "",
			want:    "::",
		},
		{
			name:    "Test with non-empty strings",
			key:     "key",
			mode:    "mode",
			profile: "profile",
			want:    "key:mode:profile",
		},
		{
			name:    "Test with mixed empty and non-empty strings",
			key:     "key",
			mode:    "",
			profile: "profile",
			want:    "key::profile",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := combineFragmentKey(tt.key, tt.mode, tt.profile); got != tt.want {
				t.Errorf("combineFragmentKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
