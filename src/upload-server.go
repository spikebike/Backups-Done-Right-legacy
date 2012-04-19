package main

import (
	
	"./bdrservice" // defines BDR related protocols
	"./tlscon"     // handles SSL connections
	"log"
	"net"
)

type Request struct{}

func (Request) Request(in *bdrservice.RequestMessage, out *bdrservice.RequestACKMessage) error {
	var records int32
	for _, blob := range in.Blobarray { 
		log.Printf("server: blobarray=%v %T", *blob.Sha256,*blob.Sha256)
		records++
	}

	
	out.Received = new(int32)
	*out.Received =  records 
	return nil
}

func RequestFunc(conn net.Conn) {
	bdrservice.ServeRequestService(conn, Request{})
}

func main() {
	fptr := RequestFunc
	tlscon.ServerTLSListen("0.0.0.0:8000", fptr)
}
