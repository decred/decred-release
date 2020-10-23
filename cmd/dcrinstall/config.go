// Copyright (c) 2016-2020 The Decred developers
// Use of this source code is governed by an ISC license that can be found in
// the LICENSE file.

package main

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"strings"
)

// override describes an entry (name) in a config file that has to be overriden
// by "content".
type override struct {
	name    string
	content string
}

func createConfigNormal(br *bufio.Reader, overrides []override) (string, error) {
	rv := ""

	for {
		line, err := br.ReadString('\n')
		if errors.Is(err, io.EOF) {
			break
		}

		for k := range overrides {
			if !strings.HasPrefix(line, overrides[k].name) {
				continue
			}
			line = strings.TrimLeft(overrides[k].name, ";# ") +
				overrides[k].content + "\n"
		}

		rv += line
	}

	return rv, nil
}

// createConfigFromFile reads a sample config file and modifies it based on the
// provided override array.
func createConfigFromFile(filename string, overrides []override) (string, error) {
	// read sample config
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	return createConfigNormal(bufio.NewReader(f), overrides)
}

// createConfigFromMemory reads a sample config file and modifies it based on
// the provided override array.
func createConfigFromMemory(conf string, overrides []override) (string, error) {
	return createConfigNormal(bufio.NewReader(bytes.NewReader([]byte(conf))),
		overrides)
}
