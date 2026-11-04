package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/fs"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

const fingerprintAnnotation = "net.panascais.docker.node.fingerprint"

func lineFingerprint(build build) (string, error) {
	values := map[string]string{
		"image":   build.image,
		"names":   strings.Join(build.names, " "),
		"pnpm":    build.pnpm,
		"version": build.version,
	}
	maps.Copy(values, build.builders)

	err := filepath.WalkDir(build.line, func(path string, entry fs.DirEntry, err error) error {
		if err != nil || entry.IsDir() {
			return err
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		sum := sha256.Sum256(content)
		values["file:"+filepath.ToSlash(path)] = hex.EncodeToString(sum[:])

		return nil
	})

	return fingerprint(values), err
}

func fingerprint(values map[string]string) string {
	hash := sha256.New()
	for _, key := range slices.Sorted(maps.Keys(values)) {
		_, _ = fmt.Fprintf(hash, "%s=%s\n", key, values[key])
	}

	return "sha256:" + hex.EncodeToString(hash.Sum(nil))
}
