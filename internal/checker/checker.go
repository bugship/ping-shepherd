// Package checker runs HTTP probes against stored targets.
package checker

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/bugship/ping-shepherd/internal/probe"
	"github.com/bugship/ping-shepherd/internal/store"
)

// Store persists check results for enabled targets.
type Store interface {
	ListTargets(context.Context) ([]store.Target, error)
	LatestCheck(context.Context, string) (store.Check, error)
	RecordCheck(context.Context, string, store.Check) (store.Check, error)
}

// Sender delivers a text alert.
type Sender interface {
	Send(context.Context, string) error
}

// Once probes every enabled target and records the result.
func Once(ctx context.Context, s Store, timeout time.Duration, alert Sender) error {
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
		prev, err := s.LatestCheck(ctx, t.ID)
		hadPrev := err == nil
		if err != nil && !errors.Is(err, store.ErrNotFound) {
			return err
		}
		res := probe.URL(ctx, t.URL, timeout)
		check := store.Check{
			Up:         res.Up,
			StatusCode: res.StatusCode,
			LatencyMS:  res.Latency.Milliseconds(),
			Error:      res.Err,
		}
		if _, err := s.RecordCheck(ctx, t.ID, check); err != nil {
			return err
		}
		if alert == nil {
			continue
		}
		if err := maybeAlert(ctx, alert, t, check, hadPrev, prev); err != nil {
			log.Printf("alert: %v", err)
		}
	}
	return nil
}

func maybeAlert(ctx context.Context, alert Sender, t store.Target, check store.Check, hadPrev bool, prev store.Check) error {
	switch {
	case !check.Up && (!hadPrev || prev.Up):
		return alert.Send(ctx, fmt.Sprintf("DOWN %s (%s) code=%d %s", t.Name, t.URL, check.StatusCode, check.Error))
	case check.Up && hadPrev && !prev.Up:
		return alert.Send(ctx, fmt.Sprintf("UP %s (%s) code=%d", t.Name, t.URL, check.StatusCode))
	default:
		return nil
	}
}

// Loop runs Once immediately, then again on each interval until ctx is done.
func Loop(ctx context.Context, s Store, interval, timeout time.Duration, alert Sender) {
	if interval <= 0 {
		interval = 30 * time.Second
	}
	run := func() {
		if err := Once(ctx, s, timeout, alert); err != nil && ctx.Err() == nil {
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
