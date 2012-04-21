package bdrsql

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"syscall"

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

func Init_db(dataBaseName string) (db *sql.DB, err error) {
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

func GetSQLFiles(db *sql.DB, dirID int64) map[string]int64 {
	var fileMap = map[string]int64{}
	var name string
	var mtime int64
	stmt, err := db.Prepare("select name,mtime from files where dirID=?")
	if err != nil {
		fmt.Println(err)
	}
	rows,err := stmt.Query(dirID)
	for rows.Next() {
		rows.Scan(&name, &mtime)
		fileMap[name] = mtime 
	}
	return fileMap
}
func GetSQLID(db *sql.DB, tablename string, field string, value string) (int64, error) {

	var dirID int64

	dirID = -1

	query := "select id from " + tablename + " where " + field + "= (?)"
	fmt.Printf("query=%s\n", query)
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Println(err)
	}
	err = stmt.QueryRow(value).Scan(&dirID)
	if err != nil {
		log.Println(err)

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

func insertSQLFile(db *sql.DB, fi os.FileInfo) error {
	e, _ := fi.Sys().(*syscall.Stat_t)
	fmt.Printf("fi %T %#v\n", fi, fi)
	fmt.Printf("fi.name %T %#v\n\n", fi.Name(), fi.Name())
	fmt.Printf("e %T %#v\n\n", e, e)

	stmt, err := db.Prepare("insert into files(name,size,mode,gid,uid,ino,dev,mtime,atime,ctime,last_seen,dirID,deleted,do_upload) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		log.Println(err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(fi.Name(), e.Size, e.Mode, e.Gid, e.Uid, e.Ino, e.Dev, e.Mtim.Sec, e.Atim.Sec, e.Ctim.Sec, "now", -1, 0, 1)
	if err != nil {
		log.Println(err)
		return err
	}
	return err
}

func main_test() {
	db, _ := Init_db("fsmeta.sql")
	id, _ := GetSQLID(db, "dirs", "path", "/home/bill/bdr/src/bdrsql")
	fmt.Printf("I found id=%d\n", id)
	d, _ := os.Open(".")
	fi, _ := d.Readdir(-1)
	for _, fi := range fi {
		//		unixStat, _ := fi.Sys().(*syscall.Stat_t)
		//		fmt.Printf("%T %#v\n",&fi,&fi)
		//		fmt.Printf("%T %#v\n\n", unixStat, unixStat)
		insertSQLFile(db, fi)
		//		fmt.Printf("%T %#v\n\n", unixStat.Ino, unixStat.Ino)
	}
}
