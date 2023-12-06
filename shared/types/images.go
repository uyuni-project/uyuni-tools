// SPDX-FileCopyrightText: 2023 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package types

type ImageFlags struct {
	Name       string `mapstructure:"image"`
	Tag        string `mapstructure:"tag"`
	PullPolicy string `mapstructure:"pullPolicy"`
}
