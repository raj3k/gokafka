package main

import (
	"context"
	"gokafka/internal"
	"gokafka/kafka"
	"log"
	"math/rand"

	kafkago "github.com/segmentio/kafka-go"
	"golang.org/x/sync/errgroup"
)

// TODO: create config

func repeatFunc[T any, K any](done <-chan K, fn func() T) <-chan T {
	stream := make(chan T, 1000)
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

func main() {

	done := make(chan bool)
	randIntMessageStream := make(<-chan kafkago.Message, 1000)
	randNumFetcher := func() kafkago.Message { return kafkago.Message{Value: internal.ItoBSlice(rand.Intn(5000000))} }
	randIntMessageStream = repeatFunc(done, randNumFetcher)
	// primeStream := primeFinder(done, randIntStream)

	// for randInt := range take(done, primeStream, 10) {
	// 	fmt.Println(randInt)
	// }

	reader := kafka.NewKafkaReader()
	writer := kafka.NewKafkaWriter()

	ctx := context.Background()
	// messages := make(chan kafkago.Message, 1000)
	// messageCommitChan := make(chan kafkago.Message, 1000)

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return writer.WriteMessages(ctx, randIntMessageStream)
	})

	g.Go(func() error {
		return reader.FetchMessages(ctx)
	})

	// g.Go(func() error {
	// 	return reader.CommitMessage(ctx, messageCommitChan)
	// })

	err := g.Wait()
	if err != nil {
		log.Fatalln(err)
	}
}
