// Copyright (c) 2020 The Decred developers
// Use of this source code is governed by an ISC license that can be found in
// the LICENSE file.

package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/cf-guardian/guardian/kernel/fileutils"
	"github.com/decred/dcrd/dcrutil"
)

var (
	bitcoinTuple = map[string]string{
		"darwin-amd64":  "osx64",
		"windows-amd64": "win64",
		"linux-amd64":   "x86_64-linux-gnu",
		"linux-arm":     "arm-linux-gnueabihf",
		"linux-arm64":   "aarch64-linux-gnu",
	}

	bitcoinVersionRE = regexp.MustCompile(`[[:digit:]]+\.[[:digit:]]+\.[[:digit:]]+`)

	bitcoinf = []decredFiles{
		{
			Name:            "bitcoin-cli",
			SupportsVersion: true,
		},
		{
			Name:            "bitcoind",
			Config:          "bitcoin.conf",
			SampleMemory:    bitcoinSampleConfig,
			SupportsVersion: true,
		},
	}
)

// downloadBitcoinBundle downloads the bitcoin bundle into the temporary
// directory. It also verifies the that the digest of the downloaded file
// matches the value in the manifest.
func downloadBitcoinBundle(digest, filename string) error {
	// Download bundle
	bitcoinBundleFilename = filepath.Join(tmpDir, filename)
	err := DownloadFile(bitcoinDownloadURI+filename, bitcoinBundleFilename)
	if err != nil {
		return fmt.Errorf("Download bitcoin bundle: %v", err)
	}

	// Verify digest
	err = sha256Verify(bitcoinBundleFilename, digest)
	if err != nil {
		return fmt.Errorf("SHA256 verification failed: %v", err)
	}

	return nil
}

// extractBitcoinBundle extracts the bitcoin bundle into the destination
// directory.
func extractBitcoinBundle() error {
	return extract(bitcoinBundleFilename, destination)
}

// bitcoinFindOS parses the bitcoin manifest and returns the digest and
// filename for the provided tuple.
func bitcoinFindOS(w, manifest string) (string, string, error) {
	which, ok := bitcoinTuple[w]
	if !ok {
		return "", "", fmt.Errorf("unsuported tuple: %v", w)
	}

	f, err := os.Open(manifest)
	if err != nil {
		return "", "", err
	}
	defer f.Close()

	// <sha256> <filename>
	br := bufio.NewReader(f)
	i := 1
	for {
		line, err := br.ReadString('\n')
		if errors.Is(err, io.EOF) {
			break
		}
		line = strings.TrimSpace(line)

		if !strings.Contains(line, which) {
			continue
		}

		a := strings.Fields(line)
		if len(a) != 2 {
			return "", "", fmt.Errorf("invalid manifest %v line %v",
				manifest, i)
		}

		// Work around windows setup. For example:
		// bitcoin-0.20.1-win64-setup.exe
		if strings.ToLower(filepath.Ext(a[1])) == ".exe" {
			continue
		}

		return a[0], a[1], nil
	}

	return "", "", fmt.Errorf("not found: %v", which)
}

