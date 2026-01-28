package main

type Message struct {
	Aftor string
	Text  string
}

/*func main() {
	messageChan1 := make(chan Message)
	messageChan2 := make(chan Message)

	go func() {
		for {
			messageChan1 <- Message{
				Aftor: "Володя гей",
				Text:  "Понюхаем кокаин!?",
			}
			time.Sleep(10000 * time.Millisecond)
		}

	}()
	go func() {
		for {

			messageChan2 <- Message{
				Aftor: "Ирина натурал",
				Text:  "Понюхаем соли!?",
			}
			time.Sleep(200 * time.Millisecond)
		}

	}()
	for {
		select {
		case msg1 := <-messageChan1:
			fmt.Println("Получено сообщение от :", msg1.Aftor, "-", msg1.Text)

		case msg2 := <-messageChan2:
			fmt.Println("Получено сообщение от :", msg2.Aftor, "-", msg2.Text)

		}
	}
}*/
