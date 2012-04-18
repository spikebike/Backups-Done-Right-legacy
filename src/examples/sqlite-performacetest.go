package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"time"
	"flag"
)

var (
	arg1 = flag.Int("values", 1000, "Number of values to be inserted")
	arg2 = flag.Int("chunks", 64, "Size of statements per commit")
)

const DataBaseName = "./foo.db"

func run() {
	var values int = *arg1
	var bufsize int = *arg2
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

	commits := values / bufsize

	fmt.Printf("values: %d\n", values)
	fmt.Printf("chunks: %d\n", bufsize)

	for a := 0; a <= commits; a++ {
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
	flag.Parse()

	t0 := time.Now()
	run()
	t1 := time.Now()

	fmt.Printf("Done!\n")
	fmt.Printf("Duration: %v\n", t1.Sub(t0))
}
