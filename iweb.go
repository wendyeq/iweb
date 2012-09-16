package main

import (
	"log"
	"net/http"
)

func main() {
	err := http.ListenAndServe(":"+GetConfig()["port"], nil)
	log.Println(err)
}

func init() {
	InitRouter()
}
