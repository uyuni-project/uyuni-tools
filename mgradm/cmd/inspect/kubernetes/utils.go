// SPDX-FileCopyrightText: 2023 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

//go:build !nok8s

package kubernetes

import (
	"encoding/json"
	"os"
	"os/exec"
	"strconv"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	inspect_shared "github.com/uyuni-project/uyuni-tools/mgradm/cmd/inspect/shared"
	shared_kubernetes "github.com/uyuni-project/uyuni-tools/shared/kubernetes"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func inspectForKubernetes(
	globalFlags *types.GlobalFlags,
	flags *kubernetesInspectFlags,
	cmd *cobra.Command,
	args []string,
) error {
	for _, binary := range []string{"kubectl", "helm"} {
		if _, err := exec.LookPath(binary); err != nil {
			log.Fatal().Err(err).Msgf("install %s before running this command", binary)
		}
	}
	serverImage, err := utils.ComputeImage(flags.Image.Name, flags.Image.Tag)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to compute image URL")
	}

	scriptDir, err := os.MkdirTemp("", "mgradm-*")
	defer os.RemoveAll(scriptDir)

	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to create temporary directory")
	}

	inspect_shared.GenerateInspectScript(scriptDir)

	command := inspect_shared.InspectOutputFile.Directory + "/" + inspect_shared.InspectScriptFilename

	podName := "inspector"

	overridesArgs := []string{"--override-type=strategic", "--overrides", `{"apiVersion":"v1","spec":{"restartPolicy":"Never","containers":[{"name":` + strconv.Quote(podName) + `,"image":` + strconv.Quote(serverImage) + `,"volumeMounts":[{"mountPath":"/var/lib/uyuni-tools","name":"var-lib-uyuni-tools"}]}],"volumes":[{"name":"var-lib-uyuni-tools","hostPath":{"path":` + strconv.Quote(scriptDir) + `,"type":"Directory"}}]}}`}

	shared_kubernetes.RunPod(podName, serverImage, command, overridesArgs...)

	shared_kubernetes.WaitForPod(podName, "Succeeded")

	inspectResult := inspect_shared.ReadInspectData(scriptDir)

	shared_kubernetes.DeletePod(podName)

	prettyInspectOutput, err := json.MarshalIndent(inspectResult, "", "  ")
	if err != nil {
		log.Fatal().Err(err).Msgf("Cannot print inspect result")
	}

	log.Info().Msgf("\n%s", string(prettyInspectOutput))
	return err
}
