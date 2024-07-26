// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package distro

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

var productMap = map[string]map[string]map[types.Arch]types.Distribution{
	"SUSE Linux Enterprise": {
		"15 SP4": {
			types.AMD64: {
				TreeLabel:    "SLES15SP4",
				InstallType:  "sles15generic",
				ChannelLabel: "sle-product-sles15-sp4-pool-x86_64",
			},
		},
		"15 SP5": {
			types.AMD64: {
				TreeLabel:    "SLES15SP5",
				InstallType:  "sles15generic",
				ChannelLabel: "sle-product-sles15-sp5-pool-x86_64",
			},
		},
		"15 SP6": {
			types.AMD64: {
				TreeLabel:    "SLES15SP6",
				InstallType:  "sles15generic",
				ChannelLabel: "sle-product-sles15-sp6-pool-x86_64",
			},
			types.AArch64: {
				TreeLabel:    "SLES15SP6",
				InstallType:  "sles15generic",
				ChannelLabel: "sle-product-sles15-sp6-pool-aarch64",
			},
		},
		"12 SP5": {
			types.AMD64: {
				TreeLabel:    "SLES12SP5",
				InstallType:  "sles12generic",
				ChannelLabel: "sles12-sp5-pool-x86_64",
			},
		},
	},

	"Red Hat Enterprise Linux": {
		"7": {
			types.AMD64: {
				TreeLabel:    "RHEL7",
				InstallType:  "rhel_7",
				ChannelLabel: "rhel7-pool-x86_64",
			},
		},
		"8": {
			types.AMD64: {
				TreeLabel:    "RHEL8",
				InstallType:  "rhel_8",
				ChannelLabel: "rhel8-pool-x86_64",
			},
		},
		"9": {
			types.AMD64: {
				TreeLabel:    "RHEL9",
				InstallType:  "rhel_9",
				ChannelLabel: "rhel9-pool-x86_64",
			},
		},
	},
}

func getDistroFromDetails(distro string, version string, arch types.Arch, flags *flagpole) (types.Distribution, error) {
	productFromConfig := flags.ProductMap
	var distribution types.Distribution
	var ok bool

	if productFromConfig[distro] != nil {
		distribution, ok = productFromConfig[distro][version][arch]
	} else if productMap[distro] != nil {
		distribution, ok = productMap[distro][version][arch]
	}

	if !ok {
		return types.Distribution{}, fmt.Errorf(L("distribution not found in product map. Please update productmap or provide channel label"))
	}
	return distribution, nil
}

func getDistroFromTreeinfo(path string, flags *flagpole) (types.Distribution, error) {
	treeinfopath := filepath.Join(path, ".treeinfo")
	log.Debug().Msgf("Reading .treeinfo %s", treeinfopath)
	treeInfoViper := viper.New()
	treeInfoViper.SetConfigType("ini")
	treeInfoViper.SetConfigName(".treeinfo")
	treeInfoViper.AddConfigPath(path)
	if err := treeInfoViper.ReadInConfig(); err != nil {
		return types.Distribution{}, fmt.Errorf(L("unable to read distribution treeinfo. Please provide distribution details and/or channel label"))
	}

	dname := treeInfoViper.GetString("release.name")
	dversion := treeInfoViper.GetString("release.version")
	darch := treeInfoViper.GetString("general.arch")
	log.Debug().Msgf("Detected distribution %s, version %s. arch %s", dname, dversion, darch)

	return getDistroFromDetails(dname, dversion, types.GetArch(darch), flags)
}

func detectDistro(path string, distroDetails types.DistributionDetails, flags *flagpole, distro *types.Distribution) error {
	treeinfopath := filepath.Join(path, ".treeinfo")
	channelLabel := flags.ChannelLabel
	if !utils.FileExists(treeinfopath) {
		log.Debug().Msgf(".treeinfo %s does not exists", treeinfopath)
		if distroDetails.Name != "" {
			if channelLabel != "" {
				log.Debug().Msg("Using channel override")
				*distro = types.Distribution{
					InstallType:  "generic_rpm",
					TreeLabel:    distroDetails.Name,
					ChannelLabel: channelLabel,
				}
				return nil
			} else if distroDetails.Version != "" && distroDetails.Arch != types.UnknownArch {
				log.Debug().Msg("Using distro details override")
				var err error
				*distro, err = getDistroFromDetails(distroDetails.Name, distroDetails.Version, distroDetails.Arch, flags)
				return err
			}
		}
		return fmt.Errorf(L("distribution treeinfo %s does not exists. Please provide distribution details and/or channel label"), treeinfopath)
	} else {
		var err error
		*distro, err = getDistroFromTreeinfo(path, flags)
		if err != nil {
			return err
		}

		// Overrides from the command line
		if distroDetails.Name != "" {
			distro.TreeLabel = distroDetails.Name
		}
		if channelLabel != "" {
			distro.ChannelLabel = channelLabel
		}
	}

	return nil
}

func getNameFromSource(source string) string {
	return strings.TrimSuffix(path.Base(source), ".iso")
}
