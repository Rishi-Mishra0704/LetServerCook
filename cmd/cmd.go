package cmd

import (
	"fmt"
	"time"

	"github.com/Rishi-Mishra0704/LetServerCook/models"
	"github.com/spf13/cobra"
)

var (
	flagURL      string
	flagMethod   string
	flagRequests int
	flagWorkers  int
	flagTimeout  time.Duration
)

var rootCmd = &cobra.Command{
	Use:   "lsc",
	Short: "LetServerCook - a blazing fast load tester in Go",
	Long: `LetServerCook is a modern CLI load testing tool.
It fires concurrent requests to stress test your API and 
then tells you if the server cooked, got cooked, or burnt the kitchen.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Build Request struct
		req := models.Request{
			URL:       flagURL,
			Method:    flagMethod,
			TotalReqs: flagRequests,
			Workers:   flagWorkers,
			TTL:       flagTimeout,
		}

		// Run benchmarks
		result := RunBenchmarks(req)

		// Print sarcastic summary
		printResult(result)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&flagURL, "url", "u", "", "Target URL to stress test (required)")
	rootCmd.Flags().StringVarP(&flagMethod, "method", "m", "GET", "HTTP method to use")
	rootCmd.Flags().IntVarP(&flagRequests, "requests", "n", 100, "Total number of requests")
	rootCmd.Flags().IntVarP(&flagWorkers, "workers", "w", 10, "Number of concurrent workers")
	rootCmd.Flags().DurationVarP(&flagTimeout, "timeout", "t", 5*time.Second, "Per-request timeout duration")

	rootCmd.MarkFlagRequired("url")
}

// printResult prints the Benchmark results with sarcastic summary
func printResult(b models.Benchmark) {
	// Decide summary text
	var summary string
	failRate := float64(b.Failures) / float64(b.Requests) * 100
	successRate := float64(b.Success) / float64(b.Requests) * 100
	if failRate < 5 {
		summary = "Server Cooked"
	} else if failRate < 30 {
		summary = "Server is Cooked"
	} else {
		summary = "Server cooked so hard, burnt the kitchen"
	}

	fmt.Println(summary)
	fmt.Println("🔥 LetServerCook Results 🔥")
	fmt.Printf("Requests:     %d\n", b.Requests)
	fmt.Printf("Concurrency:  %d\n", b.Concurrency)
	fmt.Printf("Success:      %d | Failures: %d\n", b.Success, b.Failures)
	fmt.Printf("Success Rate: %.0f percent \n", successRate)
	fmt.Printf("Avg Latency:  %s\n", b.AvgLatency)
	fmt.Printf("Max Latency:  %s\n", b.MaxLatency)
	fmt.Printf("RPS:          %.2f\n", b.RPS)
}
