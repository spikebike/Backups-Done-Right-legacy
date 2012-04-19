package main

import (
	
	"./bdrservice" // defines BDR related protocols
	"./tlscon"     // handles SSL connections
	"log"
	"net"
)

type Request struct{}

func (Request) Request(in *bdrservice.RequestMessage, out *bdrservice.RequestACKMessage) error {
	log.Printf("server: blobarray=%v %T", *in.Blobarray[0].Sha256,*in.Blobarray[0].Sha256)
	log.Printf("server: blobarray=%v %T", *in.Blobarray[0],*in.Blobarray[0])
	
	out.Received = new(int32)
	*out.Received = 5 
	return nil
}

func RequestFunc(conn net.Conn) {
	bdrservice.ServeRequestService(conn, Request{})
}

func main() {
	fptr := RequestFunc
	tlscon.ServerTLSListen("0.0.0.0:8000", fptr)
}
