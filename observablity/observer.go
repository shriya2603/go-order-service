package observablity

import (
	"sync"
	"time"
)

func init() {

}

func StartObserver(wg *sync.WaitGroup, observerAddr string, serverName string, bootTime float64) {
	defer wg.Done()
	wg.Add(1)
	go func(serverName string) {
		defer wg.Done()
		for {
			select {
			case <-time.After(30 * time.Second):
				bootTime += 30.0
			}

		}
	}(serverName)

}
