# Backups-Done-Right

Backups-Done-Right is a P2P backup program providing easy, fast and secure encrypted off-site backups.

## Features

* files encrypted with AES256 before upload
* Server does not have to have the clients encryption key
* fast filesystem walker
* possibility to run more than one fs walker at once (performance scales with number of spindles)
* simple configuration - just one config file
* simple installation (just go git)
* restores with permissions, symlinks etc.

## Potential Misfeatures

* No back doors - without the clients encryption key
* No bare metal restores - Backups Does Right depends on Go + working OS, 
  and AES256 key for a restore.  For connections to peers one of (old IP 
  address or IP address of peers) is needed.
* Some metadata is lost on cross platform restores.  Go's stat implementation
  is used.  Backups done right is only as cross platform 
  compatible as go.
* Owner of AES256 key *MUST* make copies to protect against lost.  Multiple
  Printouts and/or thumbdrives in secure locations (not in the same building) 
  are recommended to insure recovery.
                        

## Maintainers

Backups-Done-Right does have two project maintainers:

* Bill Broadley   aka spikebike	<bill@broadley.org>	(english)
* Joel Bodenmann  aka Tectu	<joel@unormal.org>	(german / english)

## Technical Overview

The walker tracks filesystem metadata in sqlite.  During each walker
invocation all configured directories are walked looking for new or changed
files.  Any new or changed files are encrypted (AES-256) and checksummed
(SHA256).  These resulting blobs are then queued for upload.

Then the walker connects to the server.  A TLS connection is opened and the
server has to prove ownership of the correct private key.  Then a list of the
updated blobs are send and the server responds with which blobs it doesn't
already have.

Server <-> server connections are used to replicated blobs to the configured
redundancy.  Only known public keys are trusted and all communications happen
over an SSL connection.

Peers are periodically challenged to prove they have the blobs they claim to
be storing.  This also allows for disaster recovery (assuming the AES256 key
is available) by waiting for challenges from your peers.

## Install and use

See doc/INSTALL

## Misc

Please see doc/* for additional informations








 
 
 
 
 
