// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package inspect

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/podman"
	adm_utils "github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	"github.com/uyuni-project/uyuni-tools/shared"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	shared_podman "github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func podmanInspect(
	globalFlags *types.GlobalFlags,
	flags *inspectFlags,
	cmd *cobra.Command,
	args []string,
) error {
	serverImage, err := utils.ComputeImage(flags.Image, flags.Tag)
	if err != nil && len(serverImage) > 0 {
		return fmt.Errorf(L("failed to determine image: %s"), err)
	}

	if len(serverImage) <= 0 {
		log.Debug().Msg("Use deployed image")

		cnx := shared.NewConnection("podman", shared_podman.ServerContainerName, "")
		serverImage, err = adm_utils.RunningImage(cnx, shared_podman.ServerContainerName)
		if err != nil {
			return fmt.Errorf(L("failed to find the image of the currently running server container: %s"))
		}
	}
	inspectResult, err := InspectPodman(serverImage, flags.PullPolicy)
	if err != nil {
		return fmt.Errorf(L("inspect command failed: %s"), err)
	}
	prettyInspectOutput, err := json.MarshalIndent(inspectResult, "", "  ")
	if err != nil {
		return fmt.Errorf(L("cannot print inspect result: %s"), err)
	}

	outputString := "\n" + string(prettyInspectOutput)
	log.Info().Msgf(outputString)

	return nil
}

// InspectPodman check values on a given image and deploy.
func InspectPodman(serverImage string, pullPolicy string) (map[string]string, error) {
	scriptDir, err := os.MkdirTemp("", "mgradm-*")
	defer os.RemoveAll(scriptDir)
	if err != nil {
		return map[string]string{}, fmt.Errorf(L("failed to create temporary directory: %s"), err)
	}

	inspectedHostValues, err := adm_utils.InspectHost()
	if err != nil {
		return map[string]string{}, fmt.Errorf(L("cannot inspect host values: %s"), err)
	}

	pullArgs := []string{}
	_, scc_user_exist := inspectedHostValues["host_scc_username"]
	_, scc_user_password := inspectedHostValues["host_scc_password"]
	if scc_user_exist && scc_user_password {
		pullArgs = append(pullArgs, "--creds", inspectedHostValues["host_scc_username"]+":"+inspectedHostValues["host_scc_password"])
	}

	preparedImage, err := shared_podman.PrepareImage(serverImage, pullPolicy, pullArgs...)
	if err != nil {
		return map[string]string{}, err
	}

	if err := adm_utils.GenerateInspectContainerScript(scriptDir); err != nil {
		return map[string]string{}, err
	}

	podmanArgs := []string{
		"-v", scriptDir + ":" + adm_utils.InspectOutputFile.Directory,
		"--security-opt", "label:disable",
	}

	err = podman.RunContainer("uyuni-inspect", preparedImage, podmanArgs,
		[]string{adm_utils.InspectOutputFile.Directory + "/" + adm_utils.InspectScriptFilename})
	if err != nil {
		return map[string]string{}, err
	}

	inspectResult, err := adm_utils.ReadInspectData(scriptDir)
	if err != nil {
		return map[string]string{}, fmt.Errorf(L("cannot inspect data: %s"), err)
	}

	return inspectResult, err
}
