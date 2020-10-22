// Copyright (c) 2016-2020 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/cf-guardian/guardian/kernel/fileutils"
	"github.com/decred/dcrd/dcrutil/v3"
)

// global context
type ctx struct {
	s *Settings

	user     string
	password string

	logFilename string
}

type binary struct {
	Name            string // binary filename
	Config          string // actual config file
	Example         string // example config file
	SupportsVersion bool   // whether or not it supports --version

	// Dex
	Dcrdex bool   // only install if user set dcrdex to true
	Conf   string // sample config
	Copy   bool   // Only copy file, bitcoin only
}

const (
	walletClientsPem = "clients.pem"
	clientPem        = "client.pem"
	clientKey        = "client-key.pem"
)

var (
	binaries = []binary{
		{
			Name:            "dcrctl",
			Config:          "dcrctl.conf",
			Example:         "sample-dcrctl.conf",
			SupportsVersion: true,
		},
		{
			Name:            "dcrd",
			Config:          "dcrd.conf",
			Example:         "sample-dcrd.conf",
			SupportsVersion: true,
		},
		{
			Name:            "dcrwallet",
			Config:          "dcrwallet.conf",
			Example:         "sample-dcrwallet.conf",
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
			Example:         "sample-dcrlnd.conf",
		},
		{
			Name:            "dcrlncli",
			SupportsVersion: true,
		},
		{
			Name:            "politeiavoter",
			Config:          "politeiavoter.conf",
			Example:         "sample-politeiavoter.conf",
			SupportsVersion: true,
		},
		{
			Name:            "gencerts",
			SupportsVersion: false,
		},
		{
			Name:   "bitcoin",
			Config: "bitcoin.conf",
			Dcrdex: true,
			Conf:   bitcoinConf,
		},
		{
			Name:   "bin/bitcoin-cli",
			Dcrdex: true,
			Copy:   true,
		},
		{
			Name:   "bin/bitcoin-qt",
			Dcrdex: true,
			Copy:   true,
		},
		{
			Name:   "bin/bitcoin-tx",
			Dcrdex: true,
			Copy:   true,
		},
		{
			Name:   "bin/bitcoin-wallet",
			Dcrdex: true,
			Copy:   true,
		},
		{
			Name:   "bin/bitcoind",
			Dcrdex: true,
			Copy:   true,
		},
		{
			Name:   "dexc",
			Config: "dexc.conf",
			Dcrdex: true,
			Conf:   dexcConf,
		},
		{
			Name:   "dexcctl",
			Config: "dexcctl.conf",
			Dcrdex: true,
			Conf:   dexcctlConf,
		},
		{
			Name:   "dexc",
			Dcrdex: true,
			Copy:   true,
		},
		{
			Name:   "dexcctl",
			Dcrdex: true,
			Copy:   true,
		},
		{
			Name:   "site",
			Dcrdex: true,
			Copy:   true,
		},
	}

	// Bitcoin core constants
	bitcoinVersion = "0.20.1"
	bitcoinLink    = "https://bitcoin.org/bin/bitcoin-core-" +
		bitcoinVersion + "/"
	bitcoinFilename  = "bitcoin-" + bitcoinVersion
	bitcoinPrefix    = bitcoinLink + bitcoinFilename
	bitcoinDownloads = map[string]string{
		"darwin-amd64":  bitcoinPrefix + "-osx64.tar.gz",
		"windows-amd64": bitcoinPrefix + "-win64.zip",
		"linux-amd64":   bitcoinPrefix + "-x86_64-linux-gnu.tar.gz",
		"linux-arm":     bitcoinPrefix + "-arm-linux-gnueabihf.tar.gz",
		"linux-arm64":   bitcoinPrefix + "-aarch64-linux-gnu.tar.gz",
	}

	// Dex
	dexVersion   = "v0.1.0"
	dexLink      = defaultURI + "/"
	dexDownloads = map[string]string{
		"darwin-amd64": dexLink + "dexc-darwin-amd64-" + dexVersion +
			".tar.gz",
		"linux-amd64": dexLink + "dexc-linux-amd64-" + dexVersion +
			".tar.gz",
		"windows-amd64": dexLink + "dexc-windows-amd64-" + dexVersion +
			".zip",
	}
)

