// Package schema contains all of the Schema Definition Files (abbreviated SDL)
// which define this API.
//
// Typically it is helpful to look at the interface separate from the implementation.
// Keeping the SDL files separate from the code and in small pieces makes it easier to
// work with.
//
// These SDL files are embedded into the application binary for ease of deployment.
// Embedding enables us to keep our deploy artifact to just one binary.
//
// Alternatively, you could package the SDL files separately (i.e. in a Docker container)
// and skip the step of embedding them.
//
// Use `go generate` to pack all *.graphql files under this directory (and sub-directories) into
// into the application binary.
package schema

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"strings"
)

// content holds all the SDL file content.
//go:embed *.graphql type/*.graphql
var content embed.FS

// String reads the .graphql schema files from the embed.FS, concatenating the
// files together into one string.
//
// If this method complains about not finding functions AssetNames() or MustAsset(),
// run `go generate` against this package to generate the functions.
func String() (string, error) {
	var buf bytes.Buffer

	fn := func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("walking dir: %w", err)
		}

		// Only add files with the .graphql extension.
		if !strings.HasSuffix(path, ".graphql") {
			return nil
		}

		b, err := content.ReadFile(path)
		if err != nil {
			return fmt.Errorf("reading file %q: %w", path, err)
		}

		// Add a newline to separate each file.
		b = append(b, []byte("\n")...)

		if _, err := buf.Write(b); err != nil {
			return fmt.Errorf("writing %q bytes to buffer: %w", path, err)
		}

		return nil
	}

	// Recursively walk this directory and append all the file contents together.
	if err := fs.WalkDir(content, ".", fn); err != nil {
		return buf.String(), fmt.Errorf("walking content directory: %w", err)
	}

	return buf.String(), nil
}
