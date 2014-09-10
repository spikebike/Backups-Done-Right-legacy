package main

import "fmt"

// return a function
func Sequence() func() int {
	n := 0
	return func() int {
		n++
		return n
	}
}

// take a funtion as argument
func Print(f func() int) {
	fmt.Printf("%d\n", f())
}

func main() {
	seq := Sequence()

	fmt.Printf("%d\n", seq())
	fmt.Printf("%d\n", seq())

	Print(seq)
}
