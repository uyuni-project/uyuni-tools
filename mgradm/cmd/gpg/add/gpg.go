// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package gpgadd

import (
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

func newCmd(globalFlags *types.GlobalFlags, run utils.CommandFunc[gpgAddFlags]) *cobra.Command {
	gpgAddKeyCmd := &cobra.Command{
		Use:   "add [URL]...",
		Short: L("Add GPG keys for 3rd party repositories"),
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags gpgAddFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, nil, run)
		},
	}

	gpgAddKeyCmd.Flags().BoolP("force", "f", false, L("Import without asking confirmation"))
	utils.AddBackendFlag(gpgAddKeyCmd)
	return gpgAddKeyCmd
}

// NewCommand import gpg keys from 3rd party repository.
func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	return newCmd(globalFlags, gpgAddKeys)
}

func gpgAddKeys(_ *types.GlobalFlags, flags *gpgAddFlags, _ *cobra.Command, args []string) error {
	cnx := shared.NewConnection(flags.Backend, podman.ServerContainerName, kubernetes.ServerFilter)
	if !cnx.TestExistenceInPod(customKeyringPath) {
		if err := adm_utils.ExecCommand(
			zerolog.InfoLevel, cnx, "mkdir", "-m", "700", "-p", filepath.Dir(customKeyringPath),
		); err != nil {
			return utils.Errorf(err, L("failed to create folder %s"), filepath.Dir(customKeyringPath))
		}
		if err := adm_utils.ExecCommand(
			zerolog.InfoLevel, cnx, "gpg", "--no-default-keyring", "--keyring", customKeyringPath, "--fingerprint",
		); err != nil {
			return utils.Errorf(err, L("failed to create keyring %s"), customKeyringPath)
		}
	}
	gpgAddCmd := []string{"gpg", "--no-default-keyring", "--import", "--import-options", "import-minimal"}

	gpgAddCmd = append(gpgAddCmd, "--keyring", customKeyringPath)

	scriptDir, cleaner, err := utils.TempDir()
	if err != nil {
		return err
	}
	defer cleaner()

	for _, keyURL := range args {
		var hostKeyPath string
		var keyname string
		if _, err := os.Stat(keyURL); err == nil {
			// gpg passed in a local file
			hostKeyPath = keyURL
			keyname = filepath.Base(hostKeyPath)
		} else {
			// Parse the URL
			parsedURL, err := url.Parse(keyURL)
			if err != nil {
				log.Error().Err(err).Msgf(L("failed to parse %s"), keyURL)
				continue
			}

			keyname = path.Base(parsedURL.Path)
			hostKeyPath = filepath.Join(scriptDir, keyname)
			if err := utils.DownloadFile(hostKeyPath, keyURL); err != nil {
				log.Error().Err(err).Msgf(L("failed to download %s"), keyURL)
				continue
			}
		}

		if err := utils.RunCmdStdMapping(zerolog.InfoLevel, "gpg", "--show-key", hostKeyPath); err != nil {
			log.Error().Err(err).Msgf(L("failed to show key %s"), hostKeyPath)
			continue
		}
		if !flags.Force {
			ret, err := utils.YesNo(L("Do you really want to trust this key"))
			if err != nil {
				return err
			}
			if !ret {
				return nil
			}
		}

		containerKeyPath := filepath.Join(filepath.Dir(customKeyringPath), keyname)

		if err := cnx.Copy(hostKeyPath, "server:"+containerKeyPath, "", ""); err != nil {
			log.Error().Err(err).Msgf(L("failed to copy %[1]s to %[2]s"), hostKeyPath, containerKeyPath)
			continue
		}
		defer func() {
			_ = adm_utils.ExecCommand(zerolog.Disabled, cnx, "rm", containerKeyPath)
		}()

		gpgAddCmd = append(gpgAddCmd, containerKeyPath)
	}

	log.Info().Msgf(L("Running %s"), strings.Join(gpgAddCmd, " "))
	if err := adm_utils.ExecCommand(zerolog.InfoLevel, cnx, gpgAddCmd...); err != nil {
		return utils.Errorf(err, L("failed to run import key"))
	}

	//this is for running import-suma-build-keys, who import customer-build-keys.gpg
	uyuniUpdateCmd := []string{"systemctl", "restart", "uyuni-update-config"}
	log.Info().Msgf(L("Running %s"), strings.Join(uyuniUpdateCmd, " "))
	if err := adm_utils.ExecCommand(zerolog.InfoLevel, cnx, uyuniUpdateCmd...); err != nil {
		return utils.Errorf(err, L("failed to restart uyuni-update-config"))
	}
	return err
}
