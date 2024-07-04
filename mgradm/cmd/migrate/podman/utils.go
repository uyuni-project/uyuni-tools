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
	"github.com/uyuni-project/uyuni-tools/shared"
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
	serverImage, err := utils.ComputeImage(flags.Image)
	if err != nil {
		return utils.Errorf(err, L("cannot compute image"))
	}

	// FIXME all this code should be centralized. Now it being called in several different places.
	inspectedHostValues, err := utils.InspectHost(false)
	if err != nil {
		return utils.Errorf(err, L("cannot inspect host values"))
	}

	pullArgs := []string{}
	_, scc_user_exist := inspectedHostValues["host_scc_username"]
	_, scc_user_password := inspectedHostValues["host_scc_password"]
	if scc_user_exist && scc_user_password {
		pullArgs = append(pullArgs, "--creds", inspectedHostValues["host_scc_username"]+":"+inspectedHostValues["host_scc_password"])
	}

	preparedImage, err := podman_utils.PrepareImage(serverImage, flags.Image.PullPolicy, pullArgs...)
	if err != nil {
		return err
	}

	// Find the SSH Socket and paths for the migration
	sshAuthSocket := migration_shared.GetSshAuthSocket()
	sshConfigPath, sshKnownhostsPath := migration_shared.GetSshPaths()

	tz, oldPgVersion, newPgVersion, err := podman.RunMigration(preparedImage, sshAuthSocket, sshConfigPath, sshKnownhostsPath, sourceFqdn, flags.User)
	if err != nil {
		return utils.Errorf(err, L("cannot run migration script"))
	}

	if oldPgVersion != newPgVersion {
		if err := podman.RunPgsqlVersionUpgrade(flags.Image, flags.DbUpgradeImage, oldPgVersion, newPgVersion); err != nil {
			return utils.Errorf(err, L("cannot run PostgreSQL version upgrade script"))
		}
	}

	schemaUpdateRequired := oldPgVersion != newPgVersion
	if err := podman.RunPgsqlFinalizeScript(preparedImage, schemaUpdateRequired); err != nil {
		return utils.Errorf(err, L("cannot run PostgreSQL finalize script"))
	}

	if err := podman.RunPostUpgradeScript(preparedImage); err != nil {
		return utils.Errorf(err, L("cannot run post upgrade script"))
	}

	if err := podman.GenerateSystemdService(tz, preparedImage, false, flags.Mirror, viper.GetStringSlice("podman.arg")); err != nil {
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

	cnx := shared.NewConnection("podman", podman_utils.ServerContainerName, "")

	if err := cnx.WaitForContainer(); err != nil {
		return err
	}

	if err := cnx.CopyCaCertificate(sourceFqdn); err != nil {
		return utils.Errorf(err, L("failed to add SSL CA certificate to host trusted certificates"))
	}

	return nil
}
