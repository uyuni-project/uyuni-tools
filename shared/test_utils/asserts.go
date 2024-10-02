// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package test_utils

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
		t.Errorf(message)
	}
}

// AssertHasAllFlags ensures that all the flags of a command are present in the args slice.
func AssertHasAllFlags(t *testing.T, cmd *cobra.Command, args []string) {
	// Some flags can be in the form --foo=bar, we only want to check the --foo part.
	noValueArgs := []string{}
	for _, arg := range args {
		noValueArgs = append(noValueArgs, strings.SplitN(arg, "=", 2)[0])
	}

	cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		flagString := "--" + flag.Name
		if !contains(noValueArgs, flagString) {
			t.Error("Untested flag " + flagString)
		}
	})
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
