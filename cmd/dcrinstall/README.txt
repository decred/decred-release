The dcrinstall tool is a tool to install and upgrade decred
installations.  It is designed to automate the most painful setup steps.
Once a machine meets the directory structure it can be perpetually be
upgraded using the tool.

The following steps are required to upgrade a currently hand installed
system to the dcrinstall expected layout.

The dcrinstall tool expects the following directory layout.  In order to
upgrade copy your current configuration files into the correct location
and ensure everything still works.

If dcrinstall detects all configuration files it'll operate in upgrade
mode.  Upgrade mode only overwrites the binaries in %HOME%\decred (or
~/decred on UNIX type OS').

The dcrinstall tool records all actions in %HOME%\decred\dcrinstall.log
(or ~/decred/dcrinstall.log on UNIX type OS').

== Windows ==

configuration files
%HOME%\AppDate\Local\Dcrctl\dcrctl.conf
%HOME%\AppDate\Local\Dcrd\dcrd.conf
%HOME%\AppDate\Local\Dcrticketbuyer\ticketbuyer.conf
%HOME%\AppDate\Local\Dcrwallet\dcrwallet.conf

binaries directory
%HOME%\decred\

== OSX ==

configuration files
~/Library/Application Support/Dcrctl/dcrctl.conf
~/Library/Application Support/Dcrd/dcrd.conf
~/Library/Application Support/Dcrticketbuyer/ticketbuyer.conf
~/Library/Application Support/Dcrwallet/dcrwallet.conf

binaries directory
~/decred

== UNIX ==

configuration files
~/.dcrctl/dcrctl.conf
~/.dcrd/dcrd.conf
~/.dcrticketbuyer/ticketbuyer.conf
~/.dcrwallet/dcrwallet.conf

binaries directory
~/decred
