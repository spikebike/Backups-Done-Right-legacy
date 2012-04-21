package main

import ( 
		"database/sql"
		"fmt"

		_ "github.com/mattn/go-sqlite3"
)	

func maptest(db *sql.DB) map[string] int64 {

		var testMap = map[string] int64 {}
		var id int64
		var path string
		rows, _:= db.Query("select id,name from files")

		for rows.Next() {
			rows.Scan(&id, &path)
			testMap[path] = id;
		}
		return testMap
}

func main() {
	db, _ := sql.Open("sqlite3", "../bdrsql/fsmeta.sql")
	m:=maptest(db)
	fmt.Printf("I found id=%v\n", m)
}
