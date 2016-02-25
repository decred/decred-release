#!/bin/bash

# Simple bash script to build basic Decred tools for all the platforms
# we support with the golang cross-compiler.
#
# Copyright (c) 2016 The Decred developers
# Use of this source code is governed by the ISC
# license.

VERSION="01"
DATE=`date +%Y%m%d`
MAINDIR=decred-$DATE-$VERSION
mkdir -p $MAINDIR
cd $MAINDIR

SYS="windows-386 windows-amd64 openbsd-386 openbsd-amd64 linux-386 linux-amd64 linux-arm darwin-386 darwin-amd64 dragonfly-amd64 freebsd-386 freebsd-amd64 freebsd-arm netbsd-386 netbsd-amd64"

for i in $SYS; do
    OS=$(echo $i | cut -f1 -d-)
    ARCH=$(echo $i | cut -f2 -d-)
    mkdir $i
    cd $i
    echo "Building:" $OS $ARCH
    env GOOS=$OS GOARCH=$ARCH go build github.com/decred/dcrd
    env GOOS=$OS GOARCH=$ARCH go build github.com/decred/dcrd/cmd/dcrctl
    env GOOS=$OS GOARCH=$ARCH go build github.com/decred/dcrwallet
    cd ..
    if [[ $OS = "windows" ]]; then
	zip -r $i-$DATE-$VERSION.zip $i
    else
	tar -cvzf $i-$DATE-$VERSION.tar.gz $i
    fi
    rm -r $i
done

sha256sum * > manifest-$DATE-$VERSION.txt
cd ..
tar -cvzf $MAINDIR.tar.gz $MAINDIR
