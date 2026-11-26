// SPDX-FileCopyrightText: 2026 Jay Prakash katara <katarajayprakash@icloud.com>
//
// SPDX-License-Identifier: Apache-2.0

package get

import (
	"testing"

	"github.com/uyuni-project/uyuni-tools/shared/testutils"
)

func TestParseFilter(t *testing.T) {
	tests := []struct {
		input        string
		wantKey      string
		wantOperator string
		wantValue    string
		wantErr      bool
	}{
		{"name=foo", "name", "=", "foo", false},
		{"extra_pkg_count>0", "extra_pkg_count", ">", "0", false},
		{"count>=10", "count", ">=", "10", false},
		{"count<=5", "count", "<=", "5", false},
		{"status!=active", "status", "!=", "active", false},
		{"id<100", "id", "<", "100", false},
		{"foo~bar", "foo", "~", "bar", false},
		{"foo!~bar", "foo", "!~", "bar", false},
		{"justkey", "justkey", "", "", false},
		{" name = foo ", "name", "=", " foo ", false},
		{"mykey |= my value", "", "", "", true},
		{"bad key=foo", "", "", "", true},
		{"=foo", "", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			key, op, value, err := parseFilter(tt.input)
			if !tt.wantErr {
				testutils.AssertNoError(t, "unexpected error", err)
				testutils.AssertEquals(t, "wrong key parsed", tt.wantKey, key)
				testutils.AssertEquals(t, "wrong operator parsed", tt.wantOperator, op)
				testutils.AssertEquals(t, "wrong value parsed", tt.wantValue, value)
			} else {
				testutils.AssertTrue(t, "expected an error, got none", err != nil)
			}
		})
	}
}
