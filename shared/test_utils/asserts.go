// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package test_utils

import "testing"

// AssertEquals ensures two values are equals and raises and error if not.
func AssertEquals[T comparable](t *testing.T, message string, expected T, actual T) {
	if actual != expected {
		t.Errorf(message+": got '%v' expected '%v'", actual, expected)
	}
}
