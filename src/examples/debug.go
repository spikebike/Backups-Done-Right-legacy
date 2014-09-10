package main

import (
	"expvar"
	"fmt"
	"net/http"
	"time"
)

// docs at http://golang.org/pkg/expvar/

const (
	httpPortTCP  = 8711
	delaySeconds = 120
)

func main() {
	go http.ListenAndServe(fmt.Sprintf(":%d", httpPortTCP), nil)
	fmt.Printf("Go to http://localhost:%d/debug/vars in the next %d seconds \n", httpPortTCP, delaySeconds)
	time.Sleep(delaySeconds * time.Second)
}

var (
	iAmAnImportantNumber = expvar.NewInt("iAmAnImportNumber")
	totalClients         = expvar.NewInt("totalClients")
	totalRequsts         = expvar.NewInt("totalRequests")
)
