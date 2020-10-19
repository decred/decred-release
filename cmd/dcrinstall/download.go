// Copyright (c) 2016 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func downloadToFile(url, filename string) error {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	if strings.HasPrefix(url, "file://") {
		localpath := url[len("file://"):]
		src, err := os.Open(localpath)
		if err != nil {
			return err
		}
		defer src.Close()
		_, err = io.Copy(f, src)
		return err
	}

	var client http.Client
	res, err := client.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP error: %v (%v)", res.Status, url)
	}

	_, err = io.Copy(f, res.Body)
	return err
}
