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

func TestFieldValue(t *testing.T) {
	item := map[string]any{
		"id":   123,
		"name": "web-server",
		"status": map[string]any{
			"state": "running",
		},
	}

	tests := []struct {
		path     string
		expected any
		found    bool
	}{
		{"id", 123, true},
		{"name", "web-server", true},
		{"status.state", "running", true},
		{"missing", nil, false},
		{"status.missing", nil, false},
	}

	for _, tt := range tests {
		result, ok := fieldValue(item, tt.path)
		if ok != tt.found {
			t.Errorf("Expected found=%v for path %q, got %v", tt.found, tt.path, ok)
		}
		if result != tt.expected {
			t.Errorf("Expected %v for path %q, got %v", tt.expected, tt.path, result)
		}
	}
}

func TestPrintTable(t *testing.T) {
	items := []map[string]any{
		{"id": float64(1), "name": "server-01"},
		{"id": float64(2), "name": "server-02"},
	}
	cols := []ColumnDef{
		{Header: "ID", Field: "id"},
		{Header: "NAME", Field: "name"},
	}

	var buf bytes.Buffer
	err := printTable(items, cols, &buf)
	if err != nil {
		t.Fatalf("printTable returned error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "ID") {
		t.Error("Expected output to contain header 'ID'")
	}
	if !strings.Contains(output, "NAME") {
		t.Error("Expected output to contain header 'NAME'")
	}
	if !strings.Contains(output, "server-01") {
		t.Error("Expected output to contain 'server-01'")
	}
	if !strings.Contains(output, "server-02") {
		t.Error("Expected output to contain 'server-02'")
	}
	if !strings.Contains(output, "1") {
		t.Error("Expected output to contain '1'")
	}
}
