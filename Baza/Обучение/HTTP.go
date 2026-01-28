/*package main

import (
	"fmt"
	"net/http"
)

func payHandler(w http.ResponseWriter, r *http.Request) {
	str := "Оплата прошла успешно изи бабки"
	x := []byte(str)
	_, err := w.Write(x)
	if err != nil {
		fmt.Println("Во время записи ответа произошла эрор:", err.Error())
	} else {
		fmt.Println("я спокойно совершил оплату")
	}
}

func cancelHandler(w http.ResponseWriter, r *http.Request) {
	str := "Возвращаю бабки за куни"
	x := []byte(str)
	_, err := w.Write(x)
	if err != nil {
		fmt.Println("Во время записи ответа произошла эрор:", err.Error())
	} else {
		fmt.Println("я спокойно отменил оплату шлюхи")
	}

}

func handler(w http.ResponseWriter, r *http.Request) {
	str := "Hello gei!"
	x := []byte(str)
	_, err := w.Write(x)
	if err != nil {
		fmt.Println("Во время записи ответа произошла эрор:", err.Error())
	} else {
		fmt.Println("я норм обработал запрос")
	}
}
*/
/*func main() {
	http.HandleFunc("/default", handler)
	http.HandleFunc("/pay", payHandler)
	http.HandleFunc("/cancel", cancelHandler)

	fmt.Println("запуск сервера летс ГООООУ")

	err := http.ListenAndServe(":228", nil)
	if err != nil {
		fmt.Println("какая то ошибка пиздеееец", err.Error())
	}
	fmt.Println("Прога кончила ты скорострел")

}*/
