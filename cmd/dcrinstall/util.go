// Copyright (c) 2016-2020 The Decred developers
// Use of this source code is governed by an ISC license that can be found in
// the LICENSE file.

package main

import (
	"archive/tar"
	"archive/zip"
	"bufio"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/decred/dcrd/dcrutil"
	humanize "github.com/dustin/go-humanize"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/clearsign"
)

// cleanAndExpandPath expands environment variables and leading ~ in the
// passed path, cleans the result, and returns it.
func cleanAndExpandPath(path string) string {
	// Nothing to do when no path is given.
	if path == "" {
		return path
	}

	// NOTE: The os.ExpandEnv doesn't work with Windows cmd.exe-style
	// %VARIABLE%, but the variables can still be expanded via POSIX-style
	// $VARIABLE.
	path = os.ExpandEnv(path)

	if !strings.HasPrefix(path, "~") {
		return filepath.Clean(path)
	}

	// Expand initial ~ to the current user's home directory, or ~otheruser
	// to otheruser's home directory.  On Windows, both forward and backward
	// slashes can be used.
	path = path[1:]

	var pathSeparators string
	if runtime.GOOS == "windows" {
		pathSeparators = string(os.PathSeparator) + "/"
	} else {
		pathSeparators = string(os.PathSeparator)
	}

	userName := ""
	if i := strings.IndexAny(path, pathSeparators); i != -1 {
		userName = path[:i]
		path = path[i:]
	}

	homeDir := ""
	var u *user.User
	var err error
	if userName == "" {
		u, err = user.Current()
	} else {
		u, err = user.Lookup(userName)
	}
	if err == nil {
		homeDir = u.HomeDir
	}
	// Fallback to CWD if user lookup fails or user has no home directory.
	if homeDir == "" {
		homeDir = "."
	}

	return filepath.Join(homeDir, path)
}

// WriteCounter keeps track of the download progress.
type WriteCounter struct {
	Total uint64
}

// Write satisfies the Writer interface for WriteCounter.
func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	wc.PrintProgress()
	return n, nil
}

// PrintProgress prints the progress of a file write
func (wc WriteCounter) PrintProgress() {
	if quiet {
		return
	}

	// Clear the line by using a character return to go back to the start
	// and remove the remaining characters by filling it with spaces
	fmt.Printf("\r%s", strings.Repeat(" ", 50))

	// Return again and print current status of download We use the
	// humanize package to print the bytes in a meaningful way (e.g. 10 MB)
	fmt.Printf("\rDownloading... %s complete", humanize.Bytes(wc.Total))
}

// DownloadFile downloads the provided URL to the filepath. If the quiet flag
// is not set it prints download progress.
func DownloadFile(url string, filepath string) error {
	log.Printf("Download file: %v -> %v", url, filepath)

	// Create the file with .tmp extension, so that we won't overwrite a
	// file until it's downloaded fully
	out, err := os.Create(filepath + ".tmp")
	if err != nil {
		return err
	}
	defer out.Close()

	// Deal with local files
	if strings.HasPrefix(url, "file://") {
		localpath := url[len("file://"):]
		src, err := os.Open(localpath)
		if err != nil {
			return err
		}
		defer src.Close()

		_, err = io.Copy(out, src)
		if err != nil {
			return err
		}
	} else {
		// Get file over HTTP
		c := http.Client{
			Transport: &http.Transport{Proxy: http.ProxyFromEnvironment},
		}
		resp, err := c.Get(url)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("%v %v", resp.StatusCode,
				http.StatusText(resp.StatusCode))
		}

		// Create our bytes counter and pass it to be used alongside
		// our writer
		counter := &WriteCounter{}
		_, err = io.Copy(out, io.TeeReader(resp.Body, counter))
		if err != nil {
			return err
		}

		// The progress use the same line so print a new line once it's
		// finished downloading
		if !quiet {
			fmt.Println()
		}
	}

	// Close file because windows
	out.Close()

	// Rename the tmp file back to the original file
	err = os.Rename(filepath+".tmp", filepath)
	if err != nil {
		return err
	}

	return nil
}

// pgpVerify verifies the signature with the provided key.
func pgpVerify(signature, manifest, key string) error {
	log.Printf("PGP verify: %v", manifest)

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

func pgpVerifyAttached(file, key string) error {
	log.Printf("PGP attached verify: %v", file)

	// open manifest signature
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	b, _ := clearsign.Decode(data)
	if b == nil {
		return fmt.Errorf("PGP attached signature failed")
	}

	// create keyring
	br := bytes.NewBufferString(key)
	keyring, err := openpgp.ReadArmoredKeyRing(br)
	if err != nil {
		return err
	}

	// verify signature
	_, err = openpgp.CheckDetachedSignature(keyring, bytes.NewReader(b.Bytes),
		b.ArmoredSignature.Body)
	if err != nil {
		return err
	}

	return nil
}

// findOS iterates over the entire manifest and plucks out the digest and
// filename of the provided os-arch tuple.  The tuple must be unique.
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
		if errors.Is(err, io.EOF) {
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
		return nil, fmt.Errorf("sha256File: %w", err)
	}

	return hasher.Sum(nil), nil
}

// sha256Verify verifies that the provided file matches the provided digest.
func sha256Verify(filename, digest string) error {
	log.Printf("Verify SHA256: %v", filename)

	d, err := sha256File(filename)
	if err != nil {
		return err
	}
	if hex.EncodeToString(d) != digest {
		return fmt.Errorf("corrupt digest")
	}
	return nil
}

// unzip unzips src to dst.
// unzip borrowed from https://golangcode.com/unzip-files-in-go/
func unzip(src string, dest string) ([]string, error) {
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
		log.Printf("Extracting: %v", f.Name)

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

// gunzip untars filename to destination.
func gunzip(filename, destination string) error {
	a, err := ioutil.ReadFile(filename)
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
		log.Printf("Extracting: %v", hdr.Name)
		target := filepath.Join(destination, hdr.Name)
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

// extract extracts the provided archive to the provided destination. It
// autodetects if it is a zip or a tar archive.
func extract(filename, dst string) error {
	log.Printf("Extracting: %v -> %v\n", filename, dst)
	var err error
	archive := filepath.Ext(filename)
	switch archive {
	case ".zip":
		_, err = unzip(filename, dst)
	case ".gz":
		err = gunzip(filename, dst)
	default:
		err = fmt.Errorf("Unknown archive type: %v", archive)
	}
	return err
}

// exists return true if the provided path exists.
func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// appFileExists return true if the provided default application filename
// exists.
func appFileExists(app, filename string) bool {
	dir := dcrutil.AppDataDir(app, false)
	return exists(filepath.Join(dir, filename))
}

// getDownloadURI returns the path portion of a URI.
func getDownloadURI(uri string) (string, error) {
	for i := len(uri) - 1; i > 0; i-- {
		if uri[i] == '/' {
			return uri[:i+1], nil
		}
	}

	return "", fmt.Errorf("invalid URI: %v", uri)
}

// runtimeTuple returns the OS and ARCH the program is running as.
func runtimeTuple() string {
	return runtime.GOOS + "-" + runtime.GOARCH
}
func printConfigError(installedConfigs, notInstalledConfigs []string) string {
	rv := "Installed configuration files:\n"
	for _, v := range installedConfigs {
		rv += "\t" + v + "\n"
	}
	rv += "\nNOT installed configuration files:\n"
	for _, v := range notInstalledConfigs {
		rv += "\t" + v + "\n"
	}
	return rv
}
