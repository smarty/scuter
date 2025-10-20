package main

import (
	"log"
	"net/http"
	"os"

	"github.com/smarty/scuter/example/internal/app"
	HTTP "github.com/smarty/scuter/example/internal/http"
)

func main() {
	address := "localhost:8080"
	logger := log.New(os.Stderr, "", log.LstdFlags)
	logger.Printf("listing on http://%s", address)
	err := http.ListenAndServe(address, HTTP.New(logger, new(app.Application)))
	if err != nil {
		logger.Panic(err)
	}
}
