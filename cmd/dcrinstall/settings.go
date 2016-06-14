package main

import (
	"flag"
	"fmt"
	"runtime"

	homedir "github.com/mitchellh/go-homedir"
)

// latestVersion and latestManifest must be updated every release.
var (
	latestManifest = "manifest-20160607-01.txt"
	defaultURI     = "https://github.com/decred/decred-release/releases/download/v0.1.5"
)

type Settings struct {
	// command line settings
	Destination  string // destination path
	Manifest     string // manifest name
	Path         string // target path for downloads
	Tuple        string // os-arch tuple
	URI          string // URI to manifest and sets
	SkipDownload bool   // requires path to files
	SkipVerify   bool   // skip TLS and signature checks, internal use only
	Verbose      bool   // loudnes
}

func parseSettings() (*Settings, error) {
	s := Settings{}

	dest := flag.String("dest", "~", "extract path")
	manifest := flag.String("manifest", latestManifest, "manifest name")
	path := flag.String("path", "", "download path")
	tuple := flag.String("tuple", runtime.GOOS+"-"+runtime.GOARCH,
		"OS-Arch tuple, e.g. windows-amd64")
	uri := flag.String("uri", defaultURI, "URI to manifest and sets")
	skip := flag.Bool("skip", false, "skip download, requires path")
	verbose := flag.Bool("verbose", false, "verbose")
	flag.Parse()

	if *tuple == "" {
		return nil, fmt.Errorf("must provide OS-Arch tuple")
	}
	if *skip && *path == "" {
		return nil, fmt.Errorf("must provide download path")
	}
	if *uri != defaultURI {
		s.SkipVerify = true
	}

	var err error
	s.Destination, err = homedir.Expand(*dest)
	if err != nil {
		return nil, err
	}

	s.Manifest = *manifest
	s.Path = *path
	s.Tuple = *tuple
	s.URI = *uri
	s.SkipDownload = *skip
	s.Verbose = *verbose

	return &s, nil
}
