package postman

import (
	"context"
	"sync"
	"time"

	"github.com/k0kubun/pp"
)

func postman(ctx context.Context, wg *sync.WaitGroup, transferPoint chan<- string, n int, mail string) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			pp.Println(" Я почтальон номер:", n, "я пошел спать я заибался нахуй уже!")
			return
		default:

			pp.Println("Я почтальон номер:", n, "Взял письмо")
			time.Sleep(1 * time.Second)
			pp.Println("Я почтальон номер:", n, "Донес письмо: ", mail, "до почты")

			transferPoint <- mail
			pp.Println("Я почтальон номер:", n, "Передал письмо:", mail, "на почту")
		}
	}
}

func PostmanPool(ctx context.Context, vsegoPostmanov int) <-chan string {
	OtdeleniePoshti := make(chan string)

	wg := &sync.WaitGroup{}

	for x := 1; x <= vsegoPostmanov; x++ {
		wg.Add(1)
		go postman(ctx, wg, OtdeleniePoshti, x, kakoyMailUPostmana(x))
	}
	go func() {
		wg.Wait()
		close(OtdeleniePoshti)
	}()
	return OtdeleniePoshti
}

func kakoyMailUPostmana(NomerPostmana int) string {
	kmup := map[int]string{
		1: "Приглашение на свингер вечеринку",
		2: "Налоговая",
		3: "Реклама трусов",
	}
	mail, ok := kmup[NomerPostmana]
	if !ok {
		return "СПАМЩИК ХУЕВ"
	}
	return mail
}
