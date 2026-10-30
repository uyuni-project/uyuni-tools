// SPDX-FileCopyrightText: 2026 Jay Prakash katara <katarajayprakash@icloud.com>
//
// SPDX-License-Identifier: Apache-2.0

package get

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type SystemGroup struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	OrgID       int    `json:"org_id"`
	SystemCount int    `json:"system_count"`
}

// systemGroupFields lists the fields supported for both filtering and
// sorting, so validation, errors and help cannot advertise different fields.
var systemGroupFields = []string{
	"id",
	"name",
	"description",
	"org_id",
	"system_count",
}

var systemGroupFilterOperators = []string{"=", "!=", ">", "<", ">=", "<="}

func systemGroupComparator(sortBy string) (func(a, b SystemGroup) bool, error) {
	switch sortBy {
	case "":
		return nil, nil
	case "id":
		return func(a, b SystemGroup) bool {
			return a.ID < b.ID
		}, nil
	case "name":
		return func(a, b SystemGroup) bool {
			return a.Name < b.Name
		}, nil
	case "description":
		return func(a, b SystemGroup) bool {
			return a.Description < b.Description
		}, nil
	case "org_id":
		return func(a, b SystemGroup) bool {
			return a.OrgID < b.OrgID
		}, nil
	case "system_count":
		return func(a, b SystemGroup) bool {
			return a.SystemCount < b.SystemCount
		}, nil
	default:
		return nil, fmt.Errorf(
			L("unsupported sort field %[1]q for system groups; supported fields: %[2]s"),
			sortBy,
			strings.Join(systemGroupFields, ", "),
		)
	}
}

// parseFilterExpr splits a single field/operator/value filter expression.
func parseFilterExpr(expr string) (string, string, string, error) {
	field, op, value := "", "", ""
	// Find the first operator character; a two-character operator takes
	// precedence over its one-character prefix.
	i := strings.IndexAny(expr, "=!<>")
	if i >= 0 && i+1 < len(expr) {
		if two := expr[i : i+2]; two == ">=" || two == "<=" || two == "!=" {
			op = two
		}
	}
	if op == "" && i >= 0 && (expr[i] == '=' || expr[i] == '>' || expr[i] == '<') {
		op = expr[i : i+1]
	}
	if op != "" {
		field = strings.TrimSpace(expr[:i])
		value = strings.TrimSpace(expr[i+len(op):])
	}
	if op == "" {
		if i >= 0 {
			return "", "", "", fmt.Errorf(
				L("unsupported filter operator %[1]q; supported operators: %[2]s"),
				string(expr[i]),
				strings.Join(systemGroupFilterOperators, ", "),
			)
		}
		return "", "", "", fmt.Errorf(
			L("invalid filter expression %[1]q: expected an operator (%[2]s)"),
			expr,
			strings.Join(systemGroupFilterOperators, ", "),
		)
	}
	if field == "" {
		return "", "", "", fmt.Errorf(L("filter field cannot be empty in %q"), expr)
	}
	for _, c := range field {
		if !unicode.IsLetter(c) && !unicode.IsDigit(c) && c != '_' && c != '-' {
			return "", "", "", fmt.Errorf(
				L("invalid filter field %q: only letters, digits, underscores and dashes are allowed"),
				field,
			)
		}
	}
	// The value may be empty: string equality can match an empty description,
	// and numericFilter reports a clearer integer error for numeric fields.
	return field, op, value, nil
}

func numericCompare(op string) (func(a, b int) bool, error) {
	switch op {
	case "=":
		return func(a, b int) bool { return a == b }, nil
	case "!=":
		return func(a, b int) bool { return a != b }, nil
	case ">":
		return func(a, b int) bool { return a > b }, nil
	case "<":
		return func(a, b int) bool { return a < b }, nil
	case ">=":
		return func(a, b int) bool { return a >= b }, nil
	case "<=":
		return func(a, b int) bool { return a <= b }, nil
	}
	return nil, fmt.Errorf(
		L("unsupported filter operator %[1]q; supported operators: %[2]s"),
		op,
		strings.Join(systemGroupFilterOperators, ", "),
	)
}

// numericFilter matches a numeric system group field against an integer
// filter value. All operators are supported.
func numericFilter(op, value string, get func(SystemGroup) int) (func(SystemGroup) bool, error) {
	compare, err := numericCompare(op)
	if err != nil {
		return nil, err
	}
	num, err := strconv.Atoi(value)
	if err != nil {
		return nil, fmt.Errorf(L("filter value %q must be an integer"), value)
	}
	return func(g SystemGroup) bool { return compare(get(g), num) }, nil
}

