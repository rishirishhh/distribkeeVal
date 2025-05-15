package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
	"time"
)

var (
	addr        = flag.String("addr", "localhost:8080", "The HTTP host port for the instance that is benchmarked.")
	iterations  = flag.Int("iterations", 1000, "The number of total iterations (requests)")
	concurrency = flag.Int("concurrency", 1, "Number of goroutines to run in parallel")
)

// Reuse a single HTTP client for all requests
var httpClient = &http.Client{
	Timeout: 5 * time.Second,
}

func writeRand() {
	key := fmt.Sprintf("key-%d", rand.Intn(1000000))
	value := fmt.Sprintf("value-%d", rand.Intn(1000000))

	values := url.Values{}
	values.Set("key", key)
	values.Set("value", value)

	resp, err := httpClient.Get("http://" + (*addr) + "/set?" + values.Encode())
	if err != nil {
		log.Printf("Error during set: %v", err)
		return
	}
	defer resp.Body.Close()
}

func benchmark(name string, fn func(), iterations int, wg *sync.WaitGroup, resultChan chan time.Duration) {
	defer wg.Done()

	for i := 0; i < iterations; i++ {
		start := time.Now()
		fn()
		duration := time.Since(start)
		resultChan <- duration
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	flag.Parse()

	fmt.Printf("Running %d iterations with concurrency level %d\n", *iterations, *concurrency)

	// Warm-up phase (10 requests to stabilize)
	for i := 0; i < 10; i++ {
		writeRand()
	}

	// Channels and sync
	var wg sync.WaitGroup
	resultChan := make(chan time.Duration, *iterations)

	iterationsPerWorker := *iterations / *concurrency
	extra := *iterations % *concurrency

	start := time.Now()

	for i := 0; i < *concurrency; i++ {
		count := iterationsPerWorker
		if i < extra {
			count++
		}
		wg.Add(1)
		go benchmark("write", writeRand, count, &wg, resultChan)
	}

	wg.Wait()
	close(resultChan)

	total := time.Since(start)

	// Analyze results
	var max time.Duration
	var min = time.Hour
	var sum time.Duration
	count := 0

	for d := range resultChan {
		if d > max {
			max = d
		}
		if d < min {
			min = d
		}
		sum += d
		count++
	}

	avg := sum / time.Duration(count)
	qps := float64(count) / total.Seconds()

	fmt.Printf("\n==== Benchmark Result ====\n")
	fmt.Printf("Total Requests: %d\n", count)
	fmt.Printf("Concurrency: %d\n", *concurrency)
	fmt.Printf("Avg: %s | Max: %s | Min: %s | QPS: %.2f\n", avg, max, min, qps)
}
