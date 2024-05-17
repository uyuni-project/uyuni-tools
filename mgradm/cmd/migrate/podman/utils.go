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

	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func migrateToPodman(globalFlags *types.GlobalFlags, flags *podmanMigrateFlags, cmd *cobra.Command, args []string) error {
	if _, err := exec.LookPath("podman"); err != nil {
		return fmt.Errorf(L("install podman before running this command"))
	}
	sourceFqdn := args[0]
	serverImage, err := utils.ComputeImage(flags.Image.Name, flags.Image.Tag)
	if err != nil {
		return utils.Errorf(err, L("cannot compute image"))
	}

	// Find the SSH Socket and paths for the migration
	sshAuthSocket := migration_shared.GetSshAuthSocket()
	sshConfigPath, sshKnownhostsPath := migration_shared.GetSshPaths()

	tz, oldPgVersion, newPgVersion, err := podman.RunMigration(serverImage, flags.Image.PullPolicy, sshAuthSocket, sshConfigPath, sshKnownhostsPath, sourceFqdn, flags.User)
	if err != nil {
		return utils.Errorf(err, L("cannot run migration script"))
	}

	if oldPgVersion != newPgVersion {
		if err := podman.RunPgsqlVersionUpgrade(flags.Image, flags.MigrationImage, oldPgVersion, newPgVersion); err != nil {
			return utils.Errorf(err, L("cannot run PostgreSQL version upgrade script"))
		}
	}

	schemaUpdateRequired := oldPgVersion != newPgVersion
	if err := podman.RunPgsqlFinalizeScript(serverImage, schemaUpdateRequired); err != nil {
		return utils.Errorf(err, L("cannot run PostgreSQL finalize script"))
	}

	if err := podman.RunPostUpgradeScript(serverImage); err != nil {
		return utils.Errorf(err, L("cannot run post upgrade script"))
	}

	if err := podman.GenerateSystemdService(tz, serverImage, false, flags.Mirror, viper.GetStringSlice("podman.arg")); err != nil {
		return utils.Errorf(err, L("cannot generate systemd service file"))
	}

	// Start the service
	if err := podman_utils.EnableService(podman_utils.ServerService); err != nil {
		return err
	}

	log.Info().Msg(L("Server migrated"))

	if err := podman_utils.EnablePodmanSocket(); err != nil {
		return utils.Errorf(err, L("cannot enable podman socket"))
	}

	return nil
}
