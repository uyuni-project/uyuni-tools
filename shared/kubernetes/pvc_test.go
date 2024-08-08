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

func TestHasPersistentVolumeClaim(t *testing.T) {
	type dataType struct {
		err      error
		out      string
		expected bool
	}
	data := []dataType{
		{nil, "persistentvolumeclaim/var-pgsql\n", true},
		{errors.New("PVC not found"), "", false},
	}

	for i, test := range data {
		runCmdOutput = func(_ zerolog.Level, _ string, _ ...string) ([]byte, error) {
			return []byte(test.out), test.err
		}
		actual := hasPersistentVolumeClaim("myns", "thepvc")
		testutils.AssertEquals(t, fmt.Sprintf("test %d: unexpected output", i), test.expected, actual)
	}
}
