package main

import (
	"encoding/json"
	"maps"
	"testing"
)

func releases(t *testing.T) []release {
	t.Helper()

	var value []release
	if err := json.Unmarshal([]byte(`[
		{"version": "v27.0.0", "lts": false},
		{"version": "v26.13.0", "lts": "Lambda"},
		{"version": "v24.22.2", "lts": "Krypton"},
		{"version": "v22.8.0", "lts": "Jod"},
		{"version": "v22.7.0", "lts": "Jod"},
		{"version": "v22.6.0", "lts": "Jod"}
	]`), &value); err != nil {
		t.Fatal(err)
	}

	return value
}

func TestNewestRelease(t *testing.T) {
	for line, expected := range map[string]string{"26": "26.13.0", "24": "24.22.2", "22": "22.6.0", "20": ""} {
		if actual := newestRelease(releases(t), line); actual != expected {
			t.Errorf("newestRelease(%s) = %q, expected %q", line, actual, expected)
		}
	}
}

func TestResolveAliases(t *testing.T) {
	for _, testCase := range []struct {
		name     string
		versions configuration
		expected configuration
	}{
		{"newest major is a line", configuration{"24": "", "26": "", "27": ""}, configuration{"latest": "27", "lts": "26"}},
		{"newest major is not a line yet", configuration{"24": "", "26": ""}, configuration{"latest": "26", "lts": "26"}},
		{"no lts line", configuration{"27": ""}, configuration{"latest": "27"}},
	} {
		if actual := resolveAliases(releases(t), testCase.versions); !maps.Equal(actual, testCase.expected) {
			t.Errorf("resolveAliases(%s) = %v, expected %v", testCase.name, actual, testCase.expected)
		}
	}
}
