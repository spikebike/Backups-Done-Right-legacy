package main

import (
	"./bdr_proto"
	"crypto/sha256"
	"fmt"
)

import M "math/rand"
import C "crypto/rand"

func main() {
	// not sure how to create more than 1 blob message within a single request.  
	// The bdr_proto.prof file allows it.  The protobufs comments at
	// http://code.google.com/p/protobuf/source/browse/trunk/src/google/protobuf/message.h
	// seem to claim something like bdr_proto.Add_RequestBlob should work 
	// the bdr_proto/bdr_proto.pb.go makes no mention of add.
	var s1 int32
	var i int32
	var intptr *int32
	var strptr *string
	var randBytes []byte
	str1 := "Hello"
	s1 = 1024

	// I had hoped just repeating the records would work.
	//
	//	t1:= &bdr_proto.RequestBlob{Sha256: &str1, Bsize: &s1, Sha256: &str2, Bsize: &s2}
	// build an array of 32 RequestBlob
	Blobarray := make([]bdr_proto.RequestBlob, 32)
	//* build an array of Requestblob pointers
	BlobarrayPtr := make([]*bdr_proto.RequestBlob, 32)
	t1 := new(bdr_proto.Request)
	t1.Blobarray = BlobarrayPtr
	for i = 0; i < 32; i++ {
		// read in 16 bytes from /dev/urandom to sha256
		randBytes = make([]byte, 16)
		_, _ = C.Read(randBytes)
		// create a new hash, and do a crypty hash of the random bytes.
		hash := sha256.New()
		_, _ = hash.Write(randBytes)
		// build a new string and store the hash in it
		strptr = new(string)
		*strptr = fmt.Sprintf("%x", hash.Sum(nil))
		/* point each array to the address of the invidual blob */
		t1.Blobarray[i] = &Blobarray[i]
		/* for that blob set the sha256 to the address of the new sha256 */
		t1.Blobarray[i].Sha256 = strptr
		/* create a new int32, store a random number in it, and assign
		it to bsize */
		intptr = new(int32)
		*intptr = M.Int31()
		t1.Blobarray[i].Bsize = intptr
		fmt.Printf("i= %2d size= %10d sha256=%s\n", i, *intptr, *strptr)
	}
	t2 := &bdr_proto.RequestBlob{Sha256: &str1, Bsize: &s1}
	fmt.Printf("t1 = %+v\n", t1)
	fmt.Printf("t2 = %+v\n", t2)
}
