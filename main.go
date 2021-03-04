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

	serveMux := http.NewServeMux()
	serveMux.Handle("/", helloHandler)
	http.ListenAndServe(":8080", serveMux)
}
