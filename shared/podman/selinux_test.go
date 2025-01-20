// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"errors"
	"fmt"
	"testing"

	"github.com/rs/zerolog"
	"github.com/uyuni-project/uyuni-tools/shared/testutils"
)

func TestIsSELinuxEnabled(t *testing.T) {
	type testType struct {
		err      error
		expected bool
	}

	cases := []testType{
		{nil, true},
		{errors.New("no such program selinuxenabled"), false},
	}

	for i, testCase := range cases {
		runCmdOutput = func(_ zerolog.Level, _ string, _ ...string) ([]byte, error) {
			return []byte(""), testCase.err
		}
		caseString := fmt.Sprintf("case %d: ", i)
		testutils.AssertEquals(t, caseString+"unexpected return value", testCase.expected, IsSELinuxEnabled())
	}
}
