/*package main

import (
	"sync"
)

var likes int = 0
var mtx sync.RWMutex

func dobavitLike(wg *sync.WaitGroup) {
	defer wg.Done()
	for x := 1; x <= 100_000; x++ {
		mtx.Lock()
		likes++
		mtx.Unlock()
	}
}

func poluchitLike(wg *sync.WaitGroup) {
	defer wg.Done()
	for x := 1; x <= 100_000; x++ {
		mtx.RLock()
		_ = likes
		mtx.RUnlock()
	}
}

func main() {
	wg := &sync.WaitGroup{}
	initTime := time.Now()

	for x := 1; x <= 10; x++ {
		wg.Add(1)
		go dobavitLike(wg)
	}
	for x := 1; x <= 10; x++ {
		wg.Add(1)
		go poluchitLike(wg)
	}
	wg.Wait()

	fmt.Println("Время выполнения:", time.Since(initTime))

}*/
