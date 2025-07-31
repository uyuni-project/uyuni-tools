// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package types

// ImageFlags represents the flags used by an image.
type ImageFlags struct {
	Registry        string `mapstructure:"registry"`
	Name            string `mapstructure:"image"`
	Tag             string `mapstructure:"tag"`
	PullPolicy      string `mapstructure:"pullPolicy"`
	SkipComputation bool   // Internal field, not a command flag, indicates if the image should be computed or not
}

// PgsqlFlags contains settings for Pgsql container.
type PgsqlFlags struct {
	Replicas  int
	Image     ImageFlags `mapstructure:",squash"`
	IsChanged bool
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

// SCCCredentials can store SCC Credentials.
type SCCCredentials struct {
	User     string
	Password string
	Registry string
}
