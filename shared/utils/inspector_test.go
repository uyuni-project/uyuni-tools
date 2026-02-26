// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"testing"

	"github.com/uyuni-project/uyuni-tools/shared/testutils"
)

func TestReadInspectData(t *testing.T) {
	// Keys must match the `mapstructure` tags defined in the structs.
	// Note: 'pg_version' maps to ContainerInspectData
	//       'db_image_pg_version' maps to DBInspectData
	content := `Timezone=Europe/Berlin
db_image_pg_version=16
pg_version=14
db_user=myuser
db_password=mysecret
db_name=mydb
db_port=1234
has_hubxmlrpc=true
`

	// Use utils.InspectData as the target type
	actual, err := ReadInspectData[InspectResult]([]byte(content))
	if err != nil {
		t.Fatalf("Unexpected failure: %s", err)
	}

	// Assuming Timezone is part of the struct (via squash)
	testutils.AssertEquals(t, "Invalid timezone", "Europe/Berlin", actual.Timezone)

	// 'PgVersion' is ambiguous (exists in multiple structs), so we must specify the path.
	testutils.AssertEquals(t, "Invalid current postgresql version", "14", actual.ContainerInspectData.PgVersion)
	testutils.AssertEquals(t, "Invalid image postgresql version", "16", actual.DBInspectData.PgVersion)

	// Unique fields (like DBUser) are promoted to the top level via 'squash'.
	testutils.AssertEquals(t, "Invalid DB user", "myuser", actual.DBUser)
	testutils.AssertEquals(t, "Invalid DB password", "mysecret", actual.DBPassword)
	testutils.AssertEquals(t, "Invalid DB name", "mydb", actual.DBName)
	testutils.AssertEquals(t, "Invalid DB port", 1234, actual.DBPort)

	testutils.AssertTrue(t, "HasHubXmlrpcApi should be true", actual.HasHubXmlrpcAPI)
}
