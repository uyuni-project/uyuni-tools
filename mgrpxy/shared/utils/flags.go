// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// ProxyImageFlags are the flags used by install proxy command.
type ProxyImageFlags struct {
	ImagesLocation string           `mapstructure:"imagesLocation"`
	Tag            string           `mapstructure:"tag"`
	PullPolicy     string           `mapstructure:"pullPolicy"`
	Httpd          types.ImageFlags `mapstructure:"httpd"`
	SaltBroker     types.ImageFlags `mapstructure:"saltBroker"`
	Squid          types.ImageFlags `mapstructure:"squid"`
	Ssh            types.ImageFlags `mapstructure:"ssh"`
	Tftpd          types.ImageFlags `mapstructure:"tftpd"`
	Tuning         Tuning           `mapstructure:"tuning"`
}

// Tuning are the custom configuration file provide by users.
type Tuning struct {
	Httpd string `mapstructure:"httpd"`
	Squid string `mapstructure:"squid"`
}

// Get the full container image name and tag for a container name.
func (f *ProxyImageFlags) GetContainerImage(name string) string {
	imageName := "proxy-" + name
	image := fmt.Sprintf("%s/%s", f.ImagesLocation, imageName)
	tag := f.Tag

	var containerImage *types.ImageFlags
	switch name {
	case "httpd":
		containerImage = &f.Httpd
	case "salt-broker":
		containerImage = &f.SaltBroker
	case "squid":
		containerImage = &f.Squid
	case "ssh":
		containerImage = &f.Ssh
	case "tftpd":
		containerImage = &f.Tftpd
	default:
		log.Warn().Msgf(L("Invalid proxy container name: %s"), name)
	}

	if containerImage != nil {
		if containerImage.Name != "" {
			image = containerImage.Name
		}
		if containerImage.Tag != "" {
			tag = containerImage.Tag
		}
	}

	imageUrl, err := utils.ComputeImage(image, tag)
	if err != nil {
		log.Fatal().Err(err).Msg(L("failed to compute image URL"))
	}
	return imageUrl
}

// AddImageFlags will add the proxy install flags to a command.
func AddImageFlags(cmd *cobra.Command) {
	cmd.Flags().String("imagesLocation", utils.DefaultNamespace,
		L("registry URL prefix containing the all the container images"))
	cmd.Flags().String("tag", utils.DefaultTag, L("image tag"))
	utils.AddPullPolicyFlag(cmd)

	addContainerImageFlags(cmd, "httpd")
	addContainerImageFlags(cmd, "saltBroker")
	addContainerImageFlags(cmd, "squid")
	addContainerImageFlags(cmd, "ssh")
	addContainerImageFlags(cmd, "tftpd")

	cmd.Flags().String("tuning-httpd", "", L("HTTPD tuning configuration file"))
	cmd.Flags().String("tuning-squid", "", L("Squid tuning configuration file"))
}

func addContainerImageFlags(cmd *cobra.Command, container string) {
	cmd.Flags().String(container+"-image", "",
		fmt.Sprintf(L("Image for %s container, overrides the namespace if set"), container))
	cmd.Flags().String(container+"-tag", "",
		fmt.Sprintf(L("Tag for %s container, overrides the global value if set"), container))
}
