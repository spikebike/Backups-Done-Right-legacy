package tlscon

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"log"
	"net"
)

func handleClient(conn net.Conn, f func(conn net.Conn)) {
	tlscon, ok := conn.(*tls.Conn)
	if ok {
		log.Print("server: conn: type assert to TLS succeedded")
		err := tlscon.Handshake()
		if err != nil {
			log.Fatalf("server: handshake failed: %s", err)
		} else {
			log.Print("server: conn: Handshake completed")
		}
		state := tlscon.ConnectionState()
		// Note we could reject clients if we don't like their public key.		
		log.Println("Server: client public key is:")
		for _, v := range state.PeerCertificates {
			log.Print(x509.MarshalPKIXPublicKey(v.PublicKey))
		}
		// Now that we have completed SSL/TLS 
		// hopefully F does the same as below
		f(conn)
		//		addservice.ServeAddService(tlscon, Add{})
	}
}

func ServerTLSListen(service string, f func(conn net.Conn)) {

	// Load x509 certificates for our private/public key, makecert.sh will
	// generate them for you.

	cert, err := tls.LoadX509KeyPair("certs/server.pem", "certs/server.key")
	if err != nil {
		log.Fatalf("server: loadkeys: %s", err)
	}
	// Note if we don't tls.RequireAnyClientCert client side certs are ignored.
	config := tls.Config{Certificates: []tls.Certificate{cert}, ClientAuth: tls.RequireAnyClientCert}
	config.Rand = rand.Reader
	listener, err := tls.Listen("tcp", service, &config)
	if err != nil {
		log.Fatalf("server: listen: %s", err)
	}
	log.Print("server: listening")
	// Keep this loop simple/fast as to be able to handle new connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("server: accept: %s", err)
			break
		}
		log.Printf("server: accepted from %s", conn.RemoteAddr())
		// Fire off go routing to handle rest of connection.
		go handleClient(conn, f)
	}
}

