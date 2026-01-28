package main

import (
	"fmt"
	"time"
)

func foo() {
	for {
		fmt.Println("zxc")
		time.Sleep(300 * time.Millisecond)
	}

}

/*func main() {
	go foo()

	go func() {
		for {
			fmt.Println("anonimus")
			time.Sleep(100 * time.Millisecond)
		}
	}()

	time.Sleep(1 * time.Second)
}
*/
