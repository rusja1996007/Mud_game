package miner

import (
	"context"
	"sync"
	"time"

	"github.com/k0kubun/pp"
)

func miner(ctx context.Context,
	wg *sync.WaitGroup,
	transferPoint chan<- int,
	n int,
	power int) {
	defer wg.Done()
	for {
		pp.Println("Я шахтер номер:", n, "пошел ебашить")
		select {
		case <-ctx.Done():
			pp.Println("Я шахтер номер:", n, "Пошел отдыхать я заибался!")
			return
		case <-time.After(1 * time.Second):
			pp.Println("Я шахтер номер:", n, "добыл тебе брат:", power, "угля")
		}
		select {
		case <-ctx.Done():
			pp.Println("Я шахтер номер:", n, "Пошел отдыхать я заибался!")
			return

		case transferPoint <- power:
			pp.Println("Я шахтер номер:", n, " передал ", power, "угля")
		}

	}
}

func MinerPool(ctx context.Context, vsegoMinerov int) <-chan int {
	skladUglya := make(chan int)

	wg := &sync.WaitGroup{}
	for x := 1; x <= vsegoMinerov; x++ {
		wg.Add(1)
		go miner(ctx, wg, skladUglya, x, x*10)

	}
	go func() {
		wg.Wait()
		close(skladUglya)
	}()

	return skladUglya
}
