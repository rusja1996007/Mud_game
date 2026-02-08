package main

import (
	"fmt"
	"study/http_server"
)

func main() {
	fmt.Println("http server started")
	err := http_server.StartHttpServer()
	if err != nil {
		fmt.Println("error http server:", err)
	} else {
		fmt.Println("http stopped.")
	}
}
9:36:49