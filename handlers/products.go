// Package classification of Product API
//
// Documentation for Product API
//
//	Schemes: http
//	BasePath: /
//	Version: 1.0.0
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
// swagger:meta
package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/iasonliu/product-api/data"
)

// A list of products returns in the response
// swagger:response productsResponse
type productsResponse struct {
	// All products in the system
	// in: body
	Body []data.Products
}

// Products is a http.Handler
type Products struct {
	l *log.Logger
}

// NewProducts creates a products handler with the given logger
func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

// swagger:route GET /products products listProducts
// Returns a list of products
// responses:
//	200: productsResponse

// GetProducts returns the products from the data store
func (p Products) GetProducts(w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle GET Products")
	// fetch the products from the data
	lp := data.GetProducts()
	// serialize the list to JSON
	err := lp.ToJSON(w)
	if err != nil {
		http.Error(w, "Unable to marshal json", http.StatusInternalServerError)
		return
	}
}

func (p *Products) AddProduct(w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST add Product")
	product := r.Context().Value(KeyProduct{}).(data.Product)

	p.l.Printf("Product: %#v", product)
	data.AddProduct(&product)
}

func (p Products) UpdateProducts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Unable to covert id", http.StatusBadRequest)
		return
	}
	p.l.Println("Handle PUT update Product", id)
	product := r.Context().Value(KeyProduct{}).(data.Product)

	err = data.UpdateProduct(id, &product)
	if err == data.ErrProductNotFound {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type KeyProduct struct{}

func (p Products) MiddlewareProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		product := data.Product{}
		err := product.FromJSON(r.Body)
		if err != nil {
			p.l.Println("[ERROR] deserializing product", err)
			http.Error(w, "Error reading product", http.StatusBadRequest)
			return
		}
		// validate the product
		err = product.Validate()
		if err != nil {
			p.l.Println("[ERROR] validate product", err)
			http.Error(
				w,
				fmt.Sprintf("Error reading product: %s", err),
				http.StatusBadRequest,
			)
			return
		}
		// add the product to the context
		ctx := context.WithValue(r.Context(), KeyProduct{}, product)
		r = r.WithContext(ctx)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func (p Products) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println(r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}