func (c *ctx) logNoTime(format string, args ...interface{}) error {
	f, err := os.OpenFile(c.logFilename, os.O_CREATE|os.O_RDWR|os.O_APPEND,
		0600)
	if err != nil {
		return err
	}
	defer f.Close()

	if c.s.Verbose {
		fmt.Printf(format, args...)
	}

	_, err = fmt.Fprintf(f, format, args...)
	return err
}

func (c *ctx) log(format string, args ...interface{}) error {
	t := time.Now().Format("15:04:05.000 ")
	return c.logNoTime(t+format, args...)
}

func (c *ctx) obtainUserName() error {
	u, err := user.Current()
	if err != nil {
		return err
	}
	c.user = u.Username
	return nil
}

func (c *ctx) obtainPassword() error {
	b := make([]byte, 24)
	_, err := io.ReadFull(rand.Reader, b[:])
	if err != nil {
		return err
	}

	// convert password to something readable
	c.password = base64.StdEncoding.EncodeToString(b)

	return nil
}

// findOS itterates over the entire manifest and plucks out the digest and
// filename of the providede os-arch tuple.  The tupple must be unique.
func findOS(which, manifest string) (string, string, error) {
	var digest, filename string

	f, err := os.Open(manifest)
	if err != nil {
		return "", "", err
	}
	defer f.Close()

	br := bufio.NewReader(f)
	i := 1
	for {
		line, err := br.ReadString('\n')
		if err == io.EOF {
			break
		}
		line = strings.TrimSpace(line)

		a := strings.Fields(line)
		if len(a) != 2 {
			return "", "", fmt.Errorf("invalid manifest %v line %v",
				manifest, i)
		}

		// add "-" to disambiguate arm and arm64
		if !strings.Contains(a[1], which+"-") {
			continue
		}

		if !(digest == "" && filename == "") {
			return "", "",
				fmt.Errorf("os-arch tuple not unique: %v", which)
		}

		digest = strings.TrimSpace(a[0])
		filename = strings.TrimSpace(a[1])
	}

	return digest, filename, nil
}

func btcFindOS(which, manifest string) (string, error) {
	f, err := os.Open(manifest)
	if err != nil {
		return "", err
	}
	defer f.Close()

	// <sha256> <filename>
	br := bufio.NewReader(f)
	i := 1
	for {
		line, err := br.ReadString('\n')
		if err == io.EOF {
			break
		}
		line = strings.TrimSpace(line)

		if !strings.Contains(line, which) {
			continue
		}

		a := strings.Fields(line)
		if len(a) != 2 {
			return "", fmt.Errorf("invalid manifest %v line %v",
				manifest, i)
		}

		return a[0], nil
	}

	return "", fmt.Errorf("not found: %v", which)
}

