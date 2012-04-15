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

type file_info_t struct {
	id int
	st_mode int
	st_ino int
	st_dev int
	st_nlink int
	st_uid int
	st_gid int
	st_size int64
	st_atime int
	st_mtime int
	st_ctime int
	name string
	dirID int
	last_seen int
	deleted int
}

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
	entry := &file_info_t{}

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
//				log.Printf("%s %d bytes %s", fi.Name(), fi.Size(),fi.ModTime())
				entry.name = fi.Name()
				entry.st_size = fi.Size()
				makeFileEntry(db, entry)
			} else {
				dirArray = append(dirArray, dirname+"/"+fi.Name())
				entry.name = fi.Name()
				makeDirEntry(db, entry)
//				log.Printf("found directory %s", fi.Name())
//				log.Println(os.Stat(fi))
			}
		}
		i++
	}
	return nil
}

func makeDirEntry(db *sql.DB, entry *file_info_t) error {
	tx, err := db.Begin()
	if err != nil {
		log.Println(err)
	}
	stmt, err := tx.Prepare("insert into dirs(name) values(?)")
	if err != nil {
		log.Println(err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(entry.name)
	if err != nil {
		log.Println(err)
		return err
	}
	tx.Commit()

	return nil
}

func makeFileEntry(db *sql.DB, entry *file_info_t) error {
	tx, err := db.Begin()
	if err != nil {
		log.Println(err)
	}
	stmt, err := tx.Prepare("insert into files(name, st_size) values(?, ?)")
	if err != nil {
		log.Println(err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(entry.name, entry.st_size)
	if err != nil {
		log.Println(err)
		return err
	}
	tx.Commit()

	return nil
}

func main() {
	flag.Parse()

	log.Printf("Loading config file from %s\n", *configFile)
	config, _ := config.ReadDefault(*configFile)

	dirList, _ := config.String("Client", "backup_dirs_secure")
	log.Printf("Backing up these directories: %s\n", dirList)

	dataBaseName, _ := config.String("Client", "sql_file")
	log.Printf("Attempting to open %s", dataBaseName)

	db, err := init_db(dataBaseName)
	err = backupDir(db, dirList)
	log.Printf("backupDir exited with %s", err)
}
