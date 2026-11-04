package main

import (
	"bytes"
	"cmp"
	"context"
	"encoding/json"
	"fmt"
	"maps"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/google/go-containerregistry/pkg/crane"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/spf13/cobra"
)

const releasesURL = "https://nodejs.org/dist/index.json"

var (
	architectures = []string{"amd64", "arm64"}
	blacklist     = []string{"22.7.0", "22.8.0"}
)

type release struct {
	Version string          `json:"version"`
	LTS     json.RawMessage `json:"lts"`
}

func (release release) number() string {
	return strings.TrimPrefix(release.Version, "v")
}

func (release release) major() string {
	major, _, _ := strings.Cut(release.number(), ".")

	return major
}

func (release release) lts() bool {
	return len(release.LTS) > 0 && string(release.LTS) != "false"
}

func updateCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "update",
		Short: "Pin the newest node release and mirrored image digest per line and refresh the aliases",
		Args:  cobra.NoArgs,
		RunE:  func(command *cobra.Command, _ []string) error { return update(command.Context()) },
	}
}

func update(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	lines, err := readLines()
	if err != nil {
		return err
	}

	before, err := readPins()
	if err != nil {
		return err
	}

	releases, err := fetchReleases(ctx)
	if err != nil {
		return err
	}

	after := pins{aliases: maps.Clone(before.aliases), builders: configuration{}, digests: configuration{}, pnpm: before.pnpm, versions: configuration{}}
	for _, line := range lines {
		if before.pnpm[line] == "" {
			return fmt.Errorf("no pnpm version for line %s in %s", line, pnpmFile)
		}

		version, digest := before.versions[line], before.digests[line]
		if candidate := newestRelease(releases, line); candidate != "" {
			pinned := ""
			if candidate == version {
				pinned = digest
			}

			resolved, err := fetchDigest(ctx, repository, candidate+"-alpine", pinned)
			if err != nil {
				return err
			}

			if resolved != "" {
				version, digest = candidate, resolved
			}
		}

		if version == "" || digest == "" {
			return fmt.Errorf("no mirrored node release found for line %s", line)
		}

		after.versions[line], after.digests[line] = version, digest
	}

	maps.Copy(after.aliases, resolveAliases(releases, after.versions))

	for image, pinned := range before.builders {
		name, tag, _ := strings.Cut(image, ":")

		digest, err := fetchDigest(ctx, mirror+name, tag, pinned)
		if err != nil {
			return err
		}

		after.builders[image] = cmp.Or(digest, pinned)
	}

	return writePins(before, after)
}

func fetchReleases(ctx context.Context) ([]release, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, releasesURL, nil)
	if err != nil {
		return nil, err
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetching %s: %s", releasesURL, response.Status)
	}

	var releases []release
	if err := json.NewDecoder(response.Body).Decode(&releases); err != nil {
		return nil, fmt.Errorf("parsing %s: %w", releasesURL, err)
	}

	return releases, nil
}

func newestRelease(releases []release, line string) string {
	for _, release := range releases {
		if release.major() == line && !slices.Contains(blacklist, release.number()) {
			return release.number()
		}
	}

	return ""
}

func resolveAliases(releases []release, versions configuration) configuration {
	aliases := configuration{}
	for _, release := range releases {
		if _, found := versions[release.major()]; !found {
			continue
		}

		if _, found := aliases["latest"]; !found {
			aliases["latest"] = release.major()
		}

		if _, found := aliases["lts"]; !found && release.lts() {
			aliases["lts"] = release.major()
		}
	}

	return aliases
}

func fetchDigest(ctx context.Context, image, tag, pinned string) (string, error) {
	head, err := crane.Head(image+":"+tag, crane.WithContext(ctx))
	if missing(err) {
		return "", nil
	}

	if err != nil || head.Digest.String() == pinned {
		return pinned, err
	}

	manifest, err := crane.Manifest(image+"@"+head.Digest.String(), crane.WithContext(ctx))
	if err != nil {
		return "", err
	}

	index, err := v1.ParseIndexManifest(bytes.NewReader(manifest))
	if err != nil {
		return "", err
	}

	for _, architecture := range architectures {
		if !slices.ContainsFunc(index.Manifests, func(entry v1.Descriptor) bool {
			return entry.Platform != nil && entry.Platform.Architecture == architecture
		}) {
			return "", nil
		}
	}

	return head.Digest.String(), nil
}
