package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	dirname := "." + string(filepath.Separator)
	d, err := os.Open(dirname)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fi, err := d.Readdir(-1)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, fi := range fi {
		if !fi.IsDir() {
			fmt.Println(fi.Name(), fi.Size(), "bytes")
		}
	}
}
