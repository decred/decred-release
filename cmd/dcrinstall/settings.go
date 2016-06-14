package main

import (
	"flag"
	"fmt"
	"runtime"

	homedir "github.com/mitchellh/go-homedir"
)

// latestVersion and latestManifest must be updated every release.
var (
	latestVersion  = "v0.1.4"
	latestManifest = "manifest-20160526-01.txt"
)

type Settings struct {
	// command line settings
	Destination  string // destination path
	Manifest     string // manifest name
	Path         string // target path for downloads
	Tuple        string // os-arch tuple
	Version      string // version used in download path
	SkipDownload bool   // requires path to files
	Verbose      bool   // loudnes
}

func parseSettings() (*Settings, error) {
	s := Settings{}

	dest := flag.String("dest", "~", "extract path")
	manifest := flag.String("manifest", latestManifest, "manifest name")
	path := flag.String("path", "", "download path")
	tuple := flag.String("tuple", runtime.GOOS+"-"+runtime.GOARCH,
		"OS-Arch tuple, e.g. windows-amd64")
	version := flag.String("version", latestVersion, "decred version to download")
	skip := flag.Bool("skip", false, "skip download, requires path")
	verbose := flag.Bool("verbose", false, "verbose")
	flag.Parse()

	if *tuple == "" {
		return nil, fmt.Errorf("must provide OS-Arch tuple")
	}
	if *skip && *path == "" {
		return nil, fmt.Errorf("must provide download path")
	}

	var err error
	s.Destination, err = homedir.Expand(*dest)
	if err != nil {
		return nil, err
	}

	s.Manifest = *manifest
	s.Path = *path
	s.Tuple = *tuple
	s.Version = *version
	s.SkipDownload = *skip
	s.Verbose = *verbose

	return &s, nil
}
