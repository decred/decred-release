# dcrinstall

dcrinstall is a tool to automate the install, upgrade, and setup
process for the decred software.

In install mode dcrinstall downloads the latest released binaries of
dcrd, dcrwallet, dcrctl, and dcrticketbuyer for your operating system
and platform, installs them, sets up the config files, and creates a
wallet for you.  In upgrade mode, dcrinstall replaces your binaries
with the latest copies but makes no changes to your configs.

## Upgrading an existing installation

The following steps are required to upgrade a system with decred that
was not installed by dcrinstall.  If you already have decred installed
you will need to follow these instructions the first time.

The dcrinstall tool expects the following directory layout.  In order
to upgrade you must copy your current configuration files into the
correct location and ensure everything still works.  You may also want
to copy you executables to the directory dcrinstall expects as well.

If dcrinstall detects all configuration files it'll operate in upgrade
mode.  Upgrade mode only overwrites the binaries in %HOMEPATH%\decred (or
~/decred on a UNIX type OS).

The dcrinstall tool records all actions in %HOMEPATH%\decred\dcrinstall.log
(or ~/decred/dcrinstall.log on a UNIX type OS).

### Windows

Configuration files:
```
%LOCALAPPDATA%\Dcrctl\dcrctl.conf
%LOCALAPPDATA%\Dcrd\dcrd.conf
%LOCALAPPDATA%\Dcrticketbuyer\ticketbuyer.conf
%LOCALAPPDATA%\Dcrwallet\dcrwallet.conf
```

Binaries directory:
```
%HOMEPATH%\decred\
```

### OSX

Configuration files:
```
~/Library/Application Support/Dcrctl/dcrctl.conf
~/Library/Application Support/Dcrd/dcrd.conf
~/Library/Application Support/Dcrticketbuyer/ticketbuyer.conf
~/Library/Application Support/Dcrwallet/dcrwallet.conf
```

Binaries directory:
```
~/decred
```

### UNIX

Configuration files:
```
~/.dcrctl/dcrctl.conf
~/.dcrd/dcrd.conf
~/.dcrticketbuyer/ticketbuyer.conf
~/.dcrwallet/dcrwallet.conf
```

Binaries directory:
```
~/decred
```

### Run the software

Now that you have the files where dcrinstall can find them you can
download and run dcrinstall

For Windows:

Open a cmd.exe window then:

```
cd %HOMEPATH%\Download
dcrinstall.exe
```

For OSX and UNIX you will also need to make the file executable before
runnning it:

```
cd Downloads/
chmod u+x dcrinstall
./dcrinstall
```

and you installation will be upgraded to the latest released version.

## Clean install

If you are doing a clean install (no existing decred configuration
files) you can just run dcrinstall and it will setup and configure all
the binaries:

For Windows open a cmd.exe window and:
```
cd %HOMEPATH%\Download
dcrinstall.exe
```

For OSX and UNIX:

```
cd Downloads/
./dcrinstall
```

You will be asked to provide a passphrase for you wallet and given the
opportunity to use and existing wallet seed if you have one.

## Log file

dcrinstall saves a log file with information on everything it did
which you may examine if you need more information.  On Windows the
file is located at:

```
%HOMEPATH%\decred\dcrinstaller.log
```

On OSX and UNIX the file is located at:

```
~/decred/dcrinstaller.log
```

## Testnet

It is possible to use dcrinstall to setup decred for testnet.  This
will ONLY work on a fresh install (as in you have no existing
configuration files on the computer).  Just add the option -net
Testnet to the command:


```
dcrinstall.exe -net testnet
```

or on OSX/UNIX:

```
./dcrinstall -net testnet
```

## Running decred programs

On Windows open cmd.exe

```
%HOMEPATH%\decred\dcrd.exe
```

One OSX and UNIX like systems:

```
cd decred/
./dcrd
```

Alternatively you can add the directory to your path.  For windows see
http://www.computerhope.com/issues/ch000549.htm  For OSX and UNIX
refer to the documentation for your shell.

## Build from source

dcrinstall can be used from the provided binaries but if you prefer to
build from source you can use these steps.  Note that dcrinstall uses
glide to manage dependencies so you cannot just use go get.  The
following instructions are for OSX and UNIX only.

```
go get -u github.com/Masterminds/glide
mkdir -p $GOPATH/src/github.com/decred
cd $GOPATH/src/github.com/decred
git clone https://github.com/decred/decred-release
cd decred-release
glide install
cd cmd/dcrinstall
go install
```

## Public Keys

The file
[cmd/dcrinstall/pubkey.go](https://github.com/decred/decred-release/blob/master/cmd/dcrinstall/pubkey.go)
contains the decred public key which is used to check the signed
manifest in the release.  You can compare the contents of this file to
what you get from a keyserver to confirm that dcrinstaller is using
the proper key.

## Notes

dcrinstall can only install decred releases v0.1.6 or later only
(although as described above it can be used to upgrade from an older
version).

dcrinstall has been tested on Windows 10, Windows 7, OSX 10.11, Bitrig current,
OpenBSD, Fedora, Ubuntu, and Raspbian.

## License

dcrinstall is licensed under the [copyfree](http://copyfree.org) ISC
License.

