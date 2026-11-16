package main

import (
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/spf13/cobra"
)

func publishCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "publish",
		Short: "Tag the pushed amd64 and arm64 digests of every built target as multi-platform images",
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
	slots := make(chan struct{}, 4)

	var group sync.WaitGroup
	for index, creation := range creations {
		group.Go(func() {
			slots <- struct{}{}
			defer func() { <-slots }()

			errs[index] = retryRateLimited(func() *exec.Cmd { return command("docker", creation...) })
		})
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
		name := strings.TrimSuffix(filepath.Base(file), ".json")

		separator := strings.LastIndex(name, "-")
		if separator < 0 {
			return nil, fmt.Errorf("no architecture in %s", file)
		}

		key, architecture := name[:separator], name[separator+1:]

		value, err := readConfiguration(file)
		if err != nil {
			return nil, err
		}

		if pushed[architecture] == nil {
			pushed[architecture] = configuration{}
		}

		pushed[architecture][key] = value[key]
	}

	return pushed, nil
}

func planManifests(pushed map[string]configuration, pins pins) ([][]string, error) {
	built := map[string]bool{}
	for _, architecture := range architectures {
		for key := range pushed[architecture] {
			built[key] = true
		}
	}

	if len(built) == 0 {
		return nil, fmt.Errorf("no pushed digests found in %s", digestsDirectory)
	}

	var creations [][]string
	for _, key := range sortedKeys(built) {
		line, variant, err := parseTarget(key)
		if err != nil {
			return nil, err
		}

		build, err := planBuild(line, variant, pins)
		if err != nil {
			return nil, err
		}

		var manifests []string
		for _, architecture := range architectures {
			digest := pushed[architecture][key]
			if digest == "" {
				return nil, fmt.Errorf("no pushed %s digest for %s", architecture, key)
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
