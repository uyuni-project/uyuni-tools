// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package types

// ImageFlags represents the flags used by an image.
type ImageFlags struct {
	Name       string `mapstructure:"image"`
	Tag        string `mapstructure:"tag"`
	PullPolicy string `mapstructure:"pullPolicy"`
}
