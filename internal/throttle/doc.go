// Package throttle provides a simple semaphore-based concurrency limiter for
// use with retryctl's runner. It caps the number of simultaneous in-flight
// attempts, preventing a burst of retries from overwhelming a downstream
// service.
//
// Usage:
//
//	th, err := throttle.New(5) // allow at most 5 concurrent attempts
//	if err != nil {
//	    log.Fatal(err)
//	}
//	th.Acquire()
//	defer th.Release()
//	// ... execute attempt ...
package throttle
