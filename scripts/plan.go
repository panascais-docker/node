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
		Short: "Print the <line>-<variant> targets whose published images are stale as a JSON array in GitHub Actions output format",
		Args:  cobra.NoArgs,
		RunE:  func(command *cobra.Command, _ []string) error { return plan(command.Context(), force) },
	}
	command.Flags().BoolVar(&force, "force", false, "treat every target as stale without reading the registry")

	return command
}

func plan(ctx context.Context, force bool) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	builds, err := plannedBuilds()
	if err != nil {
		return err
	}

	if !force && os.Getenv("GITHUB_ACTIONS") == "true" {
		if err := registries[0].login(); err != nil {
			return err
		}
	}

	expected, published := configuration{}, configuration{}
	for _, build := range builds {
		expected[build.key()] = build.fingerprint
		if force {
			continue
		}

		image := registries[0].image + ":" + build.names[0]

		fingerprint, err := publishedFingerprint(ctx, image)
		if missing(err) {
			_, _ = fmt.Fprintf(os.Stderr, "%s is missing\n", image)

			continue
		}

		if err != nil {
			return err
		}

		published[build.key()] = fingerprint
	}

	targets, err := json.Marshal(staleTargets(expected, published, force))
	if err != nil {
		return err
	}

	fmt.Println("targets=" + string(targets))

	return nil
}

func plannedBuilds() ([]build, error) {
	lines, err := readLines()
	if err != nil {
		return nil, err
	}

	pins, err := readPins()
	if err != nil {
		return nil, err
	}

	var builds []build
	for _, line := range lines {
		for _, variant := range variants {
			build, err := planBuild(line, variant, pins)
			if err != nil {
				return nil, err
			}

			builds = append(builds, build)
		}
	}

	return builds, nil
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

func staleTargets(expected, published configuration, force bool) []string {
	targets := []string{}
	for _, target := range slices.Backward(sortedKeys(expected)) {
		if force || expected[target] != published[target] {
			targets = append(targets, target)
		}
	}

	return targets
}
