package main

import (
	"fmt"
	"sync"
)

func SingleHashService(in, out chan interface{}, wg sync.WaitGroup, mutex sync.Mutex) {
	var hash1 chan string
	var hash2 chan string

	defer close(hash2)
	defer close(hash1)
	defer wg.Done()

	go func() {
		hash1 <- DataSignerCrc32(ConvToString(in))
	}()
	mutex.Lock()
	mutex.Unlock()
	go func() {
		hash2 <- DataSignerMd5(ConvToString(in))
	}()

	str := fmt.Sprintf("%s ~ %s", <-hash1, <-hash2)
	out <- str
}

func MultiHashService(in, out chan interface{}, wg sync.WaitGroup, mutex sync.Mutex) {

}
