package main

import "C"

import (
	"database/sql"
	"flag"
	"github.com/kless/goconfig/config"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"strings"
)

var (
	sqls = []string {
		"create table dirs (id INTEGER PRIMARY KEY,st_mode INT, st_ino BIGINT, st_uid INT, st_gid INT, name varchar(2048), last_seen ts, deleted INT)",
		"create table files (id INTEGER PRIMARY KEY, st_mode INT, st_ino BIGINT, st_dev BIGINT, st_nlink INT, st_uid INT, st_gid INT, st_size BIGINT, st_atime BIGINT, st_mtime BIGINT, st_ctime BIGINT, name varchar(255), dirID BIGINT, last_seen ts, deleted INT, FOREIGN KEY(dirID) REFERENCES dirs(id))",
	}

	configFile = flag.String("config", "../etc/config.cfg", "Defines where to load configuration from")
)

func init_db(dataBaseName string) (db *sql.DB, err error) {
	db, err = sql.Open("sqlite3", dataBaseName)
	if err != nil {
		log.Printf("Couldn't open database: %s", err)
		os.Exit(1)
	}
	_, rerr := db.Exec(sqls[0])
	if rerr != nil {
		log.Printf("%s", rerr)
	}
	_, rerr = db.Exec(sqls[1])
	if rerr != nil {
		log.Printf("%s", rerr)
	}

	return db, err
}

func backupDir(db *sql.DB, dirList string) error {
	var i int
	var dirname string
	i = 0
	log.Printf("backupDir received %s", dirList)
	dirArray := strings.Split(dirList," ")
	for i < len(dirArray) {
		dirname = dirArray[i]
		log.Printf("backing up dir %s", dirname)
		d, err := os.Open(dirname)
		if err != nil {
			log.Printf("Failed to open %s error=%s", dirname, err)
			os.Exit(1)
		}
		fi, err := d.Readdir(-1)
		if err != nil {
			log.Printf("Directory %s failed with error %s", dirname, err)
		}
		for _, fi := range fi {
			if !fi.IsDir() {
				log.Printf("%s %d bytes %s", fi.Name(), fi.Size(),fi.ModTime())

			} else {
				dirArray = append(dirArray, dirname+"/"+fi.Name())
				log.Printf("found directory %s", fi.Name())
//				log.Println(os.Stat(fi))
			}
		}
		i++
	}
	return nil
}

func main() {
	flag.Parse()

	log.Printf("Loading config file from %s\n", *configFile)
	config, _ := config.ReadDefault(*configFile)

	log.Printf("Backing up these directories: %s\n", dirList)
	dirList, _ := config.String("Client", "backup_dirs_secure")

	log.Printf("Attempting to open %s", dataBaseName)
	dataBaseName, _ := config.String("Client", "sql_file")

	db, err := init_db(dataBaseName)
	err = backupDir(db, dirList)
	log.Printf("backupDir exited with %s", err)
}