// download downloads the manifest, the manifest signature and the selected
// os-arch package to a temporary directory.  It returns the temporary
// directory if there is no failure.
func (c *ctx) download() (string, error) {
	// create temporary directory
	td, err := ioutil.TempDir("", "decred")
	if err != nil {
		return "", err
	}

	// download manifest
	manifestURI := c.s.URI + "/" + c.s.Manifest
	c.log("temporary directory: %v\n", td)
	c.log("downloading manifest: %v\n", manifestURI)

	manifest := filepath.Join(td, filepath.Base(manifestURI))
	err = downloadToFile(manifestURI, manifest)
	if err != nil {
		return "", err
	}

	if !c.s.SkipAsc {
		// download manifest signature
		manifestAscURI := c.s.URI + "/" + c.s.Manifest + ".asc"
		c.log("downloading manifest signatures: %v\n",
			manifestAscURI)
		manifestAsc := filepath.Join(td, filepath.Base(manifestAscURI))
		err = downloadToFile(manifestAscURI, manifestAsc)
		if err != nil {
			return "", err
		}
	}

	// determine if os-arch is supported
	_, filename, err := findOS(c.s.Tuple, manifest)
	if err != nil {
		return "", err
	}

	// download requested package
	packageURI := c.s.URI + "/" + filename
	c.log("downloading package: %v\n", packageURI)

	pkg := filepath.Join(td, filepath.Base(packageURI))
	err = downloadToFile(packageURI, pkg)
	if err != nil {
		return "", err
	}

	// Deal with bitcoin core
	if c.s.Dcrdex && c.s.BitcoinURI != "" {
		c.log("downloading bitcoin core: %v\n", c.s.BitcoinURI)
		pkg := filepath.Join(td, filepath.Base(c.s.BitcoinURI))
		err = downloadToFile(c.s.BitcoinURI, pkg)
		if err != nil {
			return "", err
		}

		// download bitcoin manifest
		btcManifestURI := bitcoinLink + "SHA256SUMS.asc"
		c.log("downloading bitcoin manifest: %v\n",
			btcManifestURI)
		btcManifestAsc := filepath.Join(td,
			filepath.Base(btcManifestURI))
		err = downloadToFile(btcManifestURI, btcManifestAsc)
		if err != nil {
			return "", err
		}
	}

	// Deal with dcrdex
	if c.s.Dcrdex && c.s.DexURI != "" {
		c.log("downloading dcrdex: %v\n", c.s.DexURI)
		pkg := filepath.Join(td, filepath.Base(c.s.DexURI))
		err = downloadToFile(c.s.DexURI, pkg)
		if err != nil {
			return "", err
		}
	}

	return td, nil
}

// verify verifies the manifest signature and the package digest.
func (c *ctx) verify() error {
	// determine if os-arch is supported
	manifest := filepath.Join(c.s.Path, c.s.Manifest)
	digest, filename, err := findOS(c.s.Tuple, manifest)
	if err != nil {
		return err
	}

	if !c.s.SkipAsc {
		// verify manifest
		c.log("verifying manifest: %v ", c.s.Manifest)
		err = pgpVerify(manifest+".asc", manifest, pubkey)
		if err != nil {
			c.logNoTime("FAIL\n")
			return fmt.Errorf("manifest PGP signature incorrect: "+
				"%v", err)
		}
		c.logNoTime("OK\n")
	}

	if c.s.Dcrdex && !c.s.SkipAsc {
		// verify bitcoin manifest
		btcManifestAsc := filepath.Join(c.s.Path, "SHA256SUMS.asc")
		c.log("verifying bitcoin manifest: %v ", "SHA256SUMS.asc")
		// The signarure is embedded in the manifest.
		err = pgpVerify(btcManifestAsc, btcManifestAsc, btcPubkey)
		if err != nil {
			// XXX cannot verify PGP signature because the manifest
			// uses an unsuported curve

			//c.logNoTime("FAIL\n")
			c.logNoTime("bitcoin signature cannot be verified: %v\n",
				fmt.Errorf("manifest PGP signature "+
					"incorrect: %v", err))
		} else {
			c.logNoTime("OK\n")
		}
	}

	// verify digest
	c.log("verifying package: %v ", filename)

	pkg := filepath.Join(c.s.Path, filename)
	d, err := sha256File(pkg)
	if err != nil {
		return err
	}

	// verify package digest
	if hex.EncodeToString(d) != digest {
		c.logNoTime("FAILED\n")
		c.log("%v %v\n", hex.EncodeToString(d), digest)

		return fmt.Errorf("corrupt digest %v", filename)
	}

	c.logNoTime("OK\n")

	// Verify bitcoind sha
	if c.s.Dcrdex {
		btcFile := filepath.Base(c.s.BitcoinURI)
		btcManifestAsc := filepath.Join(c.s.Path, "SHA256SUMS.asc")
		btcDigest, err := btcFindOS(btcFile, btcManifestAsc)
		if err != nil {
			return err
		}

		btcPkg := filepath.Join(c.s.Path, btcFile)
		c.log("verifying package: %v ", btcPkg)

		d, err := sha256File(btcPkg)
		if err != nil {
			return err
		}

		// verify package digest
		if hex.EncodeToString(d) != btcDigest {
			c.logNoTime("FAILED\n")
			c.log("%v %v\n", hex.EncodeToString(d), btcDigest)

			return fmt.Errorf("corrupt digest %v", btcFile)
		}
	}

	return nil
}

