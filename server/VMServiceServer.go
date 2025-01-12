package server

import (
	"fmt"
	"golang.org/x/time/rate"
	"log"
	"net/http"
	"sync"
)

// HUB stores all Vault structures
var HUB = Hub{}

var (
	statsLimiter = rate.NewLimiter(5, 10) // 1 request per second with a burst of 5
	statsMu      sync.Mutex
)

func statsRateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		statsMu.Lock()
		defer statsMu.Unlock()

		if !statsLimiter.Allow() {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func Run(ip string, port int) {
	fmt.Println("HELLO FROM SERVER")

	http.HandleFunc("/", handshakeHandler)
	http.Handle("/stats", statsRateLimit(http.HandlerFunc(statsHandler)))
	if err := http.ListenAndServe(fmt.Sprintf("%s:%d", ip, port), nil); err != nil {
		log.Fatal(err)
	}
}
