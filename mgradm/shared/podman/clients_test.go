// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"sort"
	"testing"

	"github.com/uyuni-project/uyuni-tools/shared/testutils"
)

func TestParseClientTrustResults(t *testing.T) {
	saltOutput := []byte(`{"m1.example.com":"OK\n","m2.example.com":"MISSING","down.example.com":"OK"}`)
	unreachable := map[string]bool{"down.example.com": true}

	migrated, notMigrated, err := parseClientTrustResults(saltOutput, unreachable)
	testutils.AssertEquals(t, "unexpected error", nil, err)

	sort.Strings(migrated)
	sort.Strings(notMigrated)
	testutils.AssertEquals(t, "wrong migrated minions", []string{"m1.example.com"}, migrated)
	testutils.AssertEquals(t, "wrong not-migrated minions", []string{"m2.example.com"}, notMigrated)
}

func TestParseClientTrustResultsEmpty(t *testing.T) {
	migrated, notMigrated, err := parseClientTrustResults([]byte("  "), nil)
	testutils.AssertEquals(t, "unexpected error", nil, err)
	testutils.AssertEquals(t, "expected no migrated minion", 0, len(migrated))
	testutils.AssertEquals(t, "expected no not-migrated minion", 0, len(notMigrated))
}
