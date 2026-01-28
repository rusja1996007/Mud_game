package main

/*func main() {

	strChan := make(chan string)
	intChan := make(chan int)
	go func() {
		time.Sleep(100 * time.Millisecond)
		intChan <- 1

	}()
	go func() {

		time.Sleep(200 * time.Millisecond)
		strChan <- "privet"

	}()
	time.Sleep(50 * time.Millisecond)
	select {
	case number := <-intChan:
		fmt.Println("числа:", number)

	case stroka := <-strChan:
		fmt.Println("строка:", stroka)
	default:
		fmt.Println("Никакой канал не готов")
	}
}
*/
