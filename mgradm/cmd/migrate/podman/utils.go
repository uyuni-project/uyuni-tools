// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/uyuni-project/uyuni-tools/mgradm/cmd/migrate/shared"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/podman"
	adm_utils "github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	podman_utils "github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func migrateToPodman(globalFlags *types.GlobalFlags, flags *podmanMigrateFlags, cmd *cobra.Command, args []string) error {
	if _, err := exec.LookPath("podman"); err != nil {
		log.Fatal().Err(err).Msg("install podman before running this command")
	}

	// Find the SSH Socket and paths for the migration
	sshAuthSocket := shared.GetSshAuthSocket()
	sshConfigPath, sshKnownhostsPath := shared.GetSshPaths()

	scriptDir, err := adm_utils.GenerateMigrationScript(args[0], false)
	if err != nil {
		return fmt.Errorf("cannot generate migration script: %s", err)
	}
	defer os.RemoveAll(scriptDir)

	extraArgs := []string{
		"--security-opt", "label:disable",
		"-e", "SSH_AUTH_SOCK",
		"-v", filepath.Dir(sshAuthSocket) + ":" + filepath.Dir(sshAuthSocket),
		"-v", scriptDir + ":/var/lib/uyuni-tools/",
	}

	if sshConfigPath != "" {
		extraArgs = append(extraArgs, "-v", sshConfigPath+":/tmp/ssh_config")
	}

	if sshKnownhostsPath != "" {
		extraArgs = append(extraArgs, "-v", sshKnownhostsPath+":/etc/ssh/ssh_known_hosts")
	}

	serverImage, err := utils.ComputeImage(flags.Image.Name, flags.Image.Tag)
	if err != nil {
		return fmt.Errorf("failed to compute image URL: %s", err)
	}
	err = podman_utils.PrepareImage(serverImage, flags.Image.PullPolicy)
	if err != nil {
		return err
	}

	log.Info().Msg("Migrating server")
	if err := podman.RunContainer("uyuni-migration", serverImage, extraArgs,
		[]string{"/var/lib/uyuni-tools/migrate.sh"}); err != nil {
		return fmt.Errorf("cannot run uyuni migration container: %s", err)
	}

	// Read the extracted data
	tz, oldPgVersion, newPgVersion, err := adm_utils.ReadContainerData(scriptDir)
	if err != nil {
		return fmt.Errorf("cannot read data from container: %s", err)
	}

	if oldPgVersion != newPgVersion {
		var migrationImage types.ImageFlags
		migrationImage.Name = flags.MigrationImage.Name
		if migrationImage.Name == "" {
			migrationImage.Name = fmt.Sprintf("%s-migration-%s-%s", flags.Image.Name, oldPgVersion, newPgVersion)
		}
		migrationImage.Tag = flags.MigrationImage.Tag
		log.Info().Msgf("Using migration image %s:%s", migrationImage.Name, migrationImage.Tag)

		image, err := utils.ComputeImage(migrationImage.Name, migrationImage.Tag)
		if err != nil {
			return fmt.Errorf("failed to compute image URL: %s", err)
		}
		err = podman_utils.PrepareImage(image, flags.Image.PullPolicy)
		if err != nil {
			return err
		}

		scriptName, err := adm_utils.GeneratePgsqlVersionUpgradeScript(scriptDir, oldPgVersion, newPgVersion, false)
		if err != nil {
			return fmt.Errorf("cannot generate postgresql database version upgrade script: %s", err)
		}

		err = podman.RunContainer("uyuni-pg-migration", image, extraArgs,
			[]string{"/var/lib/uyuni-tools/" + scriptName})
		if err != nil {
			return fmt.Errorf("cannot run uyuni-pg-migration container: %s", err)
		}
	}

	scriptName, err := adm_utils.GenerateFinalizePostgresMigrationScript(scriptDir, true, oldPgVersion != newPgVersion, true, true, false)
	if err != nil {
		return fmt.Errorf("cannot generate postgresql database migration script: %s", err)
	}

	err = podman.RunContainer("uyuni-finalize-pg", serverImage, extraArgs,
		[]string{"/var/lib/uyuni-tools/" + scriptName})
	if err != nil {
		return fmt.Errorf("cannot run uyuni-finalize-pg container: %s", err)
	}

	if err := podman.GenerateSystemdService(tz, serverImage, false, viper.GetStringSlice("podman.arg")); err != nil {
		return fmt.Errorf("cannot generate systemd service file: %s", err)
	}

	// Start the service
	if err := podman_utils.EnableService(podman_utils.ServerService); err != nil {
		return err
	}

	log.Info().Msg("Server migrated")

	if err := podman_utils.EnablePodmanSocket(); err != nil {
		return fmt.Errorf("cannot run enable podman socket: %s", err)
	}

	return nil
}
