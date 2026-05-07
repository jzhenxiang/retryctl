// Package backpressure implements a token-bucket limiter that can be used to
// apply back-pressure to retry loops.
//
// # Overview
//
// When a downstream service is struggling, hammering it with rapid retries
// often makes things worse. The Limiter hands out tokens at a controlled rate;
// each retry attempt must acquire a token before proceeding. When the bucket
// is empty the attempt blocks until either a token is refilled or the
// caller's context is cancelled.
//
// # Usage
//
//	limiter, err := backpressure.New(5, 200*time.Millisecond)
//	if err != nil { /* handle */ }
//	defer limiter.Stop()
//
//	if err := limiter.Acquire(ctx); err != nil {
//	    // context cancelled — abort retry loop
//	}
package backpressure
