package main

import (
	"sync/atomic"
)

var number atomic.Int64

/*func uvelichit(wg *sync.WaitGroup) {
	defer wg.Done()
	for x := 1; x <= 1000; x++ {
		number.Add(1)

	}
}*/

/*func main() {
	wg := &sync.WaitGroup{}
	wg.Add(10)
	go uvelichit(wg)
	go uvelichit(wg)
	go uvelichit(wg)
	go uvelichit(wg)
	go uvelichit(wg)

	go uvelichit(wg)
	go uvelichit(wg)
	go uvelichit(wg)
	go uvelichit(wg)
	go uvelichit(wg)

	wg.Wait()

	fmt.Println("number", number.Load())

}*/
