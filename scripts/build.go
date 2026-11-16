package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"math/rand/v2"
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
	variant     variant
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
		Use:   "build <line>[-<variant>]",
		Short: "Build a line's variant (alpine by default), pushing by digest on GitHub Actions and loading and smoke testing locally",
		Args:  cobra.ExactArgs(1),
		RunE:  func(_ *cobra.Command, arguments []string) error { return buildLine(arguments[0], platform) },
	}
	command.Flags().StringVar(&platform, "platform", "linux/"+runtime.GOARCH, "platform to build")

	return command
}

func buildLine(target, platform string) error {
	line, variant, err := parseTarget(target)
	if err != nil {
		return err
	}

	pins, err := readPins()
	if err != nil {
		return err
	}

	build, err := planBuild(line, variant, pins)
	if err != nil {
		return err
	}

	if os.Getenv("GITHUB_ACTIONS") == "true" {
		return push(build, platform)
	}

	return buildLocally(build, platform)
}

// parseTarget splits a <line>[-<variant>] target, defaulting to the first variant.
func parseTarget(target string) (string, variant, error) {
	line, name, found := strings.Cut(target, "-")
	if !found {
		return line, variants[0], nil
	}

	for _, variant := range variants {
		if variant.name == name {
			return line, variant, nil
		}
	}

	return "", variant{}, fmt.Errorf("unknown variant %q in %s", name, target)
}

func planBuild(line string, variant variant, pins pins) (build, error) {
	key := line + "-" + variant.name
	version, digest, pnpm := pins.versions[line], pins.digests[key], pins.pnpm[line]
	if !regexp.MustCompile(`^` + regexp.QuoteMeta(line) + `\.\d+\.\d+$`).MatchString(version) {
		return build{}, fmt.Errorf("invalid node version %q for line %s", version, line)
	}

	if !digestPattern.MatchString(digest) {
		return build{}, fmt.Errorf("invalid node digest %q for %s", digest, key)
	}

	if pnpm == "" {
		return build{}, fmt.Errorf("no pnpm version for line %s", line)
	}

	builders, err := builderImages(pins.builders, variant)
	if err != nil {
		return build{}, err
	}

	planned := build{
		line:     line,
		variant:  variant,
		version:  version,
		pnpm:     pnpm,
		image:    repository + ":" + version + "-" + variant.upstream + "@" + digest,
		builders: builders,
		names:    resolveNames(line, version, variant.name, pins.aliases),
	}

	planned.fingerprint, err = lineFingerprint(planned)

	return planned, err
}

// resolveNames lists the tags of a line's variant, canonical tag first. Every variant gets suffixed tags
// ("latest" becomes the bare variant name) and the default variant also owns the unsuffixed ones.
func resolveNames(line, version, variant string, aliases configuration) []string {
	minor, _, _ := strings.Cut(strings.TrimPrefix(version, line+"."), ".")

	bases := []string{line, line + "." + minor, version}
	for _, alias := range sortedKeys(aliases) {
		if aliases[alias] == line {
			bases = append(bases, alias)
		}
	}

	var names []string
	if variant == variants[0].name {
		names = append(names, bases...)
	}

	for _, base := range bases {
		if base == "latest" {
			names = append(names, variant)
		} else {
			names = append(names, base+"-"+variant)
		}
	}

	return names
}

// builderImages pins the base images of a variant as <NAME>_IMAGE build arguments.
func builderImages(builders configuration, variant variant) (map[string]string, error) {
	images := map[string]string{}
	for image, digest := range builders {
		name, _, found := strings.Cut(image, ":")
		if !found || !digestPattern.MatchString(digest) {
			return nil, fmt.Errorf("invalid builder %s@%s", image, digest)
		}

		if name == variant.name {
			images[strings.ToUpper(name)+"_IMAGE"] = library + image + "@" + digest
		}
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

func (build build) key() string {
	return build.line + "-" + build.variant.name
}

// source selects the variant's Dockerfile within the line's build context.
func (build build) source() []string {
	return []string{"--file", filepath.Join(build.line, "Dockerfile."+build.variant.name), build.line}
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
		build.source(),
	)
	if err := retryRateLimited(func() *exec.Cmd { return command("docker", arguments...) }); err != nil {
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

	return writeConfiguration(filepath.Join(digestsDirectory, build.key()+"-"+architecture+".json"), configuration{build.key(): digest})
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
		build.source(),
	)
	if err := command("docker", arguments...).Run(); err != nil {
		return err
	}

	return verifyImage(build)
}

func verifyImage(build build) error {
	image := registries[0].image + ":" + build.version + "-" + build.variant.name

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

func retryRateLimited(create func() *exec.Cmd) error {
	for attempt := 0; ; attempt++ {
		var stderr bytes.Buffer
		process := create()
		process.Stderr = io.MultiWriter(process.Stderr, &stderr)

		err := process.Run()
		if err == nil || attempt == 5 || !strings.Contains(stderr.String(), "429 Too Many Requests") {
			return err
		}

		delay := time.Duration(30<<attempt)*time.Second + rand.N(30*time.Second)
		_, _ = fmt.Fprintf(os.Stderr, "rate limited by a registry, retrying in %s\n", delay.Round(time.Second))
		time.Sleep(delay)
	}
}

func output(name string, arguments ...string) (string, error) {
	process := exec.Command(name, arguments...)
	process.Stderr = os.Stderr

	result, err := process.Output()

	return strings.TrimSpace(string(result)), err
}