// preconditionsBitcoinInstall determines if the tool is capable of installing
// the bitcoin bundle. It asserts that:
//   * no bitcoin daemons are running
//   * all the installed files have the same version
//   * either all or none of the config files exist
func preconditionsBitcoinInstall() error {
	if runtimeTuple() != tuple {
		log.Printf("Bitcoin installing on foreign OS, " +
			"skipping runtime checks")
		return nil
	}

	// Abort if a daemon is still running
	var isRunningList []string
	for k := range bitcoinf {
		name := bitcoinf[k].Name
		ok, err := isRunning(name)
		if err != nil {
			return fmt.Errorf("isRunning: %v", err)
		}
		if ok {
			log.Printf("Currently running: %v", name)
			isRunningList = append(isRunningList, name)
		} else {
			log.Printf("Currently NOT running: %v", name)
		}
	}
	if len(isRunningList) > 0 {
		return fmt.Errorf("Processess still running: %v",
			isRunningList)
	}

	// Determine current state
	currentlyInstalled := 0
	expectedInstalled := 0
	currentVersion := make(map[string][]string)
	for k := range bitcoinf {
		name := bitcoinf[k].Name
		filename := filepath.Join(destination, name)

		if !bitcoinf[k].SupportsVersion {
			continue
		}

		expectedInstalled++

		// Record current version
		cmd := exec.Command(filename, "--version")
		version, err := cmd.CombinedOutput()
		if err != nil {
			log.Printf("Currently not installed: %v", name)
			continue
		}
		v, err := extractSemVer(string(version))
		if err != nil {
			return fmt.Errorf("invalid version %v: %v",
				name, err)
		}
		ver := v.String()
		log.Printf("Version installed %v: %v", name, ver)
		currentlyInstalled++
		currentVersion[ver] = append(currentVersion[ver], name)
	}

	// XXX this is commented out because we mix versions.
	// Reject mixed versions
	//if len(currentVersion) > 1 {
	//	return fmt.Errorf("Multiple versions of binaries found; " +
	//		"upgrade requires human intervention")
	//}

	// If the current version is already installed error out.
	//for k := range currentVersion {
	//	log.Printf("k %v", k)
	//	// XXX check manifest version against currentVersion
	//	panic("fixme")
	//}

	// Determine if everything or nothing is installed
	if currentlyInstalled != 0 && currentlyInstalled != expectedInstalled {
		return fmt.Errorf("Currently installed binaries does not " +
			"match expected installed binaries; upgrade requires " +
			"human intervention")
	}

	// Install config files if applicable
	currentConfigFiles := 0
	expectedConfigFiles := 0
	var installedConfigs, notInstalledConfigs []string
	for k := range bitcoinf {
		if bitcoinf[k].Config == "" {
			continue
		}

		expectedConfigFiles++

		name := bitcoinf[k].Name
		dir := dcrutil.AppDataDir(name, true)
		filename := filepath.Join(dir, bitcoinf[k].Config)
		if exists(filename) {
			log.Printf("Config %s -- already installed", filename)
			currentConfigFiles++
			installedConfigs = append(installedConfigs, filename)
			continue
		}
		log.Printf("Config %s -- NOT installed", filename)
		notInstalledConfigs = append(notInstalledConfigs, filename)
	}

	if currentConfigFiles != 0 && currentConfigFiles != expectedConfigFiles {
		return fmt.Errorf("dcrinstall requires all or none of the "+
			"configuration files files to be installed. This is "+
			"to prevent improper installations or upgrades. This "+
			"upgrade/install requires human intervention.\n\n%v",
			printConfigError(installedConfigs, notInstalledConfigs))
	}

	// We can now create config files in their respective directories and
	// install the binaries into destination.

	return nil
}

