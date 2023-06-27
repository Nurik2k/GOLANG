package main

import (
	"fmt"
	"sync"
)

func ExecutePipeline(jobs ...job) {
	in := make(chan interface{})
	out := make(chan interface{})

	wg := &sync.WaitGroup{}

	for _, j := range jobs {

		go j(in, out)
		in = out
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

	for i := range in {

		go MultiHashService(i, out, wg)
		wg.Done()
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

}

func ConvToString(inter interface{}) string {
	str := fmt.Sprintf("%v", inter)
	return str
}
