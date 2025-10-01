package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Welcome to chirpy!")

	serveMux := http.NewServeMux()
	serveMux.Handle("/", http.FileServer(http.Dir(".")))

	server := http.Server{
		Handler: serveMux,
		Addr:    ":8080",
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
