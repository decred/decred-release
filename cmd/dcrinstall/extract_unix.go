// +build !windows,!plan9

package main

import (
	"fmt"
	"path/filepath"

	"github.com/docker/docker/pkg/archive"
)

// extract downloaded package.
func (c *ctx) extract() error {
	manifest := filepath.Join(c.s.Path, c.s.Manifest)
	_, filename, err := findOS(c.s.Tuple, manifest)
	if err != nil {
		return err
	}

	if c.s.Verbose {
		fmt.Printf("extracting: %v -> %v\n", filename, c.s.Destination)
	}

	err = archive.UntarPath(filepath.Join(c.s.Path, filename),
		c.s.Destination)
	if err != nil {
		return err
	}

	return nil
}
