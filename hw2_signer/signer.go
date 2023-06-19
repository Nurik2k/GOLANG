package main

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

func ExecutePipeline(jobs ...job) {
	var in, out chan interface{}
	waitGroup := &sync.WaitGroup{}

	for _, i := range jobs {
		in = out
		out = make(chan interface{}, 100)

		waitGroup.Add(1)
		go func(i job, in, out chan interface{}) {
			defer waitGroup.Done()
			defer close(out)

			i(in, out)
		}(i, in, out)
	}

	waitGroup.Wait()
}

func processStr(in, out chan interface{}, process func(string) string) {
	waitGroup := &sync.WaitGroup{}

	for i := range in {
		waitGroup.Add(1)

		go func(data string) {
			defer waitGroup.Done()

			out <- process(data)
		}(fmt.Sprintf("%v", i))
	}

	waitGroup.Wait()
}

func SingleHash(in, out chan interface{}) {
	md5Mux := &sync.Mutex{}

	processStr(in, out, func(data string) string {

		hash1 := make(chan string)
		go func() {
			defer close(hash1)

			hash1 <- DataSignerCrc32(data)
		}()

		hash2 := make(chan string)
		go func() {
			defer close(hash2)

			md5Mux.Lock()
			md5Hash := DataSignerMd5(data)
			md5Mux.Unlock()

			hash2 <- DataSignerCrc32(md5Hash)
		}()

		return fmt.Sprintf("%s~%s", <-hash1, <-hash2)
	})
}

func MultiHash(in, out chan interface{}) {
	processStr(in, out, func(data string) string {
		results := make([]string, 6)
		resultMux := &sync.Mutex{}
		resultWaitGroup := &sync.WaitGroup{}

		for th := 0; th <= 5; th++ {
			resultWaitGroup.Add(1)

			go func(th int) {
				defer resultWaitGroup.Done()

				hash := DataSignerCrc32(fmt.Sprintf("%d%s", th, data)) // crc32(th+data)

				resultMux.Lock()
				results[th] = hash
				resultMux.Unlock()
			}(th)
		}

		resultWaitGroup.Wait()
		return strings.Join(results, "")
	})
}

func CombineResults(in, out chan interface{}) {
	var results []string

	for data := range in {
		results = append(results, fmt.Sprintf("%v", data))
	}

	sort.Strings(results)
	out <- strings.Join(results, "_")
}
