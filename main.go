package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/iasonliu/working/handlers"
)

func main() {
	logger := log.New(os.Stdout, "product-api", log.LstdFlags)
	helloHandler := handlers.NewHello(logger)
	goodbyeHandler := handlers.NewGoodye(logger)

	serveMux := http.NewServeMux()
	serveMux.Handle("/", helloHandler)
	serveMux.Handle("/bye", goodbyeHandler)
	httpServer := &http.Server{
		Addr:         ":9090",
		Handler:      serveMux,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		err := httpServer.ListenAndServe()
		if err != nil {
			logger.Fatalln(err)
		}
	}()
	sigChan := make(chan os.Signal)
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
