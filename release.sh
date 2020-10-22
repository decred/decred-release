#!/bin/sh

# Copyright (c) 2019 The Decred developers
# Use of this source code is governed by the ISC
# license.

set -e

usage() {
    echo "usage: release.sh version" 1>&2
    exit 2
}

[ $# -eq 1 ] || usage
TAG=$1

set -x

GOVERSION=$(go version)
LDFLAGS="-buildid= -X main.appBuild=release"
SYS="darwin-amd64 freebsd-amd64 linux-386 linux-amd64 linux-arm linux-arm64 openbsd-amd64 windows-386 windows-amd64"

PWD=$(pwd)
PACKAGE=dcrinstall
MAINDIR=$PWD/release/$PACKAGE-$TAG
MANIFEST=dcrinstall-${TAG}-manifest.txt

[ -d ${MAINDIR} ] && rm -rf ${MAINDIR}
mkdir -p $MAINDIR

for i in $SYS; do
    OS=$(echo $i | cut -f1 -d-)
    ARCH=$(echo $i | cut -f2 -d-)
    OUT="$MAINDIR/dcrinstall-$i-$TAG"
    if [[ $OS = "windows" ]]; then
	    OUT="$OUT.exe"
    fi
    ARM=
    if [[ $ARCH = "arm" ]]; then
	ARM=7
    fi
    echo "Building:" $OS $ARCH
    env CGO_ENABLED=0 GOOS=$OS GOARCH=$ARCH GOARM=$ARM GOFLAGS= \
        go build -trimpath -tags 'safe,netgo' -o $OUT -ldflags="${LDFLAGS}" ./cmd/dcrinstall
done

(
    cd ${MAINDIR}
    openssl sha256 -r * > ${MANIFEST}
    set +x
    MANIFESTHASH=$(openssl sha256 ${MANIFEST})
    echo "${MANIFESTHASH} (${GOVERSION})"
)
