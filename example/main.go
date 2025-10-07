package main

import (
	"log"
	"net/http"
	"os"

	"github.com/mdw-go/scuter/example/internal/app"
	HTTP "github.com/mdw-go/scuter/example/internal/http"
)

func main() {
	logger := log.New(os.Stderr, "", log.LstdFlags)
	application := new(app.Application)
	http.Handle("PUT    /tasks", HTTP.NewCreateTaskShell(logger, application))
	http.Handle("DELETE /tasks", HTTP.NewDeleteTaskShell(logger, application))
	logger.Println("listing on http://localhost:8080/")
	_ = http.ListenAndServe(":8080", nil)
}
