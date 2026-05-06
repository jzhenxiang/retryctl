// Package notify defines the Notifier interface and built-in implementations
// used by retryctl to report the final outcome of a retry sequence.
//
// Built-in implementations
//
//   - LogNotifier  – writes a single structured line to any io.Writer.
//   - Multi        – fans an Event out to multiple Notifiers in order.
//
// Usage
//
//	n := notify.NewMulti(
//		notify.NewLogNotifier(os.Stdout),
//	)
//	n.Notify(notify.Event{
//		Command:  "myapp",
//		Success:  true,
//		Attempts: 2,
//		Elapsed:  300 * time.Millisecond,
//	})
package notify
