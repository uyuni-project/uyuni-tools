// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"testing"

	"github.com/uyuni-project/uyuni-tools/shared/testutils"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func TestNewCommand(t *testing.T) {
	var globalflags types.GlobalFlags
	cmd := NewCommand(&globalflags)
	testutils.AssertTrue(t, "Unexpected nil command", cmd != nil)
}
