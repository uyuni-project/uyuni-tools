// SPDX-FileCopyrightText: 2023 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package types

type Distribution struct {
	TreeLabel    string
	BasePath     string
	ChannelLabel string
	InstallType  string
	Arch         string
}
