package main

import (
	"log"
	"net/http"
)

const (
	siteDir = "./../site"
	addr    = ":8080"
)

func main() {
	log.Println("server will start at", addr)
	log.Fatal(http.ListenAndServe(addr, http.FileServer(http.Dir(siteDir))))
}
