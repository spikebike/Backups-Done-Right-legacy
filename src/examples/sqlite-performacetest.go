package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"time"
)

var (
	values int
	chunks int
	bufsize int
)

const DataBaseName = "./foo.db"

func run() {
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

	chunks := values / bufsize

	fmt.Printf("chunks: %d\n", chunks)

	for a := 0; a <= chunks; a++ {
		for i := 0; i <= bufsize; i++ {
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
	}
}

func main() {
	fmt.Printf("values: ")
	fmt.Scanf("%d", &values)
	fmt.Printf("values per Chunk: ")
	fmt.Scanf("%d", &bufsize)

	t0 := time.Now()
	run()
	t1 := time.Now()

	fmt.Printf("Done!\n")
	fmt.Printf("Duration: %v\n", t1.Sub(t0))
}
