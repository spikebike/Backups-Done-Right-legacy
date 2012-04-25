package main

import (
	"os"
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
	fmt.Println("name:  ", f.name)
	fmt.Println("uid:   ", f.uid)
	fmt.Println("gid:   ", f.gid)
	fmt.Println("size:  ", f.size)
	fmt.Println("mtime: ", f.mtime)
	fmt.Println("atime: ", f.atime)
	fmt.Println("ctime: ", f.ctime)
}

func createFile(f *file_info_t) error {
	os.Remove(f.name)

	file, err := os.Create(f.name)
	if err != nil {
		fmt.Println("couldn't create file. ERROR: %s", err)
	}

	file.Truncate(f.size)
	file.Chown(f.uid, f.gid)

	return nil
}

func readFileInfo(db *sql.DB, fname string) *file_info_t{
	f := &file_info_t{}

	f.name = fname

        stmt, err := db.Prepare("select mode, uid, gid, size, mtime, atime, ctime from files where name = ?")
        if err != nil {
                fmt.Println(err)
                return nil
        }
        defer stmt.Close()

	err = stmt.QueryRow(fname).Scan(&f.mode, &f.uid, &f.gid, &f.size, &f.mtime, &f.atime, &f.ctime)
        if err != nil {
                fmt.Println(err)
                return nil
        }

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

	fmt.Printf("file info from database:\n")
	printStruct(f)

	fmt.Printf("I restore this file now...\n")
	err = createFile(f)
	if err != nil {
		fmt.Println("couldn't restore file")
	}

	fmt.Printf("finished restoring\n")
}
