package main

import (
	"os"
	"path/filepath"
	"log"
	"flag"
)

func main() {
	flag.Parse()

	markFn := func(path string, info os.FileInfo, err error) error {
		if path == "doozerd" { // Will skip walking of directory pictures and its contents.
			log.Println("Skipping",path)
			return filepath.SkipDir
		}
		if err != nil {
			return err
		}
		log.Println(path)
		return nil
	}

	path := "."
	if flag.Arg(0) != "" {
		path = flag.Arg(0)
	}

	err := filepath.Walk(path, markFn)
	if err != nil {
		log.Fatal(err)
	}
}
