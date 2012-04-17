package main

import (
	"crypto/rand"
	"fmt"
)

func main() {
	b := make([]byte, 16)
	n, err := rand.Read(b)
	fmt.Println(n, err, b)
}
