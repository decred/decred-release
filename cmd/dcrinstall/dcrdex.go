// Copyright (c) 2022 The Decred developers
// Use of this source code is governed by an ISC license that can be found in
// the LICENSE file.

package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/decred/dcrd/dcrutil/v4"
)

var (
	dexf = []decredFiles{
		{
			Name:            "bwctl",
			Config:          "dexcctl.conf",
			SampleMemory:    dexcctlSampleConfig,
			SupportsVersion: true,
		},
		{
			Name:            "bisonw",
			Config:          "dexc.conf",
			SampleMemory:    dexcSampleConfig,
			SupportsVersion: true,
		},
	}
)

// extractDcrdexBundle extracts the dcrdex bundle into the destination
// directory.
func extractDcrdexBundle() error {
	return extract(dcrdexBundleFilename, destination)
}

// downloadDcrdexBundle downloads the dcrdex bundle into the temporary
// directory. It also verifies the that the digest of the downloaded file
// matches the value in the manifest.
func downloadDcrdexBundle(digest, filename string) error {
	// Download bundle
	dcrdexBundleFilename = filepath.Join(tmpDir, filename)
	err := DownloadFile(dcrdexDownloadURI+filename, dcrdexBundleFilename)
	if err != nil {
		return fmt.Errorf("Download dcrdex bundle: %v", err)
	}

	// Verify digest
	err = sha256Verify(dcrdexBundleFilename, digest)
	if err != nil {
		return fmt.Errorf("SHA256 verification failed: %v", err)
	}

	return nil
}

// preconditionsDcrdexInstall determines if the tool is capable of installing
// the dcrdex bundle. It asserts that:
//   - no dcrdex daemons are running
//   - all the installed files have the same version
//   - either all or none of the config files exist
func preconditionsDcrdexInstall() error {
	if runtimeTuple() != tuple {
		log.Printf("DCRDEX bundle installation on foreign OS, " +
			"skipping runtime checks")
		return nil
	}

	// Abort if a daemon is still running
	var isRunningList []string
	for k := range dexf {
		if dexf[k].Directory {
			continue
		}
		ok, err := isRunning(dexf[k].Name)
		if err != nil {
			return fmt.Errorf("isRunning: %v", err)
		}
		if ok {
			log.Printf("Currently running: %v", dexf[k].Name)
			isRunningList = append(isRunningList, dexf[k].Name)
		} else {
			log.Printf("Currently NOT running: %v", dexf[k].Name)
		}
	}
	if !allowRunning && len(isRunningList) > 0 {
		return fmt.Errorf("Processes still running: %v", isRunningList)
	}

	// Determine current state
	currentlyInstalled := 0
	expectedInstalled := 0
	currentVersion := make(map[string][]string)
	var installedBins, notInstalledBins []string
	for k := range dexf {
		filename := filepath.Join(destination, dexf[k].Name)

		if !dexf[k].SupportsVersion {
			continue
		}

		expectedInstalled++

		// Record current version
		cmd := exec.Command(filename, "--version")
		version, err := cmd.CombinedOutput()
		if err != nil {
			log.Printf("Currently not installed: %v", dexf[k].Name)
			notInstalledBins = append(notInstalledBins, filename)
			continue
		}
		v, err := extractSemVer(string(version))
		if err != nil {
			return fmt.Errorf("invalid version %v: %v",
				dexf[k].Name, err)
		}
		ver := v.String()
		log.Printf("Version installed %v: %v", dexf[k].Name, ver)
		currentlyInstalled++
		currentVersion[ver] = append(currentVersion[ver], dexf[k].Name)
		installedBins = append(installedBins, filename)
	}

	// Determine if everything or nothing is installed
	if currentlyInstalled != 0 && currentlyInstalled != expectedInstalled {
		return fmt.Errorf("dcrinstall requires all or none of the "+
			"binary files to be installed. This is "+
			"to prevent improper installations or upgrades. This "+
			"upgrade/install requires human intervention.\n\n%v",
			printConfigError(installedBins, notInstalledBins))
	}

	// Install config files if applicable
	currentConfigFiles := 0
	expectedConfigFiles := 0
	var installedConfigs, notInstalledConfigs []string
	for k := range dexf {
		if dexf[k].Config == "" {
			continue
		}

		expectedConfigFiles++

		ext := filepath.Ext(dexf[k].Config)
		name := strings.TrimSuffix(dexf[k].Config, ext)
		dir := dcrutil.AppDataDir(name, false)
		filename := filepath.Join(dir, dexf[k].Config)
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
			"configuration files to be installed. This is "+
			"to prevent improper installations or upgrades. This "+
			"upgrade/install requires human intervention.\n\n%v",
			printConfigError(installedConfigs, notInstalledConfigs))
	}

	// We can now create config files in their respective directories and
	// install the binaries into destination.

	return nil
}

