package main

import (
	"slices"
	"testing"
)

func TestResolveNames(t *testing.T) {
	aliases := configuration{"latest": "26", "lts": "26"}

	for _, testCase := range []struct {
		line, version, variant string
		expected               []string
	}{
		{"26", "26.13.0", "alpine", []string{
			"26", "26.13", "26.13.0", "latest", "lts",
			"26-alpine", "26.13-alpine", "26.13.0-alpine", "alpine", "latest-alpine", "lts-alpine",
		}},
		{"26", "26.13.0", "debian", []string{"26-debian", "26.13-debian", "26.13.0-debian", "debian", "latest-debian", "lts-debian"}},
		{"24", "24.22.2", "alpine", []string{"24", "24.22", "24.22.2", "24-alpine", "24.22-alpine", "24.22.2-alpine"}},
		{"24", "24.22.2", "debian", []string{"24-debian", "24.22-debian", "24.22.2-debian"}},
	} {
		actual := resolveNames(testCase.line, testCase.version, testCase.variant, aliases)
		if !slices.Equal(actual, testCase.expected) {
			t.Errorf("resolveNames(%s-%s) = %q, expected %q", testCase.line, testCase.variant, actual, testCase.expected)
		}
	}
}

func TestStaleTargets(t *testing.T) {
	expected := configuration{"24-alpine": "sha256:a", "26-alpine": "sha256:b", "26-debian": "sha256:c"}

	for _, testCase := range []struct {
		name      string
		published configuration
		force     bool
		targets   []string
	}{
		{"current", expected, false, []string{}},
		{"force", expected, true, []string{"26-debian", "26-alpine", "24-alpine"}},
		{"missing", configuration{"24-alpine": "sha256:a", "26-alpine": "sha256:b"}, false, []string{"26-debian"}},
		{"stale", configuration{"24-alpine": "sha256:x", "26-alpine": "sha256:b", "26-debian": "sha256:c"}, false, []string{"24-alpine"}},
	} {
		if targets := staleTargets(expected, testCase.published, testCase.force); !slices.Equal(targets, testCase.targets) {
			t.Errorf("staleTargets(%s) = %v, expected %v", testCase.name, targets, testCase.targets)
		}
	}
}
