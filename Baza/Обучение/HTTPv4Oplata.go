package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Peyment struct {
	Opisanie string `json:"opisanie"`
	USD      int    `json:"usd"`
	FullName string `json:"fullName"`
	Addres   string `json:"addres"`
	Time     time.Time
}

func (p Peyment) Println() {
	fmt.Println("Описание:", p.Opisanie)
	fmt.Println("Валюта:", p.USD)
	fmt.Println("Полное имя:", p.FullName)
	fmt.Println("Адрес:", p.Addres)
}

func (p Peyment) Validate() bool {
	if p.USD == 0 {
		return false
	}
	if p.Addres == "" {
		return false
	}
	return true

}

var mtx = sync.Mutex{}
var money = 2000
var peymentHistory = make([]Peyment, 0)

type HttpOtvet struct {
	Money          int       `json:"bbbabkki"`
	PeymentHistory []Peyment `json:"история оплаты"`
}

func payHandler(w http.ResponseWriter, r *http.Request) {
	var oplata Peyment
	if err := json.NewDecoder(r.Body).Decode(&oplata); err != nil {
		fmt.Println("Error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return

	}
	oplata.Time = time.Now()
	oplata.Println()

	mtx.Lock()
	if money-oplata.USD >= 0 {
		money -= oplata.USD
	}
	peymentHistory = append(peymentHistory, oplata)

	httpOtvet := HttpOtvet{
		Money:          money,
		PeymentHistory: peymentHistory,
	}

	x, err := json.Marshal(httpOtvet)
	if err != nil {
		fmt.Println("errror:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return

	}
	if _, err := w.Write(x); err != nil {
		fmt.Println("err", err)
		return
	}

	mtx.Unlock()

}
func main() {
	http.HandleFunc("/pay", payHandler)

	if err := http.ListenAndServe(":228", nil); err != nil {
		fmt.Println("Ошибка во  время ворка сервера исправляй", err)
	}

}
