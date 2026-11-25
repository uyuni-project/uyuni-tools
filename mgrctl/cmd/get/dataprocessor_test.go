// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package get

import (
	"testing"
	"time"

	"github.com/uyuni-project/uyuni-tools/shared/testutils"
)

type testGroup struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	OrgID       int       `json:"org_id"`
	SystemCount int       `json:"system_count"`
	CreatedAt   time.Time `json:"created_at"`
}

func testGroups() []testGroup {
	return []testGroup{
		{ID: 3, Name: "webdb", SystemCount: 1, CreatedAt: time.Date(2026, 5, 10, 12, 0, 0, 0, time.UTC)},
		{ID: 1, Name: "web", SystemCount: 5, CreatedAt: time.Date(2026, 5, 9, 23, 59, 59, 0, time.UTC)},
		{ID: 8, Name: "storage", SystemCount: 7, CreatedAt: time.Date(2026, 5, 11, 0, 0, 0, 0, time.UTC)},
		{ID: 4, Name: "app1", SystemCount: 3, CreatedAt: time.Date(2026, 5, 10, 14, 30, 0, 0, time.UTC)},
		{
			ID: 5, Name: "app2", SystemCount: 8,
			CreatedAt: time.Date(2026, 5, 10, 16, 30, 0, 0, time.FixedZone("CEST", 2*3600)),
		},
		{ID: 9, Name: "backup", SystemCount: 6, CreatedAt: time.Date(2026, 5, 12, 8, 0, 0, 0, time.UTC)},
		{ID: 2, Name: "db", SystemCount: 10, CreatedAt: time.Date(2026, 5, 10, 0, 0, 0, 0, time.UTC)},
		{ID: 6, Name: "cache", SystemCount: 2, CreatedAt: time.Date(2026, 5, 10, 14, 30, 0, 500000000, time.UTC)},
		{ID: 7, Name: "lb", SystemCount: 4, CreatedAt: time.Date(2026, 5, 10, 23, 59, 59, 0, time.UTC)},
	}
}

func groupIDs(groups []testGroup) []int {
	res := make([]int, len(groups))
	for i, g := range groups {
		res[i] = g.ID
	}
	return res
}

