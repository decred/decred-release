// Copyright (c) 2016-2020 The Decred developers
// Use of this source code is governed by an ISC license that can be found in
// the LICENSE file.

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/cf-guardian/guardian/kernel/fileutils"
	"github.com/decred/dcrd/dcrutil"
)

const (
	walletClientsPem = "clients.pem"
	clientPem        = "client.pem"
	clientKey        = "client-key.pem"
	walletDB         = "wallet.db"
	lnWalletDB       = "channel.db"
)

// decredFiles provides all needed bits to perform the install the decred
// bundle.
type decredFiles struct {
	Name            string // Binary filename
	Config          string // Actual config file
	SampleFilename  string // Sample file config file
	SampleMemory    string // Static sample config
	SupportsVersion bool   // Whether or not it supports --version
	Directory       bool   // Whether or not this filename is a directory
}

var (
	df = []decredFiles{
		{
			Name:            "dcrctl",
			Config:          "dcrctl.conf",
			SampleFilename:  "sample-dcrctl.conf",
			SupportsVersion: true,
		},
		{
			Name:            "dcrd",
			Config:          "dcrd.conf",
			SampleFilename:  "sample-dcrd.conf",
			SupportsVersion: true,
		},
		{
			Name:            "dcrwallet",
			Config:          "dcrwallet.conf",
			SampleFilename:  "sample-dcrwallet.conf",
			SupportsVersion: true,
		},
		{
			Name:            "promptsecret",
			SupportsVersion: false,
		},
		{
			Name:            "dcrlnd",
			SupportsVersion: true,
			Config:          "dcrlnd.conf",
			SampleFilename:  "sample-dcrlnd.conf",
		},
		{
			Name:            "dcrlncli",
			SupportsVersion: true,
		},
		{
			Name:            "politeiavoter",
			Config:          "politeiavoter.conf",
			SampleFilename:  "sample-politeiavoter.conf",
			SupportsVersion: true,
		},
		{
			Name:            "gencerts",
			SupportsVersion: false,
		},
	}
)

// generateClientCerts creates politeiavoter client certificates and copies
// them to the wallet directory.
func generateClientCerts() error {
	// Create certificate for politeiavoter
	gencertsExe := filepath.Join(destination,
		"decred-"+tuple+"-"+manifestDecredVersion, "gencerts")
	piDir := dcrutil.AppDataDir("politeiavoter", false)
	piClientCert := filepath.Join(piDir, clientPem)
	piClientKey := filepath.Join(piDir, clientKey)
	log.Printf("Running: %v %v %v", gencertsExe, piClientCert, piClientKey)
	o, err := exec.Command(gencertsExe, piClientCert,
		piClientKey).CombinedOutput()
	if err != nil {
		return fmt.Errorf("error: %w\noutput:\n%v", err, string(o))
	}

	// Copy certificate to dcrwallet
	dst := filepath.Join(dcrutil.AppDataDir("dcrwallet", false),
		walletClientsPem)
	fu := fileutils.New()
	if fu.Exists(dst) {
		// Shouldn't happen
		return fmt.Errorf("file already exists: %v", dst)
	}
	log.Printf("Installing: %v", dst)
	err = fu.Copy(dst, piClientCert)
	if err != nil {
		return err
	}

	return nil
}

