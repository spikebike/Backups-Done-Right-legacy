package bdr-proto

import (
	"./bdr_proto"
	"crypto/sha256"
	"fmt"
)

//import M "math/rand"
import C "crypto/rand"

const MAX = 4

func main() {
	req := &bdr_proto.Request{}
	for i := 0; i < MAX; i++ {
		// read in 16 bytes from /dev/urandom to sha256
		randBytes := make([]byte, 16)
		C.Read(randBytes)

		// get its size
		size := int32(len(randBytes))

		// create a new hash, and do a crypty hash of the random bytes.
		sha := sha256.New()
		sha.Write(randBytes)
		strhash := fmt.Sprintf("%x", sha.Sum(nil))
		fmt.Printf("i=%d sha=%s size=%d\n", i, strhash, size)

		req.Blobarray = append(req.Blobarray, &bdr_proto.RequestBlob{Sha256: &strhash, Bsize: &size})
	}
	//	fmt.Printf("%#V %T\n", req, req)
	for i := 0; i < 4; i++ {
		fmt.Printf("%s %d \n", *req.Blobarray[i].Sha256, *req.Blobarray[i].Bsize)
	}

}
