// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"testing"

	"github.com/uyuni-project/uyuni-tools/shared/testutils"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func TestSanityCheck(t *testing.T) {
	type dataType struct {
		name            string
		oldUyuniRelease string
		oldSumaRelease  string
		oldPgsqlVersion string
		newUyuniRelease string
		newSumaRelease  string
		newPgsqlVersion string
		errorPart       string
	}
	data := []dataType{
		{
			name:            "valid Uyuni upgrade",
			oldUyuniRelease: "2024.07",
			oldSumaRelease:  "",
			oldPgsqlVersion: "16",
			newUyuniRelease: "2024.13",
			newSumaRelease:  "",
			newPgsqlVersion: "17",
			errorPart:       "",
		},
		{
			name:            "valid SUSE Multi-Linux Manager upgrade",
			oldUyuniRelease: "",
			oldSumaRelease:  "5.0.1",
			oldPgsqlVersion: "16",
			newUyuniRelease: "",
			newSumaRelease:  "5.1.0",
			newPgsqlVersion: "17",
			errorPart:       "",
		},
		{
			name:            "Uyuni downgrade error",
			oldUyuniRelease: "2024.13",
			oldSumaRelease:  "",
			oldPgsqlVersion: "17",
			newUyuniRelease: "2024.07",
			newSumaRelease:  "",
			newPgsqlVersion: "16",
			errorPart:       "cannot downgrade",
		},
		{
			name:            "SUSE Multi-Linux Manager downgrade error",
			oldUyuniRelease: "",
			oldSumaRelease:  "5.1.0",
			oldPgsqlVersion: "17",
			newUyuniRelease: "",
			newSumaRelease:  "5.0.1",
			newPgsqlVersion: "16",
			errorPart:       "cannot downgrade",
		},
		{
			name:            "upgrade from Uyuni to SUSE Multi-Linux Manager unsupported",
			oldUyuniRelease: "2024.07",
			oldSumaRelease:  "",
			oldPgsqlVersion: "16",
			newUyuniRelease: "",
			newSumaRelease:  "5.1.0",
			newPgsqlVersion: "17",
			errorPart:       "Upgrade is not supported",
		},
		{
			name:            "upgrade from SUSE Multi-Linux Manager to Uyuni unsupported",
			oldUyuniRelease: "",
			oldSumaRelease:  "5.1.0",
			oldPgsqlVersion: "17",
			newUyuniRelease: "2024.07",
			newSumaRelease:  "",
			newPgsqlVersion: "16",
			errorPart:       "Upgrade is not supported",
		},
	}

	for _, test := range data {
		t.Run(test.name, func(t *testing.T) {
			// newValues represents the data inspected from the new image/server
			newValues := utils.InspectData{
				ServerInspectData: utils.ServerInspectData{
					UyuniRelease:       test.newUyuniRelease,
					SuseManagerRelease: test.newSumaRelease,
				},
				DBInspectData: utils.DBInspectData{
					PgVersion: test.newPgsqlVersion,
				},
				ContainerInspectData: utils.ContainerInspectData{
					UyuniRelease:       test.oldUyuniRelease,
					SuseManagerRelease: test.oldSumaRelease,
					PgVersion:          test.oldPgsqlVersion,
				},
			}

			// Ensure SanityCheck signature accepts (*ContainerInspectData, *InspectData)
			err := SanityCheck(&newValues)

			testutils.AssertError(t, test.errorPart, err)
		})
	}
}
