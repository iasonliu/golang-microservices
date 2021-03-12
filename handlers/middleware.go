package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/iasonliu/product-api/data"
)

type KeyProduct struct{}

// MiddlewareValidateProduct validates the product in the request and calls next if ok
func (p *Products) MiddlewareProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		product := &data.Product{}
		err := data.FromJSON(product, r.Body)
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

func (p *Products) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println(r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}
