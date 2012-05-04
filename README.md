# Backups-Done-Right

Backups-Done-Right is a P2P backup program providing easy, fast and secure encrypted off-site backups.


## Features

* files encrypted with AES256 before they go over the network
* Server does not have to have the clients encryption key
* fast filesystem walker
* posibility to run more than one fs walker at once (performance scales with number of spindles)
* simple configuration - just one config file
* simple installation (static linked build)
* restores with permissions, symlinks etc.
* open source - be sure that nobody gets your data
* backups are stored encrypted - you don't need to trust in your sysadmin
* encryption key does not get transfered - you always keep the key on your side

## Potential Misfeatures
* No back doors - without the clients AES256 key no recovery is possible
* No bare metal restores - Backups Does Right depends on Go + working unix 
          system, and AES256 key for a restore.  For connections to peers 
          one of the old IP address or the IP address of peers is needed.
* No cross client file restores - Full restore of a linux backup requires
          a linux client.  Same with Windows, and OSX.  No
          attempts are made to handle platform specific file
          system meta-data on non-native platforms.
* Owner of SHA256 key *MUST* make copies to protect against lost.  Multiple
           Printouts and thumbdrives are recommended to insure recovery.
                        

## Maintainers

Backups-Done-Right does have two project maintainers:

* Bill Broaldey   aka spikebike	<bill@broadley.org>	(english)
* Joel Bodenmann  aka Tectu	<joel@unormal.org>	(german / english)


## Build

Backups-Done-Right does depend on some go packages, which in turn depend on some software packages.  The below statement should install the dependencies for most Debian/Ubuntu based systems:


	$ apt-get install libsqlite3-dev libsqlite-dev git mercurial pkg-config

Backups Done Right depends on 3 external go packages that need to be installed.  Set GOPKG to where you want them installed.  Something like export GOPKG=/home/JoeUser/gopkg.  DO *NOT* use ~/gopkg.  To install the dependencies:

goconfig - to install, simply run:

	$ go get github.com/kless/goconfig/config


go-sqlite3 - This requires the sqlite and sql-dev packages to be installed already.  To install, simply run:

	$ go get github.com/mattn/go-sqlite3


go-rpcgen - to install, simply run:

	$ go get github.com/kylelemons/go-rpcgen/protoc-gen-go


## Run

Before you can run the software the first time, you need to create a config file which fits your needs. Please copy the example config file:

	$ cd Backups-Done-Right/etc
	$ cp config.cfg.example config.cfg

Then, edit the config file to your needs.

You do also need certificates for the SSL encryption:

	$ cd Backups-Done-Right/src/scripts
	$ ./makecerts <your_email_address>


## Technical Description

Once the filesystem walker created a database of the directories that have to be backed up, it will just update the database on every run. On each run the walker decides if the file got any changes. If yes, the file gets encrypted over AES-512 and gets uploaded to the backup server over an SSL secured TCP/IP connection. The server keeps the files encrypted.
Whenever we need a backup, we send the encrypted checksum (SHA-512) of the file, which is also stored in the database to the server. The server will send the encrypted file to the matching client over an SSL secured TCP/IP connection again. The client will then decrypt the received file and restores the complete directory tree with all the permissions, symlinks etc.


## Misc

Please see documentation/* for additional informations

