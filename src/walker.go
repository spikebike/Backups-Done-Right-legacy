package main

import "C"

import (
	"os"
	"log"
	"flag"
	"time"
	"strings"
	"syscall"
	"database/sql"
	"github.com/kless/goconfig/config"
	_ "github.com/mattn/go-sqlite3"
)

var (
	sqls = []string {
		"create table dirs (id INTEGER PRIMARY KEY, mode INT, ino BIGINT, uid INT, gid INT, name varchar(2048), last_seen ts, deleted INT)",
		"create table files (id INTEGER PRIMARY KEY, mode INT, ino BIGINT, dev BIGINT, nlink INT, uid INT, gid INT, size BIGINT, atime BIGINT, mtime BIGINT, ctime BIGINT, name varchar(255), dirID BIGINT, last_seen ts, deleted INT, FOREIGN KEY(dirID) REFERENCES dirs(id))",
	}

	configFile = flag.String("config", "../etc/config.cfg", "Defines where to load configuration from")
	newDB = flag.Bool("new-database", false, "true = creates a new database | false = use existing database")
)

type file_info_t struct {
	id int64
	mode int64
	ino int64
	dev int64
	nlink int64
	uid int64
	gid int64
	size int64
	atime int64
	mtime int64
	ctime int64
	name string
	dirID int64
	last_seen int64
	deleted int
}

func init_db(dataBaseName string) (db *sql.DB, err error) {
	if *newDB == true {
		os.Remove(dataBaseName)
	}

	db, err = sql.Open("sqlite3", dataBaseName)
	if err != nil {
		log.Printf("couldn't open database: %s", err)
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
			log.Printf("failed to open %s error : %s", dirname, err)
			os.Exit(1)
		}
		fi, err := d.Readdir(-1)
		if err != nil {
			log.Printf("directory %s failed with error %s", dirname, err)
		}
		for _, fi := range fi {
			unixStat, _ := fi.Sys().(*syscall.Stat_t)
			if !fi.IsDir() {
				entry.name = fi.Name()
				entry.size = fi.Size()
				entry.mtime = unixStat.Mtim.Sec
				entry.atime = unixStat.Atim.Sec
				entry.ctime = unixStat.Ctim.Sec
			} else {
				dirArray = append(dirArray, dirname+"/"+fi.Name())
				entry.size = 0	// VERY IMPORTANT! 
				entry.name = fi.Name()
			}

			makeEntry(db, entry)
		}
		i++
	}
	return nil
}

func makeEntry(db *sql.DB, entry *file_info_t) error {
	tx, err := db.Begin()
	if err != nil {
		log.Println(err)
	}

	if entry.size != 0 {	// is it a file or a dir entry?
		stmt, err := tx.Prepare("insert into files(name, size, mtime, atime, ctime) values(?,?,?,?,?)")
		if err != nil {
			log.Println(err)
			return err
		}
		defer stmt.Close()

		_, err = stmt.Exec(entry.name, entry.size, entry.mtime, entry.atime, entry.ctime)
		if err != nil {
			log.Println(err)
			return err
		}

	} else {
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
	}

	tx.Commit()

	return nil
}

func main() {
	flag.Parse()

	log.Printf("loading config file from %s\n", *configFile)
	config, _ := config.ReadDefault(*configFile)

	dirList, _ := config.String("Client", "backup_dirs_secure")
	log.Printf("backing up these directories: %s\n", dirList)

	dataBaseName, _ := config.String("Client", "sql_file")
	log.Printf("attempting to open %s", dataBaseName)

	db, err := init_db(dataBaseName)

	t0 := time.Now()
	log.Printf("start walking...")
	err = backupDir(db, dirList)
	t1 := time.Now()
	duration := t1.Sub(t0)

	if err != nil {
		log.Printf("Walking didn't finished successfully. Error: ", err)
	} else {
		log.Printf("walking successfully finished")
	}

	log.Printf("walking took: %v\n", duration)
}
