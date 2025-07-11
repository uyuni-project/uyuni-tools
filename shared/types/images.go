// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package types

import (
	"net/url"
	"strings"

	"github.com/rs/zerolog/log"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
)

// ImageFlags represents the flags used by an image.
type ImageFlags struct {
	Registry     string `mapstructure:"registry"`
	RegistryFQDN string `mapstructure:"registryFQDN"`
	Name         string `mapstructure:"image"`
	Tag          string `mapstructure:"tag"`
	PullPolicy   string `mapstructure:"pullPolicy"`
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
}

// RegistryFQDN return the registry FQDN
func (flags *ImageFlags) GetRegistryFQDN() string {
	reg := flags.Registry

	hasScheme := strings.Contains(reg, "://")
	toParse := reg
	if !hasScheme {
		toParse = "dummy://" + reg
	}

	u, err := url.Parse(toParse)
	if err != nil {
		log.Error().Msgf(L("Cannot extract FQDN from %s: this will be used as FQDN"))
		flags.RegistryFQDN = reg
		return flags.RegistryFQDN
	}

	if hasScheme {
		flags.RegistryFQDN = u.Scheme + "://" + u.Host
		return flags.RegistryFQDN
	}

	flags.RegistryFQDN = u.Host
	return flags.RegistryFQDN
}
