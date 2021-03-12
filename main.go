package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
	"github.com/iasonliu/product-api/handlers"
)

func main() {
	logger := log.New(os.Stdout, "[product-api] ", log.LstdFlags)
	// create the handers
	ph := handlers.NewProducts(logger)

	// create a new server mux and register hanndlers

	sm := mux.NewRouter()
	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/products", ph.GetProducts)
	getRouter.Use(ph.LoggingMiddleware)

	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/products/{id:[0-9]+}", ph.UpdateProducts)
	putRouter.Use(ph.MiddlewareProductValidation)

	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/products", ph.AddProduct)
	postRouter.Use(ph.MiddlewareProductValidation)

	deleteRouter := sm.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/{id:[0-9]+}", ph.DeleteProduct)

	// * handler for documentation
	opts := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.Redoc(opts, nil)
	getRouter.Handle("/docs", sh)
	getRouter.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	// config http server
	httpServer := &http.Server{
		Addr:         ":9090",
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// start the server
	go func() {
		logger.Println("Starting server Listen on", httpServer.Addr)
		err := httpServer.ListenAndServe()
		if err != nil {
			logger.Fatalln(err)
		}
	}()

	// trap sigterm or interupt and gracefully shotdown the http server
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)
	sig := <-sigChan
	logger.Println("Recieved terminate, graceful shutdown", sig)

	timeoutCtx, err := context.WithTimeout(context.Background(), 30*time.Second)
	if err != nil {
		logger.Println(err)
	}
	logger.Fatal(httpServer.Shutdown(timeoutCtx))
}
