// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"path"
	"testing"

	"github.com/uyuni-project/uyuni-tools/shared/test_utils"
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

	testDir, cleaner := test_utils.CreateTmpFolder(t)
	defer cleaner()

	dataPath := path.Join(testDir, "data")
	test_utils.WriteFile(t, dataPath, content)

	actual, err := ReadInspectData[InspectResult](dataPath)
	if err != nil {
		t.Fatalf("Unexpected failure: %s", err)
	}

	test_utils.AssertEquals(t, "Invalid timezone", "Europe/Berlin", actual.Timezone)
	test_utils.AssertEquals(t, "Invalid current postgresql version", "14", actual.CurrentPgVersion)
	test_utils.AssertEquals(t, "Invalid image postgresql version", "16", actual.ImagePgVersion)
	test_utils.AssertEquals(t, "Invalid DB user", "myuser", actual.DbUser)
	test_utils.AssertEquals(t, "Invalid DB password", "mysecret", actual.DbPassword)
	test_utils.AssertEquals(t, "Invalid DB name", "mydb", actual.DbName)
	test_utils.AssertEquals(t, "Invalid DB port", 1234, actual.DbPort)
	test_utils.AssertTrue(t, "HasHubXmlrpcApi should be true", actual.HasHubXmlrpcApi)
}
