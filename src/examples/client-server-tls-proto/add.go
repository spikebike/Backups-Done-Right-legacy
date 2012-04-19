package main

import (
	"./addservice"
	"./tlscon"
	"flag"
	"log"
	"net"
)

// Add is the type which will implement the addservice.AddService interface
// and can be called remotely.  In this case, it does not have any state, but
// it could.
type Add struct{}

// Add is the function that can be called remotely.  Note that this can be
// called concurrently, so if the Echo structure did have internal state,
// it should be designed for concurrent access.
func (Add) Add(in *addservice.AddMessage, out *addservice.SumMessage) error {
	log.Printf("server: X=%d", *in.X)
	log.Printf("server: Y=%d", *in.Y)
	out.Z = new(int32)
	*out.Z = *in.X + *in.Y
	log.Printf("server: Z=%d", *out.Z)
	return nil
}

func AddFunc(conn net.Conn) {
	addservice.ServeAddService(conn, Add{})
}

func main() {
	flag.Parse()
	fptr := AddFunc
	tlscon.ServerTLSListen("0.0.0.0:8000", fptr)
}
