package main

import (
	"fmt"
	"strings"
)

func main() {
	var i int
	a := strings.Split("a b c d e f", " ")
	fmt.Println(a)
	for i < rangelen(a[0]) > 0 {
		a = a[1:]
		fmt.Println(a)
	}
}
