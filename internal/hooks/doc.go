// Package hooks provides lifecycle hook support for retryctl.
//
// Hooks allow users to run arbitrary shell commands at key points during
// command execution:
//
//   - before_attempt  – fired before every attempt
//   - after_success   – fired when an attempt succeeds
//   - after_failure   – fired after each failed attempt
//   - after_final     – fired once all attempts are exhausted
//
// Example usage:
//
//	r := hooks.New([]hooks.Hook{
//		{Event: hooks.EventAfterFailure, Command: "notify-send 'retryctl: attempt failed'"},
//		{Event: hooks.EventAfterSuccess, Command: "echo done"},
//	})
//
//	if err := r.Run(hooks.EventBeforeAttempt, map[string]string{"ATTEMPT": "1"}); err != nil {
//		log.Fatal(err)
//	}
package hooks
