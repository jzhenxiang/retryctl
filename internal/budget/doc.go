// Package budget implements a sliding-window retry budget.
//
// A Budget caps the total number of retry attempts that may occur within a
// configurable time window.  Once the budget is exhausted, Allow returns
// ErrBudgetExhausted until older attempts slide out of the window or Reset
// is called explicitly.
//
// Usage:
//
//	b, err := budget.New(10, time.Minute)
//	if err != nil { /* handle */ }
//
//	if err := b.Allow(); err != nil {
//	    // too many retries in the last minute
//	}
package budget
