// SPDX-FileCopyrightText: 2023 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type ProxyInstallFlags struct {
	ImagesLocation string           `mapstructure:"imagesLocation"`
	Tag            string           `namespace:"tag"`
	PullPolicy     string           `mapstructure:"pullPolicy"`
	Httpd          types.ImageFlags `mapstructure:"httpd"`
	SaltBroker     types.ImageFlags `mapstructure:"saltBroker"`
	Squid          types.ImageFlags `mapstructure:"squid"`
	Ssh            types.ImageFlags `mapstructure:"ssh"`
	Tftpd          types.ImageFlags `mapstructure:"tftpd"`
}

// Get the full container image name and tag for a container name.
func (f *ProxyInstallFlags) GetContainerImage(name string) string {
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
		log.Warn().Msgf("Invalid proxy container name: %s", name)
	}

	if containerImage != nil {
		if containerImage.Name != "" {
			image = containerImage.Name
		}
		if containerImage.Tag != "" {
			tag = containerImage.Tag
		}
	}

	return fmt.Sprintf("%s:%s", image, tag)
}

func AddInstallFlags(cmd *cobra.Command) {
	cmd.Flags().String("imagesLocation", utils.DefaultNamespace,
		"registry URL prefix containing the all the container images")
	cmd.Flags().String("tag", utils.DefaultTag, "Tag Image")
	utils.AddPullPolicyFlag(cmd)

	addContainerImageFlags(cmd, "httpd")
	addContainerImageFlags(cmd, "saltBroker")
	addContainerImageFlags(cmd, "squid")
	addContainerImageFlags(cmd, "ssh")
	addContainerImageFlags(cmd, "tftpd")
}

func addContainerImageFlags(cmd *cobra.Command, container string) {
	cmd.Flags().String(container+"-image", "",
		fmt.Sprintf("Image for %s container, overrides the namespace if set", container))
	cmd.Flags().String(container+"-tag", "",
		fmt.Sprintf("Tag for %s container, overrides the global value if set", container))
}
