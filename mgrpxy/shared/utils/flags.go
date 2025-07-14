// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"fmt"
	"path"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// ProxyImageFlags are the flags used by install proxy command.
type ProxyImageFlags struct {
	Registry   string           `mapstructure:"registry"`
	Tag        string           `mapstructure:"tag"`
	PullPolicy string           `mapstructure:"pullPolicy"`
	Httpd      types.ImageFlags `mapstructure:"httpd"`
	SaltBroker types.ImageFlags `mapstructure:"saltBroker"`
	Squid      types.ImageFlags `mapstructure:"squid"`
	SSH        types.ImageFlags `mapstructure:"ssh"`
	Tftpd      types.ImageFlags `mapstructure:"tftpd"`
	Tuning     Tuning           `mapstructure:"tuning"`
}

// Tuning are the custom configuration file provide by users.
type Tuning struct {
	Httpd string `mapstructure:"httpd"`
	Squid string `mapstructure:"squid"`
}

func (flags *ProxyImageFlags) setTagIfMissing() {
	globalTag := utils.DefaultTag
	if flags.Tag != "" {
		globalTag = flags.Tag
	}
	if flags.Httpd.Tag == "" {
		flags.Httpd.Tag = globalTag
	}
	if flags.SSH.Tag == "" {
		flags.SSH.Tag = globalTag
	}
	if flags.SaltBroker.Tag == "" {
		flags.SaltBroker.Tag = globalTag
	}
	if flags.Squid.Tag == "" {
		flags.Squid.Tag = globalTag
	}
	if flags.Tftpd.Tag == "" {
		flags.Tftpd.Tag = globalTag
	}
}

func (flags *ProxyImageFlags) setRegistryIfMissing() {
	globalRegistry := utils.DefaultRegistry
	if flags.Registry != "" {
		globalRegistry = flags.Registry
	}
	if flags.Httpd.Registry == "" {
		flags.Httpd.Registry = globalRegistry
		flags.Httpd.GetRegistryFQDN()
	}
	if flags.SSH.Registry == "" {
		flags.SSH.Registry = globalRegistry
		flags.SSH.GetRegistryFQDN()
	}
	if flags.SaltBroker.Registry == "" {
		flags.SaltBroker.Registry = globalRegistry
		flags.SaltBroker.GetRegistryFQDN()
	}
	if flags.Squid.Registry == "" {
		flags.Squid.Registry = globalRegistry
		flags.Squid.GetRegistryFQDN()
	}
	if flags.Tftpd.Registry == "" {
		flags.Tftpd.Registry = globalRegistry
		flags.Tftpd.GetRegistryFQDN()
	}
}

// CheckParameters checks parameters for server parameters.
func (flags *ProxyImageFlags) CheckParameters() {
	flags.setRegistryIfMissing()
	flags.setTagIfMissing()
}

// GetContainerImage gets the full container image name and tag for a container name.
func (f *ProxyImageFlags) GetContainerImage(name string) string {
	var containerImage *types.ImageFlags
	switch name {
	case "httpd":
		containerImage = &f.Httpd
	case "salt-broker":
		containerImage = &f.SaltBroker
	case "squid":
		containerImage = &f.Squid
	case "ssh":
		containerImage = &f.SSH
	case "tftpd":
		containerImage = &f.Tftpd
	default:
		log.Fatal().Msgf(L("Invalid proxy container name: %s"), name)
	}
	if containerImage.Registry == "" {
		containerImage.Registry = f.Registry
	}
	if containerImage.Tag == "" {
		containerImage.Tag = f.Tag
	}

	imageURL, err := utils.ComputeImage(*containerImage)
	if err != nil {
		log.Fatal().Err(err).Msg(L("failed to compute image URL"))
	}
	return imageURL
}

// AddSCCFlag add SCC flags to a command.
func AddSCCFlag(cmd *cobra.Command) {
	cmd.Flags().String("scc-user", "",
		L("SUSE Customer Center username. It will be used to pull images from the registry"),
	)
	cmd.Flags().String("scc-password", "",
		L("SUSE Customer Center password. It will be used to pull images from the registry"),
	)

	_ = utils.AddFlagHelpGroup(cmd, &utils.Group{ID: "scc", Title: L("SUSE Customer Center Flags")})
	_ = utils.AddFlagToHelpGroupID(cmd, "scc-user", "scc")
	_ = utils.AddFlagToHelpGroupID(cmd, "scc-password", "scc")
}

// AddImageFlags will add the proxy install flags to a command.
func AddImageFlags(cmd *cobra.Command) {
	cmd.Flags().String("tag", utils.DefaultTag, L("image tag"))
	cmd.Flags().String("registry", utils.DefaultRegistry, L("Specify a registry where to pull the images from"))
	utils.AddPullPolicyFlag(cmd)

	addContainerImageFlags(cmd, "httpd", "httpd")
	addContainerImageFlags(cmd, "saltbroker", "salt-broker")
	addContainerImageFlags(cmd, "squid", "squid")
	addContainerImageFlags(cmd, "ssh", "ssh")
	addContainerImageFlags(cmd, "tftpd", "tftpd")

	cmd.Flags().String("tuning-httpd", "", L("HTTPD tuning configuration file"))
	cmd.Flags().String("tuning-squid", "", L("Squid tuning configuration file"))
}

func addContainerImageFlags(cmd *cobra.Command, paramName string, imageName string) {
	defaultImage := path.Join(utils.DefaultRegistry, "proxy-"+imageName)
	cmd.Flags().String(paramName+"-image", defaultImage,
		fmt.Sprintf(L("Image for %s container"), imageName))
	cmd.Flags().String(paramName+"-tag", "",
		fmt.Sprintf(L("Tag for %s container, overrides the global value if set"), imageName))
}
