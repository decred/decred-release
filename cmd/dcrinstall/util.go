// Copyright (c) 2016 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/moby/moby/pkg/archive"

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

	return err == nil
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
	return err
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
func (c *ctx) extract() (string, error) {
	manifest := filepath.Join(c.s.Path, c.s.Manifest)
	_, filename, err := findOS(c.s.Tuple, manifest)
	if err != nil {
		return "", err
	}

	c.log("extracting: %v -> %v\n", filename, c.s.Destination)

	src := filepath.Join(c.s.Path, filename)
	a, err := os.Open(src)
	if err != nil {
		return "", err
	}
	defer a.Close()

	options := &archive.TarOptions{NoLchown: true}
	err = archive.Untar(a, c.s.Destination, options)
	if err != nil {
		return "", err
	}

	// fish out version
	re := regexp.MustCompile("v[0-9]+.[0-9]+.[0-9]+")
	return re.FindString(filename), nil
}
