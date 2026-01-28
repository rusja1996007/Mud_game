package main

import (
	"fmt"
	"sync"
	"time"
)

func postman(wg *sync.WaitGroup, text string) {
	defer wg.Done()
	for x := 1; x <= 3; x++ {
		fmt.Println("Почтальон отнес газету с темой:", text, "в", x, "раз")
		time.Sleep(250 * time.Millisecond)
	}
}

/*func main() {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go postman(wg, "порноиндустрия")

	wg.Add(1)
	go postman(wg, "новости: геи ли пакестанцы?")

	wg.Add(1)
	go postman(wg, "Игромания")

	wg.Wait()

	fmt.Println("main кончил")

}
*/