// dcrdexDownloadAndVerify downloads, verifies and asserts that the dcrdex
// bundle can be safely upgraded. This function asserts that all preconditions
// are met before being able to proceed with the dcrdex bundle install.
func dcrdexDownloadAndVerify() error {
	// Download the dcrdex manifest
	manifestDcrdexFilename = filepath.Join(tmpDir,
		filepath.Base(dcrdexManifestURI))
	err := DownloadFile(dcrdexManifestURI, manifestDcrdexFilename)
	if err != nil {
		return fmt.Errorf("Download dcrdex manifest file: %v", err)
	}
	if dcrdexManifestDigest != "" {
		// Optional digest was set so check it
		err = sha256Verify(manifestDcrdexFilename, dcrdexManifestDigest)
		if err != nil {
			return fmt.Errorf("SHA256 of dcrdex manifest "+
				"verification failed: %v", err)
		}
	}
	dcrdexDownloadURI, err = getDownloadURI(dcrdexManifestURI)
	if err != nil {
		return fmt.Errorf("Get download URI: %v", err)
	}

	if !skipPGP {
		// Download the dcrdex manifest signature
		manifestDcrdexSignatureFilename = filepath.Join(tmpDir,
			filepath.Base(dcrdexManifestURI)+".asc")
		err = DownloadFile(dcrdexManifestURI+".asc",
			manifestDcrdexSignatureFilename)
		if err != nil {
			return fmt.Errorf("Download manifest signature file: "+
				"%v", err)
		}

		// Verify dcrdex manifest signature
		err = pgpVerify(manifestDcrdexSignatureFilename,
			manifestDcrdexFilename, dcrinstallPubkey)
		if err != nil {
			return fmt.Errorf("manifest PGP signature incorrect: "+
				"%v", err)
		}
	}

	// XXX hack around extractSemVer not working properly by feeding it the
	// filename instead of figuring it out from the URL.
	digest, filename, err := findOS(tuple, manifestDcrdexFilename)
	if err != nil {
		return fmt.Errorf("Find tuple: %v", err)
	}
	ver, err := extractSemVer(filepath.Base(filename))
	if err != nil {
		return fmt.Errorf("Extract dcrdex semver from manifest "+
			"filename %v", err)
	}
	manifestDcrdexVersion = ver.String()
	log.Printf("Attempting to upgrade to Dcrdex version: %v",
		manifestDcrdexVersion)

	// Don't download bundle if it has been extracted.
	if forceDownload || !seenBefore(filename) {
		// Download dcrdex bundle
		err = downloadDcrdexBundle(digest, filename)
		if err != nil {
			return fmt.Errorf("Download dcrdex bundle: %v", err)
		}

		err = extractDcrdexBundle()
		if err != nil {
			return fmt.Errorf("Extract dcrdex bundle: %v", err)
		}
	} else {
		log.Printf("Using cached archive: %v", filename)
	}

	err = preconditionsDcrdexInstall()
	if err != nil {
		return fmt.Errorf("Pre dcrdex install: %v", err)
	}

	return nil
}

func installDcrdexBundleConfig() error {
	if runtimeTuple() != tuple {
		log.Printf("Dcrdex bundle installation on foreign OS, " +
			"skipping configuration")
		return nil
	}

	// Install config files
	for k := range dexf {
		if dexf[k].Config == "" {
			continue
		}

		// Check if the config file is already installed.
		ext := filepath.Ext(dexf[k].Config)
		name := strings.TrimSuffix(dexf[k].Config, ext)
		dir := dcrutil.AppDataDir(name, false)
		dst := filepath.Join(dir, dexf[k].Config)
		if exists(dst) {
			continue
		}

		var overrides []override
		switch dexf[k].Name {
		default:
			overrides = []override{
				{name: "; rpc=", content: "0"},
				{name: "; rpcuser=", content: username},
				{name: "; rpcpass=", content: password},
			}
		}
		// XXX add testnet and simnet support

		// Install config file
		conf, err := createConfigFromMemory(dexf[k].SampleMemory,
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
		err = os.WriteFile(dst, []byte(conf), 0600)
		if err != nil {
			return err
		}
	}

	return nil
}

// installDcrdexBundle install all the dcrdex files. This call is only allowed
// if all dcrdex installation preconditions have been met.
func installDcrdexBundle() error {
	err := installDcrdexBundleConfig()
	if err != nil {
		return err
	}

	// Install binaries
	for k := range dexf {
		src := filepath.Join(destination,
			"bisonwallet-"+tuple+"-"+manifestDcrdexVersion, dexf[k].Name)
		dst := filepath.Join(destination, dexf[k].Name)
		// yep, this is ferrealz
		if !dexf[k].Directory && strings.HasPrefix(tuple, "windows") {
			src += ".exe"
			dst += ".exe"
		}

		//log.Printf("Installing %v -> %v\n", src, dst)
		if !fileExists(src) {
			return fmt.Errorf("file not found: %v", src)
		}
		if fileExists(dst) {
			err := os.RemoveAll(dst)
			if err != nil {
				return fmt.Errorf("Can't remove installed "+
					"file: %v", err)
			}
		}
		log.Printf("Installing: %v", dst)
		err := copyFile(dst, src)
		if err != nil {
			return err
		}

		os.Chmod(dst, 0755) // Best effort is fine
	}

	postProcess = append(postProcess, "\nDCRDEX:\n\n"+
		"Please read the release notes at https://github.com/decred/dcrdex/releases for IMPORTANT NOTICES\n\n")

	return nil
}
