# [v0.1.5](https://github.com/decred/decred-release/releases/tag/v0.1.5)

## 2016-06-07

This release contains updated binary files (dcrd, dcrctl, dcrwallet,
and dcrticketbuyer) for various platforms.

See manifest-20160607-01.txt for sha256sums of the packages and
manifest-201600607-01.txt.asc to confirm those shas.

See https://wiki.decred.org/Verifying_Binaries for more info on
verifying the files.

This release contains various fixes and improvements.

Changes include:

| Description | Pull Request |
| --- | ---- |
| peer: Unexport the mru inventory map. | [decred/dcrd#220](https://github.com/decred/dcrd/pull/220) |
| mempool/mining: Decouple and optimize priority calcs. | [decred/dcrd#223](https://github.com/decred/dcrd/pull/223) |
| database: Major redesign of database package | [decred/dcrd#91](https://github.com/decred/dcrd/pull/91) |
| mempool/mining: Introduce TxSource interface. | [decred/dcrd#225](https://github.com/decred/dcrd/pull/225) |
| mempool: Introduce mempoolConfig. | [decred/dcrd#226](https://github.com/decred/dcrd/pull/226) |
| mining: Create skeleton package. | [decred/dcrd#227](https://github.com/decred/dcrd/pull/227) |
| peer: Add DisableRelayTx to config. | [decred/dcrd#228](https://github.com/decred/dcrd/pull/228) |
| peer: Rename variable for consistency. | [decred/dcrd#229](https://github.com/decred/dcrd/pull/229) |
| Apply various upstream comment fixes. | [decred/dcrd#230](https://github.com/decred/dcrd/pull/230) |
| Merge upstream copyright date updates. | [decred/dcrd#231](https://github.com/decred/dcrd/pull/231) |
| peer: Simplify PushAddrMsg method loop. | [decred/dcrd#232](https://github.com/decred/dcrd/pull/232) |
| wire: Minor code clean up. | [decred/dcrd#233](https://github.com/decred/dcrd/pull/233) |
| txscript: Fix typo in README | [decred/dcrd#234](https://github.com/decred/dcrd/pull/234) |
| database: Merge through Implement cache layer. | [decred/dcrd#235](https://github.com/decred/dcrd/pull/235) |
| dcrjson/txscript: Merge arm-specific test updates. | [decred/dcrd#236](https://github.com/decred/dcrd/pull/236) |
| rpcserver: Optimize filteraddr code. | [decred/dcrd#237](https://github.com/decred/dcrd/pull/237) |
| Change Vin field AmountIn to display coins not int64 | [decred/dcrd#238](https://github.com/decred/dcrd/pull/238) |
| Fix median of slice of Amounts for ticketfeeinfo. | [decred/dcrd#239](https://github.com/decred/dcrd/pull/239) |
| Use atomic operations instead of mutexes. | [decred/dcrd#240](https://github.com/decred/dcrd/pull/240) |
| wire: Implement sendheaders command | [decred/dcrd#241](https://github.com/decred/dcrd/pull/241) |
| peer: Consolidate several public methods. | [decred/dcrd#242](https://github.com/decred/dcrd/pull/242) |
| server: Make consistent use of svr peer stringer. | [decred/dcrd#243](https://github.com/decred/dcrd/pull/243) |
| txscript: Comment improvements and fixes | [decred/dcrd#244](https://github.com/decred/dcrd/pull/244) |
| Implement banning based on dynamic ban scores | [decred/dcrd#245](https://github.com/decred/dcrd/pull/245) |
| wire: Export (read write)(VarInt VarBytes). | [decred/dcrd#246](https://github.com/decred/dcrd/pull/246) |
| Log block processing time in CHAN with debug on | [decred/dcrd#247](https://github.com/decred/dcrd/pull/247) |
| multi: Fix several misspellings in the comments. | [decred/dcrd#248](https://github.com/decred/dcrd/pull/248) |
| multi: Update with result of gofmt -s. | [decred/dcrd#249](https://github.com/decred/dcrd/pull/249) |
| server: Appropriately name inbound peers map in peerState. | [decred/dcrd#250](https://github.com/decred/dcrd/pull/250) |
| docs: Update READMEs with some current details. | [decred/dcrd#222](https://github.com/decred/dcrd/pull/252) |
| peer: declare QueueMessage()'s doneChan as send only. | [decred/dcrd#253](https://github.com/decred/dcrd/pull/253) |
| peer: Implement sendheaders support (BIP0130). | [decred/dcrd#254](https://github.com/decred/dcrd/pull/254) |
| server: Cleanup and optimize handleBroadcastMsg. | [decred/dcrd#255](https://github.com/decred/dcrd/pull/255) |
| config: New option --blocksonly | [decred/dcrd#256](https://github.com/decred/dcrd/pull/256) |
| Keep track of recently rejected transactions. | [decred/dcrd#257](https://github.com/decred/dcrd/pull/257) |
| mining: Export block template fields. | [decred/dcrd#258](https://github.com/decred/dcrd/pull/258) |
| server: Optimize map limiting in block manager. | [decred/dcrd#259](https://github.com/decred/dcrd/pull/259) |
| chaincfg: Register networks instead of hard coding. | [decred/dcrd#260](https://github.com/decred/dcrd/pull/260) |
| chaincfg: Consolidate tests into the chaincfg pkg. | [decred/dcrd#261](https://github.com/decred/dcrd/pull/261) |
| txscript: Correct comments on alt stack methods. | [decred/dcrd#262](https://github.com/decred/dcrd/pull/262) |
| mempool: Create and use mempoolPolicy. | [decred/dcrd#263](https://github.com/decred/dcrd/pull/263) |
| Asynchronously call TicketPoolValue to stop block manager blocking | [decred/dcrd#265](https://github.com/decred/dcrd/pull/265) |
| Add rescan and scanfrom options to importprivkey and importscript | [decred/dcrd#267](https://github.com/decred/dcrd/pull/267) |
| Bump for v0.1.5 | [decred/dcrd#268](https://github.com/decred/dcrd/pull/268) |
| Fix fee calculation for revocations and rebroastcast on start up | [decred/dcrwallet#254](https://github.com/decred/dcrwallet/pull/254) |
| rpctest behavioral test suite | [decred/dcrwallet#241](https://github.com/decred/dcrwallet/pull/241) |
| Remove unused SendRawTransaction func in StakeStore | [decred/dcrwallet#256](https://github.com/decred/dcrwallet/pull/256) |
| Remove transactions in reverse order when rolling back blocks | [decred/dcrwallet#263](https://github.com/decred/dcrwallet/pull/263) |
| Bump for v0.1.5 | [decred/dcrwallet#265](https://github.com/decred/dcrwallet/pull/265) |
| Add optional resyncing options to importscript and importprivkey | [decred/dcrwallet#264](https://github.com/decred/dcrwallet/pull/264) |
| Add gettickets to the wallet RPC client handlers | [decred/dcrrpcclient#26](https://github.com/decred/dcrrpcclient/pull/26) |
| Add rescan options for importprivkey and importscript | [decred/dcrrpcclient#27](https://github.com/decred/dcrrpcclient/pull/27) |
| Add AmountSorter, which implements the sort.Interface, for Amount. | [decred/dcrutil#12](https://github.com/decred/dcrutil/pull/12) |
| Bind to localhost only by default | [decred/dcrticketbuyer#3](https://github.com/decred/dcrticketbuyer/pull/3) |
| Fix bug where fee from difficulty window was 0 | [decred/dcrticketbuyer#5](https://github.com/decred/dcrticketbuyer/pull/5) |
| Add ability to choose which price average to use | [decred/dcrticketbuyer#6](https://github.com/decred/dcrticketbuyer/pull/6) |
| Warn the user on start up | [decred/dcrticketbuyer#7](https://github.com/decred/dcrticketbuyer/pull/7) |
| Update glide and fix unlikely simnet panic | [decred/dcrticketbuyer#8](https://github.com/decred/dcrticketbuyer/pull/8) |

## Notes

## Commits

This release was built from:

| Repository | Commit Hash |
| --- | ---- |
| decred/dcrd | b86959378985f538288c5a8d5184244d4692e0e6 |
| decred/dcrwallet | 3942d8b165842285a24973bc2e42708a65ff66ff |
| decred/dcrrpcclient | f3c620d63cb02aec0c1152a72d3c8669b92a2fb5 |
| decred/dcrutil | 4a3bdb1cb08b49811674750998363b8b8ccfd66e |
| decred/dcrticketbuyer | 65641c4458624f5a9c76116b791d48e68fe98897 |

# [v0.1.4](https://github.com/decred/decred-release/releases/tag/v0.1.4)

## 2016-05-26

This release contains updated binary files (dcrd, dcrctl, dcrwallet,
and dcrticketbuyer) for various platforms.

See manifest-20160526-01.txt for sha256sums of the packages and
manifest-201600526-01.txt.asc to confirm those shas.

See https://wiki.decred.org/Verifying_Binaries for more info on
verifying the files.

This release contains various fixes and improvements.

Changes include:

| Description | Pull Request |
| --- | ---- |
| Sync upstream through July 1, 2015 | [decred/dcrd#163](https://github.com/decred/dcrd/pull/163) |
| Sync upstream through July 22, 2015  | [decred/dcrd#164](https://github.com/decred/dcrd/pull/164) |
| Sync upstream through August 9, 2015 | [decred/dcrd#166](https://github.com/decred/dcrd/pull/166) |
| Reject very old votes from the memory pool | [decred/dcrd#168](https://github.com/decred/dcrd/pull/168) |
| Adjust getblockheader result for Decred | [decred/dcrd#170](https://github.com/decred/dcrd/pull/170) |
| Check for hidden votes by ticket hash, not vote hash | [decred/dcrd#169](https://github.com/decred/dcrd/pull/169) |
| Sync upstream through Aug 23, 2015 | [decred/dcrd#172](https://github.com/decred/dcrd/pull/172) |
| Waste less memory if sighash optimizations are on | [decred/dcrd#171](https://github.com/decred/dcrd/pull/171) |
| Sync upstream through Sep 2, 2015. | [decred/dcrd#174](https://github.com/decred/dcrd/pull/174) |
| Sync upstream through Sep 17, 2015. | [decred/dcrd#175](https://github.com/decred/dcrd/pull/175) |
| Sync upstream through Sep 24, 2015. | [decred/dcrd#177](https://github.com/decred/dcrd/pull/177) |
| Remove legacy Bitcoin addr encoding bug | [decred/dcrd#179](https://github.com/decred/dcrd/pull/179) |
| Sync upstream through Sep 28, 2015. | [decred/dcrd#178](https://github.com/decred/dcrd/pull/178) |
| wire: Use ordered Service Flags. | [decred/dcrd#182](https://github.com/decred/dcrd/pull/182) |
| rpcserver: Optimize JSON raw tx input list create. | [decred/dcrd#180](https://github.com/decred/dcrd/pull/180) |
| txscript: Sync upstream makeScriptNum tests. | [decred/dcrd#181](https://github.com/decred/dcrd/pull/181) |
| Fix VinPrevOut fields for Decred | [decred/dcrd#183](https://github.com/decred/dcrd/pull/183) |
| Add reverse order option to searchrawtransactions rpc | [decred/dcrd#185](https://github.com/decred/dcrd/pull/185) |
| main: Limit garbage collection percentage. (#686) | [decred/dcrd#187](https://github.com/decred/dcrd/pull/187) |
| Integrate a valid ECDSA signature cache. | [decred/dcrd#189](https://github.com/decred/dcrd/pull/189) |
| Add a checkpoint for block 24480 | [decred/dcrd#190](https://github.com/decred/dcrd/pull/190) |
| dcrjson: Add errors to InfoChainResult | [decred/dcrd#191](https://github.com/decred/dcrd/pull/191) |
| Use same fee policies across all networks. | [decred/dcrd#160](https://github.com/decred/dcrd/pull/160) |
| rpcserver: Correct verifymessage hash generation. | [decred/dcrd#192](https://github.com/decred/dcrd/pull/192) |
| Correct a few style related issues found by golint. | [decred/dcrd#193](https://github.com/decred/dcrd/pull/193) |
| config: New option --minrelaytxfee | [decred/dcrd#194](https://github.com/decred/dcrd/pull/194) |
| Fix magic peer initial protocol value | [decred/dcrd#195](https://github.com/decred/dcrd/pull/195) |
| peer: Refactor peer code into its own package. | [decred/dcrd#197](https://github.com/decred/dcrd/pull/197) |
| docs: Make various README.md files consistent. | [decred/dcrd#201](https://github.com/decred/dcrd/pull/201) |
| peer: Sync upstream fixes and improvements. | [decred/dcrd#202](https://github.com/decred/dcrd/pull/202) |
| Use the correct heap sorting function | [decred/dcrd#199](https://github.com/decred/dcrd/pull/199) |
| Move non-mempool specific functions to new file. | [decred/dcrd#203](https://github.com/decred/dcrd/pull/203) |
| dcrjson: Add optional locktime to createrawtransaction | [decred/dcrd#204](https://github.com/decred/dcrd/pull/204) |
| Sync upstream blockmanager comments improvements. | [decred/dcrd#205](https://github.com/decred/dcrd/pull/205) |
| Sync upstream comment and error improvements. | [decred/dcrd#152](https://github.com/decred/dcrd/pull/206) |
| chaincfg: Move DNS Seeds to chaincfg. | [decred/dcrd#209](https://github.com/decred/dcrd/pull/209) |
| peer: Fix failing test case due to wrong TimeOffset | [decred/dcrd#210](https://github.com/decred/dcrd/pull/210) |
| peer/server: various fixes from upstream | [decred/dcrd#211](https://github.com/decred/dcrd/pull/211) |
| mempool/peer: Sync upstream comment updates. | [decred/dcrd#212](https://github.com/decred/dcrd/pull/212) |
| mempool: Move checkTransactionStandard to policy. | [decred/dcrd#214](https://github.com/decred/dcrd/pull/214) |
| dcrd: do not process empty getdata messages | [decred/dcrd#215](https://github.com/decred/dcrd/pull/215) |
| Bump for v0.1.4 | [decred/dcrd#221](https://github.com/decred/dcrd/pull/221) |
| rpcserver: Add filteraddrs param to srt API. | [decred/dcrd#216](https://github.com/decred/dcrd/pull/216) |
| peer: Combine stats struct into peer struct. | [decred/dcrd#217](https://github.com/decred/dcrd/pull/217) |
| Fix dropaddrindex flag usage message | [decred/dcrd#218](https://github.com/decred/dcrd/pull/218) |
| mining: Refactor policy into its own struct. | [decred/dcrd#219](https://github.com/decred/dcrd/pull/219) |
| peer: fix panic due to err in handleVersionMsg | [decred/dcrd#222](https://github.com/decred/dcrd/pull/222) |
| Use the block timestamp on block insertion, not local | [decred/dcrwallet#240](https://github.com/decred/dcrwallet/pull/240) |
| fix spelling in comment | [decred/dcrwallet#243](https://github.com/decred/dcrwallet/pull/243) |
| Disable ticket purchase by default | [decred/dcrwallet#244](https://github.com/decred/dcrwallet/pull/244) |
| Enable stakepool for mainnet | [decred/dcrwallet#245](https://github.com/decred/dcrwallet/pull/245) |
| Change "Notifying unmined tx .." to Tracef instead of Errorf | [decred/dcrwallet#246](https://github.com/decred/dcrwallet/pull/246) |
| Enable vendor experiment earlier in travis script. | [decred/dcrwallet#247](https://github.com/decred/dcrwallet/pull/247) |
| Add offline wallet guide and movefunds utility | [decred/dcrwallet#252](https://github.com/decred/dcrwallet/pull/252) |
| Bump for v0.1.4 | [decred/dcrwallet#253](https://github.com/decred/dcrwallet/pull/253) |
| Update SearchRawTransaction calls for latest API. | [decred/dcrrpcclient#22](https://github.com/decred/dcrrpcclient/pull/22) |
| Sync upstream through Aug. 23, 2015  | [decred/dcrrpcclient#20](https://github.com/decred/dcrrpcclient/pull/20) |
| Review and fix. Mostly typos. | [decred/dcrrpcclient#21](https://github.com/decred/dcrrpcclient/pull/21) |
| Fix ticket fee info command handling | [decred/dcrrpcclient#23](https://github.com/decred/dcrrpcclient/pull/23) |
| Add optional locktime parameter to CreateRawTransaction APIs. | [decred/dcrrpcclient#24](https://github.com/decred/dcrrpcclient/pull/24) |
| Add filteraddrs param to searchrawtransactions. | [decred/dcrrpcclient#25](https://github.com/decred/dcrrpcclient/pull/25) |
| Sync upstream through July 28, 2015 | [decred/dcrutil#10](https://github.com/decred/dcrutil/pull/10) |
| Update docs for NewAmount. | [decred/dcrutil#11](https://github.com/decred/dcrutil/pull/11) |
| Add HTTP server user interface | [decred/dcrticketbuyer#1](https://github.com/decred/dcrticketbuyer/pull/1) |

## Notes

This release contains the initial version of
[dcrticketbuyer](https://github.com/decred/dcrticketbuyer).    Please
follow the link to get more information about how to run our automated
ticket buying software.

## Commits

This release was built from:

| Repository | Commit Hash |
| --- | ---- |
| decred/dcrd | 708b4007ac598e2f19fa15658b9450edd9a5f052 |
| decred/dcrwallet | c9476fab7067814497aac9692a4a8a4c98305ae8 |
| decred/dcrrpcclient | 231790f525623f78acc9a91bfd3845d52715aee5 |
| decred/dcrutil | 85fac3a15425f15408f1dcec28bfd4b18ea2f882 |
| decred/dcrticketbuyer | 471c747f656e30e951463bbca3bafbf5ecfd572f |

# [v0.1.3](https://github.com/decred/decred-release/releases/tag/v0.1.3)

## 2016-05-10

This release contains updated binary files (dcrd, dcrctl, dcrwallet)
for various platforms and is primarily a bugfix for dcrwallet.

See manifest-20160510-01.txt for sha256sums of the packages and
manifest-201600510-01.txt.asc to confirm those shas.

See https://wiki.decred.org/Verifying_Binaries for more info on
verifying the files.

This release contains various fixes and improvements.

Changes include:

| Description | Pull Request |
| --- | ---- |
| mempool: reduce lock contention | [decred/dcrd#152](https://github.com/decred/dcrd/pull/152) |
| Reject too low stake difficulty transactions and cache difficulty | [decred/dcrd#154](https://github.com/decred/dcrd/pull/154) |
| mempool: Synchronize btcd commits fixing orphan hang | [decred/dcrd#155](https://github.com/decred/dcrd/pull/155) |
| dcrd: handle signal SIGTERM (#688) | [decred/dcrd#156](https://github.com/decred/dcrd/pull/156) |
| Fix resyncing the ticket database after unexpected shutdown | [decred/dcrd#157](https://github.com/decred/dcrd/pull/157) |
| Add transaction type to listtransactions result | [decred/dcrd#158](https://github.com/decred/dcrd/pull/158) |
| Fix createrawssrtx command and logic | [decred/dcrd#159](https://github.com/decred/dcrd/pull/159) |
| Bump for v0.1.3 | [decred/dcrd#162](https://github.com/decred/dcrd/pull/162) |
| Remove btcd/wire dependency. | [decred/dcrwallet#229](https://github.com/decred/dcrwallet/pull/229) |
| Sync with upstream | [decred/dcrwallet#227](https://github.com/decred/dcrwallet/pull/227) |
| Fix glide.yaml hash in glide.lock. | [decred/dcrwallet#234](https://github.com/decred/dcrwallet/pull/234) |
| Add transaction type to listtransactions result | [decred/dcrwallet#231](https://github.com/decred/dcrwallet/pull/231) |
| Update glide repos | [decred/dcrwallet#6b2fbf8](https://github.com/decred/dcrwallet/commit/6b2fbf80a33fc52f20231fdd6e462419c2a27ff6) |
| Call the more reliable GetStakeDifficulty for ticket prices | [decred/dcrwallet#232](https://github.com/decred/dcrwallet/pull/232) |
| Fix bugs relating to reorganizations | [decred/dcrwallet#236](https://github.com/decred/dcrwallet/pull/236) |
| Update glide locks | [decred/dcrwallet#239](https://github.com/decred/dcrwallet/pull/239) |
| Bump for v0.1.3 | [decred/dcrwallet#238](https://github.com/decred/dcrwallet/pull/238) |
| Update for new createrawssrtx option | [decred/dcrrpcclient#17](https://github.com/decred/dcrrpcclient/pull/17) |
| Correct the return type for estimatestakediff | [decred/dcrrpcclient#18](https://github.com/decred/dcrrpcclient/pull/18) |
| Fix functionality of purchaseticket API | [decred/dcrrpcclient#19](https://github.com/decred/dcrrpcclient/pull/19) |

## Notes

## Commits

This release was built from:

| Repository | Commit Hash |
| --- | ---- |
| decred/dcrd | 2aec09354a7263a31f6b5d3fe5906bc534613058 |
| decred/dcrwallet | 4215ccccceee037a7835721ca59a8c6327556f62 |
| decred/dcrrpcclient | e625cc131dc06129f56e0d472061c3e378ada396 |
| decred/dcrutil | 74563ea520b1215b9c10f96507b7a9984894c0b5 |

# [v0.1.2](https://github.com/decred/decred-release/releases/tag/v0.1.2)

## 2016-05-03

This release contains updated binary files (dcrd, dcrctl, dcrwallet)
for various platforms and is primarily a bugfix for dcrwallet.

See manifest-20160503-01.txt for sha256sums of the packages and
manifest-201600503-01.txt.asc to confirm those shas.

See https://wiki.decred.org/Verifying_Binaries for more info on
verifying the files.

This release contains various fixes and improvements.

Changes include:

| Description | Pull Request |
| --- | ---- |
| Fix mempool fees variables | [decred/dcrd#141](https://github.com/decred/dcrd/pull/141) |
| Add GetStakeDifficultyResult to dcrjson so getstakedifficulty command can return more | [decred/dcrd#137](https://github.com/decred/dcrd/pull/137) |
| Remove magic number and add const maxRelayFeeMultiplier | [decred/dcrd#139](https://github.com/decred/dcrd/pull/139) |
| Add estimatestakediff RPC command | [decred/dcrd#143](https://github.com/decred/dcrd/pull/143) |
| Add ticketvwap and txfeeinfo RPC server commands | [decred/dcrd#145](https://github.com/decred/dcrd/pull/145) |
| fix sample config per issue 116 | [decred/dcrd#147](https://github.com/decred/dcrd/pull/147) |
| Add stakepooluserinfo and addticket RPC handling | [decred/dcrd#144](https://github.com/decred/dcrd/pull/144) |
| Cherry pick btcd commit that moves some functions to policy.go | [decred/dcrd#140](https://github.com/decred/dcrd/pull/140) |
| Add the constructor for AddTicketCmd | [decred/dcrd#148](https://github.com/decred/dcrd/pull/148) |
| Bump for v0.1.2 | [decred/dcrd#150](https://github.com/decred/dcrd/pull/150) |
| Fix lockup relating to channel blocking | [decred/dcrwallet#219](https://github.com/decred/dcrwallet/pull/219) |
| Add stake pool mode to the wallet | [decred/dcrwallet#192](https://github.com/decred/dcrwallet/pull/192) |
| Make purchaseticket return the correct error | [decred/dcrwallet#224](https://github.com/decred/dcrwallet/pull/224) |
| Add wallet flag for allowhighfees | [decred/dcrwallet#214](https://github.com/decred/dcrwallet/pull/214) |
| Bump for v0.1.2 | [decred/dcrwallet#225](https://github.com/decred/dcrwallet/pull/225) |
| Add RPC client pass throughs for new daemon and wallet commands | [decred/dcrrpcclient#16](https://github.com/decred/dcrrpcclient/pull/16) |

## Notes

### Added stake pool fee functionality:

We have added new config flags for dcrwallet.  Let's go over each
option to make crystal clear its usage:

#### stakepoolcoldextkey

When this option is set it turns on stake pool functionality for
wallet.  When stake pool is enabled for wallet, there are a series of
transaction checks to verify whether this wallet will vote for a
ticket that has used this stake pool's address as the ticketaddress.

This option requires the extended public key of the stake pool's cold
wallet that will receive the pool's fees.  So on simnet for instance
this option looks like this:

```
--stakepoolcoldextkey=spubVWAdividNTiSM9SdLRA5JX6LYNwt58cd51TFnpnULGQ8oqNMNskfkQwU7rjWMCY7phBguVr4XTmAWyDVRKpo2dFyjFb6QG4ihB8w64UPNuu:1000
```

The first portion (spub..., or dpub... on mainnet) is the extended
public key and the second (1000) is the number of addresses to
derive. Every user of the pool gets their own cold fee wallet address
derived, so we recommend using at least 1000 in anticipation of the
relative number of users in the stake pool.

When a vote is created by the stake pool to vote on a ticket that has
been given voting rights, it pays the pool fee to the address derived
for the cold wallet from this extended public key.

#### pooladdress

This is for use by the stake pool user.  It will be an address
provided to the user by the stake pool.  If set, this address is used
during ticket purchase and will commit to a small output in the ticket
that gives the stake pool its required fees.


#### ticketaddress

Same as the old option. This is the address that the stake pool user
is giving the ticket's voting rights to.


#### poolfees

This is the required ticket fee as requested by the stake pool.  The
value set by the user needs to be greater than or equal to that of the
pool.  The fee is a percentage based fee, based on the stake subsidy.
Here is a concrete example from simnet:

The ticket price of this ticket was 46.0551008, and the ticket relay
fees were 0.00000100 per kB. The pool fees were set to 1.00%.  The
subsidy on simnet at this block height is approximately 29.40888 Coins
per vote.  This is the ticket as purchased by the user:

```javascript
"vin": [
	... ,
],
"vout": [
	{
		"value": 46.0551008,
		"n": 0,
		"version": 0,
		"scriptPubKey": {
			... ,
			"reqSigs": 1,
			"type": "stakesubmission",
			"addresses": [
				"SsYZMHeeixdNRTkk6afzHBPL4unYDsFNd4r"
			]
		}
	},
	{
		"value": 0,
		"n": 1,
		"version": 0,
		"scriptPubKey": {
			... ,
			"type": "sstxcommitment",
			"addresses": [
				"Ssghjx8PvQVV3FM3w5FcGi9kWGvDpDkQDTV"
			],
			"commitamt": 0.17948021
		}
	},
	{
		... ,
	},
	{
		"value": 0,
		"n": 3,
		"version": 0,
		"scriptPubKey": {
			... ,
			"type": "sstxcommitment",
			"addresses": [
				"SsYUi5tbXfqHnTPgvHcajNW4yiGeSP6n7Xq"
			],
			"commitamt": 45.87562609
		}
	},
	{
		... ,
	}
],
```

And here's the vote that the stake pool created for that user's
ticket:

```javascript
"vin": [
	{
		... ,
	},
	{
		... ,
	}
],
"vout": [
	{
		... ,
	},
	{
		... ,
	},
	{
	"value": 0.2940888,
	"n": 2,
	"version": 0,
		"scriptPubKey": {
			... ,
			"type": "stakegen",
			"addresses": [
				"Ssghjx8PvQVV3FM3w5FcGi9kWGvDpDkQDTV"
			]
		}
	},
	{
		"value": 75.16989347,
		"n": 3,
		"version": 0,
		"scriptPubKey": {
			... ,
			"type": "stakegen",
			"addresses": [
				"SsYUi5tbXfqHnTPgvHcajNW4yiGeSP6n7Xq"
			]
		}
	}
]
```

As you can see '"n": 2,', the third output, is the stake pool fee
of 0.2940888.  This is 1% of the vote reward at that point
(0.2940888/29.40888). The remaining subsidy and the original coins are
returned to the take pool user in output '"n": 3,'. For more
information about stake fees, please refer to
dcrwallet/wallet/txrules/doc.go.

## Commits

This release was built from:

| Repository | Commit Hash |
| --- | ---- |
| decred/dcrd | f93cb9fd9fd7b471481e4cfb122186514f84e879 |
| decred/dcrwallet | e545bec0a3a1a3b8380224d12c9ede85bff58595 |
| decred/dcrrpcclient | a5a51f5ca4f0038e475239cfe3c635a21fd28111 |
| decred/dcrutil | 74563ea520b1215b9c10f96507b7a9984894c0b5 |
| google.golang.org/grpc | b062a3c003c22bfef58fa99d689e6a892b408f9d |

# [v0.1.1](https://github.com/decred/decred-release/releases/tag/v0.1.1)

## 2016-04-25

This release contains updated binary files (dcrd, dcrctl, dcrwallet)
for various platforms and is primarily a bugfix for dcrwallet.

See manifest-20160425-01.txt for sha256sums of the packages and
manifest-201600425-01.txt.asc to confirm those shas.

See https://wiki.decred.org/Verifying_Binaries for more info on
verifying the files.

This release contains various fixes and improvements.

Changes include:

| Description | Pull Request |
| --- | ---- |
| Catch missed error check | [decred/dcrd#127](https://github.com/decred/dcrd/pull/127) |
| fix typo | [decred/dcrd#128](https://github.com/decred/dcrd/pull/128) |
| Replace float64 and use int64 for feePerKB calculcation | [decred/dcrd#125](https://github.com/decred/dcrd/pull/125) |
| Use AllowHighFees in SendRawTransaction cmd to actually check tx fees | [decred/dcrd#124](https://github.com/decred/dcrd/pull/124) |
| Add ticketfeeinfo command | [decred/dcrd#132](https://github.com/decred/dcrd/pull/132) |
| Bump for v0.1.1 | [decred/dcrd#136](https://github.com/decred/dcrd/pull/136) |
| Regenerate walletrpc package. | [decred/dcrwallet#189](https://github.com/decred/dcrwallet/pull/189) |
| Isolate address pool to prevent excessive address creation | [decred/dcrwallet#191](https://github.com/decred/dcrwallet/pull/191) |
| Reinsert scan length variable | [decred/dcrwallet#196](https://github.com/decred/dcrwallet/pull/196) |
| Do not include zero value change outputs. | [decred/dcrwallet#193](https://github.com/decred/dcrwallet/pull/193) |
| Update help comments to show fee per kb instead of increment | [decred/dcrwallet#195](https://github.com/decred/dcrwallet/pull/195) |
| Add TicketFeeIncrementTestnet | [decred/dcrwallet#194](https://github.com/decred/dcrwallet/pull/194) |
| Allow passing an empty string for purchaseticket addresses | [decred/dcrwallet#198](https://github.com/decred/dcrwallet/pull/198) |
| Add ability to change autopurchase frequency | [decred/dcrwallet#201](https://github.com/decred/dcrwallet/pull/201) |
| Open and return wallet from CreateNewWallet. | [decred/dcrwallet#203](https://github.com/decred/dcrwallet/pull/203) |
| Avoid stdin passphrase prompt with --noinitialload. | [decred/dcrwallet#202](https://github.com/decred/dcrwallet/pull/202) |
| Regenerate JSON-RPC help descriptions. | [decred/dcrwallet#208](https://github.com/decred/dcrwallet/pull/208) |
| Bump for v0.1.1 | [decred/dcrwallet#209](https://github.com/decred/dcrwallet/pull/209) |
| use decred mainnet ports in examples | [decred/dcrrpcclient#15](https://github.com/decred/dcrrpcclient/pull/15) |

## Commits

This release was built from:

| Repository | Commit Hash |
| --- | ---- |
| decred/dcrd | 4f8ad739a231a6ecf58ae899c595fba446ffe631 |
| decred/dcrwallet | c5e47fba1608854b0c43c367b14ced6df91a6d9e |
| decred/dcrrpcclient | c69fe513f9d6beeef0cad10412e3aa804ba3fe28 |
| decred/dcrutil | 74563ea520b1215b9c10f96507b7a9984894c0b5 |
| google.golang.org/grpc | 262ed2bd6d1c8cbaa14b43c3815d2e01e4f65ca8 |

# [v0.1.0](https://github.com/decred/decred-release/releases/tag/v0.1.0)

## 2016-04-18

This release contains updated binary files (dcrd, dcrctl, dcrwallet)
for various platforms and is primarily a bugfix for dcrwallet.

See manifest-20160418-01.txt for sha256sums of the packages and
manifest-201600418-01.txt.asc to confirm those shas.

See https://wiki.decred.org/Verifying_Binaries for more info on
verifying the files.

This release contains various fixes and improvements.

Changes include:

| Description | Pull Request |
| --- | ---- |
| Fix the constructors for new RPC account commands | [decred/dcrd#106](https://github.com/decred/dcrd/pull/106) |
| TravisCI: Remove external go vet reference. (#655)  | [decred/dcrd#107](https://github.com/decred/dcrd/pull/107) |
| Clean up and fix fallthrough on invalid tx types for getrawmempool rpc request | [decred/dcrd#11](https://github.com/decred/dcrd/pull/111) |
| Pull in policy.go changes from btcd to fix issues with fee calc in dcrd | [decred/dcrd#112](https://github.com/decred/dcrd/pull/112) |
| Updated config to allow the ability to change the home directory | [decred/dcrd#109](https://github.com/decred/dcrd/pull/109) |
| Fix the mining transaction selection algorithm | [decred/dcrd#113](https://github.com/decred/dcrd/pull/113) |
| Fix rpclisten and listen port references in documentation | [decred/dcrd#118](https://github.com/decred/dcrd/pull/118) |
| Properly handle attempting reorganization to an eligible block | [decred/dcrd#117](https://github.com/decred/dcrd/pull/117) |
| Display ticket commitments in getrawtransaction | [decred/dcrd#119](https://github.com/decred/dcrd/pull/119) |
| Check to see if missingParents != nil which means isOrphan | [decred/dcrd#122](https://github.com/decred/dcrd/pull/122) |
| Modify the purchaseticket RPC command | [decred/dcrd#121](https://github.com/decred/dcrd/pull/121) |
| Bump for v0.1.0 | [decred/dcrd#123](https://github.com/decred/dcrd/pull/123) |
| Update TravisCI configs. (#409) | [decred/dcrwallet#168](https://github.com/decred/dcrwallet/pull/168) |
| Fix a bug causes wallet lockup when making transactions | [decred/dcrwallet#167](https://github.com/decred/dcrwallet/pull/167) |
| Add sweepaccount tool.  | [decred/dcrwallet#173](https://github.com/decred/dcrwallet/pull/173) |
| Add .ToCoin() to GetWalletFee return val to be consistent | [decred/dcrwallet#172](https://github.com/decred/dcrwallet/pull/172) |
| Fix bug in syncing to address index | [decred/dcrwallet#176](https://github.com/decred/dcrwallet/pull/176) |
| fix off by one when initializing a wallet | [decred/dcrwallet#177](https://github.com/decred/dcrwallet/pull/177) |
| Clean UX so it is more clear that a pass is required | [decred/dcrwallet#180](https://github.com/decred/dcrwallet/pull/180) |
| Change default relay fees | [decred/dcrwallet#182](https://github.com/decred/dcrwallet/pull/182) |
| Ticket purchasing code overhaul | [decred/dcrwallet#183](https://github.com/decred/dcrwallet/pull/183) |
| Refactor address index syncing code | [decred/dcrwallet#184](https://github.com/decred/dcrwallet/pull/184) |
| Bump for v0.1.0 | [decred/dcrwallet#185](https://github.com/decred/dcrwallet/pull/185) |
| TravisCI: Update to latest configurations. (#76) | [decred/dcrrpcclient#13](https://github.com/decred/dcrrpcclient/pull/13) |
| Add client handling for new RPC calls | [decred/dcrrpcclient#12](https://github.com/decred/dcrrpcclient/pull/12) |
| Fix the purchaseticket caller | [decred/dcrrpcclient#14](https://github.com/decred/dcrrpcclient/pull/14) |
| TravisCI: Remove external go vet reference. (#74) | [decred/dcrutil#9](https://github.com/decred/dcrutil/pull/9) |

## Commits

This release was built from:

| Repository | Commit Hash |
| --- | ---- |
| decred/dcrd | a33985293b19aab047f95d1f68d07d9625811d6d |
| decred/dcrwallet | b192834577b44602f8960bca3dcf9d35af32acb7 |
| decred/dcrrpcclient | f005c4a9466229520d7198ce1904065248f6cdd3 |
| decred/dcrutil | 74563ea520b1215b9c10f96507b7a9984894c0b5 |
| google.golang.org/grpc | 326d66361a4e305b03da4497d2c52d470f7fb584 |

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
| Set Tree field when converting wire.OutPoints. | [decred/dcrrpcclient#10](https://github.com/decred/dcrrpcclient/pull/10) |

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