// bitcoinDownloadAndVerify downloads, verifies and asserts that the bitcoin
// bundle can be safely upgraded. This function asserts that all preconditions
// are met before being able to proceed with the bitcoin bundle install.
func bitcoinDownloadAndVerify() error {
	// Download the bitcoin manifest
	manifestBitcoinFilename = filepath.Join(tmpDir,
		filepath.Base(bitcoinManifestURI))
	err := DownloadFile(bitcoinManifestURI, manifestBitcoinFilename)
	if err != nil {
		return fmt.Errorf("Download bitcoin manifest file: %v", err)
	}
	if bitcoinManifestDigest != "" {
		// Optional digest was set so check it
		err = sha256Verify(manifestBitcoinFilename,
			bitcoinManifestDigest)
		if err != nil {
			return fmt.Errorf("SHA256 of bitcoin manifest "+
				"verification failed: %v", err)
		}
	}
	bitcoinDownloadURI, err = getDownloadURI(bitcoinManifestURI)
	if err != nil {
		return fmt.Errorf("Get download URI: %v", err)
	}

	if !skipPGP {
		// Verify bitcoin manifest embedded signature
		err = pgpVerifyAttached(manifestBitcoinFilename, bitcoinPubkey)
		if err != nil {
			// XXX golang pgp does not supprt this curve so just
			// warn.

			log.Printf("Can't verify bitcoin manifest: %v", err)

			postProcess = append(postProcess, "\nThe "+
				"bitcoin signature error that was logged is "+
				"expected.\n\nThe validity of the bitcoin "+
				"archive has been validated.\n\n")
			//return fmt.Errorf("manifest PGP signature incorrect: "+
			//	"%v", err)
		}
	}

	digest, filename, err := bitcoinFindOS(tuple, manifestBitcoinFilename)
	if err != nil {
		return fmt.Errorf("Find tuple: %v", err)
	}
	ver := bitcoinVersionRE.FindString(filepath.Base(filename))
	if ver == "" {
		return fmt.Errorf("Can't Extract bitcoin version from " +
			"manifest")
	}
	_ = digest
	manifestBitcoinVersion = ver
	log.Printf("Attempting to upgrade to Bitcoin version: %v",
		manifestBitcoinVersion)

	// Download bitcoin bundle
	err = downloadBitcoinBundle(digest, filename)
	if err != nil {
		return fmt.Errorf("Download bitcoin bundle: %v", err)
	}

	err = extractBitcoinBundle()
	if err != nil {
		return fmt.Errorf("Extract bitcoin bundle: %v", err)
	}

	err = preconditionsBitcoinInstall()
	if err != nil {
		return fmt.Errorf("Pre bitcoin install: %v", err)
	}

	return nil
}

func installBitcoinBundleConfig() error {
	if runtimeTuple() != tuple {
		log.Printf("Bitcoin bundle installation on foreign OS, " +
			"skipping configuration")
		return nil
	}

	// Install config files
	for k := range bitcoinf {
		if bitcoinf[k].Config == "" {
			continue
		}

		// Check if the config file is already installed.
		name := bitcoinf[k].Name
		dir := dcrutil.AppDataDir(name, true)
		dst := filepath.Join(dir, bitcoinf[k].Config)
		if exists(dst) {
			continue
		}

		var overrides []override
		switch name {
		default:
			overrides = []override{
				{name: "#rpcuser=", content: username},
				{name: "#rpcpassword=", content: password},
				{name: "#server=", content: "1"},
				{name: "#prune=", content: "550"},
				{name: "#debug=", content: "rpc"},
			}
		}
		// XXX add testnet and simnet support

		// Install config file
		conf, err := createConfigFromMemory(bitcoinf[k].SampleMemory,
			overrides)
		if err != nil {
			return err
		}

		log.Printf("Creating directory: %v", dir)
		err = os.MkdirAll(dir, 0700)
		if err != nil {
			return err
		}

		log.Printf("Installing configuration file: %v", dst)
		err = ioutil.WriteFile(dst, []byte(conf), 0600)
		if err != nil {
			return err
		}
	}

	return nil
}

// installBitcoinBundle install all the bitcoin files. This call is only allowed
// if all bitcoin installation preconditions have been met.
func installBitcoinBundle() error {
	err := installBitcoinBundleConfig()
	if err != nil {
		return err
	}

	// Install binaries
	for k := range bitcoinf {
		name := bitcoinf[k].Name
		src := filepath.Join(destination,
			"bitcoin-"+manifestBitcoinVersion, "bin", name)
		dst := filepath.Join(destination, name)
		// yep, this is ferrealz
		if strings.HasPrefix(tuple, "windows") {
			src += ".exe"
			dst += ".exe"
		}

		//log.Printf("Installing %v -> %v\n", src, dst)
		fu := fileutils.New()
		if !fu.Exists(src) {
			return fmt.Errorf("file not found: %v", src)
		}
		if fu.Exists(dst) {
			err := os.RemoveAll(dst)
			if err != nil {
				return fmt.Errorf("Can't remove installed "+
					"file: %v", err)
			}
		}
		log.Printf("Installing: %v", dst)
		err := fu.Copy(dst, src)
		if err != nil {
			return err
		}

		os.Chmod(dst, 0755) // Best effort is fine
	}

	return nil
}
