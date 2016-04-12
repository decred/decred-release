# [v0.0.10](https://github.com/decred/decred-release/releases/tag/v0.0.10)

## 2016-04-06

This release contains updated binary files (dcrd, dcrctl, dcrwallet)
for various platforms and is primarily a bugfix for dcrwallet.

See manifest-20160406-01.txt for sha256sums of the packages and
manifest-201600406-01.txt.asc to confirm those shas.

See https://wiki.decred.org/Verifying_Binaries for more info on
verifying the files.

This release contains various fixes and improvements.

Changes include:

| Description | Pull Request |
| --- | ---- |
| Introduce a new utility to show dev premine taint. | [decred/dcrd#100](https://github.com/decred/dcrd/pull/100) |
| Bump for v0.0.10 | [decred/dcrd#101](https://github.com/decred/dcrd/pull/101) |
| Add new JSON handling for RPC commands and livetickets command | [decred/dcrd#102](https://github.com/decred/dcrd/pull/102) |
| Add stake txscript types in ListUnspent to be spendable | [decred/dcrwallet#151](https://github.com/decred/dcrwallet/pull/151) |
| Make dcrwallet pass all goclean.sh tests. | [decred/dcrwallet#155](https://github.com/decred/dcrwallet/pull/155) |
| Change initilialize to use proper index (extIdx) | [decred/dcrwallet#158](https://github.com/decred/dcrwallet/pull/158) |
| Bump for v0.0.10 | [decred/dcrwallet#159](https://github.com/decred/dcrwallet/pull/159) |
| Fix address pool syncing and add new RPC commands for the address pools | [decred/dcrwallet#161](https://github.com/decred/dcrwallet/pull/161) |
| Rollback namespace transactions when bucket is not found. | [decred/dcrwallet#163](https://github.com/decred/dcrwallet/pull/163) |
| Fix watching only wallets | [decred/dcrwallet#164](https://github.com/decred/dcrwallet/pull/164) |
| Fix case on comments | [decred/dcrwallet#165](https://github.com/decred/dcrwallet/pull/165) |

## Commits

This release was built from:

| Repository | Commit Hash |
| --- | ---- |
| decred/dcrd | 5658c503c3ad9e8b6e7eaec5183f9fe4a2e32241 |
| decred/dcrwallet | f1d9bd630188da91f7e817c49830c29d365c615d |
| decred/dcrrpcclient | b3f48780a0d68e24ef6e915e930a1c1e58b69810 |
| decred/dcrutil | 9bb7f64962cee52bb46ce588aa91ef0e6e7bb1a9 |

# [v0.0.9](https://github.com/decred/decred-release/releases/tag/v0.0.9)

## 2016-04-01

This release contains updated binary files (dcrd, dcrctl, dcrwallet)
for various platforms.

See manifest-20160401-01.txt for sha256sums of the packages and
manifest-201600401-01.txt.asc to confirm those shas.

See https://wiki.decred.org/Verifying_Binaries for more info on
verifying the files.

This release contains various fixes and improvements.

Changes include:

| Description | Pull Request |
| --- | ---- |
| txscript: New function IsUnspendable | [decred/dcrd#96](https://github.com/decred/dcrd/pull/96) |
| Get travis-ci to work again. | [decred/dcrd#97](https://github.com/decred/dcrd/pull/97) |
| peer: Remove extraneous hasTimestamp check. | [decred/dcrd#98](https://github.com/decred/dcrd/pull/98) |
| Bump to v0.0.9 for release. | [decred/dcrd#99](https://github.com/decred/dcrd/pull/99) |
| Print version string at startup. | [decred/dcrwallet#126](https://github.com/decred/dcrwallet/pull/126) |
| Sync with upstream | [decred/dcrwallet#127](https://github.com/decred/dcrwallet/pull/127) |
| Use wtxmgr for input selection. | [decred/dcrwallet#130](https://github.com/decred/dcrwallet/pull/130) |
| Fix updating the UTXO set for imported addresses | [decred/dcrwallet#133](https://github.com/decred/dcrwallet/pull/133) |
| Help prevent errors during initial sync by waiting for it to finish | [decred/dcrwallet#136](https://github.com/decred/dcrwallet/pull/136) |
| Remove voting pool package. | [decred/dcrwallet#135](https://github.com/decred/dcrwallet/pull/135) |
| Fix proportionmissed calc --> missed / missed + voted | [decred/dcrwallet#138](https://github.com/decred/dcrwallet/pull/138) |
| Refactor address pool code and automatically resync accounts from seed | [decred/dcrwallet#134](https://github.com/decred/dcrwallet/pull/134) |
| fix waddrmgr tests | [decred/dcrwallet#139](https://github.com/decred/dcrwallet/pull/139) |
| Modify the logic for password prompting | [decred/dcrwallet#142](https://github.com/decred/dcrwallet/pull/142) |
| Fix address pool panics on start up | [decred/dcrwallet#143](https://github.com/decred/dcrwallet/pull/143) |
| Add goclean.sh script from btcd. | [decred/dcrwallet#144](https://github.com/decred/dcrwallet/pull/144) |
| Bump to v0.0.9 for release. | [decred/dcrwallet#150](https://github.com/decred/dcrwallet/pull/150) |
| Update to all dcrutil tests so they successfully pass. | [decred/dcrutil#4](https://github.com/decred/dcrutil/pull/4) |
| Fix filter_test TestFilterInsertP2PubKeyOnly with correct info | [decred/dcrutil#6](https://github.com/decred/dcrutil/pull/6) |
| Fix a test output for go1.6. | [decred/dcrutil#8](https://github.com/decred/dcrutil/pull/8) |
| Enable travis-ci. | [decred/dcrutil#7](https://github.com/decred/dcrutil/pull/7) |

## Commits

This release was built from:

| Repository | Commit Hash |
| --- | ---- |
| decred/dcrd | 0ed0e815b0cd59d48380d125d47ff0de833ec43c |
| decred/dcrwallet | 4387fa379d01d125db7c9e6fcada51f8316cb0f6 |
| decred/dcrrpcclient | b3f48780a0d68e24ef6e915e930a1c1e58b69810 |
| decred/dcrutil | 9bb7f64962cee52bb46ce588aa91ef0e6e7bb1a9 |

# [v0.0.8](https://github.com/decred/decred-release/releases/tag/v0.0.8)

## 2016-03-18

This release contains updated binary files (dcrd, dcrctl, dcrwallet)
for various platforms.

See manifest-20160318-01.txt for sha256sums of the packages and
manifest-20160318-01.txt.asc to confirm those shas.

See https://wiki.decred.org/Verifying_Binaries for more info on
verifying the files.

This release contains various fixes and improvements.

Changes include:

| Description | Pull Request |
| --- | ---- |
| Update configuring_tor.md | [decred/dcrd#88](https://github.com/decred/dcrd/pull/88) |
| Add and implement the getticketpoolvalue JSON RPC command | [decred/dcrd#90](https://github.com/decred/dcrd/pull/90) |
| Add lookup of ticket commitments to addrindex | [decred/dcrd#92](https://github.com/decred/dcrd/pull/92) |
| Fix minor goclean issues. | [decred/dcrd#94](https://github.com/decred/dcrd/pull/94) |
| Add balancetomaintain rpc json parts | [decred/dcrd#93](https://github.com/decred/dcrd/pull/93) |
| Bump for 0.0.8 | [decred/dcrd#95](https://github.com/decred/dcrd/pull/95) |
| Fix a bug relating to relevantTx handling and uncaught error | [decred/dcrwallet#103](https://github.com/decred/dcrwallet/pull/103) |
| Overhaul accounts to function correctly | [decred/dcrwallet#104](https://github.com/decred/dcrwallet/pull/104) |
| Use a random address for 0-value outputs | [decred/dcrwallet#115](https://github.com/decred/dcrwallet/pull/115) |
| Fix all rpclisten references in documentation | [decred/dcrwallet#118](https://github.com/decred/dcrwallet/pull/118) |
| Fix wallet resyncing from seed and address index positioning | [decred/dcrwallet#121](https://github.com/decred/dcrwallet/pull/121) |
| Add err check for unchecked | [decred/dcrwallet#123](https://github.com/decred/dcrwallet/pull/123) |
| Catch vootingpool up with current apis | [decred/dcrwallet#122](https://github.com/decred/dcrwallet/pull/122) |
| Add new balancetomaintain rpc command | [decred/dcrwallet#120](https://github.com/decred/dcrwallet/pull/120) |
| Bump for 0.0.8 | [decred/dcrwallet#124](https://github.com/decred/dcrwallet/pull/124) |
| Set Tree field when converting wire.OutPoints. | [decred/dcrrpcclient#8](https://github.com/decred/dcrrpcclient/pull/10) |

## Commits

This release was built from:

| Repository | Commit Hash |
| --- | ---- |
| decred/dcrd | 967952c7cbf23a622cf5ada5101658037f827a2f |
| decred/dcrwallet | a981017f624e27816c6aba21b00c2086b1b5d852 |
| decred/dcrrpcclient | b3f48780a0d68e24ef6e915e930a1c1e58b69810 |
| decred/dcrutil | ae0e66b98e49e836618c01cfa4d1b3d6077e5ae7 |

# [v0.0.7](https://github.com/decred/decred-release/releases/tag/v0.0.7)

## 2016-03-09

Patched release to allow multisig votes to be properly accepted by daemons with IsStandard

Changes include:

| Description | Pull Request |
| --- | ---- |
| Fix storing the ticket database to disk on close | [decred/dcrd#80](https://github.com/decred/dcrd/pull/80) |
| Reduce likelihood of vote spam | [decred/dcrd#82](https://github.com/decred/dcrd/pull/82) |
| Optimize mining checks for various stake transactions | [decred/dcrd#83](https://github.com/decred/dcrd/pull/83) |
| Sync to upstream 0280fa0 | [decred/dcrd#78](https://github.com/decred/dcrd/pull/78) |
| Revert sync merge | [decred/dcrd#85](https://github.com/decred/dcrd/pull/85) |
| Correct the expected number of inputs for stake P2SH outputs | [decred/dcrd#86](https://github.com/decred/dcrd/pull/86) |
| Bump version for release | [decred/dcrd#87](https://github.com/decred/dcrd/pull/87) |
| Only access isClosed inside the mutex in wtxmgr. | [decred/dcrwallet#94](https://github.com/decred/dcrwallet/pull/94) |
| Fixes to work with dcrd sync to 08/11/15 | [decred/dcrwallet#91](https://github.com/decred/dcrwallet/pull/91) |
| Switch from log.Debug to log.Debugf. | [decred/dcrwallet#96](https://github.com/decred/dcrwallet/pull/96) |
| Revert "Fixes to work with dcrd sync" | [decred/dcrwallet#100](https://github.com/decred/dcrwallet/pull/100) |
| Bump version for patch release | [decred/dcrwallet#101](https://github.com/decred/dcrwallet/pull/101) |
| Fix "Established connection" log message. | [decred/dcrrpcclient#8](https://github.com/decred/dcrrpcclient/pull/9) |
| Ayp sync 1c7f05 | [decred/dcrutil#1](https://github.com/decred/dcrutil/pull/1) |
| Revert sync commit | [decred/dcrutil#2](https://github.com/decred/dcrutil/pull/2) |

## Commits

This release was built from:

| Repository | Commit Hash |
| --- | ---- |
| decred/dcrd | f2cc01cef2e58d788212dc28633c2d7b3cdf68e0 |
| decred/dcrwallet | d776d972f2f0c7b440dfbea5a10ba7ac4627cfbe |
| decred/dcrrpcclient | 4691756e416483e497d41f8883e5f432167983a2 |
| decred/dcrutil | ae0e66b98e49e836618c01cfa4d1b3d6077e5ae7 |

# [v0.0.6](https://github.com/decred/decred-release/releases/tag/v0.0.6)

## 2016-03-04

This release contains updated binary files (dcrd, dcrctl, dcrwallet)
for various platforms.

See manifest-20160304-01.txt for sha256sums of the packages and
manifest-20160304-01.txt.asc to confirm those shas.

See https://wiki.decred.org/Verifying_Binaries for more info on
verifying the files.

This release contains various fixes and improvements.

Changes include:

| Description | Pull Request |
| --- | ---- |
| Fix missing rpc help | [decred/dcrd#61](https://github.com/decred/dcrd/pull/61) |
| Fix a copy-paste error in chainsvrcmds.go | [decred/dcrd#63](https://github.com/decred/dcrd/pull/63) |
| Bug fix for checkBlockForHiddenVotes corrupting block templates | [decred/dcrd#68](https://github.com/decred/dcrd/pull/68) |
| Cherry pick commits required for wallet sync. | [decred/dcrd#69](https://github.com/decred/dcrd/pull/69) |
| Fix a panic caused by accessing unassigned pointer | [decred/dcrd#70](https://github.com/decred/dcrd/pull/70) |
| Add new RPC handlers for get/setticketfee | [decred/dcrd#71](https://github.com/decred/dcrd/pull/71) |
| Add getticketsvotebits batched command for wallet RPC | [decred/dcrd#72](https://github.com/decred/dcrd/pull/72) |
| Update default_ports.md | [decred/dcrd#75](https://github.com/decred/dcrd/pull/75) |
| Add consolidate cmd and response framework to the JSON RPC | [decred/dcrd#59](https://github.com/decred/dcrd/pull/59) |
| Add the new RPC function existsmempooltxs | [decred/dcrd#74](https://github.com/decred/dcrd/pull/74) |
| Fix bug displaying the wrong number of votes in getstakeinfo | [decred/dcrwallet#66](https://github.com/decred/dcrwallet/pull/66) |
| Merge upstream [btcsuite/btcwallet](https://github.com/btcsuite/btcwallet) code | [decred/dcrwllet#65](https://github.com/decred/dcrwallet/pull/65) |
| Add getstakeinfo online help. | [decred/dcrwallet#71](https://github.com/decred/dcrwallet/pull/71) |
| Change 'voted' in getstakeinfo to only return blockchain votes | [decred/dcrwallet#73](https://github.com/decred/dcrwallet/pull/73) |
| Get/SetTicketFee RPC and fix fee calculation in purchaseTicket | [decred/dcrwallet#72](https://github.com/decred/dcrwallet/pull/72) |
| Attempt to streamline getting/setting of fees for main/testnet | [decred/dcrwallet#74](https://github.com/decred/dcrwallet/pull/74) |
| Added getticketsvotebits functionality to the legacy RPC | [decred/dcrwallet#75](https://github.com/decred/dcrwallet/pull/75) |
| README.md: Update URL to releases | [decred/dcrwallet#78](https://github.com/decred/dcrwallet/pull/78) |
| Add wallet handling for getgenerate command. | [decred/dcrwallet#79](https://github.com/decred/dcrwallet/pull/79) |
| Correct TicketsForAddress returning pruned tickets | [decred/dcrwallet#80](https://github.com/decred/dcrwallet/pull/80) |
| Stop uses of database before closing db. | [decred/dcrwallet#87](https://github.com/decred/dcrwallet/pull/87) |
| Allow newlines and extra spaces when entering seed. | [decred/dcrwallet#88](https://github.com/decred/dcrwallet/pull/88) |
| Fix docs in grpc | [decred/dcrwallet#89](https://github.com/decred/dcrwallet/pull/89) |
| Prevent addresses from being shown more than once. | [decred/dcrwallet#89](https://github.com/decred/dcrwallet/pull/82) |
| Validate the address provided to --ticketaddress | [decred/dcrwallet#90](https://github.com/decred/dcrwallet/pull/90) |
| Add consolidate command handling to the wallet JSON RPC | [decred/dcrwallet#61](https://github.com/decred/dcrwallet/pull/61) |
| Update go versions used by travis. | [decred/dcrrpcclient#6](https://github.com/decred/dcrrpcclient/pull/6) |
| Add a hook for getticketsvotebits in wallet | [decred/dcrrpcclient#7](https://github.com/decred/dcrrpcclient/pull/7) |
| Add the existsmempooltxs command for daemon | [decred/dcrrpcclient#8](https://github.com/decred/dcrrpcclient/pull/8) |

## Commits

This release was built from:

| Repository | Commit Hash |
| --- | ---- |
| decred/dcrd | 8f0d8f2d850edef4fa684ad872512ec9c0434f20 |
| decred/dcrwallet | 3d845de5a8650459db46251883a63b78fd55d404 |
| decred/dcrrpcclient | 7181e59ba727f8e6cb2f3919bc490549f81e4d54 |
| decred/dcrutil | 025b0fb50cfb446491a6988fab4cef333830e35c |

# [v0.0.5](https://github.com/decred/decred-release/releases/tag/v0.0.5)

## 2016-02-26

This release contains updated binary files (dcrd, dcrctl, dcrwallet)
for various platforms.

See manifest-20160226-01.txt for sha256sums of the packages and
manifest-20160226-01.txt.asc to confirm those shas.

See https://wiki.decred.org/Verifying_Binaries for more info on
verifying the files.

This release is primarily to add the ability for wallet to
automatically remove old tickets and expired transactions.

Other changes include:
* Add getstakeinfo to rpcAskWallet list to return proper error
* Correct version numbers
* Fix coin supply counter to reduce work and tax subsidy based on voters
* Fix a bug that caused votes and revocations not being stored
* Add listscripts RPC command handling

## Commits

This release was built from:

| Repository | Commit Hash |
| --- | ---- |
| decred/dcrd | fbede4978022f7121f80a1ec02a217b7498c4f5b |
| decred/dcrwallet | ee2a72abe35f690fcc54c3c6234e617c79a88d19 |
| decred/dcrrpcclient | 680d8ff9cd81c017c28fd867494e20deea08e48c |
| decred/dcrutil | 025b0fb50cfb446491a6988fab4cef333830e35c |

# [v0.0.4](https://github.com/decred/decred-release/releases/tag/v0.0.4)

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

# [v0.0.3](https://github.com/decred/decred-release/releases/tag/v0.0.3)

## 2016-02-09

This wallet only release resolves an upstream wallet bug (see
decred/dcrutil 8aae5a2dacf45b7f5ee9b59c393118bc48647861).

Platform specific files are attached.

See manifest-20160209-01.txt for sha256sums of the packages and
manifest-20160209-01.txt.asc to confirm those shas.

See https://wiki.decred.org/Verifying_Binaries for more info on
verifying the files.

# [v0.0.2](https://github.com/decred/decred-release/releases/tag/v0.0.2)

## 2016-02-08

This release is an unencrypted version of the current mainnet enabled
code.

The packages below contain platform specific copies of drcd,
dcrwallet, and dcrctl.

See manifest-20160208-01.txt for sha256sums of the packages and
manifest-20160208-01.txt.asc to confirm those shas.

See https://wiki.decred.org/Verifying_Binaries for more info on
verifying the files.

# [v0.0.1](https://github.com/decred/decred-release/releases/tag/v0.0.1)

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

# [v0.0](https://github.com/decred/decred-release/releases/tag/v0.0)

## 2016-01-27

This is the testnet pre-release of Decred.

The packages below contain platform specific copies of drcd,
dcrwallet, and dcrctl.

See manifest-20160127-02.txt for sha256sums of the packages and
manifest-20160127-02.txt.asc to confirm those shas.

See https://wiki.decred.org/Verifying_Binaries for more info on
verifying the files.
