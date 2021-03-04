package main

import (
	"log"
	"net/http"
	"os"

	"github.com/iasonliu/working/handlers"
)

func main() {
	logger := log.New(os.Stdout, "product-api", log.LstdFlags)
	helloHandler := handlers.NewHello(logger)
	goodbyeHandler := handlers.NewGoodye(logger)

	serveMux := http.NewServeMux()
	serveMux.Handle("/", helloHandler)
	serveMux.Handle("/bye", goodbyeHandler)
	http.ListenAndServe(":8080", serveMux)
}
