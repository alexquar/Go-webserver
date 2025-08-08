// # Go Concurrency: Notes and Examples

// Go provides built-in support for concurrency using goroutines and channels.

// ## 1. Goroutines

// A goroutine is a lightweight thread managed by the Go runtime.

// ### Example: Starting a Goroutine

package main

import (
	"fmt"
	"time"
)

func say(s string) {
	for i := 0; i < 3; i++ {
		fmt.Println(s)
		time.Sleep(100 * time.Millisecond)
	}
}

func main() {
	go say("world") // runs concurrently
	say("hello")
}

// ## 2. Channels

// Channels are pipes that connect concurrent goroutines. You can send values into channels from one goroutine and receive those values into another goroutine.

// ### Example: Basic Channel

func main() {
	messages := make(chan string)
	// WRITES PING TO CHANNEL
	go func() { messages <- "ping" }()

	// RECEIVES FROM CHANNEL
	msg := <-messages
	fmt.Println(msg)
}

// ## 3. Buffered Channels	

// Channels can be buffered. Provide the buffer length as the second argument to make.

// ### Example: Buffered Channel

func main() {
	messages := make(chan string, 2)

	messages <- "buffered"
	messages <- "channel"

	fmt.Println(<-messages)
	fmt.Println(<-messages)
}

// ## 4. Channel Synchronization

// Channels can synchronize execution across goroutines.

// ### Example: Channel Synchronization

func worker(done chan bool) {
	fmt.Print("working...")
	time.Sleep(time.Second)
	fmt.Println("done")
	done <- true
}

func main() {
	done := make(chan bool, 1)
	go worker(done)
	<-done
}

// ## 5. Directional Channels

// You can specify if a channel is only for sending or receiving.

// ### Example: Directional Channels

func ping(pings chan<- string, msg string) {
	pings <- msg
}

func pong(pings <-chan string, pongs chan<- string) {
	msg := <-pings
	pongs <- msg
}

func main() {
	pings := make(chan string, 1)
	pongs := make(chan string, 1)
	ping(pings, "passed message")
	pong(pings, pongs)
	fmt.Println(<-pongs)
}

// ## 6. Select Statement

// The select statement lets a goroutine wait on multiple communication operations.

// ### Example: Select

func main() {
	c1 := make(chan string)
	c2 := make(chan string)

	go func() {
		time.Sleep(1 * time.Second)
		c1 <- "one"
	}()
	go func() {
		time.Sleep(2 * time.Second)
		c2 <- "two"
	}()

	for i := 0; i < 2; i++ {
		select {
		case msg1 := <-c1:
			fmt.Println("received", msg1)
		case msg2 := <-c2:
			fmt.Println("received", msg2)
		}
	}
}

// ## 7. Closing Channels

// Channels can be closed to indicate that no more values will be sent.

// ### Example: Closing Channels

func main() {
	jobs := make(chan int, 5)
	done := make(chan bool)

	go func() {
		for {
			j, more := <-jobs
			if more {
				fmt.Println("received job", j)
			} else {
				fmt.Println("received all jobs")
				done <- true
				return
			}
		}
	}()

	for j := 1; j <= 3; j++ {
		jobs <- j
		fmt.Println("sent job", j)
	}
	close(jobs)
	<-done
}

// ## 8. Range over Channels

// You can use range to iterate over values received from a channel.

// ### Example: Range over Channel

func main() {
	queue := make(chan string, 2)
	queue <- "one"
	queue <- "two"
	close(queue)
	for elem := range queue {
		fmt.Println(elem)
	}
}

// ## 9. WaitGroups

// WaitGroups are used to wait for a collection of goroutines to finish.

// ### Example: WaitGroup

import "sync"

func worker(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Worker %d starting\n", id)
	time.Sleep(time.Second)
	fmt.Printf("Worker %d done\n", id)
}

func main() {
	var wg sync.WaitGroup
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go worker(i, &wg)
	}
	wg.Wait()
}

// ## 10. Mutexes

// Mutexes provide mutual exclusion to protect shared data.

// ### Example: Mutex

import "sync"

type SafeCounter struct {
	mu sync.Mutex
	val int
}

func (c *SafeCounter) Inc() {
	c.mu.Lock()
	c.val++
	c.mu.Unlock()
}

func (c *SafeCounter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.val
}

func main() {
	c := SafeCounter{}
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			c.Inc()
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println(c.Value())
}

