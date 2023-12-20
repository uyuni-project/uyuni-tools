// SPDX-FileCopyrightText: 2023 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package distro

import (
	"fmt"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

var productMap = map[string]map[string]types.Distribution{
	"AlmaLinux": {
		"9": {
			InstallType:  "rhel_9",
			ChannelLabel: "almalinux9",
		},
		"8": {
			InstallType:  "rhel_8",
			ChannelLabel: "almalinux8",
		},
	},

	"SUSE Linux Enterprise": {
		"15 SP1": {
			InstallType:  "sles15generic",
			ChannelLabel: "sle-product-sles15-sp1-pool",
		},
		"15 SP2": {
			InstallType:  "sles15generic",
			ChannelLabel: "sle-product-sles15-sp2-pool",
		},
		"15 SP3": {
			InstallType:  "sles15generic",
			ChannelLabel: "sle-product-sles15-sp3-pool",
		},
		"15 SP4": {
			InstallType:  "sles15generic",
			ChannelLabel: "sle-product-sles15-sp4-pool",
		},
		"15 SP5": {
			InstallType:  "sles15generic",
			ChannelLabel: "sle-product-sles15-sp5-pool",
		},
		"12 SP5": {
			InstallType:  "sles12generic",
			ChannelLabel: "sles12-sp5-pool",
		},
	},

	"Red Hat Enterprise Linux": {
		"7": {
			InstallType:  "rhel_7",
			ChannelLabel: "rhel7-pool",
		},
		"8": {
			InstallType:  "rhel_8",
			ChannelLabel: "rhel8-pool",
		},
		"9": {
			InstallType:  "rhel_9",
			ChannelLabel: "rhel9-pool",
		},
	},
}

func getDistroFromDetails(distro string, version string, arch string, channeLabel string, flags *flagpole) (types.Distribution, error) {
	productFromConfig := flags.ProductMap
	var distribution types.Distribution
	var ok bool

	if productFromConfig[distro] != nil {
		distribution, ok = productFromConfig[distro][version]
	} else if productMap[distro] != nil {
		distribution, ok = productMap[distro][version]
	}

	if !ok {
		return types.Distribution{}, fmt.Errorf("unkown distribution, auto-registration is not possible")
	}

	if channeLabel != "" {
		distribution.ChannelLabel = channeLabel
	} else {
		distribution.ChannelLabel = fmt.Sprintf("%s-%s", distribution.ChannelLabel, arch)
	}

	return distribution, nil
}

func detectDistro(path string, channelLabel string, flags *flagpole, distro *types.Distribution) error {
	treeinfopath := filepath.Join(path, ".treeinfo")
	log.Trace().Msgf("Reading .treeinfo %s", treeinfopath)
	treeInfoViper := viper.New()
	treeInfoViper.SetConfigType("ini")
	treeInfoViper.SetConfigName(".treeinfo")
	treeInfoViper.AddConfigPath(path)
	if err := treeInfoViper.ReadInConfig(); err != nil {
		return err
	}

	dname := treeInfoViper.GetString("release.name")
	dversion := treeInfoViper.GetString("release.version")
	darch := treeInfoViper.GetString("general.arch")
	log.Debug().Msgf("Detected distro %s, version %s. arch %s", dname, dversion, darch)
	details, err := getDistroFromDetails(dname, dversion, darch, channelLabel, flags)
	if err != nil {
		return err
	}

	*distro = types.Distribution{
		InstallType:  details.InstallType,
		TreeLabel:    dname,
		ChannelLabel: details.ChannelLabel,
	}
	return nil
}
