// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

//go:build !nok8s

package inspect

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	adm_utils "github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	"github.com/uyuni-project/uyuni-tools/shared"
	shared_kubernetes "github.com/uyuni-project/uyuni-tools/shared/kubernetes"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func kuberneteInspect(
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

		cnx := shared.NewConnection("kubectl", "", shared_kubernetes.ServerFilter)
		serverImage, err = adm_utils.RunningImage(cnx, "uyuni")
		if err != nil {
			return fmt.Errorf(L("failed to find the image of the currently running server container: %s"))
		}
	}

	inspectResult, err := InspectKubernetes(serverImage, flags.PullPolicy)
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

// InspectKubernetes check values on a given image and deploy.
func InspectKubernetes(serverImage string, pullPolicy string) (map[string]string, error) {
	for _, binary := range []string{"kubectl", "helm"} {
		if _, err := exec.LookPath(binary); err != nil {
			return map[string]string{}, fmt.Errorf(L("install %s before running this command"), binary)
		}
	}

	scriptDir, err := os.MkdirTemp("", "mgradm-*")
	defer os.RemoveAll(scriptDir)
	if err != nil {
		return map[string]string{}, fmt.Errorf(L("failed to create temporary directory: %s"), err)
	}

	if err := adm_utils.GenerateInspectContainerScript(scriptDir); err != nil {
		return map[string]string{}, err
	}

	command := path.Join(adm_utils.InspectOutputFile.Directory, adm_utils.InspectScriptFilename)

	const podName = "inspector"

	//delete pending pod and then check the node, because in presence of more than a pod GetNode return is wrong
	if err := shared_kubernetes.DeletePod(podName, shared_kubernetes.ServerFilter); err != nil {
		return map[string]string{}, fmt.Errorf(L("cannot delete %s: %s"), podName, err)
	}

	//this is needed because folder with script needs to be mounted
	nodeName, err := shared_kubernetes.GetNode("uyuni")
	if err != nil {
		return map[string]string{}, fmt.Errorf(L("cannot find node running uyuni: %s"), err)
	}

	//generate deploy data
	deployData := types.Deployment{
		APIVersion: "v1",
		Spec: &types.Spec{
			RestartPolicy: "Never",
			NodeName:      nodeName,
			Containers: []types.Container{
				{
					Name: podName,
					VolumeMounts: append(utils.PgsqlRequiredVolumeMounts,
						types.VolumeMount{MountPath: "/var/lib/uyuni-tools", Name: "var-lib-uyuni-tools"}),
					Image: serverImage,
				},
			},
			Volumes: append(utils.PgsqlRequiredVolumes,
				types.Volume{Name: "var-lib-uyuni-tools", HostPath: &types.HostPath{Path: scriptDir, Type: "Directory"}}),
		},
	}
	//transform deploy data in JSON
	override, err := shared_kubernetes.GenerateOverrideDeployment(deployData)
	if err != nil {
		return map[string]string{}, err
	}
	err = shared_kubernetes.RunPod(podName, shared_kubernetes.ServerFilter, serverImage, pullPolicy, command, override)
	if err != nil {
		return map[string]string{}, fmt.Errorf(L("cannot run inspect pod: %s"), err)
	}

	inspectResult, err := adm_utils.ReadInspectData(scriptDir)
	if err != nil {
		return map[string]string{}, fmt.Errorf(L("cannot inspect data: %s"), err)
	}

	return inspectResult, err
}
