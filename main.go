package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/iasonliu/product-api/handlers"
)

func main() {
	logger := log.New(os.Stdout, "product-api", log.LstdFlags)
	// create the handers
	productsHandler := handlers.NewProducts(logger)

	// create a new server mux and register hanndlers
	serveMux := http.NewServeMux()
	serveMux.Handle("/", productsHandler)

	// config http server
	httpServer := &http.Server{
		Addr:         ":9090",
		Handler:      serveMux,
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
