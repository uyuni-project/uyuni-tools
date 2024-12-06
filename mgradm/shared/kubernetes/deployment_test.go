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

func TestGetRunningServerImage(t *testing.T) {
	type dataType struct {
		err      error
		out      string
		expected string
	}
	data := []dataType{
		{nil, "registry.opensuse.org/uyuni/server:latest\n", "registry.opensuse.org/uyuni/server:latest"},
		{errors.New("deployment not found"), "", ""},
	}

	for i, test := range data {
		runCmdOutput = func(_ zerolog.Level, _ string, _ ...string) ([]byte, error) {
			return []byte(test.out), test.err
		}
		actual := getRunningServerImage("myns")
		testutils.AssertEquals(t, fmt.Sprintf("test %d: unexpected result", i), test.expected, actual)
	}
}