func TestDataProcessor(t *testing.T) {
	cases := []struct {
		name string
		opts ListOptions
		want []int
		err  string
	}{
		{"sort", ListOptions{SortBy: "name"}, []int{4, 5, 9, 6, 2, 7, 8, 1, 3}, ""},
		{"sort desc", ListOptions{SortBy: "system_count", SortDescending: true}, []int{2, 5, 8, 9, 1, 7, 4, 6, 3}, ""},
		{"contains", ListOptions{FilterKey: "name", FilterOperator: "~", FilterValue: "web"}, []int{3, 1}, ""},
		{
			"not contains",
			ListOptions{FilterKey: "name", FilterOperator: "!~", FilterValue: "web"},
			[]int{8, 4, 5, 9, 2, 6, 7},
			"",
		},
		{"gte", ListOptions{FilterKey: "system_count", FilterOperator: ">=", FilterValue: "5"}, []int{1, 8, 5, 9, 2}, ""},
		{"first page", ListOptions{PageNumber: 0, PageSize: 2}, []int{3, 1}, ""},
		{"second page", ListOptions{PageNumber: 1, PageSize: 2}, []int{8, 4}, ""},
		{"page past the end", ListOptions{PageNumber: 5, PageSize: 2}, []int{}, ""},
		{"no page size", ListOptions{PageNumber: 1, PageSize: 0}, []int{3, 1, 8, 4, 5, 9, 2, 6, 7}, ""},
		{
			"page after filter+sort",
			ListOptions{
				FilterKey: "id", FilterOperator: "!=", FilterValue: "1", SortBy: "name", SortDescending: true,
				PageNumber: 0, PageSize: 1,
			},
			[]int{3},
			"",
		},
		{
			"second page after filter+sort",
			ListOptions{
				FilterKey: "id", FilterOperator: "!=", FilterValue: "1", SortBy: "name", SortDescending: true,
				PageNumber: 1, PageSize: 1,
			},
			[]int{8},
			"",
		},
		{
			"filter+sort",
			ListOptions{
				FilterKey: "id", FilterOperator: "!=", FilterValue: "1", SortBy: "name", SortDescending: true,
			},
			[]int{3, 8, 7, 2, 6, 9, 5, 4},
			"",
		},
		{
			"bad field",
			ListOptions{SortBy: "nope"},
			[]int{},
			"unsupported field \"nope\"; supported fields: id, name, description, org_id, system_count, created_at",
		},
		{
			"bad value",
			ListOptions{FilterKey: "id", FilterOperator: "=", FilterValue: "x"},
			[]int{},
			"invalid filter on \"id\": strconv.ParseInt: parsing \"x\": invalid syntax",
		},
		{
			"~ on int",
			ListOptions{FilterKey: "id", FilterOperator: "~", FilterValue: "1"},
			[]int{},
			"invalid filter on \"id\": operator ~ is only supported on string fields",
		},
		{"sort dates", ListOptions{SortBy: "created_at"}, []int{1, 2, 3, 4, 5, 6, 7, 8, 9}, ""},
		{"sort dates desc", ListOptions{SortBy: "created_at", SortDescending: true}, []int{9, 8, 7, 6, 4, 5, 3, 2, 1}, ""},
		{
			"date eq",
			ListOptions{FilterKey: "created_at", FilterOperator: "=", FilterValue: "2026-05-10"},
			[]int{3, 4, 5, 2, 6, 7},
			"",
		},
		{
			"date neq",
			ListOptions{FilterKey: "created_at", FilterOperator: "!=", FilterValue: "2026-05-10"},
			[]int{1, 8, 9},
			"",
		},
		{
			"date lt",
			ListOptions{FilterKey: "created_at", FilterOperator: "<", FilterValue: "2026-05-10"},
			[]int{1},
			"",
		},
		{
			"date lte",
			ListOptions{FilterKey: "created_at", FilterOperator: "<=", FilterValue: "2026-05-10"},
			[]int{3, 1, 4, 5, 2, 6, 7},
			"",
		},
		{
			"date gt",
			ListOptions{FilterKey: "created_at", FilterOperator: ">", FilterValue: "2026-05-10"},
			[]int{8, 9},
			"",
		},
		{
			"date gte",
			ListOptions{FilterKey: "created_at", FilterOperator: ">=", FilterValue: "2026-05-10"},
			[]int{3, 8, 4, 5, 9, 2, 6, 7},
			"",
		},
		{
			"timestamp utc eq",
			ListOptions{FilterKey: "created_at", FilterOperator: "=", FilterValue: "2026-05-10T14:30:00Z"},
			[]int{4, 5},
			"",
		},
		{
			"timestamp positive tz offset eq",
			ListOptions{FilterKey: "created_at", FilterOperator: "=", FilterValue: "2026-05-10T16:30:00+02:00"},
			[]int{4, 5},
			"",
		},
		{
			"timestamp negative tz offset eq",
			ListOptions{FilterKey: "created_at", FilterOperator: "=", FilterValue: "2026-05-10T10:30:00-04:00"},
			[]int{4, 5},
			"",
		},
		{
			"timestamp fractional seconds eq",
			ListOptions{FilterKey: "created_at", FilterOperator: "=", FilterValue: "2026-05-10T14:30:00.5Z"},
			[]int{6},
			"",
		},
		{
			"timestamp without tz eq",
			ListOptions{FilterKey: "created_at", FilterOperator: "=", FilterValue: "2026-05-10T14:30:00"},
			[]int{4, 5},
			"",
		},
		{
			"timestamp without tz minutes only eq",
			ListOptions{FilterKey: "created_at", FilterOperator: "=", FilterValue: "2026-05-10T14:30"},
			[]int{4, 5},
			"",
		},
		{
			"timestamp neq",
			ListOptions{FilterKey: "created_at", FilterOperator: "!=", FilterValue: "2026-05-10T14:30:00Z"},
			[]int{3, 1, 8, 9, 2, 6, 7},
			"",
		},
		{
			"timestamp lt",
			ListOptions{FilterKey: "created_at", FilterOperator: "<", FilterValue: "2026-05-10T12:00:00Z"},
			[]int{1, 2},
			"",
		},
		{
			"timestamp lte",
			ListOptions{FilterKey: "created_at", FilterOperator: "<=", FilterValue: "2026-05-10T12:00:00Z"},
			[]int{3, 1, 2},
			"",
		},
		{
			"timestamp gt",
			ListOptions{FilterKey: "created_at", FilterOperator: ">", FilterValue: "2026-05-10T23:59:59Z"},
			[]int{8, 9},
			"",
		},
		{
			"timestamp gte",
			ListOptions{FilterKey: "created_at", FilterOperator: ">=", FilterValue: "2026-05-10T23:59:59Z"},
			[]int{8, 9, 7},
			"",
		},
		{
			"timestamp without tz lt",
			ListOptions{FilterKey: "created_at", FilterOperator: "<", FilterValue: "2026-05-10T12:00:00"},
			[]int{1, 2},
			"",
		},
		{
			"timestamp without tz gte",
			ListOptions{FilterKey: "created_at", FilterOperator: ">=", FilterValue: "2026-05-11T00:00:00"},
			[]int{8, 9},
			"",
		},
		{
			"invalid date string",
			ListOptions{FilterKey: "created_at", FilterOperator: "=", FilterValue: "not-a-date"},
			[]int{},
			"invalid filter on \"created_at\": invalid date or timestamp \"not-a-date\" " +
				"(YYYY-MM-DDTHH:MM:SS with optional timezone offset)",
		},
		{
			"invalid date month",
			ListOptions{FilterKey: "created_at", FilterOperator: "=", FilterValue: "2026-13-40"},
			[]int{},
			"invalid filter on \"created_at\": invalid date or timestamp \"2026-13-40\" " +
				"(YYYY-MM-DDTHH:MM:SS with optional timezone offset)",
		},
		{
			"invalid timestamp hour",
			ListOptions{FilterKey: "created_at", FilterOperator: "=", FilterValue: "2026-05-10T25:00:00"},
			[]int{},
			"invalid filter on \"created_at\": invalid date or timestamp \"2026-05-10T25:00:00\" " +
				"(YYYY-MM-DDTHH:MM:SS with optional timezone offset)",
		},
		{
			"contains on date",
			ListOptions{FilterKey: "created_at", FilterOperator: "~", FilterValue: "2026-05-10"},
			[]int{},
			"invalid filter on \"created_at\": operator ~ is only supported on string fields",
		},
		{
			"not contains on date",
			ListOptions{FilterKey: "created_at", FilterOperator: "!~", FilterValue: "2026-05-10"},
			[]int{},
			"invalid filter on \"created_at\": operator !~ is only supported on string fields",
		},
		{
			"unsupported operator on date",
			ListOptions{FilterKey: "created_at", FilterOperator: "^", FilterValue: "2026-05-10"},
			[]int{},
			"invalid filter on \"created_at\": unsupported operator \"^\"",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := NewDataProcessor(testGroups(), c.opts).Sort().Filter().Get()
			if c.err == "" {
				testutils.AssertNoError(t, "unexpected error", err)
			} else {
				testutils.AssertError(t, c.err, err)
			}
			testutils.AssertEquals(t, "unexpected result", c.want, groupIDs(got))
		})
	}
}

