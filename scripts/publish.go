package main

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"github.com/spf13/cobra"
)

func publishCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "publish",
		Short: "Tag the pushed amd64 and arm64 digests of every built line as multi-platform images",
		Args:  cobra.NoArgs,
		RunE:  func(_ *cobra.Command, _ []string) error { return publish() },
	}
}

func publish() error {
	pushed, err := readPushed()
	if err != nil {
		return err
	}

	pins, err := readPins()
	if err != nil {
		return err
	}

	creations, err := planManifests(pushed, pins)
	if err != nil {
		return err
	}

	for _, registry := range registries {
		if err := registry.login(); err != nil {
			return err
		}
	}

	errs := make([]error, len(creations))

	var group sync.WaitGroup
	for index, creation := range creations {
		group.Go(func() { errs[index] = command("docker", creation...).Run() })
	}
	group.Wait()

	return errors.Join(errs...)
}

func readPushed() (map[string]configuration, error) {
	files, err := filepath.Glob(filepath.Join(digestsDirectory, "*.json"))
	if err != nil {
		return nil, err
	}

	pushed := map[string]configuration{}
	for _, file := range files {
		line, architecture, _ := strings.Cut(strings.TrimSuffix(filepath.Base(file), ".json"), "-")

		value, err := readConfiguration(file)
		if err != nil {
			return nil, err
		}

		if pushed[architecture] == nil {
			pushed[architecture] = configuration{}
		}

		pushed[architecture][line] = value[line]
	}

	return pushed, nil
}

func planManifests(pushed map[string]configuration, pins pins) ([][]string, error) {
	built := map[string]bool{}
	for _, architecture := range architectures {
		for line := range pushed[architecture] {
			built[line] = true
		}
	}

	if len(built) == 0 {
		return nil, fmt.Errorf("no pushed digests found in %s", digestsDirectory)
	}

	var creations [][]string
	for _, line := range sortedKeys(built) {
		build, err := planBuild(line, pins)
		if err != nil {
			return nil, err
		}

		var manifests []string
		for _, architecture := range architectures {
			digest := pushed[architecture][line]
			if digest == "" {
				return nil, fmt.Errorf("no pushed %s digest for %s", architecture, line)
			}

			manifests = append(manifests, digest)
		}

		for _, registry := range registries {
			creation := []string{"buildx", "imagetools", "create", "--annotation", "index:" + fingerprintAnnotation + "=" + build.fingerprint}
			for _, tag := range build.tags(registry) {
				creation = append(creation, "--tag", tag)
			}

			for _, manifest := range manifests {
				creation = append(creation, registry.image+"@"+manifest)
			}

			creations = append(creations, creation)
		}
	}

	return creations, nil
}
