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
	"os/user"
	"path/filepath"
	"strings"

	"github.com/decred/dcrd/chaincfg"
	"github.com/decred/dcrutil"
	"github.com/decred/dcrwallet/prompt"
	"github.com/decred/dcrwallet/wallet"

	_ "github.com/decred/dcrwallet/walletdb/bdb"
)

// global context
type ctx struct {
	s *Settings
}

// question answer struct
type QA struct {
	Question string
	Default  string
	Answer   string
	Visible  bool // print Default?

	Validate      func(string) (string, error) // validate and sanitize value
	ObtainDefault func() (string, error)       // obtain default value
}

type binary struct {
	Name    string // binary filename
	Config  string // actual config file
	Example string // example config file
}

var (
	binaries = []binary{
		{
			Name:    "dcrctl",
			Config:  "dcrctl.conf",
			Example: "sample-dcrctl.conf",
		},
		{
			Name:    "dcrd",
			Config:  "dcrd.conf",
			Example: "sample-dcrd.conf",
		},
		{
			Name:    "dcrticketbuyer",
			Config:  ticketbuyerConf,
			Example: "ticketbuyer-example.conf",
		},
		{
			Name:    "dcrwallet",
			Config:  "dcrwallet.conf",
			Example: "sample-dcrwallet.conf",
		},
	}
)

const (
	ticketbuyerConf = "ticketbuyer.conf"
)

func obtainUserName() (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", err
	}

	return u.Username, nil
}

