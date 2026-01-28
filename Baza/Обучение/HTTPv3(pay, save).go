package main

//var mtx = sync.Mutex{}
var babki = 1000
var bank = 0

/*func payHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	fmt.Println("http method :", r.Method)

	httpRequestBody, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		msg := "error read HTTP body:" + err.Error()
		fmt.Println(msg)
		_, err = w.Write([]byte(msg))
		if err != nil {
			fmt.Println("fail write HTTP response:", err)

		}
		return
	}
	httpRequestBodyString := string(httpRequestBody)
	paymentAmount, err := strconv.Atoi(httpRequestBodyString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		msg := "Eror convert http body to integer:" + err.Error()
		fmt.Println(msg)
		_, err = w.Write([]byte(msg))
		if err != nil {
			msg := "fail write HTTP response" + err.Error()
			fmt.Println(msg)
			_, err = w.Write([]byte(msg))
			if err != nil {
				fmt.Println("fail write HTTP response:", err)

			}

		}

		return
	}
	mtx.Lock()
	if babki-paymentAmount >= 0 {
		babki -= paymentAmount
		msg := "оплата прошла все гуд" + strconv.Itoa(babki)
		fmt.Println(msg)
		_, err = w.Write([]byte(msg))
		if err != nil {
			fmt.Println("fail write HTTP response:", err)

		}

	}
	mtx.Unlock()
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	httpRequestBody, err := io.ReadAll(r.Body)
	if err != nil {
		msg := "error read HTTP body:" + err.Error()
		fmt.Println(msg)
		_, err = w.Write([]byte(msg))
		if err != nil {
			fmt.Println("fail write HTTP response:", err)

		}
		return
	}
	httpRequestBodyString := string(httpRequestBody)
	saveAmount, err := strconv.Atoi(httpRequestBodyString)
	if err != nil {
		msg := "fail write HTTP response:" + err.Error()
		fmt.Println(msg)
		_, err = w.Write([]byte(msg))
		if err != nil {
			fmt.Println("fail write HTTP response:", err)

		}

		return
	}
	mtx.Lock()
	if babki >= saveAmount {
		babki -= saveAmount
		bank += saveAmount
		fmt.Println("новое значение переменой babki:", babki)
		fmt.Println("новое значение переменой bank:", bank)

	}
	mtx.Unlock()
}

func main() {

	http.HandleFunc("/pay", payHandler)
	http.HandleFunc("/save", saveHandler)

	err := http.ListenAndServe(":228", nil)
	if err != nil {
		fmt.Println("http server error:", err.Error())
	}

}
*/
