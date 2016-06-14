package main

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/pkg/archive"

	"golang.org/x/crypto/openpgp"
)

func answer(def string) string {
	r := bufio.NewReader(os.Stdin)
	a, _ := r.ReadString('\n')
	a = strings.TrimSpace(a)
	if len(a) == 0 {
		return def
	}
	return a
}

func yes() bool {
	r := bufio.NewReader(os.Stdin)
	a, _ := r.ReadString('\n')
	a = strings.ToUpper(strings.TrimSpace(a))
	if len(a) == 0 {
		return false
	}
	if a[0] == 'Y' {
		return true
	}
	return false
}

func exist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	}

	return true
}

// pgpVerify verifies the signature of manifest file using global pubkey.
func pgpVerify(signature, manifest string) error {
	// open manifest signature
	sf, err := os.Open(signature)
	if err != nil {
		return err
	}
	defer sf.Close()

	// open manifest
	mf, err := os.Open(manifest)
	if err != nil {
		return err
	}
	defer mf.Close()

	// create keyring
	br := bytes.NewBufferString(pubkey)
	keyring, err := openpgp.ReadArmoredKeyRing(br)
	if err != nil {
		return err
	}

	// verify signature
	_, err = openpgp.CheckArmoredDetachedSignature(keyring, mf, sf)
	if err != nil {
		return err
	}

	return nil
}

//sha256File returns the sha256 digest of the provided file.
func sha256File(filename string) ([]byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	hasher := sha256.New()
	_, err = io.Copy(hasher, f)
	if err != nil {
		return nil, fmt.Errorf("sha256: %v", err)
	}

	return hasher.Sum(nil), nil
}

// extract downloaded package.
func (c *ctx) extract() error {
	manifest := filepath.Join(c.s.Path, c.s.Manifest)
	_, filename, err := findOS(c.s.Tuple, manifest)
	if err != nil {
		return err
	}

	if c.s.Verbose {
		fmt.Printf("extracting: %v -> %v\n", filename, c.s.Destination)
	}

	err = archive.UntarPath(filepath.Join(c.s.Path, filename),
		c.s.Destination)
	if err != nil {
		return err
	}

	return nil
}
