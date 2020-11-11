package channels

import (
	"math/rand"
	"time"
)

func ping(pings chan<- int) {
	for {
		pings <- rand.Int()
		time.Sleep(1 * time.Second)
	}

}

func pong(pings <-chan int) {
	for {
		msg := <-pings
		time.Sleep(1 * time.Second)
		println(msg)
	}

}

func TestChannels() {
	pings := make(chan int, 1)
	go ping(pings)
	go pong(pings)
	for {
		time.Sleep(1 * time.Second)
	}
}
