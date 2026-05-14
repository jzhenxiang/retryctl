// Package bulkhead implements the bulkhead pattern, isolating failure domains
// by limiting the number of concurrent calls to a named partition.
package bulkhead

import (
	"errors"
	"fmt"
	"sync"
)

// ErrPartitionFull is returned when a partition has reached its concurrency limit.
var ErrPartitionFull = errors.New("bulkhead: partition is full")

// ErrUnknownPartition is returned when releasing a partition that was never acquired.
var ErrUnknownPartition = errors.New("bulkhead: unknown partition")

// Bulkhead manages isolated concurrency limits per named partition.
type Bulkhead struct {
	mu         sync.Mutex
	partitions map[string]*partition
}

type partition struct {
	max     int
	active  int
}

// New creates a Bulkhead with the given partition limits.
// The limits map must have at least one entry and all values must be positive.
func New(limits map[string]int) (*Bulkhead, error) {
	if len(limits) == 0 {
		return nil, errors.New("bulkhead: at least one partition required")
	}
	partitions := make(map[string]*partition, len(limits))
	for name, max := range limits {
		if max <= 0 {
			return nil, fmt.Errorf("bulkhead: partition %q max must be positive, got %d", name, max)
		}
		partitions[name] = &partition{max: max}
	}
	return &Bulkhead{partitions: partitions}, nil
}

// Acquire attempts to enter the named partition.
// Returns ErrPartitionFull if the limit is reached, or ErrUnknownPartition if
// the partition name was not registered at construction time.
func (b *Bulkhead) Acquire(name string) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	p, ok := b.partitions[name]
	if !ok {
		return fmt.Errorf("%w: %q", ErrUnknownPartition, name)
	}
	if p.active >= p.max {
		return fmt.Errorf("%w: %q (%d/%d)", ErrPartitionFull, name, p.active, p.max)
	}
	p.active++
	return nil
}

// Release decrements the active count for the named partition.
func (b *Bulkhead) Release(name string) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	p, ok := b.partitions[name]
	if !ok {
		return fmt.Errorf("%w: %q", ErrUnknownPartition, name)
	}
	if p.active > 0 {
		p.active--
	}
	return nil
}

// Available returns the number of remaining slots in the named partition.
// Returns -1 if the partition is unknown.
func (b *Bulkhead) Available(name string) int {
	b.mu.Lock()
	defer b.mu.Unlock()
	p, ok := b.partitions[name]
	if !ok {
		return -1
	}
	return p.max - p.active
}
