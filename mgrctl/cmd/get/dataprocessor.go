// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package get

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
)

// DataProcessor sorts and filters a slice of structs according to ListOptions.
//
// It works in place on the slice it is given: no copy of the data is made, and filtering compacts
// the slice inside its own backing array. The caller must not rely on the original slice content
// after calling Sort() or Filter(), and should use the slice returned by Get().
//
// Fields are looked up using their json tag name (e.g. "system_count").
type DataProcessor[T any] struct {
	data []T
	opts ListOptions
	err  error
}

// NewDataProcessor creates a processor for the given data and options.
func NewDataProcessor[T any](data []T, opts ListOptions) *DataProcessor[T] {
	p := &DataProcessor[T]{data: data, opts: opts}
	if reflect.TypeOf((*T)(nil)).Elem().Kind() != reflect.Struct {
		p.err = errors.New(L("data processor only supports slices of structs"))
	}
	return p
}

// Get returns the processed data and the first error met by Sort() or Filter().
// The data is then cut to the requested page, using opts.PageNumber (starting at 0) and opts.PageSize.
// Nothing is cut if the page size is not positive.
func (p *DataProcessor[T]) Get() ([]T, error) {
	if p.err != nil {
		return nil, p.err
	}
	return p.page(), nil
}

// page returns the slice of data covering the page described by opts.
func (p *DataProcessor[T]) page() []T {
	if !p.opts.hasPagination() {
		return p.data
	}
	start := max(p.opts.PageNumber, 0) * p.opts.PageSize
	if start >= len(p.data) {
		return p.data[:0]
	}
	end := min(start+p.opts.PageSize, len(p.data))
	return p.data[start:end]
}

// Sort sorts the data by opts.SortBy, in descending order if opts.SortDescending is set.
// Nothing is done if SortBy is empty.
func (p *DataProcessor[T]) Sort() *DataProcessor[T] {
	if p.err != nil || p.opts.SortBy == "" {
		return p
	}
	index, err := getFieldIndex[T](p.opts.SortBy)
	if err != nil {
		p.err = err
		return p
	}
	typ := reflect.TypeOf((*T)(nil)).Elem().Field(index).Type
	cmp, err := getComparator(typ)
	if err != nil {
		p.err = fmt.Errorf(L("sorting on field %q is not supported"), p.opts.SortBy)
		return p
	}

	data, desc := p.data, p.opts.SortDescending
	sort.SliceStable(data, func(i, j int) bool {
		// Addressing the elements avoids copying the structs.
		a := reflect.ValueOf(&data[i]).Elem().Field(index)
		b := reflect.ValueOf(&data[j]).Elem().Field(index)
		if desc {
			return cmp(b, a) < 0
		}
		return cmp(a, b) < 0
	})
	return p
}

// Filter keeps only the items matching opts.FilterKey, opts.FilterOperator and opts.FilterValue.
// Supported operators are =, !=, <, <=, >, >= and, for strings only, ~ (contains) and !~ (does not contain).
// An empty operator is treated as =. Nothing is done if there is no filter.
func (p *DataProcessor[T]) Filter() *DataProcessor[T] {
	if p.err != nil || !p.opts.hasFilter() {
		return p
	}
	index, err := getFieldIndex[T](p.opts.FilterKey)
	if err != nil {
		p.err = err
		return p
	}
	field := reflect.TypeOf((*T)(nil)).Elem().Field(index)
	match, err := buildMatcher(field.Type, p.opts.FilterOperator, p.opts.FilterValue)
	if err != nil {
		p.err = fmt.Errorf(L("invalid filter on %[1]q: %[2]w"), p.opts.FilterKey, err)
		return p
	}

	// Compact in place, reusing the backing array.
	data := p.data
	n := 0
	for i := range data {
		if match(reflect.ValueOf(&data[i]).Elem().Field(index)) {
			if n != i {
				data[n] = data[i]
			}
			n++
		}
	}
	// Zero the tail so the garbage collector can release what it references.
	var zero T
	for i := n; i < len(data); i++ {
		data[i] = zero
	}
	p.data = data[:n]
	return p
}

// getFieldIndex finds the index of the exported field of T matching name.
func getFieldIndex[T any](name string) (int, error) {
	t := reflect.TypeOf((*T)(nil)).Elem()
	names := make([]string, 0, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if !f.IsExported() {
			continue
		}
		jsonName, _, _ := strings.Cut(f.Tag.Get("json"), ",")
		if jsonName == "-" {
			jsonName = ""
		}
		if jsonName == "" {
			jsonName = f.Name
		}
		if name == jsonName {
			return i, nil
		}
		names = append(names, jsonName)
	}
	return 0, fmt.Errorf(L("unsupported field %[1]q; supported fields: %[2]s"), name, strings.Join(names, ", "))
}

type compareFunc func(a, b reflect.Value) int

