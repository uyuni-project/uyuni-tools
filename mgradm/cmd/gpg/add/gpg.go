// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package gpgadd

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	adm_utils "github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	"github.com/uyuni-project/uyuni-tools/shared"
	"github.com/uyuni-project/uyuni-tools/shared/kubernetes"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

const customKeyringPath = "/var/spacewalk/gpg/customer-build-keys.gpg"

type gpgAddFlags struct {
	Backend string
	Force   bool
}

// NewCommand import gpg keys from 3rd party repository.
func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	gpgAddKeyCmd := &cobra.Command{
		Use:   "add [URL]...",
		Short: L("Add gpg keys for 3rd party repositories"),
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags gpgAddFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, gpgAddKeys)
		},
	}

	gpgAddKeyCmd.Flags().BoolP("force", "f", false, L("Run the import"))
	utils.AddBackendFlag(gpgAddKeyCmd)
	return gpgAddKeyCmd
}

func gpgAddKeys(globalFlags *types.GlobalFlags, flags *gpgAddFlags, cmd *cobra.Command, args []string) error {
	cnx := shared.NewConnection(flags.Backend, podman.ServerContainerName, kubernetes.ServerFilter)
	if !utils.FileExists(customKeyringPath) {
		if err := adm_utils.ExecCommand(zerolog.InfoLevel, cnx, "mkdir", "-m", "700", "-p", filepath.Dir(customKeyringPath)); err != nil {
			return fmt.Errorf(L("failed to create folder %s: %s"), filepath.Dir(customKeyringPath), err)
		}
		if err := adm_utils.ExecCommand(zerolog.InfoLevel, cnx, "gpg", "--no-default-keyring", "--keyring", customKeyringPath, "--fingerprint"); err != nil {
			return fmt.Errorf(L("failed to create keyring %s: %s"), customKeyringPath, err)
		}
	}
	gpgAddCmd := []string{"gpg", "--no-default-keyring", "--import", "--import-options", "import-minimal"}

	if !flags.Force {
		gpgAddCmd = append(gpgAddCmd, "--dry-run")
	}
	gpgAddCmd = append(gpgAddCmd, "--keyring", customKeyringPath)

	scriptDir, err := os.MkdirTemp("", "mgradm-*")
	defer os.RemoveAll(scriptDir)
	if err != nil {
		return fmt.Errorf(L("failed to create temporary directory %s"), err)
	}

	for _, keyURL := range args {
		// Parse the URL
		parsedURL, err := url.Parse(keyURL)
		if err != nil {
			log.Error().Err(err).Msgf(L("failed to parse %s"), keyURL)
			continue
		}

		keyname := path.Base(parsedURL.Path)
		hostKeyPath := filepath.Join(scriptDir, keyname)
		if err := utils.DownloadFile(hostKeyPath, keyURL); err != nil {
			log.Error().Err(err).Msgf(L("failed to download %s"), keyURL)
			continue
		}

		if err := utils.RunCmdStdMapping(zerolog.InfoLevel, "gpg", "--show-key", hostKeyPath); err != nil {
			log.Error().Err(err).Msgf(L("failed to show key %s"), hostKeyPath)
			continue
		}

		containerKeyPath := filepath.Join(filepath.Dir(customKeyringPath), keyname)

		if err := cnx.Copy(hostKeyPath, "server:"+containerKeyPath, "", ""); err != nil {
			log.Error().Err(err).Msgf(L("failed to cp %s to %s"), hostKeyPath, containerKeyPath)
			continue
		}
		defer func() {
			_ = adm_utils.ExecCommand(zerolog.Disabled, cnx, "rm", containerKeyPath)
		}()

		gpgAddCmd = append(gpgAddCmd, containerKeyPath)
	}

	log.Info().Msgf(L("Running: %s"), strings.Join(gpgAddCmd, " "))
	if err := adm_utils.ExecCommand(zerolog.InfoLevel, cnx, gpgAddCmd...); err != nil {
		return fmt.Errorf(L("failed to run import key: %s"), err)
	}

	//this is for running import-suma-build-keys, who import customer-build-keys.gpg
	uyuniUpdateCmd := []string{"systemctl", "restart", "uyuni-update-config"}
	log.Info().Msgf(L("Running: %s"), strings.Join(uyuniUpdateCmd, " "))
	if err := adm_utils.ExecCommand(zerolog.InfoLevel, cnx, uyuniUpdateCmd...); err != nil {
		return fmt.Errorf(L("failed to restart uyuni-update-config: %s"), err)
	}
	return err
}
