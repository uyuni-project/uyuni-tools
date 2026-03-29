// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package get

import (
    "testing"

    "github.com/uyuni-project/uyuni-tools/shared/types"
)

func TestNewCommand(t *testing.T) {
    var globalFlags types.GlobalFlags
    cmd := NewCommand(&globalFlags)
    if cmd == nil {
        t.Fatal("NewCommand returned nil")
    }
    if cmd.Use != "get" {
        t.Errorf("expected Use 'get', got '%s'", cmd.Use)
    }
    // Verify system subcommand exists
    found := false
    for _, sub := range cmd.Commands() {
        if sub.Use == "system [name]" {
            found = true
            break
        }
    }
    if !found {
        t.Error("system subcommand not found")
    }
}

func TestPrintSystemsTable(t *testing.T) {
    systems := []System{
        {ID: 1, Name: "web-server", LastCheckin: "2026-03-26"},
    }
    // Should not error
    if err := printSystems(systems, "table"); err != nil {
        t.Errorf("unexpected error: %v", err)
    }
}

func TestPrintSystemsJSON(t *testing.T) {
    systems := []System{
        {ID: 1, Name: "web-server", LastCheckin: "2026-03-26"},
    }
    if err := printSystems(systems, "json"); err != nil {
        t.Errorf("unexpected error: %v", err)
    }
}

func TestPrintSystemsInvalidFormat(t *testing.T) {
    systems := []System{}
    if err := printSystems(systems, "invalid"); err == nil {
        t.Error("expected error for invalid format")
    }
}