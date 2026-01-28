package main

import (
	"fmt"
	"net/http"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Hello geeei!"))
	if err != nil {
		fmt.Println("Ошибка:", err.Error())
		return

	}
	fmt.Println("Handler отьебашил заебись")

}
func handlerSleep(w http.ResponseWriter, r *http.Request) {
	time.Sleep(5 * time.Second)
	_, err := w.Write([]byte("Http RESPONSE"))
	if err != nil {
		fmt.Println("Ошибка:", err.Error())
		return

	}

}

/*func main() {

	http.HandleFunc("/", handler)
	http.HandleFunc("/sleep", handlerSleep)

	err := http.ListenAndServe(":228", nil)
	if err != nil {
		fmt.Println("Error:", err.Error())
	}
}*/
