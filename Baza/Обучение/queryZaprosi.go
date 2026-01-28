package main

//localhost:228,default?foo=x&boo=y
/*func handler(w http.ResponseWriter, r *http.Request) {
	fooParam := r.URL.Query().Get("foo")
	booParam := r.URL.Query().Get("boo")

	fmt.Println("foo param:", fooParam)
	fmt.Println("boo param:", booParam)
}

func main() {
	http.HandleFunc("/default", handler)
	if err := http.ListenAndServe(":228", nil); err != nil {
		fmt.Println("failed tun server:", err)
	}
}*/