// copy installs all binaries into their final destination.
func (c *ctx) copy(version string) error {
	for _, v := range binaries {
		if v.Dcrdex {
			if !c.s.Dcrdex {
				continue
			}
			if !v.Copy {
				// Nothing to do
				continue
			}

			// Install bitcoin/dex binary
			var src string
			if strings.Contains(v.Name, "dexc") ||
				strings.Contains(v.Name, "site") {
				src = filepath.Join(c.s.Destination,
					"dexc-"+c.s.Tuple, v.Name)
			} else {
				src = filepath.Join(c.s.Destination,
					bitcoinFilename, v.Name)
			}
			dst := filepath.Join(c.s.Destination,
				filepath.Base(v.Name))

			// XXX ARRGHHHGGHGHG
			if !strings.Contains(v.Name, "site") {
				// Deal with windows.
				if runtime.GOOS == "windows" {
					src += ".exe"
					dst += ".exe"
				}
			}

			c.log("dex files installing %v -> %v\n", src, dst)
			fu := fileutils.New()
			if fu.Exists(dst) {
				err := os.RemoveAll(dst)
				if err != nil {
					return err
				}
			}
			err := fu.Copy(dst, src)
			if err != nil {
				return err
			}

			continue
		}
		src := filepath.Join(c.s.Destination,
			"decred-"+c.s.Tuple+"-"+version,
			v.Name)
		dst := filepath.Join(c.s.Destination, v.Name)

		// yep, this is ferrealz
		if runtime.GOOS == "windows" {
			src += ".exe"
			dst += ".exe"
		}

		c.log("installing %v -> %v\n", src, dst)
		srcBytes, err := ioutil.ReadFile(src)
		if err != nil {
			return err
		}
		f, err := os.OpenFile(dst, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0700)
		if err != nil {
			return err
		}
		_, err = f.Write(srcBytes)
		f.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

// validate verifies that all binaries can be executed.
func (c *ctx) validate(version string) error {
	for _, v := range binaries {
		if v.Dcrdex {
			// Skip non decred binary execution test
			continue
		}

		// not in love with this, pull this out of tar instead
		filename := filepath.Join(c.s.Destination,
			"decred-"+c.s.Tuple+"-"+version,
			v.Name)

		c.log("checking: %v ", filename)

		cmd := exec.Command(filename, "-h")
		err := cmd.Start()
		if err != nil {
			c.logNoTime("FAILED\n")
			return err
		}

		c.logNoTime("OK\n")

	}
	return nil
}

func (c *ctx) running(name string) (bool, error) {
	if c.s.DownloadOnly {
		return false, nil
	}

	return c.isRunning(name)
}

// recordCurrent iterates over binaries and records their version number in
// the log file.
func (c *ctx) recordCurrent() error {
	for _, v := range binaries {
		if !v.SupportsVersion {
			continue
		}

		// not in love with this, pull this out of tar instead
		filename := filepath.Join(c.s.Destination, v.Name)

		c.log("current version: %v ", filename)

		cmd := exec.Command(filename, "--version")
		version, err := cmd.CombinedOutput()
		if err != nil {
			c.logNoTime("NOT INSTALLED\n")
			continue
		}

		c.logNoTime("%v\n", strings.TrimSpace(string(version)))

	}

	return nil
}

func (c *ctx) createConfigNormal(b binary, br *bufio.Reader) (string, error) {
	seen := false
	rv := ""
	network := "; " + strings.ToLower(c.s.Net) + "="

	// This is a crime against humanity, but oh well.
	type override struct {
		name    string
		content string
	}
	var (
		overrides []override
		start     int
	)
	switch b.Name {
	case "dcrwallet":
		overrides = []override{
			{name: "; username=", content: c.user},
			{name: "; password=", content: c.password},
		}
		start = 2

	case "bitcoin":
		overrides = []override{
			{name: "#rpcuser=", content: c.user},
			{name: "#rpcpassword=", content: c.password},
			{name: "#server=", content: "1"},
			{name: "#prune=", content: "550"},
			{name: "#debug=", content: "rpc"},
		}
		start = 1

	case "dexc":
		overrides = []override{
			{name: "; rpc=", content: "1"},
			{name: "; rpcuser=", content: c.user},
			{name: "; rpcpass=", content: c.password},
		}
		start = 2

	case "dexctl":
		overrides = []override{
			{name: "; rpcuser=", content: c.user},
			{name: "; rpcpass=", content: c.password},
		}
		start = 2

	default:
		overrides = []override{
			{name: "; rpcuser=", content: c.user},
			{name: "; rpcpass=", content: c.password},
		}
		start = 2
	}

	for {
		line, err := br.ReadString('\n')
		if err == io.EOF {
			break
		}

		if strings.HasPrefix(line, network) {
			line = network[2:] + "1\n"
			seen = true
		}

		for k := range overrides {
			if strings.HasPrefix(line, overrides[k].name) {
				line = overrides[k].name[start:] +
					overrides[k].content + "\n"
			}
		}

		rv += line
	}

	if c.s.Net != netMain {
		if !seen {
			return "", fmt.Errorf("could not set net to %v",
				c.s.Net)
		}
	}

	return rv, nil
}

func (c *ctx) createConfig(b binary, version string) (string, error) {
	if b.Dcrdex {
		c.log("parsing: %v\n", b.Config)
		return c.createConfigNormal(b,
			bufio.NewReader(bytes.NewReader([]byte(b.Conf))))
	}

	// Regular decred files
	sample := filepath.Join(c.s.Destination,
		"decred-"+c.s.Tuple+"-"+version,
		b.Example)

	// read sample config
	f, err := os.Open(sample)
	if err != nil {
		return "", err
	}
	defer f.Close()

	c.log("parsing: %v\n", sample)

	return c.createConfigNormal(b, bufio.NewReader(f))
}

func (c *ctx) writeConfig(b binary, cf string) error {
	dir := dcrutil.AppDataDir(b.Name, false)
	conf := filepath.Join(dir, b.Config)

	c.log("writing: %v\n", conf)

	return ioutil.WriteFile(conf, []byte(cf), 0600)
}

func (c *ctx) walletDBExists() bool {
	dir := dcrutil.AppDataDir("dcrwallet", false)
	return exist(filepath.Join(dir, netMain, walletDB)) ||
		exist(filepath.Join(dir, netTest, walletDB)) ||
		exist(filepath.Join(dir, netSim, walletDB))
}

func (c *ctx) clientFileExists(app, filename string) bool {
	dir := dcrutil.AppDataDir(app, false)
	return exist(filepath.Join(dir, filename))
}

func (c *ctx) generateClientCerts(version string) error {
	// Create certificate for politeiavoter
	gencertsExe := filepath.Join(c.s.Destination,
		"decred-"+c.s.Tuple+"-"+version, "gencerts")
	if runtime.GOOS == "windows" {
		gencertsExe += ".exe"
	}
	piDir := dcrutil.AppDataDir("politeiavoter", false)
	piClientCert := filepath.Join(piDir, clientPem)
	o, err := exec.Command(gencertsExe, piClientCert,
		filepath.Join(piDir, clientKey)).CombinedOutput()
	if err != nil {
		return fmt.Errorf("error: %v\noutput:\n%v", err, string(o))
	}

	// Copy certificate to dcrwallet
	err = fileCopy(piClientCert, filepath.Join(dcrutil.AppDataDir("dcrwallet",
		false), walletClientsPem))
	if err != nil {
		return err
	}

	return nil
}

func (c *ctx) createWallet(version string) error {
	// create wallet
	c.log("creating wallet: %v\n", c.s.Net)

	dcrwalletExe := filepath.Join(c.s.Destination,
		"decred-"+c.s.Tuple+"-"+version, "dcrwallet")
	if runtime.GOOS == "windows" {
		dcrwalletExe += ".exe"
	}
	args := []string{"--create"}
	switch c.s.Net {
	case netTest:
		args = append(args, "--testnet")
	case netSim:
		args = append(args, "--simnet")
	}
	cmd := exec.Command(dcrwalletExe, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (c *ctx) main() error {
	running, err := c.running("dcrwallet")
	if err != nil {
		return err
	} else if running {
		return fmt.Errorf("dcrwallet is still running")
	}

	running, err = c.running("dcrd")
	if err != nil {
		return err
	} else if running {
		return fmt.Errorf("dcrd is still running")
	}

	running, err = c.running("dcrlnd")
	if err != nil {
		return err
	} else if running {
		return fmt.Errorf("dcrlnd is still running")
	}

	if c.s.Dcrdex {
		running, err = c.running("bitcoind")
		if err != nil {
			return err
		} else if running {
			return fmt.Errorf("bitcoind is still running")
		}
	}

	if !c.s.SkipDownload {
		c.s.Path, err = c.download()
		if err != nil {
			return err
		}
	}

	err = c.verify()
	if err != nil {
		return err
	}

	if c.s.DownloadOnly {
		// all done
		return nil
	}

	version, err := c.extract()
	if err != nil {
		return err
	}

	if c.s.Dcrdex {
		// bitcoin
		err = c.genericExtract(filepath.Base(c.s.BitcoinURI))
		if err != nil {
			return err
		}
		// dcrdex
		err = c.genericExtract(filepath.Base(c.s.DexURI))
		if err != nil {
			return err
		}
	}

	err = c.validate(version)
	if err != nil {
		return err
	}

	err = c.recordCurrent()
	if err != nil {
		return err
	}

	// prime defaults
	err = c.obtainUserName()
	if err != nil {
		return err
	}

	err = c.obtainPassword()
	if err != nil {
		return err
	}

	for _, v := range binaries {
		if v.Config == "" {
			continue
		}

		// check actual config file
		dir := dcrutil.AppDataDir(v.Name, false)
		conf := filepath.Join(dir, v.Config)
		if exist(conf) {
			c.log("skipping %s -- already installed\n", conf)
		}

		// Install config file
		config, err := c.createConfig(v, version)
		if err != nil {
			return err
		}
		c.log("creating directory: %v\n", dir)

		err = os.MkdirAll(dir, 0700)
		if err != nil {
			return err
		}

		c.log("installing configuration file: %s\n", conf)
		err = c.writeConfig(v, config)
		if err != nil {
			return err
		}
	}

	// Check client certs.
	walletCert := c.clientFileExists("dcrwallet", walletClientsPem)
	piCert := c.clientFileExists("politeiavoter", clientPem)
	piKey := c.clientFileExists("politeiavoter", clientKey)
	if walletCert && piCert && piKey {
		c.log("client certs exist, skipping cert generation.\n")
	} else if !walletCert && !piCert && !piKey {
		err = c.generateClientCerts(version)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("can't determine client certificate" +
			" state, must perform manual upgrade")
	}

	// Check client db
	if c.walletDBExists() {
		c.log("wallet.db exists, skipping wallet creation.\n")
	} else {
		err = c.createWallet(version)
		if err != nil {
			return err
		}
	}

	// install binaries in final location
	return c.copy(version)
}

func main() {
	var err error

	c := &ctx{}
	c.s, err = parseSettings()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	c.logFilename = filepath.Join(c.s.Destination, "dcrinstaller.log")

	c.logNoTime("=== dcrinstall run %v ===\n",
		time.Now().Format(time.RFC850))

	err = os.MkdirAll(c.s.Destination, 0700)
	if err != nil {
		c.log("%v\n", err)
	} else {
		err = c.main()
		if err != nil {
			c.log("%v\n", err)
		}
	}

	c.logNoTime("=== dcrinstall complete %v ===\n",
		time.Now().Format(time.RFC850))

	// exit with error set
	if err != nil {
		if !c.s.Verbose {
			// let user know something went wrong when not verbose
			fmt.Fprintf(os.Stderr, "%v\n", err)
		}
		os.Exit(1)
	}
}
