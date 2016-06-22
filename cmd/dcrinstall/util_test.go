// Copyright (c) 2016 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"path/filepath"
	"testing"

	"github.com/mitchellh/go-homedir"
)

func TestRunning(t *testing.T) {
	c := &ctx{s: &Settings{}}
	destination, err := homedir.Expand("~")
	if err != nil {
		t.Fatalf("%v", err)
	}
	c.s.Destination = filepath.Join(destination, "decred")

	r, err := c.isRunning("dcrd")
	if err != nil {
		t.Fatalf("%v", err)
	}

	t.Logf("dcrd running: %v\n", r)
}
