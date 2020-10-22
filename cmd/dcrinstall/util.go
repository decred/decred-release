// Copyright (c) 2016-2020 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"errors"
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

func exist(path string) bool {
	_, err := os.Stat(path)

	return err == nil
}

// pgpVerify verifies the signature of manifest file using global pubkey.
func pgpVerify(signature, manifest, key string) error {
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
	br := bytes.NewBufferString(key)
	keyring, err := openpgp.ReadArmoredKeyRing(br)
	if err != nil {
		return err
	}

	// verify signature
	_, err = openpgp.CheckArmoredDetachedSignature(keyring, mf, sf)
	return err
}

// sha256File returns the sha256 digest of the provided file.
func sha256File(filename string) ([]byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	hasher := sha256.New()
	_, err = io.Copy(hasher, f)
	if err != nil {
		return nil, fmt.Errorf("sha256: %w", err)
	}

	return hasher.Sum(nil), nil
}

func fileCopy(src, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE,
		sourceFileStat.Mode())
	if err != nil {
		return err
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)
	return err
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

// btcExtract extracts the downloaded bitcoin package.
func (c *ctx) genericExtract(filename string) error {
	c.log("extracting: %v -> %v\n", filename, c.s.Destination)
	var err error
	if filepath.Ext(filename) == ".zip" {
		err = c.unzip(filename)
	} else {
		err = c.gunzip(filename)
	}
	return err
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
			if errors.Is(err, io.EOF) {
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
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR,
				os.FileMode(hdr.Mode))
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
	dst := c.s.Destination
	_, err := c._unzip(src, dst)
	return err
}

// _unzip borrowed from https://golangcode.com/unzip-files-in-go/
func (c *ctx) _unzip(src string, dest string) ([]string, error) {

	var filenames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer r.Close()

	for _, f := range r.File {
		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+
			string(os.PathSeparator)) {
			return filenames, fmt.Errorf("%s: illegal file path",
				fpath)
		}
		c.log("contents of %v\n", f.Name)

		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Make File
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return filenames, err
		}

		outFile, err := os.OpenFile(fpath,
			os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return filenames, err
		}

		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration
		// of loop
		outFile.Close()
		rc.Close()

		if err != nil {
			return filenames, err
		}
	}
	return filenames, nil
}
