// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"fmt"
	"strings"
	"testing"

	"github.com/uyuni-project/uyuni-tools/shared/testutils"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func TestSanityCheck(t *testing.T) {
	type dataType struct {
		oldUyuniRelease string
		oldSumaRelease  string
		oldPgsqlVersion string
		newUyuniRelease string
		newSumaRelease  string
		newPgsqlVersion string
		errorPart       string
	}
	data := []dataType{
		{"2024.07", "", "16", "2024.13", "", "17", ""},
		{"", "5.0.1", "16", "", "5.1.0", "17", ""},
		{
			"2024.13", "", "17", "2024.07", "", "16",
			"cannot downgrade",
		},
		{
			"", "5.1.0", "17", "", "5.0.1", "16",
			"cannot downgrade",
		},
		{
			"2024.07", "", "16", "", "5.1.0", "17",
			"Upgrade is not supported",
		},
		{
			"", "5.1.0", "17", "2024.07", "", "16",
			"Upgrade is not supported",
		},
		{
			"2024.07", "", "16", "2024.13", "", "",
			"cannot fetch PostgreSQL",
		},
		{
			"2024.07", "", "", "2024.13", "", "17",
			"PostgreSQL is not installed",
		},
	}

	for i, test := range data {
		runningValues := utils.ServerInspectData{
			UyuniRelease:       test.oldUyuniRelease,
			SuseManagerRelease: test.oldSumaRelease,
		}
		newValues := utils.ServerInspectData{
			CommonInspectData: utils.CommonInspectData{
				CurrentPgVersion: test.oldPgsqlVersion,
				ImagePgVersion:   test.newPgsqlVersion,
			},
			UyuniRelease:       test.newUyuniRelease,
			SuseManagerRelease: test.newSumaRelease,
		}
		err := SanityCheck(&runningValues, &newValues, "path/to/image")
		if test.errorPart != "" {
			if err != nil {
				testutils.AssertTrue(
					t, fmt.Sprintf("test %d: Unexpected error message: %s", i+1, err),
					strings.Contains(err.Error(), test.errorPart),
				)
			} else {
				t.Errorf("test %d: expected an error, got none", i+1)
			}
		} else {
			testutils.AssertEquals(t, fmt.Sprintf("test %d: unexpected error", i+1), nil, err)
		}
	}
}
