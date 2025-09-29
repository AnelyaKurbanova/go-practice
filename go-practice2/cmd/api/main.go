package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/anelya/golang_lab2/internal/handlers"
	"github.com/anelya/golang_lab2/internal/middleware"
)

func main() {
	mux := http.NewServeMux()

	mux.Handle("/user", middleware.Auth(http.HandlerFunc(handlers.UserHandler)))

	addr := ":8080"
	fmt.Println("Сервер запущен на http://localhost" + addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
