// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"testing"

	"github.com/uyuni-project/uyuni-tools/shared/testutils"
)

func TestReadInspectData(t *testing.T) {
	content := `Timezone=Europe/Berlin
image_pg_version=16
current_pg_version=14
db_user=myuser
db_password=mysecret
db_name=mydb
db_port=1234
has_hubxmlrpc=true
`

	actual, err := ReadInspectData[InspectResult]([]byte(content))
	if err != nil {
		t.Fatalf("Unexpected failure: %s", err)
	}

	testutils.AssertEquals(t, "Invalid timezone", "Europe/Berlin", actual.Timezone)
	testutils.AssertEquals(t, "Invalid current postgresql version", "14", actual.CurrentPgVersion)
	testutils.AssertEquals(t, "Invalid image postgresql version", "16", actual.ImagePgVersion)
	testutils.AssertEquals(t, "Invalid DB user", "myuser", actual.DBUser)
	testutils.AssertEquals(t, "Invalid DB password", "mysecret", actual.DBPassword)
	testutils.AssertEquals(t, "Invalid DB name", "mydb", actual.DBName)
	testutils.AssertEquals(t, "Invalid DB port", 1234, actual.DBPort)
	testutils.AssertTrue(t, "HasHubXmlrpcApi should be true", actual.HasHubXmlrpcAPI)
}
