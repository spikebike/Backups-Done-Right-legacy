package main

import (
		"fmt"
		"math/rand"
		"time"
)

func main() {

	rand.Seed(time.Now().UnixNano())
	s := make ([]int, 64)
	for i := range s {
		s[i]=i
	}
	fmt.Printf("before shuffle = %v\n",s)
	
	for i := range s {
		j := rand.Intn(len(s) - i)
		s[i], s[j+i] = s[j+i], s[i]
	}
	fmt.Printf("post shuffle   = %v\n",s)
}
