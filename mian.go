package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func greet(w http.ResponseWriter, r *http.Request) {
	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Ooops!!!", http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "Hello %s", bs)
}

func goodbye(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "goodby!!")
}

func main() {
	http.HandleFunc("/", greet)
	http.HandleFunc("/goodbye", goodbye)
	http.ListenAndServe(":8080", nil)
}
