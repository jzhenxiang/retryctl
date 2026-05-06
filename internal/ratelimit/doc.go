// Package ratelimit implements a sliding-window rate limiter for controlling
// the frequency of retry attempts made by retryctl.
//
// # Overview
//
// A [Limiter] tracks attempt timestamps within a configurable time window and
// rejects calls to [Limiter.Allow] once the configured maximum has been
// reached.  Expired timestamps are evicted lazily on each call, so memory
// usage is bounded by Config.Max rather than by wall-clock time.
//
// # Usage
//
//	lim, err := ratelimit.New(ratelimit.Config{
//		Max:    5,
//		Window: 10 * time.Second,
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	if err := lim.Allow(); errors.Is(err, ratelimit.ErrRateLimited) {
//		// back off or abort
//	}
package ratelimit
