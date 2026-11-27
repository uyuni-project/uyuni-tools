// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"testing"

	"github.com/uyuni-project/uyuni-tools/shared/testutils"
)

func TestRingBuffer(t *testing.T) {
	tests := []struct {
		name     string
		size     int
		writes   []string
		expected string
	}{
		{
			name:     "Empty Buffer",
			size:     4,
			writes:   []string{},
			expected: "",
		},
		{
			name:     "Basic Write (Under Limit)",
			size:     4,
			writes:   []string{"123"},
			expected: "123",
		},
		{
			name:     "Exact Fill",
			size:     4,
			writes:   []string{"1234"},
			expected: "1234",
		},
		{
			name:     "Simple Wrap Around",
			size:     4,
			writes:   []string{"12", "345"},
			expected: "2345",
		},
		{
			name:     "Single Write Larger than Buffer",
			size:     4,
			writes:   []string{"1234567890"},
			expected: "7890",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rb := NewRingBuffer(tt.size)

			for _, w := range tt.writes {
				_, err := rb.Write([]byte(w))
				testutils.AssertNoError(t, "Write failed", err)
			}

			testutils.AssertEquals(t, "Bytes() mismatch", tt.expected, string(rb.Bytes()))
		})
	}
}
