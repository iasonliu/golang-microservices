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

// swagger:route DELETE /products/{id} products deleteProduct
// Returns a list of products
// responses:
//	201: noContentResponse

// DeleteProducts deletes a product from the database
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

// swagger:route POST /products products createProduct
// Create a new product
//
// responses:
//	201: productResponse
//  422: errorValidation
//  501: errorResponse

// Create handles POST requests to add new products
func (p *Products) AddProduct(w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST add Product")
	product := r.Context().Value(KeyProduct{}).(data.Product)

	p.l.Printf("[DEBUG] Inserting Product: %#v\n", product)
	data.AddProduct(&product)
}

// swagger:route PUT /products products updateProduct
// Update a products details
//
// responses:
//	201: noContentResponse
//  404: errorResponse
//  422: errorValidation
// Update handles PUT requests to update products
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
