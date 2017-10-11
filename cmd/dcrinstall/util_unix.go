// +build !windows

// Copyright (c) 2016 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"io"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
)

func (c *ctx) isRunning(name string) (bool, error) {
	var args []string

	switch runtime.GOOS {
	case "linux":
		args = []string{"aeww"}
	default:
		// BSD*
		args = []string{"Aaeww"}
	}
	cmd := exec.Command("ps", args...)
	o, err := cmd.CombinedOutput()
	if err != nil {
		return false, err
	}

	re := regexp.MustCompile("_=[[:print:]]*" + name)

	br := bytes.NewBuffer(o)
	for {
		line, err := br.ReadString('\n')
		if err == io.EOF {
			break
		}
		line = strings.TrimSpace(line)

		s := re.FindString(line)
		if s == "" {
			continue
		}
		
		if len(strings.Split(s, "=")) != 2 {
			continue
		}

		return true, nil
	}

	return false, nil
}
