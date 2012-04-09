package main

import (
	"fmt"
	"strings"
)

func main() {
	var i int
	var s string
	i = 0
	a := strings.Split("a b c d e f", " ")
	fmt.Printf("len = %d\n", len(a))
	for i < len(a) {
		s = a[i]
		fmt.Println(s)
		fmt.Println(a)
		if i == 3 {
			a = append(a, "g")
		}
		i = i + 1
	}
}
