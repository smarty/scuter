package main

import (
	"log"
	"net/http"

	"github.com/mdw-go/scuter/example/internal/app"
	HTTP "github.com/mdw-go/scuter/example/internal/http"
)

func main() {
	application := new(app.Application)
	http.Handle("PUT    /tasks", HTTP.NewCreateTaskShell(application))
	http.Handle("DELETE /tasks", HTTP.NewDeleteTaskShell(application))
	log.Println("listing on http://localhost:8080/")
	_ = http.ListenAndServe(":8080", nil)
}
