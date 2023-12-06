package utils

import (
	"github.com/cenkalti/backoff/v4"
	"time"
)

func RetryConfig() *backoff.ExponentialBackOff {
	var bo = backoff.NewExponentialBackOff()
	bo.InitialInterval = 500 * time.Millisecond
	bo.Multiplier = 1.5
	bo.RandomizationFactor = 0.5
	return bo
}
