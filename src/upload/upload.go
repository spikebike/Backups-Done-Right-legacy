package upload

import (
		"fmt"
		"../mystructs"

)

func Server(upchan chan *mystructs.Upchan_t) {
	for f := range upchan {
		fmt.Printf("Server: received rowID=%d path=%s\n", f.Rowid, f.Path)
		//      fmt.Printf("%T %#v\n",f,f)
	}
	fmt.Print("Server: Channel closed, existing\n")
}
