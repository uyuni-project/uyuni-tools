// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package kubernetes

import (
	"errors"
	"fmt"
	"testing"

	"github.com/rs/zerolog"
	"github.com/uyuni-project/uyuni-tools/shared/testutils"
)

func TestHasDeployment(t *testing.T) {
	type dataType struct {
		out      string
		err      error
		expected bool
	}

	data := []dataType{
		{"deployment.apps/traefik\n", nil, true},
		{"\n", nil, false},
		{"Some error", errors.New("Some error"), false},
	}

	for i, test := range data {
		runCmdOutput = func(_ zerolog.Level, _ string, _ ...string) ([]byte, error) {
			return []byte(test.out), test.err
		}
		testutils.AssertEquals(t, fmt.Sprintf("test %d: unexpected result", i+1), test.expected,
			HasDeployment("kube-system", "-lapp.kubernetes.io/name=traefik"),
		)
	}
}

func TestGetReplicas(t *testing.T) {
	type dataType struct {
		out      string
		err      error
		expected int
	}
	data := []dataType{
		{"2\n", nil, 2},
		{"no such deploy\n", errors.New("No such deploy"), 0},
		{"invalid output\n", nil, 0},
	}

	for i, test := range data {
		runCmdOutput = func(_ zerolog.Level, _ string, _ ...string) ([]byte, error) {
			return []byte(test.out), test.err
		}
		testutils.AssertEquals(t, fmt.Sprintf("test %d: unexpected result", i+1),
			test.expected, GetReplicas("uyuni", "uyuni-hub-api"))
	}
}
