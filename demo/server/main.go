package main

import (
	"log"
	"net/http"
)

const siteDir = "./../site"

func main() {
	log.Fatal(http.ListenAndServe(":8080", http.FileServer(http.Dir(siteDir))))
}
