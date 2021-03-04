package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Goodbye struct {
	l *log.Logger
}

func NewGoodye(l *log.Logger) *Goodbye {
	return &Goodbye{l}
}

func (h *Goodbye) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.l.Println("Goodbye~!!!")
	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Ooops!!!", http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "Goodbye %s", bs)
}
