package podman

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgradm/cmd/inspect"
	upgrade_shared "github.com/uyuni-project/uyuni-tools/mgradm/cmd/upgrade/shared"

	"github.com/uyuni-project/uyuni-tools/mgradm/shared/podman"

	adm_utils "github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	"github.com/uyuni-project/uyuni-tools/shared"
	shared_podman "github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func upgradePodman(globalFlags *types.GlobalFlags, flags *podmanUpgradeFlags, cmd *cobra.Command, args []string) error {
	serverImage, err := utils.ComputeImage(flags.Image.Name, flags.Image.Tag)
	if err != nil {
		return fmt.Errorf("Failed to compute image URL")
	}

	inspectedValues, err := inspect.InspectPodman(serverImage, flags.Image.PullPolicy)

	cnx := shared.NewConnection("podman", shared_podman.ServerContainerName, "")

	upgrade_shared.SanityCheck(cnx, inspectedValues, serverImage)

	podmanArgs := flags.Podman.Args
	if flags.MirrorPath != "" {
		podmanArgs = append(podmanArgs, "-v", flags.MirrorPath+":/mirror")
	}

	scriptDir, err := os.MkdirTemp("", "mgradm-*")
	defer os.RemoveAll(scriptDir)

	if err != nil {
		return fmt.Errorf("Failed to create temporary directory")
	}

	shared_podman.StopService(shared_podman.ServerService)
	defer shared_podman.StartService(shared_podman.ServerService)

	if inspectedValues["image_pg_version"] > inspectedValues["current_pg_version"] {
		log.Info().Msgf("Previous postgresql is %s, instead new one is %s. Performing a DB migration...", inspectedValues["current_pg_version"], inspectedValues["image_pg_version"])
		var migrationImage types.ImageFlags
		extraArgs := []string{
			"-v", scriptDir + ":/var/lib/uyuni-tools/",
		}
		migrationImage.Name = fmt.Sprintf("%s-migration-%s-%s", flags.Image.Name, inspectedValues["current_pg_version"], inspectedValues["image_pg_version"])
		shared_podman.PrepareImage(migrationImage.Name, flags.Image.PullPolicy)
		scriptName, err := adm_utils.GeneratePgMigrationScript(scriptDir, inspectedValues["current_pg_version"], inspectedValues["image_pg_version"], false)
		if err != nil {
			return fmt.Errorf("Cannot generate pg migration script %s", err)
		}

		podman.RunContainer("uyuni-upgrade-pgsql", migrationImage.Name, extraArgs,
			[]string{"/var/lib/uyuni-tools/" + scriptName})
	} else if inspectedValues["image_pg_version"] == inspectedValues["current_pg_version"] {
		log.Info().Msgf("Upgrading uyuni to %s without changing PostgreSQL version", inspectedValues["uyuni_release"])
	} else {
		return fmt.Errorf("Trying to downgrade PostgreSQL. Previous postgresql is %s, instead new one is %s", inspectedValues["current_pg_version"], inspectedValues["image_pg_version"])
	}

	extraArgs := []string{
		"-v", scriptDir + ":/var/lib/uyuni-tools/",
	}

	scriptName, err := adm_utils.GenerateFinalizePostgresMigrationScript(scriptDir, true, inspectedValues["current_pg_version"] != inspectedValues["image_pg_version"], true, true, false)
	if err != nil {
		return fmt.Errorf("Cannot generate finalize postgres script %s", err)
	}
	podman.RunContainer("uyuni-finalize-pgsql", serverImage, extraArgs,
		[]string{"/var/lib/uyuni-tools/" + scriptName})

	shared_podman.GenerateSystemdConfFile("uyuni-server", "Service", "Environment=UYUNI_IMAGE="+serverImage)
	log.Info().Msg("Waiting for the server to start...")
	shared_podman.ReloadDaemon(false)

	return nil
}
