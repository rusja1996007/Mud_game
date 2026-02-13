package http_server

import (
	"errors"
	"fmt"
	"net/http"
)

func StartHttpServer() error {
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(" обработка запроса на паттерне /ping")
		w.Write([]byte("Hello from DOcker\n"))
	})

	err := http.ListenAndServe(":5050", nil)
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	} else {
		return err
	}
}
