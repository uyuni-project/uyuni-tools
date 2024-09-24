// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package kubernetes

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/rs/zerolog"
	"github.com/uyuni-project/uyuni-tools/shared/testutils"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func TestGetCurrentDeploymentReplicaSet(t *testing.T) {
	type testType struct {
		rsOut         string
		rsErr         error
		revisionOut   string
		revisionErr   error
		expected      string
		expectedError bool
	}

	testCases := []testType{
		{
			rsOut:         "uyuni-64d597fccf,1 uyuni-66f7677dc6,2\n",
			rsErr:         nil,
			revisionOut:   "2\n",
			revisionErr:   nil,
			expected:      "uyuni-66f7677dc6",
			expectedError: false,
		},
		{
			rsOut:         "uyuni-64d597fccf,1\n",
			rsErr:         nil,
			revisionOut:   "1\n",
			revisionErr:   nil,
			expected:      "uyuni-64d597fccf",
			expectedError: false,
		},
		{
			rsOut:         "\n",
			rsErr:         nil,
			revisionOut:   "not found\n",
			revisionErr:   errors.New("not found"),
			expected:      "",
			expectedError: false,
		},
		{
			rsOut:         "get rs error\n",
			rsErr:         errors.New("get rs error"),
			revisionOut:   "1\n",
			revisionErr:   nil,
			expected:      "",
			expectedError: true,
		},
		{
			rsOut:         "uyuni-64d597fccf,1\n",
			rsErr:         nil,
			revisionOut:   "get rev error\n",
			revisionErr:   errors.New("get rev error"),
			expected:      "",
			expectedError: true,
		},
	}

	for i, test := range testCases {
		runCmdOutput = func(_ zerolog.Level, _ string, args ...string) ([]byte, error) {
			if utils.Contains(args, "rs") {
				return []byte(test.rsOut), test.rsErr
			}
			return []byte(test.revisionOut), test.revisionErr
		}
		actual, err := getCurrentDeploymentReplicaSet("uyunins", "uyuni")
		caseMsg := fmt.Sprintf("test %d: ", i+1)
		testutils.AssertEquals(t, fmt.Sprintf("%sunexpected error raised: %s", caseMsg, err),
			test.expectedError, err != nil,
		)
		testutils.AssertEquals(t, caseMsg+"unexpected result", test.expected, actual)
	}
}

func TestGetPodsFromOwnerReference(t *testing.T) {
	type testType struct {
		out      string
		err      error
		expected []string
	}

	data := []testType{
		{
			out:      "pod1 pod2 pod3\n",
			err:      nil,
			expected: []string{"pod1", "pod2", "pod3"},
		},
		{
			out:      "\n",
			err:      nil,
			expected: []string{},
		},
		{
			out:      "error\n",
			err:      errors.New("some error"),
			expected: []string{},
		},
	}

	for i, test := range data {
		runCmdOutput = func(_ zerolog.Level, _ string, _ ...string) ([]byte, error) {
			return []byte(test.out), test.err
		}
		actual, err := getPodsFromOwnerReference("myns", "owner")
		if test.err == nil {
			testutils.AssertTrue(t, "Shouldn't have raise an error", err == nil)
		} else {
			testutils.AssertTrue(t, "Unexpected error raised", strings.Contains(err.Error(), test.err.Error()))
		}
		testutils.AssertEquals(t, fmt.Sprintf("test %d: unexpected result", i+1), test.expected, actual)
	}
}
