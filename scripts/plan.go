package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"slices"
	"time"

	"github.com/google/go-containerregistry/pkg/crane"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote/transport"
	"github.com/spf13/cobra"
)

func planCommand() *cobra.Command {
	var force bool

	command := &cobra.Command{
		Use:   "plan",
		Short: "Print the lines whose published images are stale as a JSON array in GitHub Actions output format",
		Args:  cobra.NoArgs,
		RunE:  func(command *cobra.Command, _ []string) error { return plan(command.Context(), force) },
	}
	command.Flags().BoolVar(&force, "force", false, "treat every line as stale without reading the registry")

	return command
}

func plan(ctx context.Context, force bool) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	expected, err := expectedFingerprints()
	if err != nil {
		return err
	}

	published := configuration{}
	if !force {
		if os.Getenv("GITHUB_ACTIONS") == "true" {
			if err := registries[0].login(); err != nil {
				return err
			}
		}

		for _, line := range sortedKeys(expected) {
			image := registries[0].image + ":" + line

			fingerprint, err := publishedFingerprint(ctx, image)
			if missing(err) {
				_, _ = fmt.Fprintf(os.Stderr, "%s is missing\n", image)

				continue
			}

			if err != nil {
				return err
			}

			published[line] = fingerprint
		}
	}

	lines, err := json.Marshal(staleLines(expected, published, force))
	if err != nil {
		return err
	}

	fmt.Println("lines=" + string(lines))

	return nil
}

func expectedFingerprints() (configuration, error) {
	lines, err := readLines()
	if err != nil {
		return nil, err
	}

	pins, err := readPins()
	if err != nil {
		return nil, err
	}

	expected := configuration{}
	for _, line := range lines {
		build, err := planBuild(line, pins)
		if err != nil {
			return nil, err
		}

		expected[line] = build.fingerprint
	}

	return expected, nil
}

func publishedFingerprint(ctx context.Context, image string) (string, error) {
	manifest, err := crane.Manifest(image, crane.WithContext(ctx))
	if err != nil {
		return "", err
	}

	index, err := v1.ParseIndexManifest(bytes.NewReader(manifest))
	if err != nil {
		return "", err
	}

	return index.Annotations[fingerprintAnnotation], nil
}

func missing(err error) bool {
	transportError, found := errors.AsType[*transport.Error](err)
	if !found {
		return false
	}

	return transportError.StatusCode == http.StatusNotFound || slices.ContainsFunc(transportError.Errors, func(diagnostic transport.Diagnostic) bool {
		return diagnostic.Code == transport.ManifestUnknownErrorCode
	})
}

func staleLines(expected, published configuration, force bool) []string {
	lines := []string{}
	for _, line := range sortedKeys(expected) {
		if force || expected[line] != published[line] {
			lines = append(lines, line)
		}
	}

	return lines
}
