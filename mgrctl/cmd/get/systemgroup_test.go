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
	for _, want := range []string{"unsupported sort field", sortBy, strings.Join(systemGroupFields, ", ")} {
		testutils.AssertError(t, want, err)
	}
}

func TestSystemGroupListValidatesSortBeforeRequest(t *testing.T) {
	_, _, err := (&systemGroupFetcher{}).List(nil, ListOptions{SortBy: "unknown"})
	testutils.AssertError(t, "unsupported sort field", err)
}

func TestSystemGroupFilter(t *testing.T) {
	group := SystemGroup{ID: 3, Name: "Prod", Description: "prod servers", OrgID: 2, SystemCount: 10}
	tests := []struct {
		name   string
		expr   string
		want   bool
		errMsg string
	}{
		{name: "id equality", expr: "id=3", want: true},
		{name: "id inequality", expr: "id!=4", want: true},
		{name: "id greater", expr: "id>2", want: true},
		{name: "id less rejected", expr: "id<3", want: false},
		{name: "system_count greater-equal", expr: "system_count>=10", want: true},
		{name: "system_count less-equal boundary", expr: "system_count<=10", want: true},
		{name: "org_id less", expr: "org_id<3", want: true},
		{name: "name exact match", expr: "name=Prod", want: true},
		{name: "name is case sensitive", expr: "name=prod", want: false},
		{name: "name inequality", expr: "name!=prod", want: true},
		{name: "description exact match", expr: "description=prod servers", want: true},
		{name: "empty description equality", expr: "description=", want: false},
		{name: "empty description inequality", expr: "description!=", want: true},
		{name: "whitespace is trimmed", expr: " name = Prod ", want: true},
		{name: "ordering operator on string", expr: "name>abc", errMsg: "unsupported filter operator"},
		{name: "unknown field", expr: "unknown=1", errMsg: "unsupported filter field"},
		{name: "missing operator", expr: "name", errMsg: "expected an operator"},
		{name: "empty field", expr: "=3", errMsg: "filter field cannot be empty"},
		{name: "bare exclamation", expr: "name!Prod", errMsg: "unsupported filter operator"},
		{name: "empty numeric value", expr: "id=", errMsg: "must be an integer"},
		{name: "invalid field characters", expr: "na;me=Prod", errMsg: "invalid filter field"},
		{name: "non-integer numeric value", expr: "id=abc", errMsg: "must be an integer"},
		{name: "float numeric value", expr: "id=3.5", errMsg: "must be an integer"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			predicate, err := systemGroupFilter(tt.expr)
			if tt.errMsg != "" {
				testutils.AssertError(t, tt.errMsg, err)
				return
			}
			testutils.AssertNoError(t, "unexpected error", err)
			testutils.AssertEquals(t, "unexpected filter result", tt.want, predicate(group))
		})
	}
}

func TestSystemGroupFilterOperatorCharsInName(t *testing.T) {
	// Only the first operator splits field from value; the value keeps any
	// remaining operator characters.
	tests := []struct {
		expr  string
		value string
		want  bool
	}{
		{"name=prod>=test", "prod>=test", true},
		{"name=prod<=test", "prod<=test", true},
		{"name!=a=b", "a=b", false},
		{"name!=a=b", "other", true},
	}
	for _, tt := range tests {
		t.Run(tt.expr+"/"+tt.value, func(t *testing.T) {
			predicate, err := systemGroupFilter(tt.expr)
			testutils.AssertNoError(t, "unexpected error", err)
			testutils.AssertEquals(t, "unexpected filter result", tt.want,
				predicate(SystemGroup{Name: tt.value}))
		})
	}
}

func TestSystemGroupListValidatesFilterBeforeRequest(t *testing.T) {
	calls := 0
	client := &api.APIClient{
		BaseURL: "https://example.test/rhn/manager/api",
		Client: &mocks.MockClient{DoFunc: func(_ *http.Request) (*http.Response, error) {
			calls++
			return testutils.GetResponse(200, `{"success": true, "result": []}`)
		}},
	}

	tests := []struct{ filter, errMsg string }{
		{"unknown=1", "unsupported filter field"},
		{"name>abc", "unsupported filter operator"},
		{"id=abc", "must be an integer"},
		{"badfilter", "expected an operator"},
	}
	for _, tt := range tests {
		_, _, err := (&systemGroupFetcher{}).List(client, ListOptions{Filter: tt.filter})
		testutils.AssertError(t, tt.errMsg, err)
	}
	testutils.AssertEquals(t, "invalid filters must not reach the endpoint", 0, calls)
}