// stringFilter matches a string system group field with exact,
// case-sensitive equality. Ordering operators have no meaningful string
// semantics here and are rejected.
func stringFilter(op, value string, get func(SystemGroup) string) (func(SystemGroup) bool, error) {
	switch op {
	case "=":
		return func(g SystemGroup) bool { return get(g) == value }, nil
	case "!=":
		return func(g SystemGroup) bool { return get(g) != value }, nil
	}
	return nil, fmt.Errorf(
		L("unsupported filter operator %[1]q for a string field; supported operators: =, !="),
		op,
	)
}

// systemGroupFilter validates a filter expression and returns a predicate
// applied client-side to the API response.
func systemGroupFilter(expr string) (func(SystemGroup) bool, error) {
	field, op, value, err := parseFilterExpr(expr)
	if err != nil {
		return nil, err
	}
	switch field {
	case "id":
		return numericFilter(op, value, func(g SystemGroup) int { return g.ID })
	case "org_id":
		return numericFilter(op, value, func(g SystemGroup) int { return g.OrgID })
	case "system_count":
		return numericFilter(op, value, func(g SystemGroup) int { return g.SystemCount })
	case "name":
		return stringFilter(op, value, func(g SystemGroup) string { return g.Name })
	case "description":
		return stringFilter(op, value, func(g SystemGroup) string { return g.Description })
	default:
		return nil, fmt.Errorf(
			L("unsupported filter field %[1]q for system groups; supported fields: %[2]s"),
			field,
			strings.Join(systemGroupFields, ", "),
		)
	}
}

func init() {
	registerResource[SystemGroup]("systemgroup", &systemGroupFetcher{}, []string{"grp"})
}

type systemGroupFetcher struct{}

func (*systemGroupFetcher) List(client *api.APIClient, opts ListOptions) ([]SystemGroup, int, error) {
	// The filter and the comparator are validated before the group-list
	// request so invalid expressions fail without calling the endpoint.
	var match func(SystemGroup) bool
	if opts.Filter != "" {
		predicate, err := systemGroupFilter(opts.Filter)
		if err != nil {
			return nil, 0, err
		}
		match = predicate
	}

	compare, err := systemGroupComparator(opts.SortBy)
	if err != nil {
		return nil, 0, err
	}

	res, err := api.GetChecked[[]SystemGroup](client, "systemgroup/listAllGroups", "api.systemgroup.list_all_groups")
	if err != nil {
		return nil, 0, err
	}

	// The API has no filtering or pagination support: the predicate is
	// applied client-side, page options are ignored, and group counts are
	// naturally small enough.
	var groups []SystemGroup
	if opts.Filter != "" {
		groups = make([]SystemGroup, 0, len(res.Result))
		for _, group := range res.Result {
			if match(group) {
				groups = append(groups, group)
			}
		}
	} else {
		groups = res.Result
	}

	if compare != nil {
		sort.SliceStable(groups, func(i, j int) bool {
			if opts.SortDescending {
				return compare(groups[j], groups[i])
			}
			return compare(groups[i], groups[j])
		})
	}

	return groups, len(groups), nil
}

func (*systemGroupFetcher) Columns() []utils.ColumnDef {
	return []utils.ColumnDef{
		{Header: "ID", Field: "ID"},
		{Header: "NAME", Field: "Name"},
		{Header: "DESCRIPTION", Field: "Description"},
		{Header: "ORG_ID", Field: "OrgID"},
		{Header: "SYSTEM_COUNT", Field: "SystemCount"},
	}
}

func (*systemGroupFetcher) Help() ResourceHelp {
	return ResourceHelp{
		Summary: L("List system groups"),
		Details: fmt.Sprintf(
			L(`System groups support client-side filtering and sorting by:
  %[1]s
Filtering matches string fields (name, description) with exact,
case-sensitive equality (=, !=); ordering operators (>, <, >=, <=)
apply to the numeric fields only. Pagination is ignored.`),
			strings.Join(systemGroupFields, ", "),
		),
		Examples: L(`  # List all system groups
  mgrctl get systemgroup

  # Sort system groups by name
  mgrctl get systemgroup --sort-by=name

  # List system groups with ten or more systems
  mgrctl get systemgroup --filter 'system_count>=10'`),
	}
}
