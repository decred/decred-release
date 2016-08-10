// Copyright (c) 2016 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"runtime"

	homedir "github.com/marcopeereboom/go-homedir"
)

// latestVersion and latestManifest must be updated every release.
const (
	latestManifest = "manifest-v0.2.0.txt"
	defaultURI     = "https://github.com/decred/decred-binaries/releases/download/v0.2.0"

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
	SkipDownload bool   // requires path to files
	SkipVerify   bool   // skip TLS and signature checks, internal use only
	Verbose      bool   // loudnes
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

	s.Manifest = *manifest
	s.Path = *path
	s.Tuple = *tuple
	s.URI = *uri
	s.SkipDownload = *skip
	s.Verbose = *verbose

	return &s, nil
}
