#!/bin/bash

# Simple bash script to build and package Decred miner for all the platforms
# we support.
#
# Currently, ccminer cannot be crosscompiled so it must be manually
# built and zipped up on windows.  Also ccminer doesn't build on linux
# on my build machine (jcv) due to some older libs so it must also be done
# manually.  Hopefully both of these can be fixed in the near future.
#
# Copyright (c) 2016 The Decred developers
# Use of this source code is governed by the ISC
# license.

if [[ $1x = x ]]; then
    TAG=""
else
    TAG=-$1
fi

DATE=`date +%Y%m%d`
MAINDIR=miners$TAG-$DATE-$VERSION
mkdir -p $MAINDIR

cd ../cgminer/
git clean
./build_linux.sh
mkdir cgminer-decred-linux-x86_64$TAG-$DATE
cp cgminer cgminer-decred-linux-x86_64$TAG-$DATE/
cp blake256.cl cgminer-decred-linux-x86_64$TAG-$DATE/
tar -cvzf cgminer-decred-linux-x86_64$TAG-$DATE.tar.gz cgminer-decred-linux-x86_64$TAG-$DATE/
cp cgminer-decred-linux-x86_64$TAG-$DATE.tar.gz ../decred-release/$MAINDIR
./build_windows.sh
mkdir cgminer-decred-win64$TAG-$DATE
cp cgminer..exe cgminer-decred-win64$TAG-$DATE/
cp blake256.cl cgminer-decred-win64$TAG-$DATE/
zip -r cgminer-decred-win64$TAG-$DATE.zip cgminer-decred-win64$TAG-$DATE/
cp cgminer-decred-win64$TAG-$DATE.zip ../decred-release/$MAINDIR
cd ..
cp ccminer-decred-linux-x86_64$TAG-$DATE.tar.gz decred-release/$MAINDIR
cp ccminer-decred-win64$TAG-$DATE.zip decred-release/$MAINDIR
sha256sum * > manifest-miners-$DATE.txt
cd ..
tar -cvzf $MAINDIR.tar.gz $MAINDIR
