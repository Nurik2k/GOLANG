package main

import (
	"context"
	"fmt"
	"time"
)

/*
функцию тест менять нельзя

	1: таймаут контекста = 1 секунде

в этом случае мы не должны дожидаться отработки функции тест

	2: таймаут контекста = 10 секунд

в этом случае мы не должны ждать все 10 секунд, так как тест отрабатывает всего 2 секунды
*/
func main() {
	now := time.Now()
	defer func() {
		fmt.Println("finish, elapsed time", time.Since(now))
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 11*time.Second)
	defer cancel()

	c := make(chan bool)

	go func() {
		test()
		c <- true
	}()

	select {
	case <-ctx.Done():

	case <-c /*Tell finish*/ :
	}

}

func FinishTest() <-chan bool {
	c := make(chan bool)

	go func() {
		test()
		c <- true
	}()
	return c
}

func test() {
	time.Sleep(time.Second * 2)
	fmt.Println("test finished")
}
