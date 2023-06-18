package main

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

// сюда писать код
func ExecutePipeline(jobs ...job) {
	var in, out chan interface{}
	waitGroup := &sync.WaitGroup{}

	for _, i := range jobs {
		in = out
		out = make(chan interface{}, 100)

		waitGroup.Add(1)
		go func(j job, in, out chan interface{}) {
			defer waitGroup.Done()
			defer close(out)

			i(in, out)
		}(i, in, out)
	}

	waitGroup.Wait()
}

func processStr(in, out chan interface{}, process func(string) string) {
	waitGroup := &sync.WaitGroup{}

	for data := range in {
		waitGroup.Add(1)

		go func(data string) {
			defer waitGroup.Done()

			out <- process(data)
		}(fmt.Sprintln(data))
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

			hash2 <- DataSignerMd5(md5Hash)
		}()

		return fmt.Sprintf("%s~%s", <-hash1, <-hash2)
	})
}

func MultiHash(in, out chan interface{}) {
	processStr(in, out, func(data string) string {
		results := make([]string, 6)
		resultMux := &sync.Mutex{}
		resultWaitGroup := &sync.WaitGroup{}

		for i := 0; i <= 5; i++ {
			resultWaitGroup.Add(1)

			go func(i int) {
				defer resultWaitGroup.Done()

				hash := DataSignerCrc32(fmt.Sprintf("%d%s", i, data))

				resultMux.Lock()
				results[i] = hash
				resultMux.Unlock()
			}(i)
		}

		resultWaitGroup.Wait()

		return strings.Join(results, "")
	})
}

func CombineResults(in, out chan interface{}) {
	var results []string
	for i := range in {
		results = append(results, fmt.Sprintf("%d", i))
	}

	sort.Strings(results)
	out <- strings.Join(results, "_")
}
