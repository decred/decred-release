// Copyright (c) 2016-2017 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/decred/dcrd/chaincfg"
	"github.com/decred/dcrd/sampleconfig"
	"github.com/decred/dcrutil"
	"github.com/decred/dcrwallet/loader"
	"github.com/decred/dcrwallet/prompt"
	"github.com/docker/docker/pkg/archive"
	"github.com/marcopeereboom/go-homedir"

	_ "github.com/decred/dcrwallet/walletdb/bdb"
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
	ExampleGenerate bool   // whether or not to generate the example config
	SupportsVersion bool   // whether or not it supports --version
}

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
			ExampleGenerate: true,
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
	if err != nil {
		return err
	}

	return nil
}

func (c *ctx) log(format string, args ...interface{}) error {
	t := time.Now().Format("15:04:05.000 ")
	return c.logNoTime(t+format, args...)
}

func (c *ctx) obtainUserName() error {
	u, err := homedir.User()
	if err != nil {
		return err
	}
	c.user = u
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

		// XXX quirk skip if .zip XXX
		if filepath.Ext(a[1]) == ".zip" {
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
	err = downloadToFile(manifestURI, manifest, c.s.SkipVerify)
	if err != nil {
		return "", err
	}

	// download manifest signature
	manifestAscURI := c.s.URI + "/" + c.s.Manifest + ".asc"
	if c.s.SkipVerify {
		c.log("SKIPPING downloading manifest "+
			"signatures: %v\n", manifestAscURI)
	} else {
		c.log("downloading manifest signatures: %v\n",
			manifestAscURI)

		manifestAsc := filepath.Join(td, filepath.Base(manifestAscURI))
		err = downloadToFile(manifestAscURI, manifestAsc,
			c.s.SkipVerify)
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
	err = downloadToFile(packageURI, pkg, c.s.SkipVerify)
	if err != nil {
		return "", err
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

	if c.s.SkipVerify {
		c.log("SKIPPING verifying manifest: %v\n",
			c.s.Manifest)
	} else {
		// verify manifest
		c.log("verifying manifest: %v ", c.s.Manifest)

		err = pgpVerify(manifest+".asc", manifest)
		if err != nil {
			c.logNoTime("FAIL\n")
			return fmt.Errorf("manifest PGP signature incorrect: %v", err)
		}

		c.logNoTime("OK\n")
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

	return nil
}

// copy verifies that all binaries can be executed.
func (c *ctx) copy(version string) error {
	for _, v := range binaries {
		// not in love with this, pull this out of tar instead
		filename := filepath.Join(c.s.Destination,
			"decred-"+c.s.Tuple+"-"+version,
			v.Name)

		// yep, this is ferrealz
		if runtime.GOOS == "windows" {
			filename += ".exe"
		}

		c.log("installing: %v -> %v\n", filename, c.s.Destination)

		// +"/" is needed to indicate the to tar that it is a directory
		err := archive.CopyResource(filename, c.s.Destination+"/", false)
		if err != nil {
			return err
		}
	}

	return nil
}

// validate verifies that all binaries can be executed.
func (c *ctx) validate(version string) error {
	for _, v := range binaries {
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

// exists ensures that either all or none of the binary config files exist.
func (c *ctx) exists() ([]string, error) {
	x := 0
	s := ""
	found := make([]string, 0, len(binaries))
	for _, v := range binaries {
		// check actual config file
		dir := dcrutil.AppDataDir(v.Name, false)
		conf := filepath.Join(dir, v.Config)

		if !exist(conf) {
			continue
		}

		found = append(found, filepath.Base(conf))
		s += filepath.Base(conf) + " "
		x++
	}

	if x != 0 {
		return found, fmt.Errorf("%valready exists", s)
	}

	return nil, nil
}

func (c *ctx) createConfigNormal(b binary, f *os.File) (string, error) {
	seen := false
	rv := ""
	usr := "; rpcuser="
	pwd := "; rpcpass="
	network := "; " + strings.ToLower(c.s.Net) + "="
	if b.Name == "dcrwallet" {
		usr = "; username="
		pwd = "; password="
	}

	br := bufio.NewReader(f)
	for {
		line, err := br.ReadString('\n')
		if err == io.EOF {
			break
		}

		if strings.HasPrefix(line, usr) {
			line = usr[2:] + c.user + "\n"
		}
		if strings.HasPrefix(line, pwd) {
			line = pwd[2:] + c.password + "\n"
		}
		if strings.HasPrefix(line, network) {
			line = network[2:] + "1\n"
			seen = true
		}

		rv += line
	}

	if c.s.Net != netMain {
		if seen == false {
			return "", fmt.Errorf("could not set net to %v\n",
				c.s.Net)
		}
	}

	return rv, nil
}

func (c *ctx) createConfig(b binary, version string) (string, error) {
	sample := filepath.Join(c.s.Destination,
		"decred-"+c.s.Tuple+"-"+version,
		b.Example)

	// write sample config if needed
	if b.ExampleGenerate {
		switch b.Name {
		case "dcrd":
			err := ioutil.WriteFile(sample, []byte(sampleconfig.FileContents), 0644)
			if err != nil {
				return "", fmt.Errorf("unable to write sample config to %v: %v",
					sample, err)
			}
		}
	}

	// read sample config
	f, err := os.Open(sample)
	if err != nil {
		return "", err
	}
	defer f.Close()

	c.log("parsing: %v\n", sample)

	return c.createConfigNormal(b, f)
}

func (c *ctx) writeConfig(b binary, cf string) error {
	dir := dcrutil.AppDataDir(b.Name, false)
	conf := filepath.Join(dir, b.Config)

	c.log("writing: %v\n", conf)

	err := ioutil.WriteFile(conf, []byte(cf), 0600)
	if err != nil {
		return err
	}

	return nil
}

func (c *ctx) walletDBExists() bool {
	dir := dcrutil.AppDataDir("dcrwallet", false)
	if exist(filepath.Join(dir, netMain, walletDB)) ||
		exist(filepath.Join(dir, netTest, walletDB)) ||
		exist(filepath.Join(dir, netSim, walletDB)) {
		return true
	}

	return false
}

func (c *ctx) createWallet() error {
	// create wallet
	c.log("creating wallet: %v\n", c.s.Net)

	r := bufio.NewReader(os.Stdin)
	privPass, pubPass, seed, err := prompt.Setup(r)
	if err != nil {
		return err
	}

	var chainParams *chaincfg.Params
	switch c.s.Net {
	case netMain:
		chainParams = &chaincfg.MainNetParams
	case netTest:
		chainParams = &chaincfg.TestNet2Params
	case netSim:
		chainParams = &chaincfg.SimNetParams
	default:
		return fmt.Errorf("invalid wallet type: %v",
			c.s.Net)
	}

	dbDir := filepath.Join(dcrutil.AppDataDir("dcrwallet",
		false), chainParams.Name)
	loader := loader.NewLoader(chainParams, dbDir,
		new(loader.StakeOptions), 0, false, 0)
	w, err := loader.CreateNewWallet(pubPass, privPass, seed)
	if err != nil {
		return err
	}
	_ = w

	err = loader.UnloadWallet()
	if err != nil {
		return err
	}
	return nil
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

	err = c.validate(version)
	if err != nil {
		return err
	}

	err = c.recordCurrent()
	if err != nil {
		return err
	}

	found, err := c.exists()
	if err != nil {
		c.log("--- Performing upgrade ---\n")
	} else if len(found) == 0 {
		c.log("--- Performing install ---\n")

		// prime defaults
		err = c.obtainUserName()
		if err != nil {
			return err
		}

		err = c.obtainPassword()
		if err != nil {
			return err
		}

		// lay down config files with parsed answers only if a Config
		// was defined
		for _, v := range binaries {
			if v.Config != "" {
				config, err := c.createConfig(v, version)
				if err != nil {
					return err
				}

				dir := dcrutil.AppDataDir(v.Name, false)
				c.log("creating directory: %v\n", dir)

				err = os.MkdirAll(dir, 0700)
				if err != nil {
					return err
				}

				err = c.writeConfig(v, config)
				if err != nil {
					return err
				}
			}
		}

		if c.walletDBExists() {
			c.log("wallet.db exists, skipping wallet creation.\n")
		} else {
			err = c.createWallet()
			if err != nil {
				return err
			}
		}
	}

	// install binaries in final location
	err = c.copy(version)
	if err != nil {
		return err
	}

	return nil
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
