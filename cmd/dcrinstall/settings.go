// Copyright (c) 2016-2020 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
)

// latestVersion and latestManifest must be updated every release.
const (
	latestVersion  = "v1.6.0-rc1"
	latestManifest = "decred-" + latestVersion + "-manifest.txt"
	defaultURI     = "https://github.com/decred/decred-binaries/releases/download/" + latestVersion

	netMain  = "mainnet"
	netTest  = "testnet3"
	netSim   = "simnet"
	walletDB = "wallet.db" // start using wallet package one
)

// Settings defines the settings
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
	SkipAsc      bool   // disables downloading and verifying the pgp signatures
	Dcrdex       bool   // install dcrdex
	BitcoinURI   string // bitcoin core uri
	DexURI       string // dcrdex uri
	Quiet        bool   // quiet
	Verbose      bool   // loudness
	Version      bool   // show version.
}

func parseSettings() (*Settings, error) {
	defaultTuple := runtime.GOOS + "-" + runtime.GOARCH
	defaultBitcoinURI := bitcoinDownloads[defaultTuple]
	defaultDexURI := dexDownloads[defaultTuple]

	s := Settings{}
	u, err := user.Current()
	if err != nil {
		return nil, err
	}
	dest := flag.String("dest", filepath.Join(u.HomeDir, "decred"), "extract path")
	manifest := flag.String("manifest", latestManifest, "manifest name")
	net := flag.String("net", netMain, "decred net "+netMain+", "+netTest+
		" or "+netSim)
	path := flag.String("path", "", "download path")
	tuple := flag.String("tuple", defaultTuple,
		"OS-Arch tuple, e.g. windows-amd64")
	uri := flag.String("uri", defaultURI, "URI to manifest and sets")
	download := flag.Bool("downloadonly", false, "download binaries "+
		"but don't install")
	skip := flag.Bool("skip", false, "skip download, requires path")
	skipAsc := flag.Bool("skipasc", false, "skip download and verification"+
		" of pgp signatures, requires path")
	dcrdex := flag.Bool("dcrdex", false, "install dcrdex")
	bitcoinURI := flag.String("bitcoinuri", defaultBitcoinURI,
		"bitcoin download path")
	dexURI := flag.String("dexuri", defaultDexURI, "dcrdex download path")
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

	if *dest == "" {
		return nil, fmt.Errorf("destination not set")
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

	switch *net {
	case netMain, netTest, netSim:
	default:
		return nil, fmt.Errorf("invalid net, please use %v, %v or %v",
			netMain, netTest, netSim)
	}
	s.Net = *net
	s.Destination = filepath.Clean(*dest)

	if *verbose {
		s.Verbose = true
	}
	if *quiet {
		s.Verbose = false
	}

	// Check to see if dcrdex is supported on this platform.
	if *dcrdex {
		if *dexURI != "" {
			s.DexURI = *dexURI
		} else {
			return nil, fmt.Errorf("dcrdex cannot be installed "+
				"because it does not support %v", *tuple)
		}
		if *bitcoinURI != "" {
			s.BitcoinURI = *bitcoinURI
		} else {
			return nil, fmt.Errorf("dcrdex cannot be installed "+
				"because bitcoin core does not support %v",
				*tuple)
		}
	}

	s.Manifest = *manifest
	s.Path = *path
	s.Tuple = *tuple
	s.URI = *uri
	s.SkipDownload = *skip
	s.Dcrdex = *dcrdex
	s.SkipAsc = *skipAsc
	s.DownloadOnly = *download

	return &s, nil
}
