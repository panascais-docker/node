package main

import "testing"

func TestAbsoluteLinks(t *testing.T) {
	for markdown, expected := range map[string]string{
		"[guidelines](contributing.md)":         "[guidelines](" + sourceURL + "contributing.md)",
		"[Variants](#variants)":                 "[Variants](#variants)",
		"[docs](https://pnpm.io/docker)":        "[docs](https://pnpm.io/docker)",
		"[(silas@x)](mailto:silas@x)":           "[(silas@x)](mailto:silas@x)",
		"![size](https://img.shields.io/a?b=c)": "![size](https://img.shields.io/a?b=c)",
		"[config](configuration/versions.json)": "[config](" + sourceURL + "configuration/versions.json)",
	} {
		if actual := absoluteLinks(markdown); actual != expected {
			t.Errorf("absoluteLinks(%q) = %q, expected %q", markdown, actual, expected)
		}
	}
}
