package main

import (
	"expvar"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// docs at http://golang.org/pkg/expvar/

const (
	httpPortTCP  = 8711
	delaySeconds = 10
)

func HandleSignals(sig chan os.Signal) {
	// this loop runs once per signal
	for ss := range sig {
		fmt.Printf("received: %s\n", ss)
		fmt.Printf("iAmAnImportantNumber=%s\n", iAmAnImportantNumber.String())
		fmt.Printf("totalClients=%s\n", totalClients.String())
		fmt.Printf("totalRequests=%s\n\n", totalRequests.String())
	}
}

func main() {
	// launch webserver to serve out debugging variables
	go http.ListenAndServe(fmt.Sprintf(":%d", httpPortTCP), nil)

	// make channel for sending signals
	sig := make(chan os.Signal, 1)
	// send a signal on channel cig for each SIGUSR1
	signal.Notify(sig, syscall.SIGUSR1)
	// run signal handler in gorouting
	go HandleSignals(sig)

	fmt.Printf("\nGo to http://localhost:%d/debug/vars in the next %d seconds \n", httpPortTCP, delaySeconds)
	fmt.Printf("my PID=%d, kill -SIGUSR1 %d if you want\n\n", syscall.Getpid(), syscall.Getpid())
	time.Sleep(delaySeconds * time.Second)
	fmt.Print("waited 10, sending SIGUSR1, waiting another 10 seconds\n")
	syscall.Kill(syscall.Getpid(), syscall.SIGUSR1)
	time.Sleep(delaySeconds * time.Second)
	fmt.Print("Done waiting, goodbye\n")
}

var (
	iAmAnImportantNumber = expvar.NewInt("iAmAnImportNumber")
	totalClients         = expvar.NewInt("totalClients")
	totalRequests        = expvar.NewInt("totalRequests")
)
