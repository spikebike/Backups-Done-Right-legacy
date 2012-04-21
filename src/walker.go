package main

import "C"

import (
	"database/sql"
	"flag"
	"log"
	"os"
	"strings"
//	"syscall"
	"time"

	"./bdrsql"

	"github.com/kless/goconfig/config"
	_ "github.com/mattn/go-sqlite3"
)

var (
	sqls = []string{
		"create table dirs (id INTEGER PRIMARY KEY, mode INT, ino BIGINT, uid INT, gid INT, path varchar(2048), last_seen BIGINT, deleted INT)",
		"create table files (id INTEGER PRIMARY KEY, mode INT, ino BIGINT, dev BIGINT, uid INT, gid INT, size BIGINT, atime BIGINT, mtime BIGINT, ctime BIGINT, name varchar(255), dirID BIGINT, last_seen BIGINT, deleted INT, do_upload INT, FOREIGN KEY(dirID) REFERENCES dirs(id))",
	}

	configFile = flag.String("config", "../etc/config.cfg", "Defines where to load configuration from")
	newDB      = flag.Bool("new-db", false, "true = creates a new database | false = use existing database")
	debug      = flag.Bool("debug", false, "activates debug mode")
)

func backupDir(db *sql.DB, dirList string, bufsize int) error {
	var i int
	var dirname string

//	entry := &file_info_t{}
	i = 0
//	start:=time.Now().Unix()
	log.Printf("backupDir received %s", dirList)
	dirArray := strings.Split(dirList, " ")
	for i < len(dirArray) {
		dirname = dirArray[i]
		// does dirname exist in the dir table
		dirID,err := bdrsql.GetSQLID(db,"dirs","path",dirname)
//      add code here
//		sqlFiles = "select FI from files,directories where dir.id=$dirID and files.deleted=false
		Fmap:=bdrsql.GetSQLFiles(db,1)
		log.Printf("map=%T map=%#v dirID=%d\n",Fmap,Fmap,dirID)
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
			// add code here
			// return zero if fi.Name not in sqlfiles
			// sqlModTime = getModTime(sqlFiles,fi.Name())
//			unixStat, _ := fi.Sys().(*syscall.Stat_t)
			// if it's a file not a dir
			if !fi.IsDir() {
				log.Printf("found %s file\n",fi.Name())
				// and it's been modified since last backup
//				if fi.Modfile<=sqlModTime { //already backed up
//					"update files set Last_seen=now() where name=fi.Name and deted=fales"
//			} else { // Either fi is newer or modified
//					makeSQLEntry(db, fi)
//			} else {
//				Fullpath:=dirname+"/"+fi.Name()
//				if Fullpath not in dirArray { // directory needs walked
//					dirArray = append(dirArray, dirname+"/"+fi.Name())
//				}
			}
		}
		i++
	}
	// once done walking any file not seen must have been deleted.
//	Update files set deleted=True where last_seen < $start;
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

	db, err := bdrsql.Init_db(dataBaseName)

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
