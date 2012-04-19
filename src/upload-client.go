package main

import (
	"./bdrservice" // defines BDR related protocols
	"./tlscon"     // handles SSL connections
	C "crypto/rand"
	"crypto/sha256"
	"fmt"
	"log"
)

func main() {
	conn, err := tlscon.OpenTLSClient("127.0.0.1:8000")
	if err != nil {
		log.Fatalf("dial: %s", err)
	}

	request := bdrservice.NewRequestServiceClient(conn)
	req := &bdrservice.RequestMessage{}
	ack := &bdrservice.RequestACKMessage{}
	for i := 0; i < 4; i++ {
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

		req.Blobarray = append(req.Blobarray, &bdrservice.RequestMessageBlob{Sha256: &strhash, Bsize: &size})
	}
	if err := request.Request(req, ack); err != nil {
		log.Fatalf("Add failed with: %s", err)
	}
	fmt.Printf("Server ACKs %d records\n\n", *ack.Received)
}
