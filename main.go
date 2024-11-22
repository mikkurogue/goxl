package main

import (
	"fmt"
	"goxl/bootstrap"
	"net/http"
)

func main() {
	mux := bootstrap.Bootstrap()

	fmt.Println("Server running on http://localhost:9292")
	err := http.ListenAndServe(":9292", mux)
	if err != nil {
		fmt.Println("Error starting server", err)
	}
}
