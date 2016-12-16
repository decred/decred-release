// Copyright (c) 2016 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	homedir "github.com/marcopeereboom/go-homedir"
)

// latestVersion and latestManifest must be updated every release.
const (
	latestManifest = "manifest-v0.7.0.txt"
	defaultURI     = "https://github.com/decred/decred-binaries/releases/download/v0.7.0"

	netMain  = "mainnet"
	netTest  = "testnet"
	netSim   = "simnet"
	walletDB = "wallet.db" // start using wallet package one
)

type Settings struct {
	// command line settings
	Destination  string // destination path
	Manifest     string // manifest name
	Net          string // which network to use
	Path         string // target path for downloads
	Tuple        string // os-arch tuple
	URI          string // URI to manifest and sets
	DownloadOnly bool   // download files only
	SkipDownload bool   // requires path to files
	SkipVerify   bool   // skip TLS and signature checks, internal use only
	Quiet        bool   // quiet
	Verbose      bool   // loudness
	Version      bool   // show version.
}

func parseSettings() (*Settings, error) {
	s := Settings{}

	dest := flag.String("dest", "~/decred", "extract path")
	manifest := flag.String("manifest", latestManifest, "manifest name")
	net := flag.String("net", netMain, "decred net "+netMain+", "+netTest+
		" or "+netSim)
	path := flag.String("path", "", "download path")
	tuple := flag.String("tuple", runtime.GOOS+"-"+runtime.GOARCH,
		"OS-Arch tuple, e.g. windows-amd64")
	uri := flag.String("uri", defaultURI, "URI to manifest and sets")
	download := flag.Bool("downloadonly", false, "download binaries "+
		"but don't install")
	skip := flag.Bool("skip", false, "skip download, requires path")
	ver := flag.Bool("version", false, "display version")
	// for backwards compat
	quiet := flag.Bool("quiet", false, "quiet (default false)")
	verbose := flag.Bool("verbose", true, "verbose (deprecated in favor of quiet)")
	flag.Parse()

	if *ver {
		// Show the version and exit if the version flag was specified.
		appName := filepath.Base(os.Args[0])
		appName = strings.TrimSuffix(appName, filepath.Ext(appName))
		fmt.Println(appName, "version", version())
		os.Exit(0)
	}

	if *tuple == "" {
		return nil, fmt.Errorf("must provide OS-Arch tuple")
	}
	if *skip && *path == "" {
		return nil, fmt.Errorf("must provide download path")
	}
	if *skip && *download {
		return nil, fmt.Errorf("downloadonly and skip are mutually exclusive")
	}
	if *uri != defaultURI {
		s.SkipVerify = true
	}

	switch *net {
	case netMain, netTest, netSim:
	default:
		return nil, fmt.Errorf("invalid net, please use %v, %v or %v",
			netMain, netTest, netSim)
	}
	s.Net = *net

	destination, err := homedir.Expand(*dest)
	if err != nil {
		return nil, err
	}
	s.Destination = destination

	if *verbose == true {
		s.Verbose = true
	}
	if *quiet == true {
		s.Verbose = false
	}

	s.Manifest = *manifest
	s.Path = *path
	s.Tuple = *tuple
	s.URI = *uri
	s.SkipDownload = *skip
	s.DownloadOnly = *download

	return &s, nil
}
