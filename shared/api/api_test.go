// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"testing"

	"github.com/uyuni-project/uyuni-tools/shared/testutils"
)

func TestRedactHeaders(t *testing.T) {
	data := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "JSESSIONID cookie",
			input:    `"JSESSIONID=supersecret; Path=/; Secure; HttpOnly; HttpOnly;HttpOnly;Secure"`,
			expected: `"JSESSIONID=<REDACTED>; Path=/; Secure; HttpOnly; HttpOnly;HttpOnly;Secure"`,
		},
		{
			name:     "pxt-session-cookie cookie",
			input:    `"pxt-session-cookie=supersecret; Max-Age=0;"`,
			expected: `"pxt-session-cookie=<REDACTED>; Max-Age=0;"`,
		},
	}

	for _, tt := range data {
		t.Run(tt.name, func(t *testing.T) {
			actual := redactHeaders(tt.input)
			testutils.AssertEquals(t, "Redact failed", tt.expected, actual)
		})
	}
}
