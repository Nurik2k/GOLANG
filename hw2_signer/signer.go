package main

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

func ExecutePipeline(jobs ...job) {
	in := make(chan interface{})
	out := make(chan interface{})
	wg := &sync.WaitGroup{}

	for _, value := range jobs {
		in = out
		out = make(chan interface{}, MaxInputDataLen)

		wg.Add(1)

		go func(job2 job, input, output chan interface{}) {
			defer wg.Done()
			defer close(output)

			job2(input, output)
		}(value, in, out)
	}

	wg.Wait()
}

func SingleHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}
	mutex := &sync.Mutex{}

	for i := range in {
		wg.Add(1)
		go SingleHashService(i, out, wg, mutex)
	}

	wg.Wait()
}

func MultiHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}

	for value := range in {
		wg.Add(1)
		go MultiHashService(value, out, wg)
	}
	wg.Wait()
}

func main() {

	jobs := []job{
		job(func(in, out chan interface{}) {
			out <- 1
			out <- 2
			out <- 3
		}),
		job(func(in, out chan interface{}) {
			for value := range in {
				fmt.Println("This is JOB 2, sending data to out channel...")
				out <- value
			}
		}),
		job(func(in, out chan interface{}) {
			for value := range in {
				fmt.Println("This is JOB 3, sending data to out channel...")
				out <- value
			}
		}),
		job(func(in, out chan interface{}) {
			for value := range in {
				fmt.Println("FINAL RESULT: ", value)
			}
		}),
	}
	ExecutePipeline(jobs...)
}

func CombineResults(in, out chan interface{}) {
	inputValue := make([]string, MaxInputDataLen)

	for value := range in {
		inputValue = append(inputValue, ConvToString(value))
	}

	sort.SliceIsSorted(inputValue, func(i, j int) bool {
		return inputValue[i] < inputValue[j]
	})

	out <- strings.Join(inputValue, "_")
}

func ConvToString(inter interface{}) string {
	str := fmt.Sprintf("%v", inter)
	return str
}
