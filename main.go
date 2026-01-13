package main

import (
	"fmt"
	"study/feature1"
	"study/feature2"
	"study/feature_postgres/simple_connect"
)

func main() {
	fmt.Println("Hello gei!")
	feature1.Feauture1()
	feature2.Feature2()
	simple_connect.CheckConnect()
}
