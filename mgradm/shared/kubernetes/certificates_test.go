// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

//go:build !nok8s

package kubernetes

import (
	"errors"
	"fmt"
	"testing"

	"github.com/rs/zerolog"
	"github.com/uyuni-project/uyuni-tools/shared/testutils"
)

func TestHasIssuer(t *testing.T) {
	type testType struct {
		out      string
		err      error
		expected bool
	}

	data := []testType{
		{
			out:      "issuer.cert-manager.io/someissuer\n",
			err:      nil,
			expected: true,
		},
		{
			out:      "any error\n",
			err:      errors.New("Any error"),
			expected: false,
		},
	}

	for i, test := range data {
		runCmdOutput = func(_ zerolog.Level, _ string, _ ...string) ([]byte, error) {
			return []byte(test.out), test.err
		}
		testutils.AssertEquals(t, fmt.Sprintf("test %d: unexpected result", i+1), test.expected,
			HasIssuer("somens", "someissuer"),
		)
	}
}
