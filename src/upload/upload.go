package upload

import (
	"../bdrupload"
	"bufio"
	"crypto/sha256"
	"fmt"
	"os"
	"time"
)

func Server(upchan chan *bdrupload.Upchan_t) {
	var count int
	var size int64
	buffer := make([]byte, 16384)

	for f := range upchan {
		fmt.Printf("Server: received rowID=%d path=%s\n", f.Rowid, f.Path)
		//      fmt.Printf("%T %#v\n",f,f)
		size = 0
		file, err := os.Open(f.Path)
		t0 := time.Now().UnixNano()
		reader := bufio.NewReader(file)
		h := sha256.New() // h is a hash.Hash
		for {
			if count, err = reader.Read(buffer); err != nil {
				break
			}
			size=size+int64(count)
			h.Write(buffer[:count])
		}
		t1 := time.Now().UnixNano()
		seconds :=float64(t1-t0)/1000000000
		fmt.Printf("%x %s %4.2f MB/sec\n", h.Sum(nil), f.Path,float64(size)/(1024*1024*seconds))

	}
	fmt.Print("Server: Channel closed, existing\n")
}
