// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import "testing"

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
		if actual != expected {
			t.Errorf("%s: expected %v got %v", value, expected, actual)
		}
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
		if actual != expected {
			t.Errorf("%s: expected %v got %v", value, expected, actual)
		}
	}
}
