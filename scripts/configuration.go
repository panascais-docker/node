package main

import (
	"bytes"
	"cmp"
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

const (
	aliasesFile  = "configuration/aliases.json"
	buildersFile = "configuration/builders.json"
	digestsFile  = "configuration/digests.json"
	pnpmFile     = "configuration/pnpm.json"
	versionsFile = "configuration/versions.json"
	mirror       = "mirror.gcr.io/library/"
	repository   = mirror + "node"
)

var linePattern = regexp.MustCompile(`^\d+$`)

type configuration map[string]string

type pins struct {
	aliases  configuration
	builders configuration
	digests  configuration
	pnpm     configuration
	versions configuration
}

func (value configuration) equal(other configuration) bool {
	return maps.Equal(value, other)
}

func readConfiguration(path string) (configuration, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var value configuration
	if err := json.Unmarshal(content, &value); err != nil {
		return nil, fmt.Errorf("parsing %s: %w", path, err)
	}

	return value, nil
}

func writeConfiguration(path string, value configuration) error {
	var entries []string
	for _, key := range sortedKeys(value) {
		entries = append(entries, fmt.Sprintf("%q:%q", key, value[key]))
	}

	var indented bytes.Buffer
	if err := json.Indent(&indented, []byte("{"+strings.Join(entries, ",")+"}"), "", "    "); err != nil {
		return err
	}

	indented.WriteByte('\n')

	return os.WriteFile(path, indented.Bytes(), 0o644)
}

func readPins() (pins, error) {
	var value pins
	for path, target := range map[string]*configuration{
		aliasesFile:  &value.aliases,
		buildersFile: &value.builders,
		digestsFile:  &value.digests,
		pnpmFile:     &value.pnpm,
		versionsFile: &value.versions,
	} {
		content, err := readConfiguration(path)
		if err != nil {
			return pins{}, err
		}

		*target = content
	}

	return value, nil
}

func writePins(before, after pins) error {
	for path, values := range map[string][2]configuration{
		aliasesFile:  {before.aliases, after.aliases},
		buildersFile: {before.builders, after.builders},
		digestsFile:  {before.digests, after.digests},
		pnpmFile:     {before.pnpm, after.pnpm},
		versionsFile: {before.versions, after.versions},
	} {
		if values[0].equal(values[1]) {
			continue
		}

		if err := writeConfiguration(path, values[1]); err != nil {
			return err
		}
	}

	return nil
}

func readLines() ([]string, error) {
	entries, err := os.ReadDir(".")
	if err != nil {
		return nil, err
	}

	var lines []string
	for _, entry := range entries {
		if !entry.IsDir() || !linePattern.MatchString(entry.Name()) {
			continue
		}

		if _, err := os.Stat(filepath.Join(entry.Name(), "Dockerfile")); err == nil {
			lines = append(lines, entry.Name())
		}
	}

	slices.SortFunc(lines, compareKeys)

	return lines, nil
}

func sortedKeys[Value any](record map[string]Value) []string {
	return slices.SortedFunc(maps.Keys(record), compareKeys)
}

func compareKeys(left, right string) int {
	return cmp.Or(slices.Compare(keyNumbers(left), keyNumbers(right)), strings.Compare(left, right))
}

func keyNumbers(key string) []int {
	var numbers []int
	for part := range strings.SplitSeq(key, ".") {
		number, err := strconv.Atoi(part)
		if err != nil {
			return nil
		}

		numbers = append(numbers, number)
	}

	return numbers
}
