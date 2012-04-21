package main

import "C"

import (
	"database/sql"
	"flag"
	"log"
	"os"
	"strings"
	//	"syscall"
	"fmt"
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

func checkPath(dirArray []string, dir string) bool {
	for _, i := range dirArray {
		if i == dir {
			return true
		}
	}
	return false
}

func backupDir(db *sql.DB, dirList string) error {
	var dirname string
	var i int
	var fileC int64
	var dirC int64
	var dFile int64
	var dDir int64
	fileC = 0
	dirC = 0
	dFile = 0
	dDir = 0
	start := time.Now().Unix()
	dirArray := strings.Split(dirList, " ")
	i = 0
	for i < len(dirArray) {
		dirname = dirArray[i]
		// get dirID of dirname, even if it needs inserted.
		dirID, err := bdrsql.GetSQLID(db, "dirs", "path", dirname)
		// get a map for filename -> modified time
		SQLmap := bdrsql.GetSQLFiles(db, dirID)
		fmt.Printf("Scanning dir %s ", dirname)
		d, err := os.Open(dirname)
		if err != nil {
			log.Printf("failed to open %s error : %s", dirname, err)
			os.Exit(1)
		}
		fi, err := d.Readdir(-1)
		if err != nil {
			log.Printf("directory %s failed with error %s", dirname, err)
		}
		Fmap := map[string]int64{}
		// Iterate over the entire directory
		dFile = 0
		dDir = 0
		for _, f := range fi {
			if !f.IsDir() {
				fileC++
				dFile++
				// and it's been modified since last backup
				if f.ModTime().Unix() <= SQLmap[f.Name()] {
					//					log.Printf("NO backup needed for %s \n",f.Name())
					Fmap[f.Name()] = f.ModTime().Unix()
				} else {
					//					log.Printf("backup needed for %s \n",f.Name())
					bdrsql.InsertSQLFile(db, f, dirID)
				}
			} else { // is directory
				dirC++
				dDir++
				fullpath := dirname + "/" + f.Name()
				// avoid an infinite loop 
				if !checkPath(dirArray, fullpath) {
					dirArray = append(dirArray, fullpath)
				}
			}
		}
		// All files that we've seen, set last_seen
		t1 := time.Now().UnixNano()
		bdrsql.SetSQLSeen(db, Fmap, dirID)
		t2 := time.Now().UnixNano()
		fmt.Printf("files=%d dirs=%d duration=%dms\n", dFile, dDir, (t2-t1)/1000000)
		i++
	}
	// if we have seen the files since start it must have been deleted.
	bdrsql.SetSQLDeleted(db, start)
	log.Printf("fileC=%d dirC=%d\n", fileC, dirC)
	return nil
}

func main() {
	flag.Parse()

	log.Printf("loading config file from %s\n", *configFile)
	configF, err := config.ReadDefault("../etc/config.cfg")
	if err != nil {
		log.Fatalf("ERROR: %s", err)
	}
	dirList, err := configF.String("Client", "backup_dirs_secure")
	if err != nil {
		log.Fatalf("ERROR: %s", err)
	}
	log.Printf("backing up these directories: %s\n", dirList)

	dataBaseName, err := configF.String("Client", "sql_file")
	if err != nil {
		log.Fatalf("ERROR: %s", err)
	}
	log.Printf("attempting to open %s", dataBaseName)

	db, err := bdrsql.Init_db(dataBaseName)
	if err != nil {
		log.Printf("could not open %s, error: %s", dataBaseName, err)
	} else {
		log.Printf("Opened database %v\n", db)
	}

	log.Printf("start walking...")
	t0 := time.Now()
	err = backupDir(db, dirList)
	t1 := time.Now()
	duration := t1.Sub(t0)

	if err != nil {
		log.Printf("Walking didn't finished successfully. Error: %s", err)
	} else {
		log.Printf("walking successfully finished")
	}

	log.Printf("walking took: %v\n", duration)
}
