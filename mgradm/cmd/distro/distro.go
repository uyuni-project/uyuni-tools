// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package distro

import (
	"strings"
	"text/template"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
	"gopkg.in/yaml.v2"
)

type flagpole struct {
	Backend           string
	ChannelLabel      string                `mapstructure:"channel"`
	ProductMap        types.ProductMap      `mapstructure:"ProductMap"`
	ConnectionDetails api.ConnectionDetails `mapstructure:"api"`
}

type productMapTemplateData struct {
	DefaultProductMapRender string
	ProductMapRender        string
}

func prettyPrint(productMap types.ProductMap) string {
	if prettyPrintedProductMapBytes, err := yaml.Marshal(map[string]interface{}{"ProductMap": productMap}); err == nil {
		return string(prettyPrintedProductMapBytes)
	}
	return ""
}

func getProductMapHelp() string {
	return L(`Auto installation distribution product mapping.

For distribution to be registered by the server it is important to map distribution to the correct software channel.

Software channels can be named differently without any corellation to distribution name;
it is then needed to allow custom distribution name to software channel mapping.

One way to set software channel is by flag --channel to the distribution copy command.

For frequent usage it is possible to write custom product mapping to the mgradm configuration file as follows:

ProductMap:
  <distribution name>:
    <distribution version>:
      <distribution architecture>:
        ChannelLabel: <channel label>
        InstallType: <one of rhel_7|rhel_8|rhel_9|sles12generic|sles15generic|generic_rpm>
        TreeLabel: <custom distribution name>

Where
* <distribution name> is the name of the distribution, by default taken from '.treeinfo' file from the media.
  If '.treeinfo' is not found or available, command line option is required and used.
* <distribution version> is the version of the distribution, by default taken from '.treeinfo' file from the media.
  If'.treeinfo' is not found, command line option is required and used.
* <distribution architecture> is distribution architecture, by default taken from '.treeinfo' file from the media.
  If '.treeinfo' is not found, command line option is required and used.
* ChannelLabel is the channel label from Uyuni server and which is to be used for this distribution;
  can be overridden by command line flag.
* InstallType is used when installer is known (for autoyast or kickstart) or use 'generic_rpm'.
* TreeLabel is how the distribution will be presented in the Uyuni server UI. If not set <distribution name> is used.

Build-in product map:

{{ .DefaultProductMapRender }}
`)
}

func showHelp(_ *types.GlobalFlags,
	flags *flagpole,
	_ *cobra.Command,
	_ []string,
) error {
	mergedMaps := make(types.ProductMap, len(defaultProductMap))
	for k, v := range defaultProductMap {
		mergedMaps[k] = v
	}
	for distro, versions := range flags.ProductMap {
		if _, ok := mergedMaps[distro]; ok {
			for version, archs := range versions {
				if _, ok := mergedMaps[distro][version]; ok {
					for arch, distroDetail := range archs {
						// product map from config file has prio, overwrite
						mergedMaps[distro][version][arch] = distroDetail
					}
				} else {
					mergedMaps[distro][version] = archs
				}
			}
		} else {
			mergedMaps[distro] = versions
		}
	}
	print(prettyPrint(mergedMaps))
	return nil
}

func newCmd(globalFlags *types.GlobalFlags, run utils.CommandFunc[flagpole]) (*cobra.Command, error) {
	var flags flagpole

	distroCmd := &cobra.Command{
		Use:     "distribution",
		GroupID: "tool",
		Short:   L("Distributions management"),
		Long:    L("Tools for autoinstallation distributions management"),
		Aliases: []string{"distro"},
	}

	cpCmd := &cobra.Command{
		Use:   "copy path-to-source [distribution-name [version arch]]",
		Short: L("Copy distribution files from iso to the container"),
		Long: L(`Takes a path to source iso file or directory with mounted iso and copies it into the container.

Optional parameters 'distribution-name', 'version' and 'arch' specifies custom distribution.
If not set, distribution name is attempted to be autodetected:

- use name from '.treeinfo' file if exists
- use name of the ISO or passed directory

Note: API details are required for auto registration.`),
		Aliases: []string{"cp"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return utils.CommandHelper(globalFlags, cmd, args, &flags, nil, run)
		},
	}
	cpCmd.Flags().String("channel", "", L("Set parent channel for the distribution."))

	cpCmdHelp := &cobra.Command{
		Use:   "productmap",
		Short: L("Show distribution product map"),
		RunE: func(cmd *cobra.Command, args []string) error {
			return utils.CommandHelper(globalFlags, cmd, args, &flags, nil, showHelp)
		},
	}

	t := template.Must(template.New("help").Parse(getProductMapHelp()))
	var helpBuilder strings.Builder
	if err := t.Execute(&helpBuilder, productMapTemplateData{
		DefaultProductMapRender: prettyPrint(defaultProductMap),
		ProductMapRender:        prettyPrint(flags.ProductMap),
	}); err != nil {
		log.Fatal().Err(err).Msg(L("failed to compute config help command"))
	}
	cpCmdHelp.SetHelpTemplate(helpBuilder.String())

	api.AddAPIFlags(distroCmd)
	distroCmd.AddCommand(cpCmd)
	distroCmd.AddCommand(cpCmdHelp)

	return distroCmd, nil
}

// NewCommand command for distribution management.
func NewCommand(globalFlags *types.GlobalFlags) (*cobra.Command, error) {
	return newCmd(globalFlags, distroCp)
}
