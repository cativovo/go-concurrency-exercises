//////////////////////////////////////////////////////////////////////
//
// Given is a producer-consumer scenario, where a producer reads in
// tweets from a mockstream and a consumer is processing the
// data. Your task is to change the code so that the producer as well
// as the consumer can run concurrently
//

package main

import (
	"fmt"
	"sync"
	"time"
)

func producer(stream Stream) chan Tweet {
	tweetChan := make(chan Tweet)

	go func() {
		for {
			tweet, err := stream.Next()
			if err == ErrEOF {
				close(tweetChan)
				return
			}

			tweetChan <- *tweet
		}
	}()

	return tweetChan
}

func consumer(tweetChan <-chan Tweet, wg *sync.WaitGroup) {
	for t := range tweetChan {
		wg.Add(1)

		go func(t Tweet) {
			defer wg.Done()
			if t.IsTalkingAboutGo() {
				fmt.Println(t.Username, "\ttweets about golang")
			} else {
				fmt.Println(t.Username, "\tdoes not tweet about golang")
			}
		}(t)
	}
}

func main() {
	start := time.Now()
	stream := GetMockStream()

	// Producer
	tweetChan := producer(stream)

	var wg sync.WaitGroup
	// Consumer
	consumer(tweetChan, &wg)

	wg.Wait()

	fmt.Printf("Process took %s\n", time.Since(start))
}
