package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type State int

const (
	UnknownState State = iota
	FailureState
	SuccessState
)

type Counter interface {
	Count(State)
	ConsecutiveFailures() uint32
	LastActivity() time.Time
}
type Circuit func(context.Context) error

type counter struct {
	mu                 sync.Mutex
	consecutiveFailure uint32
	lastActivity       time.Time
}

func NewCounter() Counter {
	return &counter{}
}

func (c *counter) Count(s State) {
	c.mu.Lock()
	defer c.mu.Unlock()

	switch s {
	case SuccessState:
		c.consecutiveFailure = 0
	case FailureState:
		c.consecutiveFailure++
	}

	c.lastActivity = time.Now()
}

func (c *counter) ConsecutiveFailures() uint32 {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.consecutiveFailure
}

func (c *counter) LastActivity() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.lastActivity
}

func Breaker(c Circuit, failureThreshold uint32) Circuit {
	cnt := NewCounter()

	return func(ctx context.Context) error {
		if cnt.ConsecutiveFailures() >= failureThreshold {
			canRetry := func(cnt Counter) bool {
				backoffLevel := cnt.ConsecutiveFailures() - failureThreshold

				// Calculates when should the circuit breaker resume propagating requests
				// to the service
				shouldRetryAt := cnt.LastActivity().Add(time.Second * 2 << backoffLevel)

				return time.Now().After(shouldRetryAt)
			}

			if !canRetry(cnt) {
				// Fails fast instead of propagating requests to the circuit since
				// not enough time has passed since the last failure to retry
				return fmt.Errorf("service unavailable")
			}
		}

		// Unless the failure threshold is exceeded the wrapped service mimics the
		// old behavior and the difference in behavior is seen after consecutive failures
		if err := c(ctx); err != nil {
			cnt.Count(FailureState)
			return err
		}

		cnt.Count(SuccessState)
		return nil
	}
}

func main() {
	// Define a circuit that fails after 3 attempts
	circuit := func(ctx context.Context) error {
		attempts := 0
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				attempts++
				if attempts > 3 {
					return fmt.Errorf("circuit failed after %d attempts", attempts)
				}
				fmt.Println("circuit attempt", attempts)
				time.Sleep(time.Second)
			}
		}
	}

	// Wrap the circuit with a breaker that allows up to 2 consecutive failures
	breaker := Breaker(circuit, 2)

	// Call the circuit through the breaker
	for i := 0; i < 5; i++ {
		err := breaker(context.Background())
		if err != nil {
			fmt.Println("breaker error:", err)
		} else {
			fmt.Println("breaker success")
		}
		time.Sleep(time.Second)
	}
}
