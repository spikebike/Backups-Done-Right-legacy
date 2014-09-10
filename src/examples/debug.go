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

func TestStress(sig chan os.Signal) {
	for ss := range sig {
		fmt.Printf("received: %s\n", ss)
		fmt.Printf("iAmAnImportantNumber=%s\n", iAmAnImportantNumber.String())
	}
}

func main() {
	go http.ListenAndServe(fmt.Sprintf(":%d", httpPortTCP), nil)
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGUSR1)
	go TestStress(sig)

	fmt.Printf("Go to http://localhost:%d/debug/vars in the next %d seconds \n", httpPortTCP, delaySeconds)
	fmt.Printf("my PID=%d, kill -SIGUSR1 %d if you want\n", syscall.Getpid(), syscall.Getpid())
	time.Sleep(delaySeconds * time.Second)
	fmt.Print("waited 10, sending SIGUSR1, waiting another 10 seconds\n")
	syscall.Kill(syscall.Getpid(), syscall.SIGUSR1)
	time.Sleep(delaySeconds * time.Second)
}

var (
	iAmAnImportantNumber = expvar.NewInt("iAmAnImportNumber")
	totalClients         = expvar.NewInt("totalClients")
	totalRequsts         = expvar.NewInt("totalRequests")
)
