package main

import (
	"fmt"
	"sync"
)

func SingleHashService(in interface{}, out chan interface{}, wg *sync.WaitGroup, mutex *sync.Mutex) {
	hash1 := make(chan string)
	hash2 := make(chan string)
	defer func() {
		close(hash1)
		close(hash2)
		wg.Done()
	}()

	dataCrc32 := DataSignerCrc32(ConvToString(in))

	go func(h1 chan string, d2 string) {
		h1 <- d2
	}(hash1, dataCrc32)

	mutex.Lock()
	dataMd5 := DataSignerMd5(ConvToString(in))
	mutex.Unlock()

	go func(h2 chan string, d2 string) {
		h2 <- d2
	}(hash2, dataMd5)

	str := fmt.Sprintf("%s ~ %s", <-hash1, <-hash2)
	out <- str
}

func MultiHashService(in interface{}, out chan interface{}, wg *sync.WaitGroup) {
	hash := make(chan string)

	defer close(hash)
	defer wg.Done()

	go func() {
		hash <- DataSignerCrc32(ConvToString(in))
	}()

	out <- in
}
