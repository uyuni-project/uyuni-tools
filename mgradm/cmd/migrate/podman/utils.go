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
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/coco"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/hub"
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
	serverImage, err := utils.ComputeImage(globalFlags.Registry, utils.DefaultTag, flags.Image)
	if err != nil {
		return utils.Errorf(err, L("cannot compute image"))
	}

	authFile, cleaner, err := podman_utils.PodmanLogin()
	if err != nil {
		return utils.Errorf(err, L("failed to login to registry.suse.com"))
	}
	defer cleaner()

	preparedImage, err := podman_utils.PrepareImage(authFile, serverImage, flags.Image.PullPolicy)
	if err != nil {
		return err
	}

	// Find the SSH Socket and paths for the migration
	sshAuthSocket := migration_shared.GetSshAuthSocket()
	sshConfigPath, sshKnownhostsPath := migration_shared.GetSshPaths()

	extractedData, err := podman.RunMigration(preparedImage, sshAuthSocket, sshConfigPath, sshKnownhostsPath, sourceFqdn, flags.User)
	if err != nil {
		return utils.Errorf(err, L("cannot run migration script"))
	}

	oldPgVersion := extractedData.CurrentPgVersion
	newPgVersion := extractedData.ImagePgVersion

	if oldPgVersion != newPgVersion {
		if err := podman.RunPgsqlVersionUpgrade(
			authFile, globalFlags.Registry, flags.Image, flags.DbUpgradeImage, oldPgVersion, newPgVersion,
		); err != nil {
			return utils.Errorf(err, L("cannot run PostgreSQL version upgrade script"))
		}
	}

	schemaUpdateRequired := oldPgVersion != newPgVersion
	if err := podman.RunPgsqlFinalizeScript(preparedImage, schemaUpdateRequired, true); err != nil {
		return utils.Errorf(err, L("cannot run PostgreSQL finalize script"))
	}

	if err := podman.RunPostUpgradeScript(preparedImage); err != nil {
		return utils.Errorf(err, L("cannot run post upgrade script"))
	}

	if err := podman.GenerateSystemdService(
		extractedData.Timezone, preparedImage, false, flags.Mirror, viper.GetStringSlice("podman.arg"),
	); err != nil {
		return utils.Errorf(err, L("cannot generate systemd service file"))
	}

	// Start the service
	if err := podman_utils.EnableService(podman_utils.ServerService); err != nil {
		return err
	}

	// Prepare confidential computing containers
	if err = coco.Upgrade(
		authFile, globalFlags.Registry, flags.Coco.Image, flags.Image,
		extractedData.DbPort, extractedData.DbName,
		extractedData.DbUser, extractedData.DbPassword,
	); err != nil {
		return utils.Errorf(err, L("cannot setup confidential computing attestation service"))
	}

	if flags.Coco.Replicas > 0 {
		err := podman_utils.ScaleService(flags.Coco.Replicas, podman_utils.ServerAttestationService)
		if err != nil {
			return err
		}
	}

	if err := hub.SetupHubXmlrpc(
		authFile, globalFlags.Registry, flags.Image.PullPolicy, flags.Image.Tag, flags.HubXmlrpc.Image,
	); err != nil {
		return err
	}

	if err := hub.EnableHubXmlrpc(flags.HubXmlrpc.Replicas); err != nil {
		return err
	}

	log.Info().Msg(L("Server migrated"))

	if err := podman_utils.EnablePodmanSocket(); err != nil {
		return utils.Errorf(err, L("cannot enable podman socket"))
	}

	cnx := shared.NewConnection("podman", podman_utils.ServerContainerName, "")
	if err := cnx.CopyCaCertificate(sourceFqdn); err != nil {
		return utils.Errorf(err, L("failed to add SSL CA certificate to host trusted certificates"))
	}

	return nil
}
