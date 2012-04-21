package main

import (
	"fmt"
	//	"io/ioutil"
	"log"
	"os"
	//	"path"
	"path/filepath"
	"syscall"
)

func main() {

	t, err := os.Readlink("bar")
	fmt.Printf("t: %+v\n", t)

	t2, err := filepath.EvalSymlinks("bar")
	fmt.Printf("t2: %+v\n", t2)

	fi, err := os.Lstat("foo")
	if fi.Mode()&os.ModeSymlink != 0 {
		fmt.Printf("we haave a link!!!\n")
	}

	fi, err = os.Stat("inode.go")
	if fi.Mode()&os.ModeSymlink != 0 {
		fmt.Printf("we haave a link!!!\n")
	}
	fmt.Printf("%T %#v\n", fi, fi)

	if err != nil {
		log.Fatal(err)
	}
	// Check that it's a Unix file.
	unixStat, ok := fi.Sys().(*syscall.Stat_t)
	if !ok {
		log.Fatal("hello.go: not a Unix file")
	}

	fmt.Printf("full stat: %+v\n", unixStat)
	fmt.Printf("file i-number: %d\n", unixStat.Ino)
	fmt.Printf("mode: %+v\n", unixStat.Mode)
	fmt.Printf("mode: %+v\n", unixStat.Mode&0777)
	fmt.Printf("mode: %+v\n", fi.Mode()&os.ModePerm)
	fmt.Printf("UID: %+v\n", unixStat.Uid)
	fmt.Printf("GID: %+v\n", unixStat.Gid)
	fmt.Printf("Atim: %+v\n", unixStat.Atim)
	fmt.Printf("Mtim: %+v\n", unixStat.Mtim)
	fmt.Printf("Ctim: %+v\n", unixStat.Ctim)
}
