package main

import (
	"bufio"
	"crypto/sha256"
	"flag"
	"fmt"
	"os"

//	"bytes"
//    "hash"
)

var plaintext string

func main() {

	var (
		filename string
		file     *os.File
		count    int
		err      error
	)
	flag.Parse()
	if flag.NArg() > 0 {
		filename = flag.Arg(0)
	} else {
		filename = "1MB"
	}

	if file, err = os.Open(filename); err != nil {
		return
	}
	reader := bufio.NewReader(file)
	buffer := make([]byte, 16384)
	h := sha256.New() // h is a hash.Hash
	for {
		if count, err = reader.Read(buffer); err != nil {
			break
		}
		h.Write(buffer[:count])
	}
	fmt.Printf("%x  %s\n", h.Sum(nil), filename)
}
