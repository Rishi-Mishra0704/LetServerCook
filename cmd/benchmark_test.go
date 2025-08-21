package cmd_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Rishi-Mishra0704/LetServerCook/cmd"
	"github.com/Rishi-Mishra0704/LetServerCook/models"
	"github.com/stretchr/testify/assert"
)

func TestRunBenchmarks_AllSuccess(t *testing.T) {
	// Local server that always returns 200
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("ok"))
	}))
	defer ts.Close()

	r := models.Request{
		URL:       ts.URL,
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

func TestRunBenchmarks_AllFail(t *testing.T) {
	// Local server that always returns 500
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		w.Write([]byte("fail"))
	}))
	defer ts.Close()

	r := models.Request{
		URL:       ts.URL,
		Method:    "GET",
		TotalReqs: 10,
		Workers:   2,
		TTL:       1 * time.Second,
	}

	bench := cmd.RunBenchmarks(r)

	assert.Equal(t, 0, bench.Success)
	assert.Equal(t, r.TotalReqs, bench.Failures)
}

func TestWorkerHandlesTimeout(t *testing.T) {
	// Server sleeps longer than TTL to trigger timeout
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(200)
	}))
	defer ts.Close()

	r := models.Request{
		URL:       ts.URL,
		Method:    "GET",
		TotalReqs: 5,
		Workers:   1,
		TTL:       50 * time.Millisecond, // shorter than server sleep
	}

	bench := cmd.RunBenchmarks(r)

	// All requests should fail due to timeout
	assert.Equal(t, 0, bench.Success)
	assert.Equal(t, r.TotalReqs, bench.Failures)
}
