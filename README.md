# Backups-Done-Right

Backups-Done-Right is a P2P backup program providing easy, fast and secure encrypted off-site backups.


## Features

* file transfer is always encrypted
* very fast file system walker
* posibility to run more than one fs walker at once (huge speed increasing on multi HDD systems)
* simple configuration - just one config file
* simple installation (static linked build)
* open source


## Maintainers

Backups-Done-Right does have two project maintainers:

* Bill Broaldey   aka spikebike	<bill@broadley.org>	(english)
* Joel Bodenmann  aka Tectu	<joel@unormal.org>	(german / english)


## Build

Backups Done Right depends on 3 external go packages that need to be installed:

goconfig - to install, simply run:

	$ go get github.com/kless/goconfig/gonfig


goconfig - to install, simply run:

	$ go get github.com/mattn/go-sqlite3


go-rpcgen - to install, simply run:

	$ go get github.com/kylelemons/go-rpcgen/protoc-gen-go


Please note that Backups-Done-Right does also require sqlite3 which is not a part of go or Backups-Done-Right itself. 
Installing sqlite3 depends on your system. Therefor we cannot give you an installing howto. If you don't know how to install sqlite3 on your system, we recommend to use google to find out.


## Run

Before you can run the software the first time, you need to create a config file which fits your needs. Please copy the example config file:

	$ cd Backups-Done-Right/etc
	$ cp config.cfg.example config.cfg

Then, edit the config file to your needs.

You do also need certificates for the encryption:

	$ cd Backups-Done-Right/src
	$ ./makecerts <your_email_address>


## Misc

Please see documentation/* for additional informations

