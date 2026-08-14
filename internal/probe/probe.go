// Package probe performs HTTP GET checks.
package probe

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// Result is the outcome of one HTTP GET.
type Result struct {
	Up         bool
	StatusCode int
	Latency    time.Duration
	Err        string
}

// URL GETs rawURL. Status codes 2xx and 3xx count as up.
func URL(ctx context.Context, rawURL string, timeout time.Duration) Result {
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return Result{Err: err.Error()}
	}
	req.Header.Set("User-Agent", "ping-shepherd/0")

	client := &http.Client{
		Timeout: timeout,
		CheckRedirect: func(_ *http.Request, via []*http.Request) error {
			if len(via) >= 5 {
				return fmt.Errorf("too many redirects")
			}
			return nil
		},
	}

	start := time.Now()
	res, err := client.Do(req)
	latency := time.Since(start)
	if err != nil {
		return Result{Latency: latency, Err: err.Error()}
	}
	defer res.Body.Close()

	up := res.StatusCode >= 200 && res.StatusCode < 400
	return Result{
		Up:         up,
		StatusCode: res.StatusCode,
		Latency:    latency,
	}
}
