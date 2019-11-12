// Copyright (c) 2016 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"archive/tar"
	"archive/zip"
	"bufio"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/crypto/openpgp"
)

var relRE = regexp.MustCompile(`(v|release-v)?(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(-(0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(\.(0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*)?(\+[0-9a-zA-Z-]+(\.[0-9a-zA-Z-]+)*)?`)

type semVerInfo struct {
	Major      uint32
	Minor      uint32
	Patch      uint32
	PreRelease string
	Build      string
}

func extractSemVer(s string) (*semVerInfo, error) {
	matches := relRE.FindStringSubmatch(s)
	if len(matches) == 0 {
		return nil, fmt.Errorf("version string %q does not follow semantic "+
			"versioning requirements", s)
	}

	major, err := strconv.ParseInt(matches[2], 10, 32)
	if err != nil {
		return nil, err
	}
	minor, err := strconv.ParseInt(matches[3], 10, 32)
	if err != nil {
		return nil, err
	}
	patch, err := strconv.ParseInt(matches[4], 10, 32)
	if err != nil {
		return nil, err
	}

	return &semVerInfo{
		Major:      uint32(major),
		Minor:      uint32(minor),
		Patch:      uint32(patch),
		PreRelease: matches[6],
		Build:      matches[9],
	}, nil
}

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
	if filepath.Ext(filename) == ".zip" {
		err = c.unzip(filename)
	} else {
		err = c.gunzip(filename)
	}
	if err != nil {
		return "", err
	}

	// fish out version
	info, err := extractSemVer(filename)
	if err != nil {
		return "", err
	}

	version := fmt.Sprintf("v%v.%v.%v", info.Major, info.Minor, info.Patch)
	if info.PreRelease != "" {
		version += "-" + info.PreRelease
	}

	return version, nil
}

func (c *ctx) gunzip(filename string) error {
	src := filepath.Join(c.s.Path, filename)

	a, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}
	ab := bytes.NewReader(a)
	gz, err := gzip.NewReader(ab)
	if err != nil {
		return err
	}
	defer gz.Close()
	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				break // end of archive
			}
			return err
		}
		if hdr == nil {
			continue
		}
		target := filepath.Join(c.s.Destination, hdr.Name)
		c.log("contents of %v\n", target)
		switch hdr.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(hdr.Mode))
			if err != nil {
				return err
			}

			// copy to file
			if _, err := io.Copy(f, tr); err != nil {
				f.Close()
				return err
			}

			if err := f.Close(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *ctx) unzip(filename string) error {
	src := filepath.Join(c.s.Path, filename)

	uz, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer uz.Close()

	for _, file := range uz.File {
		c.log("contents of %v\n", file.Name)
		rc, err := file.Open()
		if err != nil {
			return err
		}

		target := filepath.Join(c.s.Destination, file.Name)
		if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
			rc.Close()
			return err
		}

		f, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY, os.FileMode(0755))
		if err != nil {
			rc.Close()
			return err
		}
		if _, err = io.Copy(f, rc); err != nil {
			f.Close()
			rc.Close()
			return err
		}
		if err = f.Close(); err != nil {
			return err
		}
		if err = rc.Close(); err != nil {
			return err
		}
	}
	return nil
}