func TestSystemGroupListFiltersAndSortsAPIResponse(t *testing.T) {
	client := &api.APIClient{
		BaseURL: "https://example.test/rhn/manager/api",
		Client: &mocks.MockClient{DoFunc: func(_ *http.Request) (*http.Response, error) {
			return testutils.GetResponse(200, `{
				"success": true,
				"result": [
					{"id": 1, "name": "zulu", "system_count": 3},
					{"id": 2, "name": "Prod", "system_count": 15},
					{"id": 3, "name": "alpha", "system_count": 10},
					{"id": 4, "name": "beta", "system_count": 10}
				]
			}`)
		}},
	}

	tests := []struct {
		name      string
		opts      ListOptions
		wantIDs   []int
		wantTotal int
	}{
		{
			name:      "string filter",
			opts:      ListOptions{Filter: "name=Prod"},
			wantIDs:   []int{2},
			wantTotal: 1,
		},
		{
			name:      "numeric filter and ascending sort",
			opts:      ListOptions{Filter: "system_count>=10", SortBy: "name"},
			wantIDs:   []int{2, 3, 4},
			wantTotal: 3,
		},
		{
			name:      "filter and descending sort",
			opts:      ListOptions{Filter: "system_count<=10", SortBy: "system_count", SortDescending: true},
			wantIDs:   []int{3, 4, 1},
			wantTotal: 3,
		},
		{
			name:      "no matches",
			opts:      ListOptions{Filter: "id=99"},
			wantIDs:   []int{},
			wantTotal: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			groups, total, err := (&systemGroupFetcher{}).List(client, tt.opts)
			testutils.AssertNoError(t, "unexpected error", err)
			testutils.AssertEquals(t, "wrong total", tt.wantTotal, total)
			testutils.AssertEquals(t, "wrong number of returned items", tt.wantTotal, len(groups))
			gotIDs := make([]int, len(groups))
			for i, group := range groups {
				gotIDs[i] = group.ID
			}
			testutils.AssertEquals(t, "wrong returned IDs", tt.wantIDs, gotIDs)
		})
	}
}

func TestSystemGroupListEmptyAPIResponse(t *testing.T) {
	client := &api.APIClient{
		BaseURL: "https://example.test/rhn/manager/api",
		Client: &mocks.MockClient{DoFunc: func(_ *http.Request) (*http.Response, error) {
			return testutils.GetResponse(200, `{"success": true, "result": []}`)
		}},
	}

	groups, total, err := (&systemGroupFetcher{}).List(client, ListOptions{SortBy: "name"})
	testutils.AssertNoError(t, "unexpected error", err)
	testutils.AssertEquals(t, "wrong total", 0, total)
	if groups == nil {
		t.Error("empty result must be an empty slice, not nil")
	}
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

func TestSystemGroupListSortsDescending(t *testing.T) {
	client := &api.APIClient{
		BaseURL: "https://example.test/rhn/manager/api",
		Client: &mocks.MockClient{DoFunc: func(_ *http.Request) (*http.Response, error) {
			return testutils.GetResponse(200, `{
				"success": true,
				"result": [
					{"id": 1, "name": "bar"},
					{"id": 2, "name": "foo"}
				]
			}`)
		}},
	}

	groups, _, err := (&systemGroupFetcher{}).List(client, ListOptions{SortBy: "name", SortDescending: true})
	testutils.AssertNoError(t, "unexpected error", err)
	gotIDs := make([]int, len(groups))
	for i, group := range groups {
		gotIDs[i] = group.ID
	}
	testutils.AssertEquals(t, "wrong returned IDs for descending sort", []int{2, 1}, gotIDs)
}
