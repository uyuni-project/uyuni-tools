// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package cmd

import "testing"

func TestShouldLogToConsoleForInteractiveSQL(t *testing.T) {
	rootCmd, err := NewUyuniadmCommand()
	if err != nil {
		t.Fatalf("failed to create command: %v", err)
	}

	sqlCmd, _, err := rootCmd.Find([]string{"support", "sql"})
	if err != nil {
		t.Fatalf("failed to find sql command: %v", err)
	}

	if err := sqlCmd.ParseFlags([]string{"--interactive"}); err != nil {
		t.Fatalf("failed to parse flags: %v", err)
	}

	if shouldLogToConsole(sqlCmd, "debug") {
		t.Error("interactive sql should not log to console in debug mode")
	}
	if shouldLogToConsole(sqlCmd, "trace") {
		t.Error("interactive sql should not log to console in trace mode")
	}
	if !shouldLogToConsole(sqlCmd, "info") {
		t.Error("interactive sql should still log to console in info mode")
	}
	if !shouldLogToConsole(sqlCmd, "") {
		t.Error("interactive sql should still log to console with default log level")
	}
}

func TestShouldLogToConsoleForNonInteractiveSQL(t *testing.T) {
	rootCmd, err := NewUyuniadmCommand()
	if err != nil {
		t.Fatalf("failed to create command: %v", err)
	}

	sqlCmd, _, err := rootCmd.Find([]string{"support", "sql"})
	if err != nil {
		t.Fatalf("failed to find sql command: %v", err)
	}

	if shouldLogToConsole(sqlCmd, "debug") == false {
		t.Error("non-interactive sql should keep console logging")
	}
}
