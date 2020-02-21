// Copyright (c) 2016-2020 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"os/user"
	"path/filepath"
	"testing"
)

func TestRunning(t *testing.T) {
	c := &ctx{s: &Settings{}}
	u, err := user.Current()
	if err != nil {
		t.Fatalf("%v", err)
	}
	c.s.Destination = filepath.Join(u.HomeDir, "decred")

	r, err := c.isRunning("dcrd")
	if err != nil {
		t.Fatalf("%v", err)
	}

	t.Logf("dcrd running: %v\n", r)
}
