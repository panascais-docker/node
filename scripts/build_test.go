package main

import (
	"slices"
	"testing"
)

func TestResolveNames(t *testing.T) {
	aliases := configuration{"latest": "26", "lts": "26"}

	for _, testCase := range []struct {
		line, version string
		expected      []string
	}{
		{"26", "26.13.0", []string{"26", "26.13", "26.13.0", "latest", "lts"}},
		{"24", "24.22.2", []string{"24", "24.22", "24.22.2"}},
	} {
		if actual := resolveNames(testCase.line, testCase.version, aliases); !slices.Equal(actual, testCase.expected) {
			t.Errorf("resolveNames(%s) = %q, expected %q", testCase.line, actual, testCase.expected)
		}
	}
}

func TestStaleLines(t *testing.T) {
	expected := configuration{"24": "sha256:a", "26": "sha256:b"}

	for _, testCase := range []struct {
		name      string
		published configuration
		force     bool
		lines     []string
	}{
		{"current", expected, false, []string{}},
		{"force", expected, true, []string{"24", "26"}},
		{"missing", configuration{"24": "sha256:a"}, false, []string{"26"}},
		{"stale", configuration{"24": "sha256:x", "26": "sha256:b"}, false, []string{"24"}},
	} {
		if lines := staleLines(expected, testCase.published, testCase.force); !slices.Equal(lines, testCase.lines) {
			t.Errorf("staleLines(%s) = %v, expected %v", testCase.name, lines, testCase.lines)
		}
	}
}
