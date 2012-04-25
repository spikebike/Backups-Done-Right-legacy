package bdrsql

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"log"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

var (
	sqls = []string{
		"create table dirs (id INTEGER PRIMARY KEY, mode INT, ino BIGINT, uid INT, gid INT, path varchar(2048), last_seen BIGINT, deleted INT)",
		"create table files (id INTEGER PRIMARY KEY, mode INT, ino BIGINT, dev BIGINT, uid INT, gid INT, size BIGINT, atime BIGINT, mtime BIGINT, ctime BIGINT, name varchar(254), dirID BIGINT, last_seen BIGINT, deleted INT, do_upload INT, FOREIGN KEY(dirID) REFERENCES dirs(id))",
		"create index ctimeindex on files (ctime)",
		"create index pathindex on dirs (path)",
	}
)

func Init_db(dataBaseName string, newDB bool, debug bool) (db *sql.DB, err error) {
	if newDB == true {
		os.Remove(dataBaseName)
	}

	db, err = sql.Open("sqlite3", dataBaseName)
	if err != nil {
		log.Printf("couldn't open database: %s", err)
		os.Exit(1)
	}
	// Allow commits to be buffered, MUCH faster.  
	// debug = true makes database writes synchronous and much slower,
	if debug == false {
		_, err = db.Exec("PRAGMA synchronous=OFF")
		if err != nil {
			log.Printf("%s", err)
		}
	}
	return db, err
}

func CopyFile(dstName, srcName string) (written int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		return
	}
	defer src.Close()
	//log.Printf("dest=%#v\n", dstName)

	dst, err := os.Create(dstName)
	if err != nil {
		return
	}
	defer dst.Close()

	return io.Copy(dst, src)
}

func BackupDB(db *sql.DB, dbname string) (*sql.DB, error) {
	time := time.Now().UnixNano()
	tmpDBName := fmt.Sprintf("%s.%d", dbname, time)
	log.Printf("backing up %s to %s\n", dbname, tmpDBName)
	db.Close()                  // Insure that ALL writes are completed.
	CopyFile(tmpDBName, dbname) // make a snapshot
	fi, err := os.Stat(tmpDBName)
	directory, _ := filepath.Split(tmpDBName)
	db2, err := Init_db(dbname, false, false)
	if err != nil {
		log.Fatalf("Init_DB failed with: %s\n", err)
	}
	dirID, err := GetSQLID(db2, "dirs", "path", directory)
	InsertSQLFile(db2, fi, dirID)
	return db2, err
}

func CreateBDRTables(db *sql.DB, debug bool) error {
	var err error

	for _, sql := range sqls {
		_, err = db.Exec(sql)
		if err != nil {
			log.Printf("%s", err)
		}
	}

	return err
}

func GetSQLFiles(db *sql.DB, dirID int64) map[string]int64 {
	var fileMap = map[string]int64{}
	var name string
	var mtime int64
	stmt, err := db.Prepare("select name,mtime from files where dirID=? and deleted=0")
	if err != nil {
		fmt.Printf("GetSQLFiles prepare of select failed: %s\n", err)
	}
	rows, err := stmt.Query(dirID)
	if err != nil {
		fmt.Printf("GetSQLFiles query failed: %s\n", err)
	}
	for rows.Next() {
		rows.Scan(&name, &mtime)
		fileMap[name] = mtime
	}
	return fileMap
}

func SetSQLSeen(db *sql.DB, fmap map[string]int64, dirID int64) {
	now := time.Now().Unix()
	//	tx,_ := db.Begin()
	update := fmt.Sprintf("update files set last_seen=%d where name=? and dirID=%d and deleted=0 and ctime=?", now, dirID)
	stmt, _ := db.Prepare(update)
	for i, _ := range fmap {
		//		log.Printf("file = %v dirID=%d\n",i,dirID)
		stmt.Exec(i, fmap[i])
	}
	//	tx.Commit()
}

func SetSQLDeleted(db *sql.DB, now int64) {
	stmt, err := db.Prepare("update files set deleted=1 where last_seen<?")
	if err != nil {
		log.Println(err)
	}
	stmt.Exec(now)
}

func GetSQLID(db *sql.DB, tablename string, field string, value string) (int64, error) {

	var dirID int64

	dirID = -1

	query := "select id from " + tablename + " where " + field + "= (?)"
	//	fmt.Printf("query=%s\n", query)
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Printf("GetSQLID: prepare select failed: %s\n", err)
	}
	err = stmt.QueryRow(value).Scan(&dirID)
	if err != nil {
		log.Printf("GetSQLID: missing %s, error %s\n", value, err)

		insert := "insert into " + tablename + "(" + field + ") values(?);"
		stmt, err := db.Prepare(insert)
		if err != nil {
			log.Println(err)
			return -1, err
		}
		result, err := stmt.Exec(value)
		dirID, err = result.LastInsertId()
	}
	return dirID, err
}

func InsertSQLFile(db *sql.DB, fi os.FileInfo, dirID int64) error {
	now := time.Now().Unix()
	e, _ := fi.Sys().(*syscall.Stat_t)

	stmt, err := db.Prepare("insert into files(name,size,mode,gid,uid,ino,dev,mtime,atime,ctime,last_seen,dirID,deleted,do_upload) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		log.Printf("InsertSQL prepare: %s\n", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(fi.Name(), e.Size, e.Mode, e.Gid, e.Uid, e.Ino, e.Dev, e.Mtim.Sec, e.Atim.Sec, e.Ctim.Sec, now, dirID, 0, 1)
	if err != nil {
		log.Printf("InsertSQL Exec: %s\n", err)
		return err
	}
	return err
}

func main_test() {
	db, _ := Init_db("fsmeta.sql", true, false)
	id, _ := GetSQLID(db, "dirs", "path", "/home/bill/bdr/src/bdrsql")
	fmt.Printf("I found id=%d\n", id)
	d, _ := os.Open(".")
	fi, _ := d.Readdir(-1)
	for _, fi := range fi {
		//		unixStat, _ := fi.Sys().(*syscall.Stat_t)
		//		fmt.Printf("%T %#v\n",&fi,&fi)
		//		fmt.Printf("%T %#v\n\n", unixStat, unixStat)
		InsertSQLFile(db, fi, -1)
		//		fmt.Printf("%T %#v\n\n", unixStat.Ino, unixStat.Ino)
	}
}
