# Backups-Done-Right

Backups-Done-Right is an easy way to do fast, simple and secure backups.


## Features

* file transfer is always encrypted
* very fast file system walker
* simple configuration - just one config file
* simple installation (static linked build)
* open source


## Maintainers

Backups-Done-Right does have two project maintainers:

* Bill Broaldey   aka spikebike	<bill@broadley.org>
* Joel Bodenmann  aka Tectu	<joel@unormal.org>


## Build

Backups Done Right depends on 3 external go packages that need installed:

goconfig - to install, simply run:

	$ go get github.com/kless/goconfig/gonfig


goconfig - to install, simply run:

	$ go get github.com/mattn/go-sqlite3


go-rpcgen - to install, simply run:

	$ go get github.com/kylelemons/go-rpcgen/protoc-gen-go


Then to make certificates run:

	$ cd Backups-Done-Right/src
	$ ./makecerts <your email address>


Please note that Backups-Done-Right does also require sqlite3 which is not a part of go. 


## Run

Before you can run the software the first time, you need to create
a config file which fits your needs:

	$ cd Backups-Done-Right/etc
	$ cp config.cfg.example config.cfg

Then, edit the config file to your needs.

To run


## Misc

Please see documentation/* for feature informations