func getComparator(typ reflect.Type) (compareFunc, error) {
	if typ == reflect.TypeOf(time.Time{}) {
		return func(a, b reflect.Value) int {
			ta := a.Interface().(time.Time)
			tb := b.Interface().(time.Time)
			return ta.Compare(tb)
		}, nil
	}

	switch typ.Kind() {
	case reflect.String:
		return func(a, b reflect.Value) int { return strings.Compare(a.String(), b.String()) }, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return func(a, b reflect.Value) int { return cmpOrdered(a.Int(), b.Int()) }, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return func(a, b reflect.Value) int { return cmpOrdered(a.Uint(), b.Uint()) }, nil
	case reflect.Float32, reflect.Float64:
		return func(a, b reflect.Value) int { return cmpOrdered(a.Float(), b.Float()) }, nil
	case reflect.Bool:
		return func(a, b reflect.Value) int {
			switch {
			case a.Bool() == b.Bool():
				return 0
			case !a.Bool():
				return -1
			}
			return 1
		}, nil
	}
	return nil, fmt.Errorf(L("unsupported type %s"), typ)
}

func cmpOrdered[V int64 | uint64 | float64](a, b V) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	}
	return 0
}

// buildMatcher parses the filter value once and returns the predicate to apply on a field.
func buildMatcher(typ reflect.Type, op string, value string) (func(reflect.Value) bool, error) {
	if op == "" {
		op = "="
	}

	if op == "~" || op == "!~" {
		if typ.Kind() != reflect.String {
			return nil, fmt.Errorf(L("operator %s is only supported on string fields"), op)
		}
		want := op == "~"
		return func(v reflect.Value) bool { return strings.Contains(v.String(), value) == want }, nil
	}

	var test func(c int) bool
	switch op {
	case "=":
		test = func(c int) bool { return c == 0 }
	case "!=":
		test = func(c int) bool { return c != 0 }
	case "<":
		test = func(c int) bool { return c < 0 }
	case "<=":
		test = func(c int) bool { return c <= 0 }
	case ">":
		test = func(c int) bool { return c > 0 }
	case ">=":
		test = func(c int) bool { return c >= 0 }
	default:
		return nil, fmt.Errorf(L("unsupported operator %q"), op)
	}

	if typ == reflect.TypeOf(time.Time{}) {
		return buildTimeMatcher(op, value, test)
	}

	switch typ.Kind() {
	case reflect.String:
		return func(v reflect.Value) bool { return test(strings.Compare(v.String(), value)) }, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		w, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil, err
		}
		return func(v reflect.Value) bool { return test(cmpOrdered(v.Int(), w)) }, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		w, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return nil, err
		}
		return func(v reflect.Value) bool { return test(cmpOrdered(v.Uint(), w)) }, nil
	case reflect.Float32, reflect.Float64:
		w, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, err
		}
		return func(v reflect.Value) bool { return test(cmpOrdered(v.Float(), w)) }, nil
	case reflect.Bool:
		w, err := strconv.ParseBool(value)
		if err != nil {
			return nil, err
		}
		if op != "=" && op != "!=" {
			return nil, fmt.Errorf(L("operator %s is not supported on boolean fields"), op)
		}
		return func(v reflect.Value) bool { return (v.Bool() == w) == (op == "=") }, nil
	}
	return nil, fmt.Errorf(L("unsupported type %s"), typ)
}

// buildTimeMatcher parses an ISO-8601 date or timestamp (with or without timezone offset)
// and builds a matcher for time.Time fields.
func buildTimeMatcher(op string, value string, test func(int) bool) (func(reflect.Value) bool, error) {
	parsed, isDateOnly, err := parseISO8601(value)
	if err != nil {
		return nil, err
	}

	if isDateOnly {
		// Treat a date-only filter as the full 24-hour calendar day [start, end).
		start := parsed
		end := parsed.AddDate(0, 0, 1)
		return func(v reflect.Value) bool {
			tv := v.Interface().(time.Time)
			switch op {
			case "=":
				return !tv.Before(start) && tv.Before(end)
			case "!=":
				return tv.Before(start) || !tv.Before(end)
			case "<":
				return tv.Before(start)
			case "<=":
				return tv.Before(end)
			case ">":
				return !tv.Before(end)
			case ">=":
				return !tv.Before(start)
			}
			return false
		}, nil
	}

	return func(v reflect.Value) bool {
		tv := v.Interface().(time.Time)
		return test(tv.Compare(parsed))
	}, nil
}

// parseISO8601 parses ISO-8601 date and timestamp strings with or without a timezone offset.
// Inputs without a timezone offset default to UTC.
func parseISO8601(value string) (time.Time, bool, error) {
	// Date-only: YYYY-MM-DD
	if t, err := time.ParseInLocation(time.DateOnly, value, time.UTC); err == nil {
		return t, true, nil
	}

	// Timestamp with timezone offset (2026-05-10T14:30:00Z or +02:00)
	if t, err := time.Parse(time.RFC3339Nano, value); err == nil {
		return t, false, nil
	}

	// Timestamp without timezone offset (with 'T' or space separator, optional fractional seconds)
	noTZLayouts := []string{
		"2006-01-02T15:04:05",
		"2006-01-02T15:04",
	}
	for _, layout := range noTZLayouts {
		if t, err := time.ParseInLocation(layout, value, time.UTC); err == nil {
			return t, false, nil
		}
	}

	return time.Time{}, false,
		fmt.Errorf(L("invalid date or timestamp %q (YYYY-MM-DDTHH:MM:SS with optional timezone offset)"), value)
}
