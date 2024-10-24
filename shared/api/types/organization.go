// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package types

// Organization describe an organization in the API.
type Organization struct {
	ID                    int
	Name                  string
	ActiveUsers           int `mapstructure:"active_users"`
	Systems               int
	Trusts                int
	SystemGroups          int  `mapstructure:"system_groups"`
	ActivationKeys        int  `mapstructure:"activation_keys"`
	KickstartProfiles     int  `mapstructure:"kickstart_profiles"`
	ConfigurationChannels int  `mapstructure:"configuration_channels"`
	StagingContentEnabled bool `mapstructure:"staging_content_enabled"`
}
