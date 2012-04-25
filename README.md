# Backups-Done-Right

Backups-Done-Right is an easy way to do fast, simple and secure backups.


## Features

* file transfer is always encrypted
* simple installation (static linked build)
* open source


## Maintainers

Backups-Done-Right does have two project maintainers:

* Bill Broaldey   aka spikebike 	<bill@broadley.org>
* Joel Bodenmann  aka Tectu       <joel@unormal.org>


## Build & Run

Backups Done Right depends on 3 external packages that need installed:

goconfig - to install, simply run:

	$ go get github.com/kless/goconfig/config'


goconfig - to install, simply run:

	$ go get github.com/mattn/go-sqlite3'


go-rpcgen - to install, simply run:

	$ go get github.com/kylelemons/go-rpcgen/protoc-gen-go'


Then to make certificates run:

	$ cd Backups-Done-Right/src
	$ ./makecerts <your email address>


Before you can run the software the first time, you need to create
a config file which fits your needs:

	$ cd Backups-Done-Right/etc
	$ cp config.cfg.example config.cf

Then, edit the config file to your needs.


## Misc

Please see documentation/* for feature informations

