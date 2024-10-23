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

var systemd podman_utils.Systemd = podman_utils.SystemdImpl{}

func migrateToPodman(globalFlags *types.GlobalFlags, flags *podmanMigrateFlags, cmd *cobra.Command, args []string) error {
	if _, err := exec.LookPath("podman"); err != nil {
		return fmt.Errorf(L("install podman before running this command"))
	}
	sourceFqdn, err := utils.GetFqdn(args)
	if err != nil {
		return err
	}

	serverImage, err := utils.ComputeImage(flags.Image.Registry, utils.DefaultTag, flags.Image)
	if err != nil {
		return utils.Errorf(err, L("cannot compute image"))
	}

	hostData, err := podman_utils.InspectHost()
	if err != nil {
		return err
	}

	authFile, cleaner, err := podman_utils.PodmanLogin(hostData, flags.SCC)
	if err != nil {
		return utils.Errorf(err, L("failed to login to registry.suse.com"))
	}
	defer cleaner()

	preparedImage, err := podman_utils.PrepareImage(authFile, serverImage, flags.Image.PullPolicy, true)
	if err != nil {
		return err
	}

	// Find the SSH Socket and paths for the migration
	sshAuthSocket := migration_shared.GetSshAuthSocket()
	sshConfigPath, sshKnownhostsPath := migration_shared.GetSshPaths()

	extractedData, err := podman.RunMigration(
		preparedImage, sshAuthSocket, sshConfigPath, sshKnownhostsPath, sourceFqdn, flags.User, flags.Prepare,
	)
	if err != nil {
		return utils.Errorf(err, L("cannot run migration script"))
	}
	if flags.Prepare {
		log.Info().Msg(L("Migration prepared. Run the 'migrate' command without '--prepare' to finish the migration."))
		return nil
	}

	oldPgVersion := extractedData.CurrentPgVersion
	newPgVersion := extractedData.ImagePgVersion

	if oldPgVersion != newPgVersion {
		if err := podman.RunPgsqlVersionUpgrade(
			authFile, flags.Image.Registry, flags.Image, flags.DbUpgradeImage, oldPgVersion, newPgVersion,
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
		systemd, extractedData.Timezone, preparedImage, false, flags.Mirror, viper.GetStringSlice("podman.arg"),
	); err != nil {
		return utils.Errorf(err, L("cannot generate systemd service file"))
	}

	// Start the service
	if err := systemd.EnableService(podman_utils.ServerService); err != nil {
		return err
	}

	// Prepare confidential computing containers
	if flags.Coco.Replicas > 0 {
		if err = coco.Upgrade(
			systemd, authFile, flags.Image.Registry, flags.Coco, flags.Image,
			extractedData.DbPort, extractedData.DbName,
			extractedData.DbUser, extractedData.DbPassword,
		); err != nil {
			return utils.Errorf(err, L("cannot setup confidential computing attestation service"))
		}

		err := systemd.ScaleService(flags.Coco.Replicas, podman_utils.ServerAttestationService)
		if err != nil {
			return err
		}
	}

	hubReplicas := flags.HubXmlrpc.Replicas
	if extractedData.HasHubXmlrpcApi {
		log.Info().Msg(L("Enabling Hub XML-RPC API since it is enabled on the migrated server"))
		hubReplicas = 1
	}
	if hubReplicas > 0 {
		if err := hub.SetupHubXmlrpc(
			systemd, authFile, flags.Image.Registry, flags.Image.PullPolicy, flags.Image.Tag, flags.HubXmlrpc,
		); err != nil {
			return err
		}
		if err := hub.EnableHubXmlrpc(systemd, hubReplicas); err != nil {
			return err
		}
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
