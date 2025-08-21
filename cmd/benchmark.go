package cmd

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/Rishi-Mishra0704/LetServerCook/models"
)

type WorkerReq struct {
	id       int
	wg       *sync.WaitGroup
	tasks    <-chan int
	r        models.Request
	mu       *sync.Mutex
	success  *int
	failures *int
	client   *http.Client
}

type RequestResult struct {
	index   int
	latency time.Duration
	success bool
}

func RunBenchmarks(r models.Request) models.Benchmark {
	fmt.Println("Starting LetServerCook benchmark...")
	fmt.Printf("URL: %s | Method: %s | Requests: %d | Workers: %d | Timeout: %s\n",
		r.URL, r.Method, r.TotalReqs, r.Workers, r.TTL)

	tasks := make(chan int, r.TotalReqs)
	results := make(chan RequestResult, r.TotalReqs) // Channel for results
	var wg sync.WaitGroup

	var successCount, failureCount int
	var mu sync.Mutex
	var latencies []time.Duration

	fmt.Println("Spawning workers...")

	// Improved HTTP client configuration
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: r.Workers + 10, // More connections per host
			MaxConnsPerHost:     r.Workers + 10,
			IdleConnTimeout:     30 * time.Second,
			DisableKeepAlives:   false, // Keep connections alive
		},
		Timeout: r.TTL, // Set client-level timeout
	}

	// Start result collector goroutine
	go func() {
		latencies = make([]time.Duration, 0, r.TotalReqs)
		for result := range results {
			mu.Lock()
			latencies = append(latencies, result.latency)
			if result.success {
				successCount++
			} else {
				failureCount++
			}
			mu.Unlock()
		}
	}()

	for i := 0; i < r.Workers; i++ {
		wg.Add(1)
		wr := WorkerReq{
			id:       i,
			wg:       &wg,
			tasks:    tasks,
			r:        r,
			success:  &successCount,
			failures: &failureCount,
			mu:       &mu,
			client:   client,
		}
		go WorkerFixed(wr, results)
	}

	fmt.Println("Enqueuing tasks...")
	for i := 0; i < r.TotalReqs; i++ {
		tasks <- i
	}
	close(tasks)

	startBenchmark := time.Now()
	fmt.Println("Running benchmark...")
	wg.Wait()
	close(results) // Signal result collector to finish

	totalDuration := time.Since(startBenchmark)

	fmt.Println("Calculating results...")

	// Wait a bit for result collector to finish
	time.Sleep(10 * time.Millisecond)

	mu.Lock()
	var totalLatency time.Duration
	var maxLatency time.Duration
	validLatencies := 0

	for _, l := range latencies {
		if l > 0 { // Only count non-zero latencies
			totalLatency += l
			validLatencies++
			if l > maxLatency {
				maxLatency = l
			}
		}
	}

	var avgLatency time.Duration
	if validLatencies > 0 {
		avgLatency = totalLatency / time.Duration(validLatencies)
	}

	rps := float64(successCount) / totalDuration.Seconds()
	mu.Unlock()

	fmt.Println("Benchmark completed!")
	return models.Benchmark{
		Requests:    r.TotalReqs,
		Concurrency: r.Workers,
		Success:     successCount,
		Failures:    failureCount,
		AvgLatency:  avgLatency,
		MaxLatency:  maxLatency,
		RPS:         rps,
	}
}

func WorkerFixed(wr WorkerReq, results chan<- RequestResult) {
	defer wr.wg.Done()
	count := 0

	for idx := range wr.tasks {
		ctx, cancel := context.WithTimeout(context.Background(), wr.r.TTL)
		start := time.Now()

		req, err := http.NewRequestWithContext(ctx, wr.r.Method, wr.r.URL, nil)
		if err != nil {
			cancel()
			results <- RequestResult{index: idx, latency: 0, success: false}
			continue
		}

		resp, err := wr.client.Do(req)
		elapsed := time.Since(start)
		cancel()

		success := err == nil && resp != nil && resp.StatusCode < 400

		// Send result through channel (thread-safe)
		results <- RequestResult{
			index:   idx,
			latency: elapsed,
			success: success,
		}

		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}

		count++
		if count%10 == 0 {
			fmt.Printf("Worker %d processed %d requests...\n", wr.id, count)
		}
	}
}
