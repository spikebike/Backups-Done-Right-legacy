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
	"syscall"
	"time"
)

var (
	sqls = []string{
		"create table dirs (id INTEGER PRIMARY KEY, mode INT, ino BIGINT, uid INT, gid INT, path varchar(2048), last_seen BIGINT, deleted INT)",
		"create table files (id INTEGER PRIMARY KEY, mode INT, ino BIGINT, dev BIGINT, uid INT, gid INT, size BIGINT, atime BIGINT, mtime BIGINT, ctime BIGINT, name varchar(255), dirID BIGINT, last_seen BIGINT, deleted INT, do_upload INT, FOREIGN KEY(dirID) REFERENCES dirs(id))",
	}

	configFile = flag.String("config", "../etc/config.cfg", "Defines where to load configuration from")
	newDB      = flag.Bool("new-db", false, "true = creates a new database | false = use existing database")
	debug	   = flag.Bool("debug", false, "activates debug mode")
)

type file_info_t struct {
	id        int64
	mode      uint32
	ino       uint64
	dev       uint64
	nlink     int64
	uid       uint32
	gid       uint32
	size      int64
	atime     int64
	mtime     int64
	ctime     int64
	name      string
	path      string
	dirID     int
	last_seen int64
	deleted   int
	do_upload int
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

func backupDir(db *sql.DB, dirList string, bufsize int) error {
	var i int
	var dirname string
	entry := &file_info_t{}
	i = 0

	log.Printf("backupDir received %s", dirList)
	dirArray := strings.Split(dirList, " ")
	for i < len(dirArray) {
		dirname = dirArray[i]
		if *debug == true {
			log.Printf("backing up dir %s", dirname)
		}
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
				entry.deleted = 0
				entry.do_upload = 0
				entry.name = fi.Name()
				entry.size = fi.Size()
				entry.gid = unixStat.Gid
				entry.uid = unixStat.Uid
				entry.ino = unixStat.Ino
				entry.dev = unixStat.Dev
				entry.mode = unixStat.Mode
				entry.mtime = unixStat.Mtim.Sec
				entry.atime = unixStat.Atim.Sec
				entry.ctime = unixStat.Ctim.Sec
				entry.last_seen = time.Now().Unix()
			} else {
				dirArray = append(dirArray, dirname+"/"+fi.Name())
				entry.size = 0 // VERY IMPORTANT!
				entry.deleted = 0
				entry.gid = unixStat.Gid
				entry.uid = unixStat.Uid
				entry.ino = unixStat.Ino
				entry.mode = unixStat.Mode
				entry.last_seen = time.Now().Unix()
				entry.path = dirname + "/" + fi.Name()
			}

			queue_file(db, fi, bufsize)
			makeEntry(db, entry)
		}
		i++
	}
	return nil
}

func queue_file(db *sql.DB, fi os.FileInfo, bufsize int) error {

	return nil
}

func makeEntry(db *sql.DB, e *file_info_t) error {
	tx, err := db.Begin()
	if err != nil {
		log.Println(err)
	}

	if e.size != 0 { // is it a file or a dir entry?
		stmt, err := db.Prepare("select id from dirs where path = ?")
		if err != nil {
			log.Println(err)
			return err
		}
		defer stmt.Close()
		err = stmt.QueryRow(e.path).Scan(&e.dirID)
		if err != nil {
			log.Println(err)
			return err
		}

		stmt, err = tx.Prepare("insert into files(name,size,mode,gid,uid,ino,dev,mtime,atime,ctime,last_seen,dirID,deleted,do_upload) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
		if err != nil {
			log.Println(err)
			return err
		}
		defer stmt.Close()

		_, err = stmt.Exec(e.name, e.size, e.mode, e.gid, e.uid, e.ino, e.dev, e.mtime, e.atime, e.ctime, e.last_seen, e.dirID, e.deleted, e.do_upload)
		if err != nil {
			log.Println(err)
			return err
		}
	} else {
		stmt, err := tx.Prepare("insert into dirs(path,mode,gid,uid,ino,last_seen,deleted) values(?,?,?,?,?,?,?)")
		if err != nil {
			log.Println(err)
			return err
		}
		defer stmt.Close()

		_, err = stmt.Exec(e.path, e.mode, e.gid, e.uid, e.ino, e.last_seen, e.deleted)
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

	bufsize, _ := config.Int("Client", "buffer_size")

	t0 := time.Now()
	log.Printf("start walking...")
	err = backupDir(db, dirList, bufsize)
	t1 := time.Now()
	duration := t1.Sub(t0)

	if err != nil {
		log.Printf("Walking didn't finished successfully. Error: ", err)
	} else {
		log.Printf("walking successfully finished")
	}

	log.Printf("walking took: %v\n", duration)
}
