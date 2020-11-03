// Copyright (c) 2016-2020 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"regexp"
	"strconv"
)

var relRE = regexp.MustCompile(`(v|release-v)?(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(-(0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(\.(0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*)?(\+[0-9a-zA-Z-]+(\.[0-9a-zA-Z-]+)*)?`)

type semVerInfo struct {
	Major      uint32
	Minor      uint32
	Patch      uint32
	PreRelease string
	Build      string
}

// extractSemVer peels a semver out of a string.
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

// String satisfies the Stringer interface for semVerInfo.
func (s semVerInfo) String() string {
	var pre string
	if s.PreRelease != "" {
		pre = "-" + s.PreRelease
	}
	return fmt.Sprintf("v%v.%v.%v%v", s.Major, s.Minor, s.Patch, pre)
}
