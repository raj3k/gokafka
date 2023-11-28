package main

import (
	"context"
	"encoding/binary"
	"log"
	"math/bits"
	"math/rand"
	"time"

	"github.com/segmentio/kafka-go"
)

const (
	topic     = "randomIntStream"
	partition = 0
)

func repeatFunc[T any, K any](done <-chan K, fn func() T) <-chan T {
	stream := make(chan T)
	go func() {
		defer close(stream)
		for {
			select {
			case <-done:
				return
			case stream <- fn():
			}
		}
	}()
	return stream
}

func take[T any, K any](done <-chan K, stream <-chan T, n int) <-chan T {
	taken := make(chan T)
	go func() {
		defer close(taken)
		for i := 0; i < n; i++ {
			select {
			case <-done:
				return
			case taken <- <-stream:
			}
		}
	}()
	return taken
}

func primeFinder(done <-chan bool, randIntStream <-chan int) <-chan int {
	isPrime := func(randomInt int) bool {
		for i := randomInt - 1; i > 1; i-- {
			if randomInt%i == 0 {
				return false
			}
		}
		return true
	}

	primes := make(chan int)
	go func() {
		defer close(primes)
		for {
			select {
			case <-done:
			case randomInt := <-randIntStream:
				if isPrime(randomInt) {
					primes <- randomInt
				}
			}
		}
	}()
	return primes
}

func encodeUint(x uint64) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, x)
	return buf[bits.LeadingZeros64(x)>>3:]
}

func main() {

	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9093", topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

	done := make(chan bool)
	randNumFetcher := func() uint64 { return rand.Uint64() }
	randIntStream := repeatFunc(done, randNumFetcher)
	// primeStream := primeFinder(done, randIntStream)

	// for randInt := range take(done, primeStream, 10) {
	// 	fmt.Println(randInt)
	// }

	for randInt := range randIntStream {
		_, err := conn.Write(encodeUint(randInt))
		if err != nil {
			log.Fatal("failed to write message:", err)
		}
	}
	if err := conn.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}
}
