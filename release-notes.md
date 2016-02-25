# v0.0.4

## 2016-02-24

This release contains updated binary files (dcrd, dcrctl, dcrwallet)
for various platforms.

See manifest-20160224-01.txt for sha256sums of the packages and
manifest-20160224-01.txt.asc to confirm those shas.

See https://wiki.decred.org/Verifying_Binaries for more info on
verifying the files.

This release includes a number of fixes for both wallet and daemon as
well as several new rpc calls.

This includes (but is not limited to):
* Added getcoinsupply, get/setticketvotebits, existslivetickets, getstakeinfo
* First checkpoint added
* Several fee related issues
* Disable unsafe RPC calls on mainnet
* Corrected fee estimation for general transactions
* Allow wallet to accept hex or words as seed
* Other bug fixes and cleanups

# v0.0.3

## 2016-02-09

This wallet only release resolves an upstream wallet bug (see
decred/dcrutil 8aae5a2dacf45b7f5ee9b59c393118bc48647861).

Platform specific files are attached.

See manifest-20160209-01.txt for sha256sums of the packages and
manifest-20160209-01.txt.asc to confirm those shas.

See https://wiki.decred.org/Verifying_Binaries for more info on
verifying the files.

# v0.0.2

## 2016-02-08

This release is an unencrypted version of the current mainnet enabled
code.

The packages below contain platform specific copies of drcd,
dcrwallet, and dcrctl.

See manifest-20160208-01.txt for sha256sums of the packages and
manifest-20160208-01.txt.asc to confirm those shas.

See https://wiki.decred.org/Verifying_Binaries for more info on
verifying the files.

# v0.0.1

## 2016-02-07

This is the initial mainnet binaries for Decred.

The packages below contain platform specific copies of drcd,
dcrwallet, and dcrctl.

See manifest-20160207-01.txt for sha256sums of the packages and
manifest-20160207-01.txt.asc to confirm those shas.

See https://wiki.decred.org/Verifying_Binaries for more info on
verifying the files.

The platform specific archives have been encrypted with 7zip. The key
will be made available when mainnet is launched.

To unencrypt on the command line you can do:

```bash
7za e FILENAME
```

then provide the password when asked.

Mainnet binary decryption key (password): yqJgFJUmQODUOWP2jJez5gt1

# v0.0

## 2016-01-27

This is the testnet pre-release of Decred.

The packages below contain platform specific copies of drcd,
dcrwallet, and dcrctl.

See manifest-20160127-02.txt for sha256sums of the packages and
manifest-20160127-02.txt.asc to confirm those shas.

See https://wiki.decred.org/Verifying_Binaries for more info on
verifying the files.
