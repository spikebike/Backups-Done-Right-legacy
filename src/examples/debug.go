package main

import (
	"expvar"
	"fmt"
	"net/http"
	"time"
)

const (
	httpPortTCP  = 8711
	delaySeconds = 120
)

func main() {
	go http.ListenAndServe(fmt.Sprintf(":%d", httpPortTCP), nil)
	fmt.Printf("To go http://localhost:%d/debug/vars in the next %d seconds \n", httpPortTCP, delaySeconds)
	time.Sleep(delaySeconds * time.Second)
}

var (
	iAmAnImportantNumber = expvar.NewInt("iAmAnImportNumber")
	totalClients         = expvar.NewInt("totalClients")
	totalRequsts         = expvar.NewInt("totalRequests")
)
