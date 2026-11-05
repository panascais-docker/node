package main

import (
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"slices"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

const digestsDirectory = "digests"

type registry struct {
	host    string
	image   string
	secrets string
}

var (
	registries = []registry{
		{host: "ghcr.io", image: "ghcr.io/panascais-docker/node/node", secrets: "CONTAINER"},
		{image: "panascais/node", secrets: "DOCKER"},
		{host: "quay.io", image: "quay.io/panascais/node", secrets: "QUAY"},
	}
	digestPattern = regexp.MustCompile(`^sha256:[0-9a-f]{64}$`)
)

type build struct {
	line        string
	version     string
	pnpm        string
	image       string
	builders    map[string]string
	names       []string
	fingerprint string
}

func buildCommand() *cobra.Command {
	var platform string

	command := &cobra.Command{
		Use:   "build <line>",
		Short: "Build a line, pushing by digest on GitHub Actions and loading and smoke testing locally",
		Args:  cobra.ExactArgs(1),
		RunE:  func(_ *cobra.Command, arguments []string) error { return buildLine(arguments[0], platform) },
	}
	command.Flags().StringVar(&platform, "platform", "linux/"+runtime.GOARCH, "platform to build")

	return command
}

func buildLine(line, platform string) error {
	pins, err := readPins()
	if err != nil {
		return err
	}

	build, err := planBuild(line, pins)
	if err != nil {
		return err
	}

	if os.Getenv("GITHUB_ACTIONS") == "true" {
		return push(build, platform)
	}

	return buildLocally(build, platform)
}

func planBuild(line string, pins pins) (build, error) {
	version, digest, pnpm := pins.versions[line], pins.digests[line], pins.pnpm[line]
	if !regexp.MustCompile(`^` + regexp.QuoteMeta(line) + `\.\d+\.\d+$`).MatchString(version) {
		return build{}, fmt.Errorf("invalid node version %q for line %s", version, line)
	}

	if !digestPattern.MatchString(digest) {
		return build{}, fmt.Errorf("invalid node digest %q for line %s", digest, line)
	}

	if pnpm == "" {
		return build{}, fmt.Errorf("no pnpm version for line %s", line)
	}

	builders, err := builderImages(pins.builders)
	if err != nil {
		return build{}, err
	}

	planned := build{
		line:     line,
		version:  version,
		pnpm:     pnpm,
		image:    repository + ":" + version + "-alpine@" + digest,
		builders: builders,
		names:    resolveNames(line, version, pins.aliases),
	}

	planned.fingerprint, err = lineFingerprint(planned)

	return planned, err
}

func resolveNames(line, version string, aliases configuration) []string {
	minor, _, _ := strings.Cut(strings.TrimPrefix(version, line+"."), ".")

	names := []string{line, line + "." + minor, version}
	for _, alias := range sortedKeys(aliases) {
		if aliases[alias] == line {
			names = append(names, alias)
		}
	}

	return names
}

func builderImages(builders configuration) (map[string]string, error) {
	images := map[string]string{}
	for image, digest := range builders {
		name, _, found := strings.Cut(image, ":")
		if !found || !digestPattern.MatchString(digest) {
			return nil, fmt.Errorf("invalid builder %s@%s", image, digest)
		}

		images[strings.ToUpper(name)+"_IMAGE"] = library + image + "@" + digest
	}

	return images, nil
}

func (build build) tags(registry registry) []string {
	tags := make([]string, len(build.names))
	for index, name := range build.names {
		tags[index] = registry.image + ":" + name
	}

	return tags
}

func (build build) arguments(revision string) []string {
	values := map[string]string{
		"BUILD_DATE":   time.Now().UTC().Format(time.RFC3339),
		"NODE_IMAGE":   build.image,
		"NODE_VERSION": build.version,
		"PNPM_VERSION": build.pnpm,
		"VCS_REF":      revision,
	}
	maps.Copy(values, build.builders)

	var arguments []string
	for _, key := range slices.Sorted(maps.Keys(values)) {
		arguments = append(arguments, "--build-arg", key+"="+values[key])
	}

	return arguments
}

func push(build build, platform string) error {
	revision, err := output("git", "rev-parse", "--short", "HEAD")
	if err != nil {
		return err
	}

	var images []string
	for _, registry := range registries {
		if err := registry.login(); err != nil {
			return err
		}

		images = append(images, registry.image)
	}

	metadataFile := filepath.Join(os.TempDir(), "node-metadata.json")
	arguments := slices.Concat(
		[]string{"buildx", "build", "--platform", platform, "--progress=plain", "--metadata-file", metadataFile},
		[]string{"--output", `type=image,"name=` + strings.Join(images, ",") + `",push-by-digest=true,name-canonical=true,push=true`},
		build.arguments(revision),
		[]string{build.line},
	)
	if err := command("docker", arguments...).Run(); err != nil {
		return err
	}

	digest, err := readPushedDigest(metadataFile)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(digestsDirectory, 0o755); err != nil {
		return err
	}

	_, architecture, _ := strings.Cut(platform, "/")

	return writeConfiguration(filepath.Join(digestsDirectory, build.line+"-"+architecture+".json"), configuration{build.line: digest})
}

func readPushedDigest(metadataFile string) (string, error) {
	content, err := os.ReadFile(metadataFile)
	if err != nil {
		return "", err
	}

	var metadata struct {
		Digest string `json:"containerimage.digest"`
	}
	if err := json.Unmarshal(content, &metadata); err != nil || !digestPattern.MatchString(metadata.Digest) {
		return "", fmt.Errorf("no pushed digest in %s", metadataFile)
	}

	return metadata.Digest, nil
}

func buildLocally(build build, platform string) error {
	builder, err := output("docker", "context", "show")
	if err != nil {
		return err
	}

	var tags []string
	for _, registry := range registries {
		for _, tag := range build.tags(registry) {
			tags = append(tags, "--tag", tag)
		}
	}

	arguments := slices.Concat(
		[]string{"buildx", "build", "--builder", builder, "--platform", platform, "--progress=plain", "--load"},
		tags,
		build.arguments("local"),
		[]string{build.line},
	)
	if err := command("docker", arguments...).Run(); err != nil {
		return err
	}

	return verifyImage(build)
}

func verifyImage(build build) error {
	image := registries[0].image + ":" + build.version

	version, err := output("docker", "run", "--rm", image, "node", "--version")
	if err != nil {
		return err
	}

	if version != "v"+build.version {
		return fmt.Errorf("%s reports node %s, expected v%s", image, version, build.version)
	}

	return command("docker", "run", "--rm", image, "pnpm", "--version").Run()
}

func (registry registry) login() error {
	arguments := []string{"login"}
	if registry.host != "" {
		arguments = append(arguments, registry.host)
	}

	login := command("docker", append(arguments, "-u", os.Getenv(registry.secrets+"_REGISTRY_USERNAME"), "--password-stdin")...)
	login.Stdout = os.Stderr
	login.Stdin = strings.NewReader(os.Getenv(registry.secrets + "_REGISTRY_TOKEN"))

	return login.Run()
}

func command(name string, arguments ...string) *exec.Cmd {
	process := exec.Command(name, arguments...)
	process.Stdout, process.Stderr = os.Stdout, os.Stderr

	return process
}

func output(name string, arguments ...string) (string, error) {
	process := exec.Command(name, arguments...)
	process.Stderr = os.Stderr

	result, err := process.Output()

	return strings.TrimSpace(string(result)), err
}
