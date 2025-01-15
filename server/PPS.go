package server

import (
	"sync/atomic"
	"time"
)

// this file implements a service that keeps track of packet per second statistic metric

var (
	inPacketCount        uint64
	inPacketCounterChan  = make(chan int, 500) // buffered channel for packet counting, used to hopefully bypass slowdowns with overusing mutexes
	inPPS                uint64                // stores the actual packet per second metric
	inMaxPPS             uint64                // runtime max of pps
	outPacketCount       uint64
	outPacketCounterChan = make(chan int, 500)
	outPPS               uint64
	outMaxPPS            uint64
)

func PPSInit() {
	go processInPacketCounter()
	go processOutPacketCounter()
	go calculatePPS()
}

// processInPacketCounter runs as a goroutine and continuously processes packet counter buffer
func processInPacketCounter() {
	for inc := range inPacketCounterChan {
		if inc == 1 {
			atomic.AddUint64(&inPacketCount, 1)
		}
	}
}

func processOutPacketCounter() {
	for inc := range inPacketCounterChan {
		if inc == 1 {
			atomic.AddUint64(&outPacketCount, 1)
		}
	}
}

func calculatePPS() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		currentInCount := atomic.SwapUint64(&inPacketCount, 0)
		if inMaxPPS < currentInCount {
			atomic.SwapUint64(&inMaxPPS, currentInCount)
		}
		atomic.StoreUint64(&inPPS, currentInCount)

		currentOutCount := atomic.SwapUint64(&outPacketCount, 0)
		if outMaxPPS < currentOutCount {
			atomic.SwapUint64(&outMaxPPS, currentOutCount)
		}
		atomic.StoreUint64(&outPPS, currentOutCount)
	}
}

func GetInPPS() uint64 {
	return inPPS
}

func GetInMaxPPS() uint64 {
	return inMaxPPS
}

func GetOutPPS() uint64 {
	return outPPS
}

func GetOutMaxPPS() uint64 {
	return outMaxPPS
}
