package models

import "time"

type Benchmark struct {
	Requests    int
	Concurrency int
	Success     int
	Failures    int
	AvgLatency  time.Duration
	MaxLatency  time.Duration
	RPS         float64
}
