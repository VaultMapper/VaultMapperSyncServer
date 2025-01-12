package server

import (
	"fmt"
	"log"
	"net/http"
)

// HUB stores all Vault structures
var HUB = Hub{}

func Run(ip string, port int) {
	fmt.Println("HELLO FROM SERVER")

	http.HandleFunc("/", handshakeHandler)
	http.Handle("/stats", rateLimit(http.HandlerFunc(statsHandler)))
	if err := http.ListenAndServe(fmt.Sprintf("%s:%d", ip, port), nil); err != nil {
		log.Fatal(err)
	}
}
