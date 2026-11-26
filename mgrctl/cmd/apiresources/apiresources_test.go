// SPDX-FileCopyrightText: 2026 Jay Prakash katara <katarajayprakash@icloud.com>
//
// SPDX-License-Identifier: Apache-2.0

package apiresources

import (
	"bytes"
	"testing"

	"github.com/uyuni-project/uyuni-tools/mgrctl/cmd/get"
)

func TestPrintResources(t *testing.T) {
	resources := []get.ResourceInfo{{
		Name:    "system",
		Aliases: []string{"systems"},
		Help:    get.ResourceHelp{Summary: "Managed systems"},
	}}

	var output bytes.Buffer
	if err := printResources(&output, resources); err != nil {
		t.Fatalf("printResources() error = %v", err)
	}

	want := "NAME      ALIAS NAMES    DESCRIPTION\nsystem    systems        Managed systems\n"
	if output.String() != want {
		t.Errorf("printResources() = %q, want %q", output.String(), want)
	}
}

func TestPrintResourcesEmpty(t *testing.T) {
	var output bytes.Buffer
	if err := printResources(&output, nil); err != nil {
		t.Fatalf("printResources() error = %v", err)
	}

	if got, want := output.String(), "No resources registered.\n"; got != want {
		t.Errorf("printResources() = %q, want %q", got, want)
	}
}
