package main

import (
	"fmt"
	"log"
	"net/http"
	"umbrellacorp/handlers"
	"umbrellacorp/router"
)

func main() {
	initialize()
	fmt.Printf("\nStarting Server\n")
	log.Fatal(http.ListenAndServe(":8080", router.NewRouter()))
}

func initialize() {
	handlers.Init()
}
