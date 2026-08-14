// Package checker runs HTTP probes against stored targets.
package checker

import (
	"context"
	"log"
	"time"

	"github.com/bugship/ping-shepherd/internal/probe"
	"github.com/bugship/ping-shepherd/internal/store"
)

// Store persists check results for enabled targets.
type Store interface {
	ListTargets(context.Context) ([]store.Target, error)
	RecordCheck(context.Context, string, store.Check) (store.Check, error)
}

// Once probes every enabled target and records the result.
func Once(ctx context.Context, s Store, timeout time.Duration) error {
	targets, err := s.ListTargets(ctx)
	if err != nil {
		return err
	}
	for _, t := range targets {
		if !t.Enabled {
			continue
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}
		res := probe.URL(ctx, t.URL, timeout)
		_, err := s.RecordCheck(ctx, t.ID, store.Check{
			Up:         res.Up,
			StatusCode: res.StatusCode,
			LatencyMS:  res.Latency.Milliseconds(),
			Error:      res.Err,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// Loop runs Once immediately, then again on each interval until ctx is done.
func Loop(ctx context.Context, s Store, interval, timeout time.Duration) {
	if interval <= 0 {
		interval = 30 * time.Second
	}
	run := func() {
		if err := Once(ctx, s, timeout); err != nil && ctx.Err() == nil {
			log.Printf("check pass: %v", err)
		}
	}
	run()
	tick := time.NewTicker(interval)
	defer tick.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-tick.C:
			run()
		}
	}
}
