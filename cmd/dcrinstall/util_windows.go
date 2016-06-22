// +build windows

// Copyright (c) 2016 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
)

func (c *ctx) isRunning(name string) (bool, error) {
	filename := filepath.Join(c.s.Destination, name)
	f, err := os.OpenFile(filename+".exe", os.O_RDWR, 0600)
	if err != nil {
		if os.IsNotExist(err) {
			// file doesn't exist so it can't be running
			return false, nil
		}

		// try to see if file was locked
		x, ok := err.(*os.PathError)
		if !ok {
			return false, fmt.Errorf("invalid type")
		}
		e, ok := x.Err.(syscall.Errno)
		if !ok {
			return false, fmt.Errorf("invalid error type")
		}
		if e == 0x20 {
			return true, nil
		}

		return false, err
	}
	defer f.Close()

	return false, nil
}
