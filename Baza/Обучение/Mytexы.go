package main

/*var slice []int
var mtx sync.Mutex

func uvelichit(wg *sync.WaitGroup) {
	defer wg.Done()
	for x := 1; x <= 1000; x++ {
		mtx.Lock()
		slice = append(slice, x)
		mtx.Unlock()
	}
}

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

	fmt.Println("slice len:", len(slice))

}
*/