func TestParseISO8601(t *testing.T) {
	cases := []struct {
		name         string
		input        string
		wantTime     time.Time
		wantDateOnly bool
		wantErr      string
	}{
		{
			"date only",
			"2026-05-10",
			time.Date(2026, 5, 10, 0, 0, 0, 0, time.UTC),
			true,
			"",
		},
		{
			"timestamp UTC with Z",
			"2026-05-10T14:30:00Z",
			time.Date(2026, 5, 10, 14, 30, 0, 0, time.UTC),
			false,
			"",
		},
		{
			"timestamp with positive offset",
			"2026-05-10T16:30:00+02:00",
			time.Date(2026, 5, 10, 14, 30, 0, 0, time.UTC),
			false,
			"",
		},
		{
			"timestamp with negative offset",
			"2026-05-10T10:30:00-04:00",
			time.Date(2026, 5, 10, 14, 30, 0, 0, time.UTC),
			false,
			"",
		},
		{
			"timestamp with fractional seconds and offset",
			"2026-05-10T14:30:00.123456789Z",
			time.Date(2026, 5, 10, 14, 30, 0, 123456789, time.UTC),
			false,
			"",
		},
		{
			"timestamp without offset seconds",
			"2026-05-10T14:30:05",
			time.Date(2026, 5, 10, 14, 30, 5, 0, time.UTC),
			false,
			"",
		},
		{
			"timestamp without offset minutes",
			"2026-05-10T14:30",
			time.Date(2026, 5, 10, 14, 30, 0, 0, time.UTC),
			false,
			"",
		},
		{
			"empty string",
			"",
			time.Time{},
			false,
			"invalid date or timestamp \"\" (YYYY-MM-DDTHH:MM:SS with optional timezone offset)",
		},
		{
			"invalid text",
			"not-a-date",
			time.Time{},
			false,
			"invalid date or timestamp \"not-a-date\" (YYYY-MM-DDTHH:MM:SS with optional timezone offset)",
		},
		{
			"invalid month",
			"2026-13-01",
			time.Time{},
			false,
			"invalid date or timestamp \"2026-13-01\" (YYYY-MM-DDTHH:MM:SS with optional timezone offset)",
		},
		{
			"invalid day",
			"2026-05-32",
			time.Time{},
			false,
			"invalid date or timestamp \"2026-05-32\" (YYYY-MM-DDTHH:MM:SS with optional timezone offset)",
		},
		{
			"invalid separator slash",
			"2026/05/10",
			time.Time{},
			false,
			"invalid date or timestamp \"2026/05/10\" (YYYY-MM-DDTHH:MM:SS with optional timezone offset)",
		},
		{
			"invalid hour",
			"2026-05-10T25:00:00",
			time.Time{},
			false,
			"invalid date or timestamp \"2026-05-10T25:00:00\" (YYYY-MM-DDTHH:MM:SS with optional timezone offset)",
		},
		{
			"space separator unsupported",
			"2026-05-10 14:30:00",
			time.Time{},
			false,
			"invalid date or timestamp \"2026-05-10 14:30:00\" (YYYY-MM-DDTHH:MM:SS with optional timezone offset)",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			gotTime, isDateOnly, err := parseISO8601(c.input)
			if c.wantErr == "" {
				testutils.AssertNoError(t, "unexpected error", err)
				testutils.AssertTrue(t, "wrong parsed time", gotTime.Equal(c.wantTime))
				testutils.AssertEquals(t, "wrong isDateOnly", c.wantDateOnly, isDateOnly)
			} else {
				testutils.AssertError(t, c.wantErr, err)
			}
		})
	}
}
