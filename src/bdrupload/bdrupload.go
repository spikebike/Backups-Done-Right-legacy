package bdrupload

import (
	"fmt"
)

type Upchan_t struct {
        Rowid int64
        Path string
}

type Downchan_t struct {
        Rowid int
        Err error
}

func HelloWorld() {
	fmt.Println("Hello World!")
}

func main_test() {
	HelloWorld()
}
