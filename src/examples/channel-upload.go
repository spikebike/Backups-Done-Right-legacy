package main

import (
	"fmt"
	"time"
)

func server(UpChan chan string, done chan bool) {
	for f := range UpChan {
		fmt.Printf("Server: received %s\n", f)
	}
	fmt.Print("Server: Channel closed, existing\n")
}

func client(UpChan chan string, done chan bool) {
	for i := 0; i < 10; i++ {
		fmt.Print("Client: sending file\n")
		UpChan <- "/home/bill/test"
		time.Sleep(400 * time.Millisecond)
	}
	fmt.Print("Client: announcing done and existing\n")
	done <- true
}

func main() {
	done := make(chan bool)
	UpChan := make(chan string)
	go server(UpChan, done)
	go client(UpChan, done)
	<-done
	fmt.Print("Main: Client finished closing server\n")
	close(UpChan)
}
