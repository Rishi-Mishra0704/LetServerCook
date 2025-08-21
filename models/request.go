package models

import "time"

type Request struct {
	URL       string
	Method    string
	Workers   int
	TotalReqs int
	TTL       time.Duration
}
