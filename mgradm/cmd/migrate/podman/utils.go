// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"fmt"
	"os/exec"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	migration_shared "github.com/uyuni-project/uyuni-tools/mgradm/cmd/migrate/shared"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/podman"
	podman_utils "github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"

	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func migrateToPodman(globalFlags *types.GlobalFlags, flags *podmanMigrateFlags, cmd *cobra.Command, args []string) error {
	if _, err := exec.LookPath("podman"); err != nil {
		log.Fatal().Err(err).Msg("install podman before running this command")
	}
	sourceFqdn := args[0]
	serverImage, err := utils.ComputeImage(flags.Image.Name, flags.Image.Tag)
	if err != nil {
		return fmt.Errorf("cannot compute image: %s", err)
	}

	// Find the SSH Socket and paths for the migration
	sshAuthSocket := migration_shared.GetSshAuthSocket()
	sshConfigPath, sshKnownhostsPath := migration_shared.GetSshPaths()

	tz, oldPgVersion, newPgVersion, err := podman.RunMigration(serverImage, flags.Image.PullPolicy, sshAuthSocket, sshConfigPath, sshKnownhostsPath, sourceFqdn)
	if err != nil {
		return fmt.Errorf("cannot run migration script: %s", err)
	}

	if oldPgVersion != newPgVersion {
		if err := podman.RunPgsqlVersionUpgrade(flags.Image, flags.MigrationImage, oldPgVersion, newPgVersion); err != nil {
			return fmt.Errorf("cannot run PostgreSQL version upgrade script: %s", err)
		}
	}

	schemaUpdateRequired := oldPgVersion != newPgVersion
	if err := podman.RunPgsqlFinalizeScript(serverImage, schemaUpdateRequired); err != nil {
		return fmt.Errorf("cannot run PostgreSQL finalize script: %s", err)
	}

	if err := podman.RunPostUpgradeScript(serverImage); err != nil {
		return fmt.Errorf("cannot run post upgrade script: %s", err)
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
