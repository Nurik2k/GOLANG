package main

import (
	"sync"
)

func SingleHashService(in, out chan interface{}, wg sync.WaitGroup, mutex sync.Mutex) {
	var hash1 chan string
	var hash2 chan string

	go func() {
		hash1 <- DataSignerCrc32(ConvToString(in))
	}()

	go func() {
		hash2 <- DataSignerMd5(ConvToString(in))
	}()
}

func MultiHashService(in, out chan interface{}, wg sync.WaitGroup, mutex sync.Mutex) {

}
