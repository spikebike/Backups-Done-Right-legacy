package main

import (
	"fmt"
	//  "io/ioutil"
	"log"
	"os"
	//  "path"
	"path/filepath"
	"syscall"
	"time"
)

const (
	timeFmt = "Jan _2 15:04"
)



func main() {
	var dirname string
	dirname = "."
	d, err := os.Open(dirname)
	fi, err := d.Readdir(-1)
	if err != nil {
		log.Printf("Directory %s failed with error %s", dirname, err)
	}
	for _, fi := range fi {
		unixStat, ok := fi.Sys().(*syscall.Stat_t)
		if !ok {
			fmt.Printf("%s is not a Unix file",fi.Name())
		}
		if fi.Mode()&os.ModeSymlink != 0 {
			t, _:= os.Readlink(fi.Name())
			fmt.Printf("%s is a link, destination is %s!!!!!\n",fi.Name(),t)
			t2, _:= filepath.EvalSymlinks(fi.Name())
			fmt.Printf("EvalSymLinks claims destination is %s\n",t2)
		} else {
			fmt.Printf("%s is a not a link\n",fi.Name())
		}
		if !fi.IsDir() {
			fmt.Printf("%s %d bytes Mode=%+v Perm=%+v UID=%+v GID=%+v\n", fi.Name(), fi.Size(), unixStat.Mode,fi.Mode()&os.ModePerm,unixStat.Uid,unixStat.Gid)
			t := time.Unix(unixStat.Atim.Sec,unixStat.Atim.Nsec)
			fmt.Printf("atime: %+v\n",t)
			t = time.Unix(unixStat.Ctim.Sec,unixStat.Ctim.Nsec)
			fmt.Printf("ctime: %+v\n",t)
			t = time.Unix(unixStat.Mtim.Sec,unixStat.Mtim.Nsec)
			fmt.Printf("mtime: %+v\n",t)
		} else {
			fmt.Printf("found directory %s\n", fi.Name())
		}
		fmt.Printf("\n");
	}
}
