// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"errors"
	"testing"

	"github.com/rs/zerolog"
	"github.com/uyuni-project/uyuni-tools/shared/testutils"
)

func TestIsSELinuxEnabled(t *testing.T) {
	type testType struct {
		name     string
		err      error
		expected bool
	}

	cases := []testType{
		{
			name:     "selinux enabled",
			err:      nil,
			expected: true,
		},
		{
			name:     "selinux not enabled",
			err:      errors.New("no such program selinuxenabled"),
			expected: false,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			runCmdOutput = func(_ zerolog.Level, _ string, _ ...string) ([]byte, error) {
				return []byte(""), testCase.err
			}
			testutils.AssertEquals(t, "unexpected return value", testCase.expected, IsSELinuxEnabled())
		})
	}
}
