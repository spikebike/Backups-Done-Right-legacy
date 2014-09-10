package main

import "fmt"
import "time"

func getID(id chan int64) {
	var counter int64
	counter = 0
	for {
		id <- counter
		counter = counter + 1
	}
}

func consumer(id chan int64, done chan bool) {
	var sum int64
	sum = 0
	t0 := time.Now().UnixNano()
	for i := 0; i < 100000000; i++ {
		<-id
		sum = sum + 1
	}
	t1 := time.Now().UnixNano()
	fmt.Println("time: ", float64(t1-t0)/1000000000, "seconds  sum: ", sum)
	done <- true
}

func main() {
	id := make(chan int64, 256)
	done := make(chan bool, 1)
	go getID(id)
	go consumer(id, done)
	go consumer(id, done)

	// Block till goroutines finish
	<-done
	<-done
}
