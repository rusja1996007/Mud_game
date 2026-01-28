package main

import (
	"fmt"
	"time"
)

func mine(punktPeredachi chan int, n int) {
	fmt.Println("Попер в шахту в", n, "раз")
	time.Sleep(2 * time.Second)
	fmt.Println("Поход номер", n, "закончен")
	punktPeredachi <- 30
	fmt.Println("Поход номер", n, "передал уголь")
}

/*func main() {
	SkladUglya := 0
	punktPeredachi := make(chan int, 3)
	t := time.Now()

	go mine(punktPeredachi, 1)
	go mine(punktPeredachi, 2)
	go mine(punktPeredachi, 3)

	SkladUglya += <-punktPeredachi

	SkladUglya += <-punktPeredachi

	SkladUglya += <-punktPeredachi

	SkladUglya += <-punktPeredachi

	fmt.Println("Всего угля:", SkladUglya)
	fmt.Println("Прошло времения:", time.Since(t))

}*/
