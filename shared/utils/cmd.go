// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

// Default path where to look for locale files.
//
// On SUSE distros this should be overridden with /usr/share/locale.
var LocaleRoot = "locale"

// DefaultNamespace represents the default name used for image.
var DefaultNamespace = "registry.opensuse.org/uyuni"

// DefaultTag represents the default tag used for image.
var DefaultTag = "latest"

// This variable needs to be set a build time using git tags.
var Version = "0.0.0"

// CommandFunc is a function to be executed by a Cobra command.
type CommandFunc[F interface{}] func(*types.GlobalFlags, *F, *cobra.Command, []string) error

// CommandHelper parses the configuration file into the flags and runs the fn function.
// This function should be passed to Command's RunE.
func CommandHelper[T interface{}](
	globalFlags *types.GlobalFlags,
	cmd *cobra.Command,
	args []string,
	flags *T,
	fn CommandFunc[T],
) error {
	viper, err := ReadConfig(cmd, GlobalConfigFilename, globalFlags.ConfigPath)
	if err != nil {
		return err
	}
	if err := viper.Unmarshal(&flags); err != nil {
		log.Error().Err(err).Msg(L("failed to unmarshall configuration"))
		return Errorf(err, L("failed to unmarshall configuration"))
	}
	return fn(globalFlags, flags, cmd, args)
}

// AddBackendFlag add the flag for setting the backend ('podman', 'podman-remote', 'kubectl').
func AddBackendFlag(cmd *cobra.Command) {
	cmd.Flags().String("backend", "", L("tool to use to reach the container. Possible values: 'podman', 'podman-remote', 'kubectl'. Default guesses which to use."))
}

// AddPullPolicyFlag adds the --pullPolicy flag to a command.
//
// Since podman doesn't have such a concept of pull policy like kubernetes,
// the values need some explanations for it:
//   - Never: just check and fail if needed
//   - IfNotPresent: check and pull
//   - Always: pull without checking
//
// For kubernetes the value is simply passed to the helm charts.
func AddPullPolicyFlag(cmd *cobra.Command) {
	cmd.Flags().String("pullPolicy", "IfNotPresent",
		L("set whether to pull the images or not. The value can be one of 'Never', 'IfNotPresent' or 'Always'"))
}

// AddPullPolicyFlag adds the --pullPolicy flag to an upgrade command.
func AddPullPolicyUpgradeFlag(cmd *cobra.Command) {
	cmd.Flags().String("pullPolicy", "Always",
		L("set whether to pull the images or not. The value can be one of 'Never', 'IfNotPresent' or 'Always'"))
}

// AddPTFFlag add PTF flag to a command.
func AddPTFFlag(cmd *cobra.Command) {
	cmd.Flags().String("ptf", "", L("PTF ID"))
	cmd.Flags().String("test", "", L("Test package ID"))
	cmd.Flags().String("user", "", L("SCC user"))
}

// PurgeFlags defined what has te be removed in an uninstall command.
type PurgeFlags struct {
	Volumes bool
}

// UninstallFlags are the common flags for uninstall commands.
type UninstallFlags struct {
	Backend string
	Force   bool
	Purge   PurgeFlags
}

// AddUninstallFlags adds the common flags for uninstall commands.
func AddUninstallFlags(cmd *cobra.Command, withBackend bool) {
	cmd.Flags().BoolP("force", "f", false, L("Actually remove the server"))
	cmd.Flags().Bool("purge-volumes", false, L("Also remove the volumes"))

	if withBackend {
		AddBackendFlag(cmd)
	}
}
