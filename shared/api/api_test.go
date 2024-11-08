// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package api

import "testing"

func TestRedactHeaders(t *testing.T) {
	data := [][]string{
		{
			`"JSESSIONID=supersecret; Path=/; Secure; HttpOnly; HttpOnly;HttpOnly;Secure"`,
			`"JSESSIONID=<REDACTED>; Path=/; Secure; HttpOnly; HttpOnly;HttpOnly;Secure"`,
		},
		{
			`"pxt-session-cookie=supersecret; Max-Age=0;"`,
			`"pxt-session-cookie=<REDACTED>; Max-Age=0;"`,
		},
	}

	for i, testCase := range data {
		input := testCase[0]
		expected := testCase[1]

		actual := redactHeaders(input)

		if actual != expected {
			t.Errorf("Testcase %d: Expected %s got %s when redacting  %s", i, expected, actual, input)
		}
	}
}
