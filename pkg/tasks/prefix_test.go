package tasks

import (
	"bytes"
	"testing"
)

func TestCommandPrefixer_Write(t *testing.T) {
	// Test writing an empty slice
	t.Run("empty slice", func(t *testing.T) {
		var dest bytes.Buffer
		prefixer := NewPrefixer("prefix: ", &dest)
		n, err := prefixer.Write([]byte(""))
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if n != 0 {
			t.Errorf("Expected 0 bytes written, got %d", n)
		}
		if dest.String() != "" {
			t.Errorf("Expected empty string, got %s", dest.String())
		}
	})

	// Test writing a slice with multiple lines
	t.Run("multiple lines", func(t *testing.T) {
		var dest bytes.Buffer
		prefixer := NewPrefixer("prefix: ", &dest)
		n, err := prefixer.Write([]byte("line1\nline2\n"))
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if n != 12 {
			t.Errorf("Expected 12 bytes written, got %d", n)
		}
		if dest.String() != "prefix: line1\nprefix: line2\n" {
			t.Errorf("Unexpected output: %s", dest.String())
		}
	})

	t.Run("multiple lines multiple writes", func(t *testing.T) {
		var dest bytes.Buffer
		prefixer := NewPrefixer("prefix: ", &dest)
		n1, err := prefixer.Write([]byte("line1\nline2"))
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if n1 != 11 {
			t.Errorf("Expected 12 bytes written in first call, got %d", n1)
		}

		n2, err := prefixer.Write([]byte("line2end\nline3\n"))
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if n2 != 15 {
			t.Errorf("Expected 15 bytes written in second call, got %d", n2)
		}
		if dest.String() != "prefix: line1\nprefix: line2line2end\nprefix: line3\n" {
			t.Errorf("Unexpected output: %s", dest.String())
		}
	})

	// Test writing a slice with escape sequences
	t.Run("escape sequences", func(t *testing.T) {
		var dest bytes.Buffer
		prefixer := &CommandPrefixer{
			dest:   &dest,
			prefix: []byte("prefix: "),
		}
		n, err := prefixer.Write([]byte("\x1b[2Kline1\n\x1b[1Gline2\n"))
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if n != 20 {
			t.Errorf("Expected 20 bytes written, got %d", n)
		}
		if dest.String() != "prefix: line1\nprefix: line2\n" {
			t.Errorf("Unexpected output: %s", dest.String())
		}
	})
}
