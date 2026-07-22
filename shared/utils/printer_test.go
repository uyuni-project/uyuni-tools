// SPDX-FileCopyrightText: 2026 Jayprakash
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"bytes"
	"reflect"
	"strings"
	"testing"
)

func TestParseCustomColumns(t *testing.T) {
	spec := "ID:.id,NAME:.name,Nested Field:.status.state"
	expected := []ColumnDef{
		{Header: "ID", Field: "id"},
		{Header: "NAME", Field: "name"},
		{Header: "Nested Field", Field: "status.state"},
	}

	result := parseCustomColumns(spec)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

type status struct {
	State string
}

type testItem struct {
	ID     int
	Name   string
	Status status
}

func TestFieldValue(t *testing.T) {
	item := testItem{
		ID:   123,
		Name: "web-server",
		Status: status{
			State: "running",
		},
	}

	tests := []struct {
		description string
		path        string
		expected    any
		found       bool
	}{
		{"Field with an int value", "ID", 123, true},
		{"Field with a string value", "Name", "web-server", true},
		{"Nested field", "Status.State", "running", true},
		{"Non existent field", "Missing", nil, false},
		{"Non existent nested field", "Status.missing", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			result, ok := fieldValue(item, tt.path)
			if ok != tt.found {
				t.Errorf("Expected found=%v for path %q, got %v", tt.found, tt.path, ok)
			}
			if result != tt.expected {
				t.Errorf("Expected %v for path %q, got %v", tt.expected, tt.path, result)
			}
		})
	}
}

func TestPrintTable(t *testing.T) {
	items := []testItem{
		{ID: 1, Name: "server-01", Status: status{State: "running"}},
		{ID: 2, Name: "server-02", Status: status{State: "stopped"}},
	}
	cols := []ColumnDef{
		{Header: "ID", Field: "ID"},
		{Header: "NAME", Field: "Name"},
		{Header: "STATUS", Field: "Status.State"},
	}

	var buf bytes.Buffer
	err := printTable(items, cols, &buf)
	if err != nil {
		t.Fatalf("printTable returned error: %v", err)
	}

	output := buf.String()
	for _, word := range []string{"ID", "NAME", "server-01", "server-02", "1", "running", "stopped"} {
		if !strings.Contains(output, word) {
			t.Errorf("Expected output to contain header '%s'", word)
		}
	}
}
