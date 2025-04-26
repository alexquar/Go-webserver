package main

import (
	"fmt"
)

func main() {
	// c1 := make(chan string)
	// c2 := make(chan string)

	// go func() {
	// 	for {
	// 		c1 <- "Ever 500ms"
	// 		time.Sleep(time.Millisecond * 500)
	// 	}
	// }()

	// go func() {
	// 	for {
	// 		c2 <- "Ever 2s"
	// 		time.Sleep(time.Second * 2)
	// 	}
	// }()
	// for {
	// 	select {
	// 	case msg1 := <-c1:
	// 		fmt.Println(msg1)
	// 	case msg2 := <-c2:
	// 		fmt.Println(msg2)
	// 	}

	// }]
	jobs := make(chan int, 10)
	result := make(chan int, 10)
	go worker(jobs, result)
	for i := 0; i < 10; i++ {
		jobs <- i
	}
	close(jobs)
	for i := 0; i < 10; i++ {
		fmt.Println(<-result)
	}

}

// func count(thing string, c chan string) {
// 	for i := 1; i <= 5; i++ {
// 		c <- thing
// 		time.Sleep(time.Millisecond * 500)
// 	}
// 	close(c)
// }

func worker(jobs <-chan int, result chan<- int) {
	for n := range jobs {
		result <- fib(n)
	}

}

func fib(n int) int {
	if n <= 1 {
		return n
	}
	return fib(n-1) + fib(n-2)
}
