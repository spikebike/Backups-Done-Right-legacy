package main

import (
	"fmt"
	"../bdrsql"
	"database/sql"
	)

const DataBaseName = "./example-db.sql"

type file_info_t struct {
	id int
	mode int
	ino int
	dev int
	nlink int
	uid int
	gid int
	size int64
	atime int
	mtime int
	ctime int
	name string
	dirID int
	last_seen int
	deleted int
}

func printStruct(f *file_info_t) {
	fmt.Println("name: ", f.name)
	fmt.Println("uid: ", f.uid)
	fmt.Println("gid: ", f.gid)
	fmt.Println("size: ", f.size)
	fmt.Println("mtime: ", f.mtime)
	fmt.Println("atime: ", f.atime)
	fmt.Println("ctime: ", f.ctime)
}

func readFileInfo(db *sql.DB, fname string) *file_info_t{
	f := &file_info_t{}

        stmt, err := db.Prepare("select mode, ino, uid, gid, size, mtime, atime, ctime from files where name = ?")
        if err != nil {
                fmt.Println(err)
                return nil
        }
        defer stmt.Close()

	err = stmt.QueryRow(fname).Scan(&f.mode, &f.ino, &f.uid, &f.gid, &f.size, &f.mtime, &f.atime, &f.ctime)
        if err != nil {
                fmt.Println(err)
                return nil
        }
	f.name = fname

	return f
}

func main() {
	db, err := bdrsql.Init_db(DataBaseName, false);
	if err != nil {
		fmt.Println("coudln't open database. Error: %s", err)
	} else {
		fmt.Println("opened database")
	}
	defer db.Close()

	f := readFileInfo(db, "avr.odt")
	printStruct(f)
}
