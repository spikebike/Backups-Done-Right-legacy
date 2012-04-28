package main

import "C"

import (
	"database/sql"
	"flag"
	"log"
	"os"
	"strings"
	"fmt"
	"time"
	"runtime"
	"./bdrsql"
	"./bdrupload"
	"path/filepath"
	"github.com/kless/goconfig/config"
	_ "github.com/mattn/go-sqlite3"
)

var (
	configFile = flag.String("config", "../etc/config.cfg", "Defines where to load configuration from")
	newDB      = flag.Bool("new-db", false, "true = creates a new database | false = use existing database")
	debug      = flag.Bool("debug", false, "activates debug mode")
	pool_flag  = flag.Int("threads", 0, "overwrites threads in [Client] section in config.cfg")

	upchan = make(chan *bdrupload.Upchan_t, 100)
	downchan = make(chan *bdrupload.Downchan_t, 100)
	done = make(chan bool)

	pool int
)

func checkPath(dirArray []string, dir string) bool {
	for _, i := range dirArray {
		if i == dir {
			return true
		}
	}
	return false
}

func backupDir(db *sql.DB, dirList string, excludeList string, dataBaseName string, debug bool) error {
	var dirname string
	var i int
	var fileC int64
	var backupFileC int64
	var dirC int64
	var dFile int64
	var dDir int64
	fileC = 0
	dirC = 0
	backupFileC = 0
	dFile = 0
	dDir = 0
	start := time.Now().Unix()
	dirArray := strings.Split(dirList, " ")
	excludeArray := strings.Split(excludeList, " ")
	i = 0
	for i < len(dirArray) {
		dirname = dirArray[i]
		// get dirID of dirname, even if it needs inserted.
		dirID, err := bdrsql.GetSQLID(db, "dirs", "path", dirname)
		// get a map for filename -> modified time
		SQLmap := bdrsql.GetSQLFiles(db, dirID)
		if debug == true {
			fmt.Printf("Scanning dir %s ", dirname)
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
		Fmap := map[string]int64{}
		// Iterate over the entire directory
		dFile = 0
		dDir = 0
		for _, f := range fi {
			if !f.IsDir() {
				fileC++ //track files per backup
				dFile++ //trace files per directory
				// and it's been modified since last backup
				if f.ModTime().Unix() <= SQLmap[f.Name()] {
					// log.Printf("NO backup needed for %s \n",f.Name())
					Fmap[f.Name()] = f.ModTime().Unix()
				} else {
					// log.Printf("backup needed for %s \n",f.Name())
					backupFileC++
					bdrsql.InsertSQLFile(db, f, dirID)
				}
			} else { // is directory
				dirC++ //track directories per backup
				dDir++ //track subdirs per directory
				fullpath := filepath.Join(dirname, f.Name())

				for _, excludeEntry := range excludeArray {
					// avoid an infinite loop 
					if !checkPath(dirArray, fullpath) && !strings.Contains(fullpath, excludeEntry){
						dirArray = append(dirArray, fullpath)
					}
				}
			}
		}
		// All files that we've seen, set last_seen
		t1 := time.Now().UnixNano()
		bdrsql.SetSQLSeen(db, Fmap, dirID)
		if debug == true {
			t2 := time.Now().UnixNano()
			fmt.Printf("files=%d dirs=%d duration=%dms\n", dFile, dDir, (t2-t1)/1000000)
		}
		i++
	}
	// if we have not seen the files since start it must have been deleted.
	bdrsql.SetSQLDeleted(db, start)

	log.Printf("Scanned %d files and %d directories\n", fileC, dirC)
	log.Printf("%d files scheduled for backup\n", backupFileC)

	return nil
}

func main() {
	flag.Parse()

	log.Printf("loading config file from %s\n", *configFile)
	configF, err := config.ReadDefault("../etc/config.cfg")
	if err != nil {
		log.Fatalf("ERROR: %s", err)
	}

	pool_config, err := configF.Int("Client", "threads")
	if err != nil {
		log.Fatalf("ERROR: %s", err)
	}
	if *pool_flag != 0 {
		pool = *pool_flag
	} else {
		pool = pool_config
	}

	runtime.GOMAXPROCS(pool)

	dirList, err := configF.String("Client", "backup_dirs_secure")
	if err != nil {
		log.Fatalf("ERROR: %s", err)
	}

	excludeList, err := configF.String("Client", "exclude_dirs")
	if err != nil {
		log.Fatalf("ERROR: %s", err)
	}

	dataBaseName, err := configF.String("Client", "sql_file")
	if err != nil {
		log.Fatalf("ERROR: %s", err)
	}
	log.Printf("attempting to open %s", dataBaseName)

	db, err := bdrsql.Init_db(dataBaseName, *newDB, *debug)
	if err != nil {
		log.Printf("could not open %s, error: %s", dataBaseName, err)
	} else {
		log.Printf("opened database %v\n", dataBaseName)
	}

	err = bdrsql.CreateBDRTables(db)
	if err != nil {
		log.Printf("couldn't create tables: %s", err)
	} else {
		log.Printf("created tables\n")
	}

	log.Printf("backing up these directories: %s\n", dirList)
	log.Printf("start walking...")
	t0 := time.Now()
	err = backupDir(db, dirList, excludeList, dataBaseName, *debug)
	t1 := time.Now()
	duration := t1.Sub(t0)
	if err != nil {
		log.Printf("walking didn't finished successfully. Error: %s", err)
	} else {
		log.Printf("walking successfully finished")
	}
	log.Printf("walking took: %v\n", duration)

	// shutdown database, make a copy, open it, backup copy of db
	// db, _ = bdrsql.BackupDB(db,dataBaseName)
	// launch server to receive uploads
	for i := 0; i<pool; i++ {
		go bdrupload.Uploader(upchan, done, *debug)
	}
	log.Printf("started %d uploaders\n", pool)
	// send all files to be uploaded to server.

	log.Printf("started sending files to uploaders...\n")
	bdrsql.SQLUpload(db, upchan)
	for i:=0;i<pool;i++ {
		<- done
	}
	log.Printf("uploading successfully finished\n")
}
