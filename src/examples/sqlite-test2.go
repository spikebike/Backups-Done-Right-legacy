package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
)

const DataBaseName = "./foo.db"

func main() {
	var tx *sql.Tx

	os.Remove(DataBaseName)

	db, err := sql.Open("sqlite3", DataBaseName)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	_, err = db.Exec("create table foo (id INT)")
	if err != nil {
		fmt.Printf("%q: %s\n", err)
		return
	}

	tx, err = db.Begin()
	if err != nil {
		fmt.Println(err)
	}

	for i := 0; i < 42; i++ {
		stmt, err := tx.Prepare("insert into foo(id) values(?)")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer stmt.Close()
		_, err = stmt.Exec(i)
		if err != nil {
			fmt.Println(err)
			return
		}

	}
	tx.Commit()

	fmt.Println("Done!")
}
