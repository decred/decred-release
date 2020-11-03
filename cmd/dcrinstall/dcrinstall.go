// Copyright (c) 2016-2020 The Decred developers
// Use of this source code is governed by an ISC license that can be found in
// the LICENSE file.

package main

import (
	"bufio"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

var (
	// Generated values such as file and directory names
	username                        string // Username used in config files
	password                        string // Password used in config files
	manifestDecredVersion           string // Decred manifest version
	manifestDecredFilename          string // Decred manifest filename
	manifestDecredSignatureFilename string // Decred manifest signature filename
	decredBundleFilename            string // Decred bundle that is downloaded
	decredDownloadURI               string // Decred bundle download URI

	// Generated dcrdex stuff
	manifestDcrdexVersion           string // Dcrdex manifest version
	manifestDcrdexFilename          string // Dcrdex manifest filename
	manifestDcrdexSignatureFilename string // Dcrdex manifest signature filename
	dcrdexBundleFilename            string // Dcrdex bundle that is downloaded
	dcrdexDownloadURI               string // Dcrdex bundle download URI

	// Generated bitcoin stuff
	manifestBitcoinVersion  string // Bitcoin manifest version
	manifestBitcoinFilename string // Bitcoin manifest filename
	bitcoinBundleFilename   string // Bitcoin bundle that is downloaded
	bitcoinDownloadURI      string // Bitcoin bundle download URI

	postProcess []string // Things to tell the user after installation
)

// dcrinstall performs the install/ugrade for everything.
func dcrinstall() error {
	log.Printf("=== dcrinstall start ===")

	// create temporary directory
	var err error
	tmpDir, err = ioutil.TempDir("", "dcrinstall")
	if err != nil {
		return fmt.Errorf("Create temporary file: %v", err)
	}
	log.Printf("Download directory: %v", tmpDir)

	// Decred pre conditions
	err = decredDownloadAndVerify()
	if err != nil {
		return fmt.Errorf("Decred download and verify: %v", err)
	}

	// Dcrdex pre conditions
	if dcrdex {
		err = dcrdexDownloadAndVerify()
		if err != nil {
			return fmt.Errorf("Dcrdex download and verify: %v", err)
		}
	}

	// Bitcoin pre conditions
	if dcrdex {
		err = bitcoinDownloadAndVerify()
		if err != nil {
			return fmt.Errorf("Bitcoin download and verify: %v", err)
		}
	}

	// Install decred
	err = installDecredBundle()
	if err != nil {
		return fmt.Errorf("Decred install: %v", err)
	}

	// Install dcrdex
	if dcrdex {
		err = installDcrdexBundle()
		if err != nil {
			return fmt.Errorf("Dcrdex install: %v", err)
		}
	}

	// Install bitcoin
	if dcrdex {
		err = installBitcoinBundle()
		if err != nil {
			return fmt.Errorf("Bitcoin install: %v", err)
		}
	}

	log.Printf("=== dcrinstall complete ===")

	postProcess = append(postProcess,
		fmt.Sprintf("\nAll binaries have been installed to %v\n\n"+
			"For example, to run dcrd use the following command: '%v'\n\n"+
			"The subdirectories that exist in %v are backups of installation artifacts."+
			" Please do not remove or use them unless directed to.\n\n",
			destination, filepath.Join(destination, "dcrd"),
			destination))

	if len(postProcess) > 0 {
		fmt.Println()
		for k := range postProcess {
			fmt.Printf("=== Post installation message %v ===\n", k)
			fmt.Printf("%v", postProcess[k])
		}
	}

	return nil
}

var (
	defaultLatestManifestURI = "https://raw.githubusercontent.com/decred/decred-release/master/latest"

	defaultTuple                  = runtime.GOOS + "-" + runtime.GOARCH
	defaultDecredManifestVersion  = "v1.6.0-rc2"
	defaultDecredManifestFilename = "decred-" + defaultDecredManifestVersion +
		"-manifest.txt"
	defaultDecredManifestURI = "https://github.com/decred/decred-binaries" +
		"/releases/download/" + defaultDecredManifestVersion + "/" +
		defaultDecredManifestFilename

	// dcrdex
	defaultDcrdexManifestVersion  = "v0.1.0"
	defaultDcrdexManifestFilename = "dexc-" + defaultDcrdexManifestVersion +
		"-manifest.txt"
	defaultDcrdexManifestURI = "https://github.com/decred/decred-binaries" +
		"/releases/download/" + defaultDecredManifestVersion + "/" +
		defaultDcrdexManifestFilename // Yes defaultDecredManifestVersion

	// bitcoin
	defaultBitcoinManifestVersion  = "0.20.1"
	defaultBitcoinManifestFilename = "SHA256SUMS.asc"
	defaultBitcoinManifestURI      = "https://bitcoincore.org/bin/" +
		"bitcoin-core-" + defaultBitcoinManifestVersion + "/" +
		defaultBitcoinManifestFilename

	// dcrinstall itself
	defaultDcrinstallManifestVersion  = "v1.6.0-rc1"
	defaultDcrinstallManifestFilename = "dcrinstall-" +
		defaultDcrinstallManifestVersion + "-manifest.txt"
	defaultDcrinstallManifestURI = "https://github.com/decred/decred-release/releases/download/" +
		defaultDcrinstallManifestVersion + "/" +
		defaultDcrinstallManifestFilename
	//	https://github.com/decred/decred-release/releases/download/v1.6.0-rc1/dcrinstall-v1.6.0-manifest.txt

	// Settings
	tmpDir                string // Directory where files are downloaded to
	destination           string // Base directory where all files land
	latestManifestURI     string // Manifest of manifests filename
	decredManifestURI     string // Decred manifest filename
	decredManifestDigest  string // Decred manifest digest, if used
	dcrdexManifestURI     string // Dcrdex manifest filename
	dcrdexManifestDigest  string // Dcrdex manifest digest, if used
	bitcoinManifestURI    string // Bitcoin manifest filename
	bitcoinManifestDigest string // Bitcoin manifest digest, if used
	tuple                 string // Download tuple
	network               string // Installing for network
	dcrdex                bool   // Install dcrdex
	skipPGP               bool   // Don't download and verify PGP signatures
	quiet                 bool   // Don't output anything but errors

	// Regexp
	decredRE     = regexp.MustCompile(`decred-v[[:digit:]]\.[[:digit:]]\.[[:digit:]][[:print:]]*-manifest\.txt`)
	dexcRE       = regexp.MustCompile(`dexc-v[[:digit:]]\.[[:digit:]]\.[[:digit:]][[:print:]]*-manifest\.txt`)
	bitcoinRE    = regexp.MustCompile(`\/SHA256SUMS.asc`)
	dcrinstallRE = regexp.MustCompile(`dcrinstall-v[[:digit:]]\.[[:digit:]]\.[[:digit:]][[:print:]]*-manifest\.txt`)
)

// downloadManifest downloads the latest manifest and verifies them.
func downloadManifest() error {
	f, err := ioutil.TempFile("", "dcrinstall")
	if err != nil {
		return err
	}
	f.Close()

	// Download latest manifest
	err = DownloadFile(latestManifestURI, f.Name())
	if err != nil {
		return err
	}

	// Check sig
	if !skipPGP {
		err = pgpVerifyAttached(f.Name(), dcrinstallPubkey)
		if err != nil {
			return err
		}
	}

	// Pluck out links
	f, err = os.Open(f.Name())
	if err != nil {
		return err
	}
	defer f.Close()

	var dcrinstallURI, dcrinstallDigest string
	// <sha256> <filename>
	br := bufio.NewReader(f)
	i := 1
	for {
		line, err := br.ReadString('\n')
		if errors.Is(err, io.EOF) {
			break
		}

		var uri, digest *string
		switch {
		case decredRE.MatchString(line):
			uri = &decredManifestURI
			digest = &decredManifestDigest
		case dexcRE.MatchString(line):
			uri = &dcrdexManifestURI
			digest = &dcrdexManifestDigest
		case bitcoinRE.MatchString(line):
			uri = &bitcoinManifestURI
			digest = &bitcoinManifestDigest
		case dcrinstallRE.MatchString(line):
			uri = &dcrinstallURI
			digest = &dcrinstallDigest
		default:
			continue
		}

		a := strings.Fields(line)
		if len(a) != 2 {
			return fmt.Errorf("invalid manifest %v line %v",
				latestManifestURI, i)
		}

		*digest = a[0]
		*uri = a[1]
	}

	if dcrinstallURI == "" || dcrinstallDigest == "" {
		return fmt.Errorf("Invalid dcrinstall, contact maintainers")
	}
	// Deal with dcrinstall versions
	if defaultDcrinstallManifestURI != dcrinstallURI {
		log.Printf("=== dcrinstall must be updated ===")
		log.Println()
		log.Printf("A new version of dcrinstall was detected. " +
			"Dcrinstall must upgraded before continuing")
		log.Println()
		log.Printf("The latest version can be found on 'decred.org'. " +
			"This tool does not print the link for security " +
			"reasons.")
		log.Println()
		log.Printf("Please see 'https://github.com/decred/decred-release'" +
			" for more information")

		return fmt.Errorf("Please update dcrinstall before continuing")
	}

	return nil
}

func _main() error {
	// Username
	var err error
	u, err := user.Current()
	if err != nil {
		return err
	}
	username = u.Username

	// Password
	b := make([]byte, 24)
	_, err = io.ReadFull(rand.Reader, b)
	if err != nil {
		return err
	}
	password = base64.StdEncoding.EncodeToString(b)

	// Flags
	destF := flag.String("dest", filepath.Join(u.HomeDir, "decred"),
		"extract path")
	latestManifestURIF := flag.String("manifest",
		defaultLatestManifestURI, "latest manifest URI")
	decredManifestURIF := flag.String("decredmanifest", "",
		"Decred manifest URI override")
	dcrdexManifestURIF := flag.String("dcrdexmanifest", "",
		"dcrdex manifest URI override")
	bitcoinManifestURIF := flag.String("bitcoinmanifest", "",
		"bitcoin manifest URI override")
	tupleF := flag.String("tuple", defaultTuple,
		"OS-Arch tuple, e.g. windows-amd64")
	dcrdexF := flag.Bool("dcrdex", false, "Install Dcrdex")
	skipPGPF := flag.Bool("skippgp", false, "skip download and "+
		"verification of pgp signatures")
	quietF := flag.Bool("quiet", false, "quiet (default false)")
	flag.Parse()

	// Prepare environment
	destination = cleanAndExpandPath(*destF)
	tuple = *tupleF
	dcrdex = *dcrdexF
	skipPGP = *skipPGPF
	quiet = *quietF

	// Deal with manifest logic
	if *latestManifestURIF == "" {
		// Manifest was cleared so use defaults
		decredManifestURI = defaultDecredManifestURI
		dcrdexManifestURI = defaultDcrdexManifestURI
		bitcoinManifestURI = defaultBitcoinManifestURI
	} else {
		// Download manifest but let cli options override
		latestManifestURI = *latestManifestURIF
		err = downloadManifest()
		if err != nil {
			return err
		}

		if *decredManifestURIF != "" {
			decredManifestURI = *decredManifestURIF
		}
		if *dcrdexManifestURIF != "" {
			dcrdexManifestURI = *dcrdexManifestURIF
		}
		if *bitcoinManifestURIF != "" {
			bitcoinManifestURI = *bitcoinManifestURIF
		}

		log.Printf("Decred manifest URI: %v\n", decredManifestURI)
		log.Printf("DCRDEX manifest URI: %v\n", dcrdexManifestURI)
		log.Printf("Bitcoin manifest URI: %v\n", bitcoinManifestURI)
	}

	// XXX this needs to be come a flag and tested.
	network = "mainnet"

	err = os.MkdirAll(destination, 0700)
	if err != nil {
		return err
	}

	// Setup logging
	lw, err := os.Create(filepath.Join(destination, "dcrinstall.log"))
	if err != nil {
		return err
	}
	if quiet {
		log.SetOutput(lw)
	} else {
		log.SetOutput(io.MultiWriter(os.Stdout, lw))
	}

	return dcrinstall()
}

func main() {
	if err := _main(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