// createWallet creates a wallet.
func createWallet(net string) error {
	// create wallet
	log.Printf("Creating lightning wallet: %v", net)

	dcrwalletExe := filepath.Join(destination,
		"decred-"+tuple+"-"+manifestDecredVersion, "dcrwallet")
	args := []string{"--create"}
	switch net {
	case "testnet":
		args = append(args, "--testnet")
	case "simnet":
		args = append(args, "--simnet")
	}
	cmd := exec.Command(dcrwalletExe, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// lnCreateWallet creates a lightning wallet.
func lnCreateWallet(net string) error {
	// create wallet
	log.Printf("Creating lightning wallet: %v", net)

	dcrwalletExe := filepath.Join(destination,
		"decred-"+tuple+"-"+manifestDecredVersion, "dcrlncli")
	args := []string{"create"}
	switch net {
	case "testnet":
		args = append(args, "--testnet")
	case "simnet":
		args = append(args, "--simnet")
	}
	cmd := exec.Command(dcrwalletExe, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// walletDBExists return true if the decred wallet is already created.
func walletDBExists(net string) bool {
	dir := dcrutil.AppDataDir("dcrwallet", false)
	return exists(filepath.Join(dir, net, walletDB))
}

// lnWalletDB return true if the decred lightning wallet is already created.
func lnWalletDBExists(net string) bool {
	dir := dcrutil.AppDataDir("dcrlnd", false)
	return exists(filepath.Join(dir, "data", "graph", net, lnWalletDB))
}

// extractDecredBundle extracts the decred bundle into the destination
// directory.
func extractDecredBundle() error {
	return extract(decredBundleFilename, destination)
}

// downloadDecredBundle downloads the decred bundle into the temporary
// directory. It also verifies the that the digest of the downloaded file
// matches the value in the manifest.
func downloadDecredBundle(digest, filename string) error {
	// Download bundle
	decredBundleFilename = filepath.Join(tmpDir, filename)
	err := DownloadFile(decredDownloadURI+filename, decredBundleFilename)
	if err != nil {
		return fmt.Errorf("Download decred bundle: %v", err)
	}

	// Verify digest
	err = sha256Verify(decredBundleFilename, digest)
	if err != nil {
		return fmt.Errorf("SHA256 verification failed: %v", err)
	}

	return nil
}

// preconditionsDecredInstall determines if the tool is capable of installing
// the decred bundle. It asserts that:
//   * no decred daemons are running
//   * all the installed files have the same version
//   * either all or none of the config files exist
func preconditionsDecredInstall() error {
	if runtimeTuple() != tuple {
		log.Printf("Decred bundle installation on foreign OS, " +
			"skipping runtime checks")
		return nil
	}

	// Abort if a daemon is still running
	var isRunningList []string
	for k := range df {
		if df[k].Directory {
			continue
		}
		ok, err := isRunning(df[k].Name)
		if err != nil {
			return fmt.Errorf("isRunning: %v", err)
		}
		if ok {
			log.Printf("Currently running: %v", df[k].Name)
			isRunningList = append(isRunningList, df[k].Name)
		} else {
			log.Printf("Currently NOT running: %v", df[k].Name)
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
	var installedBins, notInstalledBins []string
	for k := range df {
		filename := filepath.Join(destination, df[k].Name)

		if !df[k].SupportsVersion {
			continue
		}

		expectedInstalled++

		// Record current version
		cmd := exec.Command(filename, "--version")
		version, err := cmd.CombinedOutput()
		if err != nil {
			log.Printf("Currently not installed: %v", df[k].Name)
			notInstalledBins = append(notInstalledBins, filename)
			continue
		}
		v, err := extractSemVer(string(version))
		if err != nil {
			return fmt.Errorf("invalid version %v: %v",
				df[k].Name, err)
		}
		ver := v.String()
		log.Printf("Version installed %v: %v", df[k].Name, ver)
		currentlyInstalled++
		currentVersion[ver] = append(currentVersion[ver], df[k].Name)
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
	for k := range df {
		if df[k].Config == "" {
			continue
		}

		expectedConfigFiles++

		dir := dcrutil.AppDataDir(df[k].Name, false)
		filename := filepath.Join(dir, df[k].Config)
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

// decredDownloadAndVerify downloads, verifies and asserts that the decred
// bundle can be safely upgraded. This function asserts that all preconditions
// are met before being able to proceed with the decred bundle install.
func decredDownloadAndVerify() error {
	// Download the decred manifest
	manifestDecredFilename = filepath.Join(tmpDir,
		filepath.Base(decredManifestURI))

	err := DownloadFile(decredManifestURI, manifestDecredFilename)
	if err != nil {
		return fmt.Errorf("Download manifest file: %v", err)
	}
	if decredManifestDigest != "" {
		// Optional digest was set so check it
		err = sha256Verify(manifestDecredFilename, decredManifestDigest)
		if err != nil {
			return fmt.Errorf("SHA256 of decred manifest "+
				"verification failed: %v", err)
		}
	}
	decredDownloadURI, err = getDownloadURI(decredManifestURI)
	if err != nil {
		return fmt.Errorf("Get download URI: %v", err)
	}

	if !skipPGP {
		// Download the decred manifest signature
		manifestDecredSignatureFilename = filepath.Join(tmpDir,
			filepath.Base(decredManifestURI)+".asc")
		err = DownloadFile(decredManifestURI+".asc",
			manifestDecredSignatureFilename)
		if err != nil {
			return fmt.Errorf("Download manifest signature file: "+
				"%v", err)
		}

		// Verify decred manifest signature
		err = pgpVerify(manifestDecredSignatureFilename,
			manifestDecredFilename, dcrinstallPubkey)
		if err != nil {
			return fmt.Errorf("manifest PGP signature incorrect: "+
				"%v", err)
		}
	}

	// XXX hack around extractSemVer not working properly by feeding it the
	// filename instead of figuring it out from the URL.
	digest, filename, err := findOS(tuple, manifestDecredFilename)
	if err != nil {
		return fmt.Errorf("Find tuple: %v", err)
	}
	ver, err := extractSemVer(filepath.Base(filename))
	if err != nil {
		return fmt.Errorf("Extract decred semver from manifest "+
			"filename %v", err)
	}
	manifestDecredVersion = ver.String()
	log.Printf("Attempting to upgrade to Decred version: %v",
		manifestDecredVersion)

	// Don't download bundle if it has been extracted.
	if forceDownload || !seenBefore(filename) {
		// Download decred bundle
		err = downloadDecredBundle(digest, filename)
		if err != nil {
			return fmt.Errorf("Download decred bundle: %v", err)
		}

		err = extractDecredBundle()
		if err != nil {
			return fmt.Errorf("Extract decred bundle: %v", err)
		}
	} else {
		log.Printf("Using cached archive: %v", filename)
	}

	err = preconditionsDecredInstall()
	if err != nil {
		return fmt.Errorf("Pre decred install: %v", err)
	}

	return nil
}

func installDecredBundleConfig() error {
	if runtimeTuple() != tuple {
		log.Printf("Decred bundle installation on foreign OS, " +
			"skipping configuration")
		return nil
	}

	// Install config files
	for k := range df {
		if df[k].Config == "" {
			continue
		}

		// Check if the config file is already installed.
		dir := dcrutil.AppDataDir(df[k].Name, false)
		dst := filepath.Join(dir, df[k].Config)
		if exists(dst) {
			continue
		}

		var overrides []override
		switch df[k].Name {
		case "dcrwallet":
			overrides = []override{
				{name: "; username=", content: username},
				{name: "; password=", content: password},
			}
		case "dcrlnd":
			overrides = []override{
				{name: "; dcrd.rpcuser=", content: username},
				{name: "; dcrd.rpcpass=", content: password},
			}
		default:
			overrides = []override{
				{name: "; rpcuser=", content: username},
				{name: "; rpcpass=", content: password},
			}
		}
		// XXX add testnet and simnet support

		// Install config file
		src := filepath.Join(destination,
			"decred-"+tuple+"-"+manifestDecredVersion,
			df[k].SampleFilename)
		conf, err := createConfigFromFile(src, overrides)
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

	// Check client certs.
	walletCert := appFileExists("dcrwallet", walletClientsPem)
	piCert := appFileExists("politeiavoter", clientPem)
	piKey := appFileExists("politeiavoter", clientKey)
	switch {
	case walletCert && piCert && piKey:
		log.Printf("Client certs exist, skipping client cert " +
			"generation.")
	case !walletCert && !piCert && !piKey:
		err := generateClientCerts()
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("Can't determine client certificate" +
			" state, must perform manual upgrade")
	}

	// Check if wallet exists.
	if walletDBExists(network) {
		log.Printf("Wallet exists, skipping creation.")
	} else {
		err := createWallet(network)
		if err != nil {
			return fmt.Errorf("Can't create wallet: %v", err)
		}
	}

	// Check if lin wallet exists.
	if lnWalletDBExists(network) {
		log.Printf("Lightning wallet exists, skipping creation.")
	} else {
		// XXX can't create a lightning wallet without dcrlnd running.
		// Tell the user at the end of the install what commands to
		// run.
		log.Printf("Lightning wallet does not exist.")

		lndw := filepath.Join(destination, "dcrlncli")
		postProcess = append(postProcess, fmt.Sprintf("\nThe lightning "+
			"wallet could not be automatically created.\n\n"+
			"To create a lightning wallet:\n"+
			"* Start dcrlnd\n"+
			"* Run '%v create'\n\n", lndw))

		//err := lnCreateWallet(network)
		//if err != nil {
		//	return fmt.Errorf("Can't create wallet: %v", err)
		//}
	}

	return nil
}

// installDecredBundle install all the decred files. This call is only allowed
// if all decred installation preconditions have been met.
func installDecredBundle() error {
	err := installDecredBundleConfig()
	if err != nil {
		return err
	}

	// Install binaries
	for k := range df {
		src := filepath.Join(destination,
			"decred-"+tuple+"-"+manifestDecredVersion, df[k].Name)
		dst := filepath.Join(destination, df[k].Name)
		// yep, this is ferrealz
		if !df[k].Directory && strings.HasPrefix(tuple, "windows") {
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
