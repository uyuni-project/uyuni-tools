// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgradm/cmd/inspect"
	upgrade_shared "github.com/uyuni-project/uyuni-tools/mgradm/cmd/upgrade/shared"

	"github.com/uyuni-project/uyuni-tools/mgradm/shared/podman"

	"github.com/uyuni-project/uyuni-tools/shared"
	shared_podman "github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func upgradePodman(globalFlags *types.GlobalFlags, flags *podmanUpgradeFlags, cmd *cobra.Command, args []string) error {
	serverImage, err := utils.ComputeImage(flags.Image.Name, flags.Image.Tag)
	if err != nil {
		return fmt.Errorf("failed to compute image URL")
	}

	inspectedValues, err := inspect.InspectPodman(serverImage, flags.Image.PullPolicy)
	if err != nil {
		return fmt.Errorf("cannot inspect podman values: %s", err)
	}

	cnx := shared.NewConnection("podman", shared_podman.ServerContainerName, "")

	if err := upgrade_shared.SanityCheck(cnx, inspectedValues, serverImage); err != nil {
		return err
	}

	if err := shared_podman.StopService(shared_podman.ServerService); err != nil {
		return fmt.Errorf("cannot stop service %s", err)
	}

	defer func() {
		err = shared_podman.StartService(shared_podman.ServerService)
	}()
	if inspectedValues["image_pg_version"] > inspectedValues["current_pg_version"] {
		if err := podman.RunPgsqlVersionUpgrade(flags.Image, flags.MigrationImage, inspectedValues["current_pg_version"], inspectedValues["image_pg_version"]); err != nil {
			return fmt.Errorf("cannot run PostgreSQL version upgrade script: %s", err)
		}
	} else if inspectedValues["image_pg_version"] == inspectedValues["current_pg_version"] {
		log.Info().Msgf("Upgrading to %s without changing PostgreSQL version", inspectedValues["uyuni_release"])
	} else {
		return fmt.Errorf("trying to downgrade postgresql from %s to %s", inspectedValues["current_pg_version"], inspectedValues["image_pg_version"])
	}

	schemaUpdateRequired := inspectedValues["current_pg_version"] != inspectedValues["image_pg_version"]
	if err := podman.RunPgsqlFinalizeScript(serverImage, schemaUpdateRequired); err != nil {
		return fmt.Errorf("cannot run PostgreSQL version upgrade script: %s", err)
	}

	if err := podman.RunPostUpgradeScript(serverImage); err != nil {
		return fmt.Errorf("cannot run post upgrade script: %s", err)
	}

	if err := shared_podman.GenerateSystemdConfFile("uyuni-server", "Service", "Environment=UYUNI_IMAGE="+serverImage); err != nil {
		return err
	}
	log.Info().Msg("Waiting for the server to start...")
	return shared_podman.ReloadDaemon(false)
}
