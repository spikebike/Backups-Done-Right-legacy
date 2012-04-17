package main

import (
	"./bdr_proto"
	"fmt"
	"crypto/sha256"
	)

import M "math/rand"
import C "crypto/rand"

func main() {
// not sure how to create more than 1 blob message within a single request.  
// The bdr_proto.prof file allows it.  The protobufs comments at
// http://code.google.com/p/protobuf/source/browse/trunk/src/google/protobuf/message.h
// seem to claim something like bdr_proto.Add_RequestBlob should work 
// the bdr_proto/bdr_proto.pb.go makes no mention of add.
	var s1 int32;
	var i int32;
	var intptr *int32;
	var strptr *string;
	randBytes := make([]byte,16)
	_,_ = C.Read(randBytes)
	str1:="Hello"
	s1=1024

// I had hoped just repeating the records would work.
//
//	t1:= &bdr_proto.RequestBlob{Sha256: &str1, Bsize: &s1, Sha256: &str2, Bsize: &s2}
	Blobarray := make([]bdr_proto.RequestBlob,32)
	BlobarrayPtr := make([]*bdr_proto.RequestBlob,32)
	t1:= new(bdr_proto.Request)
	t1.Blobarray=BlobarrayPtr;
	for i =0; i<32; i++ {
		fmt.Printf("i=%d\n",i)
		hash:=sha256.New()
		_,_ = hash.Write(randBytes)
		_,_ = C.Read(randBytes)
		strptr=new(string)
		*strptr=fmt.Sprintf("i=%v",hash.Sum(nil))
		fmt.Printf("sha256 for i=%d is %s\n",i,strptr);
//		BlobarrayPtr[i]{Sha256: &str1, Bsize: &s1}
		t1.Blobarray[i]=&Blobarray[i];
		t1.Blobarray[i].Sha256=strptr;
		intptr=new(int32)
		*intptr=M.Int31();
		t1.Blobarray[i].Bsize=intptr;
	//	t1.Blobarray[i].RequestBlob{Sha256: &str1, Bsize: &i}
	}
//	t1.Blobarray = &Blobarray
	t2:= &bdr_proto.RequestBlob{Sha256: &str1, Bsize: &s1}
	fmt.Printf("t1 = %+v\n",t1)
	fmt.Printf("t2 = %+v\n",t2)
}
