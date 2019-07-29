package main

import (
	"log"
	"net/http"
)

func startServer() {
	rp := NewReverseProxyPool()
	err := http.ListenAndServe("127.0.0.1:8888", rp)
	if err != nil {
		log.Fatalln("ListenAndServe: ", err)
	}
}

func main() {
	startServer()
}
