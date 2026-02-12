// SPDX-FileCopyrightText: 2024-2025 SUSE LLC and contributors
//
// SPDX-License-Identifier: Apache-2.0

package types

// GlobalFlags represents the flags used by all commands.
type GlobalFlags struct {
	ConfigPath  string
	LogLevel    string
	KeepTempDir bool
}
