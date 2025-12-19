// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package testutils

import (
	"reflect"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// AssertEquals ensures two values are equals and raises and error if not.
func AssertEquals[T any](t *testing.T, message string, expected T, actual T) {
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf(message+": got '%v' expected '%v'", actual, expected)
	}
}

// AssertTrue ensures a value is true and raises and error if not.
func AssertTrue(t *testing.T, message string, actual bool) {
	if !actual {
		t.Error(message)
	}
}

// AssertNoError ensures error was not produced.
func AssertNoError(t *testing.T, message string, err error) {
	if err != nil {
		t.Errorf(message+"err: %v", err)
	}
}

// AssertError ensures error mesasge was produced.
func AssertError(t *testing.T, message string, err error) {
	t.Helper() // Important: Marks this function as a test helper

	if err == nil {
		t.Fatal("Expected error but got success")
	}

	if message != "" && !strings.Contains(err.Error(), message) {
		t.Errorf("Expected error message to contain %q, got %q", message, err.Error())
	}
}

// AssertHasAllFlagsIgnores ensures that all but the ignored flags are present in the args slice.
func AssertHasAllFlagsIgnores(t *testing.T, cmd *cobra.Command, args []string, ignored []string) {
	// Some flags can be in the form --foo=bar, we only want to check the --foo part.
	noValueArgs := []string{}
	for _, arg := range args {
		noValueArgs = append(noValueArgs, strings.SplitN(arg, "=", 2)[0])
	}

	cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		flagString := "--" + flag.Name
		if !contains(ignored, flagString) && !contains(noValueArgs, flagString) {
			t.Error("Untested flag " + flagString)
		}
	})
}

// AssertHasAllFlags ensures that all the flags of a command are present in the args slice.
func AssertHasAllFlags(t *testing.T, cmd *cobra.Command, args []string) {
	AssertHasAllFlagsIgnores(t, cmd, args, []string{})
}

// AssertContains ensures a slice contains the expected value.
func AssertContains(t *testing.T, message string, actual []string, expected string) {
	if !contains(actual, expected) {
		t.Error(message)
	}
}

// AssertNotContains ensures a slice contains the expected value.
func AssertNotContains(t *testing.T, message string, actual []string, expected string) {
	if contains(actual, expected) {
		t.Error(message)
	}
}

// contains is copied from utils to avoid to dependency loop.
func contains(slice []string, needle string) bool {
	for _, item := range slice {
		if item == needle {
			return true
		}
	}
	return false
}
