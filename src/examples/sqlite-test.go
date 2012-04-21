package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
)

const DataBaseName = "./foo.db"

func main() {
	os.Remove(DataBaseName)

	db, err := sql.Open("sqlite3", DataBaseName)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	sqls := []string{
		"create table foo (name text)",
		"delete from foo",
	}
	for _, sql := range sqls {
		_, err = db.Exec(sql)
		if err != nil {
			fmt.Printf("%q: %s\n", err, sql)
			return
		}
	}

	tx, err := db.Begin()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("tx = %#v tx=%T\n", tx, tx)
	stmt, err := tx.Prepare("insert into foo (name) values(?)")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer stmt.Close()
	for i := 0; i < 12; i++ {
		result, err := stmt.Exec(fmt.Sprintf("Hello World! %03d", i))
		id, err := result.LastInsertId()
		r, err := result.RowsAffected()

		fmt.Printf("last inserted ID = %d RowsAffected=%d\n", id, r)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	tx.Commit()

	rows, err := db.Query("select rowid,name from foo")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var rowid int64
		var name string
		rows.Scan(&rowid, &name)
		println(rowid, name)
	}
	rows.Close()

	stmt, err = db.Prepare("select name from foo where rowid = ?")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer stmt.Close()
	var name string
	err = stmt.QueryRow("3").Scan(&name)
	if err != nil {
		fmt.Println(err)
		return
	}
	println(name)
}
