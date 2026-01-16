// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"sync"
)

// RingBuffer is a thread-safe, fixed-size buffer that preserves the most recent data.
type RingBuffer struct {
	mu   sync.Mutex
	data []byte
	size int
	pos  int  // The write pointer
	full bool // Whether we have wrapped around at least once
}

// NewRingBuffer creates a buffer fixed at 'size' bytes.
func NewRingBuffer(size int) *RingBuffer {
	return &RingBuffer{
		data: make([]byte, size),
		size: size,
	}
}

// Write writes data to the buffer, overwriting the oldest data if full.
// It is thread-safe (handles concurrent writes from Stdout/Stderr).
func (r *RingBuffer) Write(p []byte) (n int, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	n = len(p)
	if n == 0 {
		return 0, nil
	}

	// Case 1: The new data is bigger than the entire buffer.
	// We just keep the last 'size' bytes of the new data.
	if n > r.size {
		p = p[n-r.size:]
	}

	// Case 2: Standard circular write
	remaining := r.size - r.pos
	if len(p) > remaining {
		// Overflow: Write to the end, then wrap to the beginning
		copy(r.data[r.pos:], p[:remaining])
		copy(r.data[0:], p[remaining:])
		r.pos = len(p) - remaining
		r.full = true
	} else {
		// No Overflow: Simple copy
		copy(r.data[r.pos:], p)
		r.pos += len(p)
		if r.pos == r.size {
			r.pos = 0 // Exactly hit the end
			r.full = true
		}
	}

	// Always return the original input length so io.MultiWriter
	// doesn't think the write failed.
	return n, nil
}

// Bytes returns the contents of the buffer as a byte slice.
func (r *RingBuffer) Bytes() []byte {
	r.mu.Lock()
	defer r.mu.Unlock()

	// 1. If never full, just return a copy of what we have
	if !r.full {
		out := make([]byte, r.pos)
		copy(out, r.data[:r.pos])
		return out
	}

	// 2. If full, stitch the "Oldest" (tail) + "Newest" (head)
	out := make([]byte, r.size)
	// Copy oldest data first (from pos to end)
	copy(out, r.data[r.pos:])
	// Copy newest data second (from 0 to pos)
	copy(out[r.size-r.pos:], r.data[:r.pos])

	return out
}
