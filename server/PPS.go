package server

import (
	"log"
	"sync/atomic"
	"time"
)

// this file implements a service that keeps track of packet per second statistic metric

var (
	packetCount       uint64
	packetCounterChan = make(chan bool, 500) // buffered channel for packet counting, used to hopefully bypass slowdowns with overusing mutexes
	pps               uint64                 // stores the actual packet per second metric
	maxPPS            uint64                 // runtime max of pps
)

func PPSInit() {
	go processPacketCounter()
	go calculatePPS()
}

// processPacketCounter runs as a goroutine and continuously processes packet counter buffer
func processPacketCounter() {
	for inc := range packetCounterChan {
		if inc {
			atomic.AddUint64(&packetCount, 1)
		}
	}
}

func calculatePPS() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		currentCount := atomic.SwapUint64(&packetCount, 0)
		if pps < currentCount {
			atomic.SwapUint64(&maxPPS, currentCount)
		}
		atomic.StoreUint64(&pps, currentCount)
		log.Println(pps)
	}
}

func GetPPS() uint64 {
	return pps
}

func GetMaxPPS() uint64 {
	return maxPPS
}
