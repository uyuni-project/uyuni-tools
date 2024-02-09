// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package types

// Distribution contains information about the distribution.
type Distribution struct {
	TreeLabel    string
	BasePath     string
	ChannelLabel string
	InstallType  string
	Arch         string
}
