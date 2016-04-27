# [v0.1.0_miners](https://github.com/decred/decred-release/releases/tag/v0.1.0_miners)

## 2016-04-27

This is a binary release of cgminer and ccminer for decred.  Binaries
for 64bit Linux and Windows are provided.

The attached manifest file can be used to confirm the sha of the
download.

See [README.md](./README.md) for more info on verifying the files.

Changes include:

| Description | Pull Request |
| --- | ---- |
| blake kernel optimisation for nvidia cards | [decred/cgminer#24](https://github.com/decred/cgminer/pull/24) |
| Fix broken intensity. | [decred/cgminer#25](https://github.com/decred/cgminer/pull/25) |
| Fixed high CPU usage and time-too-new bug. | [decred/cgminer#27](https://github.com/decred/cgminer/pull/27) |
| .gitignore the rest of the generated files. | [decred/cgminer#28](https://github.com/decred/cgminer/pull/28) |
| Use decred specific version number. | [decred/cgminer#30](https://github.com/decred/cgminer/pull/30) |
| Add script to make mostly static linux builds. | [decred/cgminer#31](https://github.com/decred/cgminer/pull/31) |
| Have git ignore more temporary files. | [decred/cgminer#32](https://github.com/decred/cgminer/pull/32) |
| Add message at the end of build script. | [decred/cgminer#33](https://github.com/decred/cgminer/pull/33) |
| Add cross-compiling for windows build script. | [decred/cgminer#34](https://github.com/decred/cgminer/pull/34) |
| Bump for v0.1.0 | [decred/cgminer#35](https://github.com/decred/cgminer/pull/35) |
| ~10% speedup | [decred/ccminer#1](https://github.com/decred/ccminer/pull/1) |
| Use decred specific version number. | [decred/ccminer#2](https://github.com/decred/ccminer/pull/2) |
| Add script to make mostly static linux builds. | [decred/ccminer#3](https://github.com/decred/ccminer/pull/3) |
| Add message at the end of build script. | [decred/ccminer#4](https://github.com/decred/ccminer/pull/4) |
| Bump for v0.1.0 | [decred/ccminer#5](https://github.com/decred/ccminer/pull/5) |

## Commits

This release was built from:

| Repository | Commit Hash |
| --- | ---- |
| cgminer | 2616e256b99924de072bfde177abe3c470cd1b32 |
| ccminer | 45aaecfcdf5ae23ab1e9090cb3140d758a287bd4 |

# [v0.0.5_cgminer](https://github.com/decred/decred-release/releases/tag/v0.0.5_cgminer)

## 2016-02-10

This is a build of cgminer for 32bit windows (with no TLS support).

The attached manifest file can be used to confirm the sha of the
download.

See https://wiki.decred.org/Verifying_Binaries for more info on
verifying the files.

# [v0.0.4_cgminer](https://github.com/decred/decred-release/releases/tag/v0.0.4_cgminer)

## 2016-02-08

Bugfix for incorrectly displaying difficulty.

See manifest-cgminer-20160208.txt for sha256sums of the packages and
manifest-cgminer-20160208.txt.asc to confirm those shas.

See https://wiki.decred.org/Verifying_Binaries for more info on
verifying the files.

# [v0.0.3_cgminer](https://github.com/decred/decred-release/releases/tag/v0.0.3_cgminer)

## 2016-02-04

This is the testnet pre-release of cgminer for Decred.

This is an update to the Windows build ONLY to address missing
libraries in the previous release. There are no code changes.

See manifest-cgminer-20160204-2.txt for sha256sums of the packages and
manifest-cgminer-20160204-2.txt.asc to confirm those shas.

See https://wiki.decred.org/Verifying_Binaries for more info on
verifying the files.

# [v0.0.2_cgminer](https://github.com/decred/decred-release/releases/tag/v0.0.2_cgminer)

## 2016-02-04

This is the testnet pre-release of cgminer for Decred.

The packages below contain platform specific copies of cgminer.

Bug fixes, rebuilt with support for the AMD Display Library, TLS-self
signed certificates can now be specified with --cert

See manifest-cgminer-20160204.txt for sha256sums of the packages and
manifest-cgminer-20160204.txt.asc to confirm those shas.

See https://wiki.decred.org/Verifying_Binaries for more info on
verifying the files.

# [v0.0.1_cgminer](https://github.com/decred/decred-release/releases/tag/v0.0.1_cgminer)

## 2016-01-28

This is the testnet pre-release of cgminer for Decred.

The packages below contain platform specific copies of cgminer.

The windows binary remains the same as in the previous release. The
linux package has been updated.

See manifest-cgminer-20160128.txt for sha256sums of the packages and
manifest-cgminer-20160128.txt.asc to confirm those shas.

See https://wiki.decred.org/Verifying_Binaries for more info on
verifying the files.

# [v0.0_cgminer](https://github.com/decred/decred-release/releases/tag/v0.0.0_cgminer)

## 2016-01-27

This is the testnet pre-release of cgminer for Decred.

The packages below contain platform specific copies of cgminer.

See manifest-cgminer-20160127.txt for sha256sums of the packages and
manifest-cgminer-20160127.txt.asc to confirm those shas.

See https://wiki.decred.org/Verifying_Binaries for more info on
verifying the files

