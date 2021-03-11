package handlers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/iasonliu/product-api/data"
)

func (p *Products) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	// this will always covert because of the router
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	p.l.Println("Handle DELETE Product", id)
	err := data.DeleteProduct(id)
	if err == data.ErrProductNotFound {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Product not found", http.StatusInternalServerError)
		return
	}
}
