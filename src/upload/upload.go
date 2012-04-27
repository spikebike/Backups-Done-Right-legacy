package upload

import (
	"../mystructs"
	"crypto/sha256"
	"bufio"
	"fmt"
	"os"
)

func Server(upchan chan *mystructs.Upchan_t) {
	var count int
	buffer := make([]byte, 16384)

	for f := range upchan {
		fmt.Printf("Server: received rowID=%d path=%s\n", f.Rowid, f.Path)
		//      fmt.Printf("%T %#v\n",f,f)
		file, err := os.Open(f.Path)
		reader := bufio.NewReader(file)
		h := sha256.New() // h is a hash.Hash
		for {
				if count, err = reader.Read(buffer); err != nil {
				break
			}
			h.Write(buffer[:count])
		}
		fmt.Printf("%x  %s\n", h.Sum(nil), f.Path)

	}
	fmt.Print("Server: Channel closed, existing\n")
}