func obtainPassword() (string, error) {
	b := make([]byte, 24)
	_, err := io.ReadFull(rand.Reader, b[:])
	if err != nil {
		return "", err
	}

	// convert password to something readable
	password := base64.StdEncoding.EncodeToString(b)

	return password, nil
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

		if !strings.Contains(a[1], which) {
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
	if c.s.Verbose {
		fmt.Printf("temporary directory: %v\n", td)
		fmt.Printf("downloading manifest: %v\n", manifestURI)
	}

	manifest := filepath.Join(td, filepath.Base(manifestURI))
	err = downloadToFile(manifestURI, manifest, c.s.SkipVerify)
	if err != nil {
		return "", err
	}

	// download manifest signature
	manifestAscURI := c.s.URI + "/" + c.s.Manifest + ".asc"
	if c.s.SkipVerify {
		if c.s.Verbose {
			fmt.Printf("SKIPPING downloading manifest "+
				"signatures: %v\n", manifestAscURI)
		}
	} else {
		if c.s.Verbose {
			fmt.Printf("downloading manifest signatures: %v\n",
				manifestAscURI)
		}

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
	if c.s.Verbose {
		fmt.Printf("downloading package: %v\n", packageURI)
	}

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
		if c.s.Verbose {
			fmt.Printf("SKIPPING verifying manifest: %v\n",
				c.s.Manifest)
		}
	} else {
		// verify manifest
		if c.s.Verbose {
			fmt.Printf("verifying manifest: %v ", c.s.Manifest)
		}

		err = pgpVerify(manifest+".asc", manifest)
		if err != nil {
			if c.s.Verbose {
				fmt.Printf("fail\n")
			}
			return fmt.Errorf("manifest PGP signature incorrect: %v", err)
		}

		if c.s.Verbose {
			fmt.Printf("OK\n")
		}
	}

	// verify digest
	if c.s.Verbose {
		fmt.Printf("verifying package: %v ", filename)
	}

	pkg := filepath.Join(c.s.Path, filename)
	d, err := sha256File(pkg)
	if err != nil {
		return err
	}

	// verify package digest
	if hex.EncodeToString(d) != digest {
		if c.s.Verbose {
			fmt.Printf("failed\n")
		}
		fmt.Printf("%v %v\n", hex.EncodeToString(d), digest)

		return fmt.Errorf("corrupt digest %v", filename)
	}

	if c.s.Verbose {
		fmt.Printf("OK\n")
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

		if c.s.Verbose {
			fmt.Printf("checking: %v ", filename)
		}

		cmd := exec.Command(filename, "-h")
		err := cmd.Start()
		if err != nil {
			if c.s.Verbose {
				fmt.Printf("failed\n")
			}
			return err
		}

		if c.s.Verbose {
			fmt.Printf("OK\n")
		}

	}
	return nil
}

// exists ensures that either all or none of the binary config files exst.
func (c *ctx) exists() error {
	x := 0
	s := ""
	for _, v := range binaries {
		// check actual config file
		dir := dcrutil.AppDataDir(v.Name, false)
		conf := filepath.Join(dir, v.Config)

		if !exist(conf) {
			continue
		}

		s += filepath.Base(conf) + " "
		x++
	}

	if x != 0 {
		return fmt.Errorf("%valready exists", s)
	}

	return nil
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
			username, err := obtainUserName()
			if err != nil {
				return "", err
			}
			line = usr[2:] + username + "\n"
		}
		if strings.HasPrefix(line, pwd) {
			password, err := obtainPassword()
			if err != nil {
				return "", err
			}
			line = pwd[2:] + password + "\n"
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

func (c *ctx) createConfigTicketbuyer(b binary, f *os.File) (string, error) {
	username, err := obtainUserName()
	if err != nil {
		return "", err
	}

	password, err := obtainPassword()
	if err != nil {
		return "", err
	}

	seen := false
	rv := ""
	br := bufio.NewReader(f)
	for {
		line, err := br.ReadString('\n')
		if err == io.EOF {
			break
		}

		switch {
		case strings.HasPrefix(line, "dcrduser"):
			line = fmt.Sprintf("dcrduser=%v\n", username)

		case strings.HasPrefix(line, "dcrwuser"):
			line = fmt.Sprintf("dcrwuser=%v\n", username)

		case strings.HasPrefix(line, "dcrdpass"):
			line = fmt.Sprintf("dcrdpass=%v\n", password)

		case strings.HasPrefix(line, "dcrwpass"):
			line = fmt.Sprintf("dcrwpass=%v\n", password)

		case strings.HasPrefix(line, "httpsvrport"):
			// use default from config file

		case strings.HasPrefix(line, "httpuipath"):
			dir := filepath.Join(c.s.Destination, "webui")
			line = fmt.Sprintf("httpuipath=%v\n", dir)

		case strings.HasPrefix(line, "simnet"):
			a := "0"
			if c.s.Net == netSim {
				a = "1"
			}
			line = fmt.Sprintf("simnet=%v\n", a)
			seen = true

		case strings.HasPrefix(line, "testnet"):
			a := "0"
			if c.s.Net == netTest {
				a = "1"
			}
			line = fmt.Sprintf("testnet=%v\n", a)
			seen = true

		case strings.HasPrefix(line, "\n"):
			// do nothing

		case !strings.HasPrefix(line, "#"):
			// comment out
			line = "#" + line
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
	// read sample config
	sample := filepath.Join(c.s.Destination,
		"decred-"+c.s.Tuple+"-"+version,
		b.Example)

	f, err := os.Open(sample)
	if err != nil {
		return "", err
	}
	defer f.Close()

	if c.s.Verbose {
		fmt.Printf("parsing: %v\n", sample)
	}

	if b.Config == ticketbuyerConf {
		return c.createConfigTicketbuyer(b, f)
	}

	return c.createConfigNormal(b, f)
}

func (c *ctx) writeConfig(b binary, cf string) error {
	dir := dcrutil.AppDataDir(b.Name, false)
	conf := filepath.Join(dir, b.Config)

	if c.s.Verbose {
		fmt.Printf("writing: %v\n", conf)
	}

	err := ioutil.WriteFile(conf, []byte(cf), 0600)
	if err != nil {
		return err
	}

	return nil
}

func _main() error {
	var err error

	c := &ctx{}
	c.s, err = parseSettings()
	if err != nil {
		return err
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

	version, err := c.extract()
	if err != nil {
		return err
	}

	err = c.validate(version)
	if err != nil {
		return err
	}

	err = c.exists()
	if err != nil {
		return err
	}

	// lay down config files with parsed answers
	for _, v := range binaries {
		config, err := c.createConfig(v, version)
		if err != nil {
			return err
		}

		dir := dcrutil.AppDataDir(v.Name, false)
		if c.s.Verbose {
			fmt.Printf("creating directory: %v\n", dir)
		}

		err = os.MkdirAll(dir, 0700)
		if err != nil {
			return err
		}

		err = c.writeConfig(v, config)
		if err != nil {
			return err
		}
	}

	// create wallet
	if c.s.Verbose {
		fmt.Printf("creating wallet: %v\n", c.s.Net)
	}

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
		chainParams = &chaincfg.TestNetParams
	case netSim:
		chainParams = &chaincfg.SimNetParams
	default:
		return fmt.Errorf("invalid wallet type: %v", c.s.Net)
	}

	dbDir := filepath.Join(dcrutil.AppDataDir("dcrwallet", false),
		chainParams.Name)
	loader := wallet.NewLoader(chainParams, dbDir, new(wallet.StakeOptions),
		false, false, 0, false)
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

func main() {
	err := _main()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
