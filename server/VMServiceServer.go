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
	statsLimiter  = rate.NewLimiter(5, 10) // 1 request per second with a burst of 5
	statsMu       sync.Mutex
	viewerLimiter = rate.NewLimiter(1, 1) // 1 request per second with a burst of 5
	viewerMu      sync.Mutex
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

func viewerRateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		viewerMu.Lock()
		defer viewerMu.Unlock()

		if !viewerLimiter.Allow() {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func Run(ip string, port int) {
	fmt.Println("HELLO FROM SERVER")

	// start up pps service
	PPSInit()

	http.HandleFunc("/", handshakeHandler)
	http.Handle("/stats", statsRateLimit(http.HandlerFunc(statsHandler)))
	http.Handle("/view", viewerRateLimit(http.HandlerFunc(viewerHandshakeHandler)))
	if err := http.ListenAndServe(fmt.Sprintf("%s:%d", ip, port), nil); err != nil {
		log.Fatal(err)
	}
}
