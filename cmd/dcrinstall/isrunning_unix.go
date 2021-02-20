//go:build !windows
// +build !windows

// Copyright (c) 2016 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"runtime"
)

// isRunning returns true if the provided name appears in ps.
func isRunning(name string) (bool, error) {
	var args []string

	// Darwin allows both forms
	switch runtime.GOOS {
	case "linux":
		args = []string{"-Aaww"}
	default:
		// BSD*
		args = []string{"Aaww"}
	}
	cmd := exec.Command("ps", args...)
	o, err := cmd.CombinedOutput()
	if err != nil {
		return false, err
	}

	res := fmt.Sprintf(`(?:^|\W)%s(?:$|\W)`, name)
	re, err := regexp.Compile(res)
	if err != nil {
		return false, err
	}

	br := bytes.NewBuffer(o)
	for {
		line, err := br.ReadString('\n')
		if errors.Is(err, io.EOF) {
			break
		}

		if re.MatchString(line) {
			return true, nil
		}
	}

	return false, nil
}
