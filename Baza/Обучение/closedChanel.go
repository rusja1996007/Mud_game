package main

/*func main() {
	transPoint := make(chan int)

	go func() {
		//3 4 5 6(3 + 0 или 1 или 2 или 3)
		iterat := 3 + rand.Intn(4)
		fmt.Println("итерация:", iterat)
		for x := 1; x <= iterat; x++ {
			transPoint <- 10
			time.Sleep(1 * time.Second)
		}
		close(transPoint)
	}()

	skladUglya := 0

	for y := range transPoint {
		skladUglya += y
		fmt.Println("добавили в склад:", skladUglya)
	}

	fmt.Println("на складе всего стало :", skladUglya)

}*/
