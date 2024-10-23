// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package kubernetes

import (
	"errors"
	"fmt"
	"testing"

	"github.com/rs/zerolog"
	"github.com/uyuni-project/uyuni-tools/shared/test_utils"
)

func TestHasVolume(t *testing.T) {
	type dataType struct {
		err      error
		out      string
		expected bool
	}
	data := []dataType{
		{nil, "Bound\n", true},
		{nil, "Pending\n", false},
		{errors.New("PVC not found"), "", false},
	}

	for i, test := range data {
		runCmdOutput = func(logLevel zerolog.Level, command string, args ...string) ([]byte, error) {
			return []byte(test.out), test.err
		}
		actual := HasVolume("myns", "thepvc")
		test_utils.AssertEquals(t, fmt.Sprintf("test %d: unexpected output", i), test.expected, actual)
	}
}
