package main

import (
	"fmt"
	"sync"
)

func ExecutePipeline(jobs ...job) {
	wg := &sync.WaitGroup{}

	in := make(chan interface{})
	out := make(chan interface{})

	for _, i := range jobs {
		wg.Add(1)
		go func(job job, in, out chan interface{}) {
			defer wg.Done()

			i(in, out)
		}(i, in, out)
	}
	wg.Wait()
}

func SingleHash(in, out chan interface{}) {

}

func MultiHash(in, out chan interface{}) {

}

func CombineResults(in, out chan interface{}) {

}

func main() {
	jobs := []job{
		job(func(in, out chan interface{}) {
			out <- 1
			//out <- 2
			//out <- 3
		}),
		//job(func(in, out chan interface{}) {
		//	for value := range in {
		//		fmt.Println("This is JOB 2, sending data to out channel...")
		//		out <- value
		//	}
		//}),
		//job(func(in, out chan interface{}) {
		//	for value := range in {
		//		fmt.Println("This is JOB 3, sending data to out channel...")
		//		out <- value
		//	}
		//}),
		job(func(in, out chan interface{}) {
			for value := range in {
				fmt.Println("FINAL RESULT: ", value)
			}
		}),
	}

	ExecutePipeline(jobs...)
}

func ConvToString(inter chan interface{}) string {
	str := fmt.Sprintf("%v", inter)
	return str
}
