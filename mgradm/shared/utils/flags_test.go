// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"testing"

	"github.com/uyuni-project/uyuni-tools/shared/testutils"
)

func TestIdChecker(t *testing.T) {
	data := map[string]bool{
		"foo":       true,
		"foo bar":   false,
		"\u798f":    false,
		"foo123._-": true,
		"foo+":      false,
		"foo&":      false,
		"foo'":      false,
		"foo\"":     false,
		"foo`":      false,
		"foo=":      false,
		"foo#":      false,
	}
	for value, expected := range data {
		actual := idChecker(value)
		testutils.AssertEquals(t, value, expected, actual)
	}
}

func TestEmailChecker(t *testing.T) {
	data := map[string]bool{
		"root@localhost":           true,
		"joe.hacker@foo.bar.com":   true,
		"<joe.hacker@foo.bar.com>": false,
		"fooo":                     false,
	}
	for value, expected := range data {
		actual := emailChecker(value)
		testutils.AssertEquals(t, value, expected, actual)
	}
}
