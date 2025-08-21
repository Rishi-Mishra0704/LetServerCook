package cmd_test

import (
	"testing"
	"time"

	"github.com/Rishi-Mishra0704/LetServerCook/cmd"
	"github.com/Rishi-Mishra0704/LetServerCook/models"
	"github.com/stretchr/testify/assert"
)

func TestRunBenchmarks_AllSuccess(t *testing.T) {
	r := models.Request{
		URL:       "http://example.com",
		Method:    "GET",
		TotalReqs: 10,
		Workers:   2,
		TTL:       1 * time.Second,
	}

	bench := cmd.RunBenchmarks(r)

	assert.Equal(t, r.TotalReqs, bench.Success)
	assert.Equal(t, 0, bench.Failures)
	assert.True(t, bench.AvgLatency > 0)
	assert.True(t, bench.MaxLatency > 0)
	assert.True(t, bench.RPS > 0)
}

func TestWorkerHandlesTimeout(t *testing.T) {
	r := models.Request{
		URL:       "http://example.com",
		Method:    "GET",
		TotalReqs: 5,
		Workers:   1,
		TTL:       50 * time.Millisecond,
	}

	bench := cmd.RunBenchmarks(r)

	assert.Equal(t, 0, bench.Success)
	assert.Equal(t, r.TotalReqs, bench.Failures)
}
