package main

import (
	"context"
	"fmt"
	"time"
)

// родительский контекст
func fuu(ctx context.Context, n int) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("fUU завершилась чтото 1000-7", n)
			return
		default:
			fmt.Println("fUUUUUUUU ДАльше ебашит", n)

		}
		time.Sleep(100 * time.Millisecond)
	}

}

// дочерний контекст
func boo(ctx context.Context, n int) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("bOO завершилась чтото 1000-7", n)
			return
		default:
			fmt.Println("bOOOOOOO ДАльше ебашит", n)

		}
		time.Sleep(100 * time.Millisecond)
	}
}

/*func main() {
	prntContext, prntCancel := context.WithCancel(context.Background())
	childContext, childCancel := context.WithCancel(prntContext)
	go fuu(prntContext, 1)
	go fuu(prntContext, 2)
	go fuu(prntContext, 3)

	go boo(childContext, 1)
	go boo(childContext, 2)
	go boo(childContext, 3)

	time.Sleep(1 * time.Second)
	childCancel()

	time.Sleep(1 * time.Second)
	prntCancel()

	time.Sleep(3 * time.Second)
	fmt.Println("Main завершен")

}*/
