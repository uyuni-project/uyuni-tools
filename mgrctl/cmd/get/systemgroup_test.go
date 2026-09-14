// SPDX-FileCopyrightText: 2026 Jay Prakash katara <katarajayprakash@icloud.com>
//
// SPDX-License-Identifier: Apache-2.0

package get

import (
	"net/http"
	"sort"
	"strings"
	"testing"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/mocks"
	"github.com/uyuni-project/uyuni-tools/shared/testutils"
)

func TestSystemGroupComparator(t *testing.T) {
	groups := []SystemGroup{
		{ID: 3, Name: "same", Description: "charlie", OrgID: 2, SystemCount: 10},
		{ID: 1, Name: "alpha", Description: "bravo", OrgID: 3, SystemCount: 20},
		{ID: 2, Name: "same", Description: "alpha", OrgID: 1, SystemCount: 5},
	}
	tests := []struct {
		name    string
		sortBy  string
		wantIDs []int
	}{
		{name: "no sorting", wantIDs: []int{3, 1, 2}},
		{name: "ID", sortBy: "id", wantIDs: []int{1, 2, 3}},
		{name: "name preserves ties", sortBy: "name", wantIDs: []int{1, 3, 2}},
		{name: "description", sortBy: "description", wantIDs: []int{2, 1, 3}},
		{name: "organization ID", sortBy: "org_id", wantIDs: []int{2, 3, 1}},
		{name: "system count", sortBy: "system_count", wantIDs: []int{2, 3, 1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compare, err := systemGroupComparator(tt.sortBy)
			testutils.AssertNoError(t, "unexpected error", err)

			got := append([]SystemGroup(nil), groups...)
			if compare != nil {
				sort.SliceStable(got, func(i, j int) bool {
					return compare(got[i], got[j])
				})
			}
			gotIDs := make([]int, len(got))
			for i, group := range got {
				gotIDs[i] = group.ID
			}
			testutils.AssertEquals(t, "unexpected sorted IDs", tt.wantIDs, gotIDs)
		})
	}
}

func TestSystemGroupComparatorRejectsUnknownField(t *testing.T) {
	const sortBy = "unknown"
	_, err := systemGroupComparator(sortBy)
	for _, want := range []string{"unsupported sort field", sortBy, strings.Join(systemGroupSortFields, ", ")} {
		testutils.AssertError(t, want, err)
	}
}

func TestSystemGroupListValidatesSortBeforeRequest(t *testing.T) {
	_, _, err := (&systemGroupFetcher{}).List(nil, ListOptions{SortBy: "unknown"})
	testutils.AssertError(t, "unsupported sort field", err)
}

func TestSystemGroupListRejectsFilterBeforeRequest(t *testing.T) {
	_, _, err := (&systemGroupFetcher{}).List(nil, ListOptions{Filter: "name=example"})
	testutils.AssertError(t, "filtering is not supported", err)
}

func TestSystemGroupListSortsAPIResponse(t *testing.T) {
	client := &api.APIClient{
		BaseURL: "https://example.test/rhn/manager/api",
		Client: &mocks.MockClient{DoFunc: func(req *http.Request) (*http.Response, error) {
			testutils.AssertEquals(t, "wrong system group endpoint", "/rhn/manager/api/systemgroup/listAllGroups", req.URL.Path)
			return testutils.GetResponse(200, `{
				"success": true,
				"result": [
					{"id": 1, "name": "zulu"},
					{"id": 2, "name": "alpha"}
				]
			}`)
		}},
	}

	groups, total, err := (&systemGroupFetcher{}).List(client, ListOptions{SortBy: "name"})
	testutils.AssertNoError(t, "unexpected error", err)
	testutils.AssertEquals(t, "wrong number of returned items", 2, total)
	gotIDs := make([]int, len(groups))
	for i, group := range groups {
		gotIDs[i] = group.ID
	}
	testutils.AssertEquals(t, "wrong returned IDs", []int{2, 1}, gotIDs)
}
