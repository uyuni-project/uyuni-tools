// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// ProxyImageFlags are the flags used by install proxy command.
type ProxyImageFlags struct {
	Registry   types.Registry       `mapstructure:"registry"`
	Tag        string               `mapstructure:"tag"`
	PullPolicy string               `mapstructure:"pullPolicy"`
	Httpd      types.ImageFlags     `mapstructure:"httpd"`
	SaltBroker types.ImageFlags     `mapstructure:"saltBroker"`
	Squid      types.ImageFlags     `mapstructure:"squid"`
	SSH        types.ImageFlags     `mapstructure:"ssh"`
	Tftpd      types.ImageFlags     `mapstructure:"tftpd"`
	Tuning     Tuning               `mapstructure:"tuning"`
	SCC        types.SCCCredentials `mapstructure:"scc"`
}

// Tuning are the custom configuration file provide by users.
type Tuning struct {
	Httpd string `mapstructure:"httpd"`
	Squid string `mapstructure:"squid"`
	SSH   string `mapstructure:"ssh"`
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

	imageURL, err := utils.ComputeImage(f.Registry.Host, f.Tag, *containerImage)
	if err != nil {
		log.Fatal().Err(err).Msg(L("failed to compute image URL"))
	}
	return imageURL
}

// AddSCCFlag add SCC flags to a command.
func AddSCCFlag(cmd *cobra.Command) {
	cmd.Flags().String("scc-user", "",
		L("SUSE Customer Center username. It will be used to pull images from SCC registry"),
	)
	cmd.Flags().String("scc-password", "",
		L("SUSE Customer Center password. It will be used to pull images from SCC registry"),
	)
	_ = utils.AddFlagHelpGroup(cmd, &utils.Group{ID: "scc", Title: L("SUSE Customer Center Flags")})
	_ = utils.AddFlagToHelpGroupID(cmd, "scc-user", "scc")
	_ = utils.AddFlagToHelpGroupID(cmd, "scc-password", "scc")
}

// AddImageFlags will add the proxy install flags to a command.
func AddImageFlags(cmd *cobra.Command) {
	cmd.Flags().String("tag", utils.DefaultTag, L("image tag"))
	utils.AddPullPolicyFlag(cmd)

	addContainerImageFlags(cmd, "httpd", "httpd")
	addContainerImageFlags(cmd, "saltbroker", "salt-broker")
	addContainerImageFlags(cmd, "squid", "squid")
	addContainerImageFlags(cmd, "ssh", "ssh")
	utils.AddTFTPDFlags(cmd, false, "")

	cmd.Flags().String("tuning-httpd", "", L("HTTPD tuning configuration file"))
	cmd.Flags().String("tuning-squid", "", L("Squid tuning configuration file"))
	cmd.Flags().String("tuning-ssh", "", L("SSH server tuning configuration file"))
	utils.AddRegistryFlag(cmd)
}

func addContainerImageFlags(cmd *cobra.Command, paramName string, imageName string) {
	utils.AddContainerImageFlags(cmd, paramName, imageName, "", "proxy-"+imageName)
}
