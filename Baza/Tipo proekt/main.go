package main

import (
	"concurrency/miner"
	"concurrency/postman"
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/k0kubun/pp"
)

func main() {
	var vsegoUglia atomic.Int64
	mtx := sync.Mutex{}
	var vsePisma []string

	minerContext, minerCancel := context.WithCancel(context.Background())
	postmanContext, postmanCancel := context.WithCancel(context.Background())

	go func() {
		time.Sleep(3 * time.Second)
		fmt.Println("ШАХТЕРЫ КОНЧИЛИ!!!")
		minerCancel()
	}()

	go func() {
		time.Sleep(6 * time.Second)
		fmt.Println("ПОЧТАЛЬОНЫ КОНЧИЛИ!")
		postmanCancel()
	}()

	skladUglya := miner.MinerPool(minerContext, 3)
	OtdeleniePoshti := postman.PostmanPool(postmanContext, 3)

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for v := range skladUglya {
			vsegoUglia.Add(int64(v))
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for v := range OtdeleniePoshti {
			mtx.Lock()
			vsePisma = append(vsePisma, v)
			mtx.Unlock()
		}
	}()

	wg.Wait()

	pp.Println("Всего собрано угля:", vsegoUglia.Load())

	mtx.Lock()
	pp.Println("Всего притащили писем:", len(vsePisma))
	mtx.Unlock()

}
