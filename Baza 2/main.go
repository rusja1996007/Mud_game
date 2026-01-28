package main

import (
	"fmt"
	"restapi/http"
	"restapi/todo"
)

func main() {
	tdList := todo.NewList()
	httpHandler := http.NewHTTPHandlers(tdList)
	httpServer := http.NewHTTPServer(httpHandler)

	if err := httpServer.StartServer(); err != nil {
		fmt.Println("failed to start http server", err)
	}

}
