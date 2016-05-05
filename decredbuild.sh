#!/bin/bash

# Simple bash script to build basic Decred tools for all the platforms
# we support with the golang cross-compiler.
#
# Copyright (c) 2016 The Decred developers
# Use of this source code is governed by the ISC
# license.

if [[ $1x = x ]]; then
    TAG=""
else
    TAG=-$1
fi

VERSION="01"
DATE=`date +%Y%m%d`
MAINDIR=decred$TAG-$DATE-$VERSION
mkdir -p $MAINDIR
cd $MAINDIR

SYS="windows-386 windows-amd64 openbsd-386 openbsd-amd64 linux-386 linux-amd64 linux-arm linux-arm64 darwin-386 darwin-amd64 dragonfly-amd64 freebsd-386 freebsd-amd64 freebsd-arm netbsd-386 netbsd-amd64"

# Use the first element of $GOPATH in the case where GOPATH is a list
# (something that is totally allowed).
GPATH=$(echo $GOPATH | cut -f1 -d:)

for i in $SYS; do
    OS=$(echo $i | cut -f1 -d-)
    ARCH=$(echo $i | cut -f2 -d-)
    mkdir decred-$i
    cd decred-$i
    echo "Building:" $OS $ARCH
    env GOOS=$OS GOARCH=$ARCH go build github.com/decred/dcrd
    env GOOS=$OS GOARCH=$ARCH go build github.com/decred/dcrd/cmd/dcrctl
    env GOOS=$OS GOARCH=$ARCH go build github.com/decred/dcrwallet
    cp $GPATH/src/github.com/decred/dcrd/sample-dcrd.conf .
    cp $GPATH/src/github.com/decred/dcrd/cmd/dcrctl/sample-dcrctl.conf .
    cp $GPATH/src/github.com/decred/dcrwallet/sample-dcrwallet.conf .
    cd ..
    if [[ $OS = "windows" ]]; then
	zip -r decred$TAG-$i-$DATE-$VERSION.zip decred-$i
    else
	tar -cvzf decred$TAG-$i-$DATE-$VERSION.tar.gz decred-$i
    fi
    rm -r decred-$i
done

sha256sum * > manifest-$DATE-$VERSION.txt
cd ..
tar -cvzf $MAINDIR.tar.gz $MAINDIR
