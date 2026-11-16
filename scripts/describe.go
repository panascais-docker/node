package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"slices"
	"time"

	"github.com/spf13/cobra"
)

const (
	hubURL     = "https://hub.docker.com"
	readmeFile = "readme.md"
	sourceURL  = "https://github.com/panascais-docker/node/blob/master/"
	// overviewLimit is the byte limit Docker Hub enforces on a repository overview.
	overviewLimit = 25000
)

// relativeLinkPattern matches markdown link targets without a scheme or leading anchor, such as ](license.md).
var relativeLinkPattern = regexp.MustCompile(`\]\(([^)\s:#][^)\s:]*)\)`)

func describeCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "describe",
		Short: "Publish readme.md as the Docker Hub repository overview",
		Args:  cobra.NoArgs,
		RunE:  func(command *cobra.Command, _ []string) error { return describe(command.Context()) },
	}
}

func describe(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	content, err := os.ReadFile(readmeFile)
	if err != nil {
		return err
	}

	overview := absoluteLinks(string(content))
	if len(overview) > overviewLimit {
		return fmt.Errorf("%s is %d bytes, Docker Hub accepts at most %d", readmeFile, len(overview), overviewLimit)
	}

	hub := registries[slices.IndexFunc(registries, func(registry registry) bool { return registry.host == "" })]

	var login struct {
		AccessToken string `json:"access_token"`
	}
	if err := hubRequest(ctx, http.MethodPost, "/v2/auth/token", "", map[string]string{
		"identifier": os.Getenv(hub.secrets + "_REGISTRY_USERNAME"),
		"secret":     os.Getenv(hub.secrets + "_REGISTRY_TOKEN"),
	}, &login); err != nil {
		return err
	}

	if login.AccessToken == "" {
		return errors.New("docker hub returned no access token")
	}

	return hubRequest(ctx, http.MethodPatch, "/v2/repositories/"+hub.image, login.AccessToken, map[string]string{"full_description": overview}, nil)
}

// absoluteLinks points relative links at the GitHub repository, since Docker Hub cannot resolve them.
func absoluteLinks(markdown string) string {
	return relativeLinkPattern.ReplaceAllString(markdown, "]("+sourceURL+"$1)")
}

func hubRequest(ctx context.Context, method, path, token string, body, result any) error {
	payload, err := json.Marshal(body)
	if err != nil {
		return err
	}

	request, err := http.NewRequestWithContext(ctx, method, hubURL+path, bytes.NewReader(payload))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")
	if token != "" {
		request.Header.Set("Authorization", "Bearer "+token)
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode > 299 {
		detail, _ := io.ReadAll(io.LimitReader(response.Body, 512))

		return fmt.Errorf("%s %s: %s %s", method, path, response.Status, bytes.TrimSpace(detail))
	}

	if result == nil {
		return nil
	}

	return json.NewDecoder(response.Body).Decode(result)
}
