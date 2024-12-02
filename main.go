package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/nurovic/hmall/api"
)

func main() {
	router := api.NewRouter()

	port := "8080"
	fmt.Printf("Sunucu %s portunda çalışıyor...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
