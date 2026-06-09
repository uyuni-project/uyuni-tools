// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package types

// System describes an Uyuni registered system/minion in the API.
type System struct {
	ID          int    `json:"id" mapstructure:"id"`
	Name        string `json:"name" mapstructure:"name"`
	LastCheckin string `json:"last_checkin" mapstructure:"last_checkin"`
	Created     string `json:"created" mapstructure:"created"`
	LastBoot    string `json:"last_boot" mapstructure:"last_boot"`
}
