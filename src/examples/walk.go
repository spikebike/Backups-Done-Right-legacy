package main

import (
	"os"
	"path/filepath"
	"log"
)

func main() {
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
	err := filepath.Walk(".", markFn)
	if err != nil {
		log.Fatal(err)
	}
}
