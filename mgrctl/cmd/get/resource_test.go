// SPDX-FileCopyrightText: 2026 Jay Prakash katara <katarajayprakash@icloud.com>
//
// SPDX-License-Identifier: Apache-2.0

package get

import (
	"fmt"
	"testing"

	"github.com/go-viper/mapstructure/v2"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/testutils"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func TestSortByFlagConfiguration(t *testing.T) {
	const expected = "system_count"

	cmd := NewCommand(&types.GlobalFlags{})
	testutils.AssertNoError(t, "unexpected error setting --sort-by flag", cmd.Flags().Set("sort-by", expected))

	config, err := utils.ReadConfig(cmd, t.TempDir()+"/missing.yaml")
	testutils.AssertNoError(t, "unexpected error reading command configuration", err)

	var flags getFlags
	testutils.AssertNoError(
		t,
		"unexpected error decoding command configuration",
		mapstructure.Decode(config.AllSettings(), &flags),
	)
	testutils.AssertEquals(t, "wrongly decoded sort field", expected, flags.Sort.By)
}

func TestResourceTypesNoDuplicates(t *testing.T) {
	var seen []string

	for name, res := range resourceTypes {
		testutils.AssertNotContains(t, "Duplicate resource key found", seen, name)
		seen = append(seen, name)

		for _, alias := range res.Aliases {
			testutils.AssertNotContains(
				t,
				fmt.Sprintf("Duplicate resource alias found: %s (in resource %s)", alias, name),
				seen,
				alias,
			)
			seen = append(seen, alias)
		}
	}
}

func TestParseFilter(t *testing.T) {
	tests := []struct {
		input     string
		wantKey   string
		wantValue string
		wantErr   bool
	}{
		{"name=foo", "name", "=foo", false},
		{"extra_pkg_count>0", "extra_pkg_count", ">0", false},
		{"count>=10", "count", ">=10", false},
		{"count<=5", "count", "<=5", false},
		{"status!=active", "status", "!=active", false},
		{"id<100", "id", "<100", false},
		{"justkey", "justkey", "", false},
		{" name = foo ", "name", "= foo", false},
		{"mykey |= my value", "", "", true},
		{"bad key=foo", "", "", true},
		{"=foo", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			key, value, err := parseFilter(tt.input)
			if !tt.wantErr {
				testutils.AssertNoError(t, "unexpected error", err)
				testutils.AssertEquals(t, "wrong key parsed", tt.wantKey, key)
				testutils.AssertEquals(t, "wrong value parsed", tt.wantValue, value)
			} else {
				testutils.AssertTrue(t, "expected an error, got none", err != nil)
			}
		})
	}
}

type lookupTestItem struct {
	ID   int
	Name string
}

type lookupTestFetcher struct{}

func (*lookupTestFetcher) List(*api.APIClient, ListOptions) ([]lookupTestItem, int, error) {
	return []lookupTestItem{{ID: 1, Name: "one"}}, 1, nil
}

func (*lookupTestFetcher) Columns() []utils.ColumnDef {
	return []utils.ColumnDef{{Header: "ID", Field: "ID"}}
}

func (*lookupTestFetcher) Help() ResourceHelp {
	return ResourceHelp{
		Summary:  "test resource help",
		Details:  "test resource details",
		Examples: "  mgrctl get lookuptest",
	}
}

// TestRegisterAndLookupResource checks registerResource + lookupResource with name and aliases.
func TestRegisterAndLookupResource(t *testing.T) {
	const name = "lookuptest"
	aliases := []string{"lt", "lookup-test"}
	wantColumns := []utils.ColumnDef{{Header: "ID", Field: "ID"}}
	wantHelp := (&lookupTestFetcher{}).Help()

	registerResource[lookupTestItem](name, &lookupTestFetcher{}, aliases)
	t.Cleanup(func() { delete(resourceTypes, name) })

	for _, input := range append([]string{name}, aliases...) {
		res, err := lookupResource(input)
		testutils.AssertNoError(t, "lookupResource("+input+")", err)
		testutils.AssertEquals(t, "unexpected Help for "+input, wantHelp, res.Help())
		testutils.AssertEquals(t, "unexpected Aliases for "+input, aliases, res.Aliases)
		testutils.AssertEquals(t, "unexpected Columns for "+input, wantColumns, res.Columns())
		testutils.AssertTrue(t, "ListAndPrint should be set for "+input, res.ListAndPrint != nil)
	}

	for _, want := range append([]string{name}, aliases...) {
		testutils.AssertTrue(t, "registeredTypes should include "+want, utils.Contains(registeredTypes(), want))
	}
	testutils.AssertStringContains(t, "missing resource summary", getResourceHelpSummaries(), wantHelp.Summary)
	testutils.AssertStringContains(t, "missing resource details", getResourceHelpDetails(), wantHelp.Details)
	testutils.AssertStringContains(t, "missing resource examples", getResourceHelpExamples(), wantHelp.Examples)

	_, err := lookupResource("not-a-lookuptest")
	testutils.AssertError(t, "unknown resource type", err)
}

func TestResourceHelpSections(t *testing.T) {
	summaries := []string{
		"system (alias: sys)",
		"List systems",
		"systemgroup (alias: grp)",
		"List system groups",
	}
	for _, expected := range summaries {
		testutils.AssertStringContains(t, "missing resource summary", getResourceHelpSummaries(), expected)
	}
	details := []string{
		"Systems support server-side filtering, sorting, and pagination.",
		"System groups support client-side sorting by:",
	}
	for _, expected := range details {
		testutils.AssertStringContains(t, "missing resource details", getResourceHelpDetails(), expected)
	}
	examples := []string{
		"mgrctl get system --page 0 --page-size 10",
		"mgrctl get systemgroup --sort-by=name",
	}
	for _, expected := range examples {
		testutils.AssertStringContains(t, "missing resource example", getResourceHelpExamples(), expected)
	}
}

func TestLookupResource(t *testing.T) {
	tests := []struct {
		input   string
		wantErr string
	}{
		{"system", ""},
		{"sys", ""},
		{"systemgroup", ""},
		{"grp", ""},
		{"", "unknown resource type"},
		{"unknown", "unknown resource type"},
		{"SYSTEM", "unknown resource type"},
	}

	for _, tt := range tests {
		res, err := lookupResource(tt.input)
		if tt.wantErr != "" {
			testutils.AssertError(t, tt.wantErr, err)
			continue
		}
		testutils.AssertNoError(t, "lookupResource("+tt.input+")", err)
		testutils.AssertTrue(t, "Help summary should not be empty for "+tt.input, res.Help().Summary != "")
		testutils.AssertTrue(t, "Columns should not be empty for "+tt.input, len(res.Columns()) > 0)
		testutils.AssertTrue(t, "ListAndPrint should be set for "+tt.input, res.ListAndPrint != nil)
	}
}
