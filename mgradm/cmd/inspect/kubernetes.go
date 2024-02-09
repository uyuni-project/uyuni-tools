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

	inspect_shared "github.com/uyuni-project/uyuni-tools/mgradm/cmd/inspect/shared"
	shared_kubernetes "github.com/uyuni-project/uyuni-tools/shared/kubernetes"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

var kubernetesBuilt = true

func InspectKubernetes(serverImage string, pullPolicy string) (map[string]string, error) {
	for _, binary := range []string{"kubectl", "helm"} {
		if _, err := exec.LookPath(binary); err != nil {
			return map[string]string{}, fmt.Errorf("install %s before running this command. %s", binary, err)
		}
	}

	scriptDir, err := os.MkdirTemp("", "mgradm-*")
	defer os.RemoveAll(scriptDir)
	if err != nil {
		return map[string]string{}, fmt.Errorf("Failed to create temporary directory. %s", err)
	}

	if err := inspect_shared.GenerateInspectScript(scriptDir); err != nil {
		return map[string]string{}, err
	}

	command := path.Join(inspect_shared.InspectOutputFile.Directory, inspect_shared.InspectScriptFilename)

	const podName = "inspector"

	//delete pending pod and then check the node, because in presence of more than a pod GetNode return is wrong
	shared_kubernetes.DeletePod(podName)

	//this is needed because folder with script needs to be mounted
	nodeName, err := shared_kubernetes.GetNode("uyuni")
	if err != nil {
		return map[string]string{}, fmt.Errorf("Cannot find node for app uyuni %s", err)
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
					VolumeMounts: append(shared_kubernetes.PgsqlRequiredVolumeMounts,
						types.VolumeMount{MountPath: "/var/lib/uyuni-tools", Name: "var-lib-uyuni-tools"}),
					Image: serverImage,
				},
			},
			Volumes: append(shared_kubernetes.PgsqlRequiredVolumes,
				types.Volume{Name: "var-lib-uyuni-tools", HostPath: &types.HostPath{Path: scriptDir, Type: "Directory"}}),
		},
	}
	//transform deploy data in JSON
	override, err := shared_kubernetes.GenerateOverrideDeployment(deployData)
	if err != nil {
		return map[string]string{}, err
	}
	err = shared_kubernetes.RunPod(podName, serverImage, pullPolicy, command, override)
	if err != nil {
		return map[string]string{}, fmt.Errorf("Cannot run inspect pod %s", err)
	}

	inspectResult, err := inspect_shared.ReadInspectData(scriptDir)
	if err != nil {
		return map[string]string{}, fmt.Errorf("Cannot inspect data. %s", err)
	}

	prettyInspectOutput, err := json.MarshalIndent(inspectResult, "", "  ")
	if err != nil {
		return map[string]string{}, fmt.Errorf("Cannot print inspect result. %s", err)
	}

	log.Info().Msgf("\n%s", string(prettyInspectOutput))
	return inspectResult, err
}
