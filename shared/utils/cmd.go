// SPDX-FileCopyrightText: 2024 SUSE LLC
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
)

// Default path where to look for locale files.
//
// On SUSE distros this should be overridden with /usr/share/locale.
var LocaleRoot = "locale"

// DefaultRegistryServer represents the default registry FQDN for images.
var DefaultRegistryServer = "registry.opensuse.org"

// DefaultRegistryPath represents the default registry path used for images.
var DefaultRegistryPath = "/uyuni"

// DefaultRegistryPath represents the default registry path used for helm charts.
// The value is the same as a default here, but could be configured with different in the spec file.
var DefaultRegistryHelmPath = "/uyuni"

// DefaultTag represents the default tag used for image.
var DefaultTag = "latest"

// DefaultPullP represents the default pull policy used for image.
var DefaultPullPolicy = "Always"

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

// AddRegistryFlags adds the flags setting the registry server and path.
func AddRegistryFlags(cmd *cobra.Command) {
	cmd.Flags().String("registry-server", "",
		fmt.Sprintf(
			L(`Server FQDN or IP and optional port for the container images registry. (default "%s")`),
			DefaultRegistryServer))
	cmd.Flags().String("registry-path", "",
		fmt.Sprintf(
			L(`Path to the container images in the registry. (default "%s")`),
			DefaultRegistryPath))
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
	cmd.Flags().String("pullPolicy", DefaultPullPolicy,
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
	Images  bool
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
	cmd.Flags().Bool("purge-images", false, L("Also remove the container images"))

	if withBackend {
		AddBackendFlag(cmd)
	}
}

// Get the contatenated server and path using default values if needed.
func GetRegistryPath(registry *types.RegistryFlags) string {
	server := DefaultRegistryServer
	if registry.Server != "" {
		server = registry.Server
	}

	imagesPath := DefaultRegistryPath
	if registry.Path != "" {
		imagesPath = registry.Path
	}

	return path.Join(server, imagesPath)
}
