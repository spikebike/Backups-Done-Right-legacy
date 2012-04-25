package main

import (
	"fmt"
	"../bdrsql"
	)

const DataBaseName = "./example-db.sql"

func main() {
	db, err := bdrsql.Init_db(DataBaseName, false);
	if err != nil {
		fmt.Println("coudln't open database. Error: %s", err)
	} else {
		fmt.Println("opened database")
	}
	defer db.Close()

	stmt, err := db.Prepare("select name from files where id = ?")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer stmt.Close()

	var name string
	err = stmt.QueryRow("1").Scan(&name)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("id 1 name %s\n", name)
}
