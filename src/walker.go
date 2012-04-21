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

func backupDir(db *sql.DB, dirList string) error {
	var i int
	var dirname string

	//	entry := &file_info_t{}
	i = 0
	//	start:=time.Now().Unix()
	log.Printf("backupDir received %s", dirList)
	dirArray := strings.Split(dirList, " ")
	for i < len(dirArray) {
		dirname = dirArray[i]
		// get dirID of dirname, even if it needs inserted.
		dirID, err := bdrsql.GetSQLID(db, "dirs", "path", dirname)
		// get a map for filename -> modified time
		Fmap := bdrsql.GetSQLFiles(db, dirID)
		log.Printf("map=%T map=%#v dirID=%d\n", Fmap, Fmap, dirID)
		log.Printf("backing up dir %s id=", dirname, dirID)
		d, err := os.Open(dirname)
		if err != nil {
			log.Printf("failed to open %s error : %s", dirname, err)
			os.Exit(1)
		}
		fi, err := d.Readdir(-1)
		if err != nil {
			log.Printf("directory %s failed with error %s", dirname, err)
		}

		for _, f := range fi {
			if !f.IsDir() {
				// and it's been modified since last backup
				if f.ModTime().Unix() <= Fmap[f.Name()] {
					log.Printf("NO backup needed for %s \n",f.Name())
				} else {
					log.Printf("backup needed for %s \n",f.Name())
					bdrsql.InsertSQLFile(db,f,dirID)
				}
					 //already backed up
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

	t0 := time.Now()
	log.Printf("start walking...")
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
