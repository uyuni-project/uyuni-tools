// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package types

// Registry represents the flags used to set the registry.
type RegistryFlags struct {
	Server string `mapstructure:"server"`
	Path   string `mapstructure:"path"`
}

// ImageFlags represents the flags used by an image.
type ImageFlags struct {
	Name       string `mapstructure:"image"`
	Tag        string `mapstructure:"tag"`
	PullPolicy string `mapstructure:"pullPolicy"`
}

// ImageMetadata represents the image metadata of an RPM image.
type ImageMetadata struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
	File string   `json:"file"`
}

// Metadata represents the metadata of an RPM image.
type Metadata struct {
	Image ImageMetadata `json:"image"`
}
