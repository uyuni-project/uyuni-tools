// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package kubernetes

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

const k3sTraefikConfigPath = "/var/lib/rancher/k3s/server/manifests/k3s-traefik-config.yaml"

// InstallK3sTraefikConfig install K3s Traefik configuration.
func InstallK3sTraefikConfig(tcpPorts []types.PortMap, udpPorts []types.PortMap) {
	log.Info().Msg(L("Installing K3s Traefik configuration"))

	data := K3sTraefikConfigTemplateData{
		TcpPorts: tcpPorts,
		UdpPorts: udpPorts,
	}
	if err := utils.WriteTemplateToFile(data, k3sTraefikConfigPath, 0600, false); err != nil {
		log.Fatal().Err(err).Msgf(L("Failed to write K3s Traefik configuration"))
	}

	// Wait for traefik to be back
	log.Info().Msg(L("Waiting for Traefik to be reloaded"))
	for i := 0; i < 60; i++ {
		out, err := utils.RunCmdOutput(zerolog.TraceLevel, "kubectl", "get", "job", "-n", "kube-system",
			"-o", "jsonpath={.status.completionTime}", "helm-install-traefik")
		if err == nil {
			completionTime, err := time.Parse(time.RFC3339, string(out))
			if err == nil && time.Since(completionTime).Seconds() < 60 {
				break
			}
		}
	}
}

// UninstallK3sTraefikConfig uninstall K3s Traefik configuration.
func UninstallK3sTraefikConfig(dryRun bool) {
	utils.UninstallFile(k3sTraefikConfigPath, dryRun)
}

// InspectKubernetes check values on a given image and deploy.
func InspectKubernetes(namespace string, serverImage string, pullPolicy string) (*utils.ServerInspectData, error) {
	for _, binary := range []string{"kubectl", "helm"} {
		if _, err := exec.LookPath(binary); err != nil {
			return nil, fmt.Errorf(L("install %s before running this command"), binary)
		}
	}

	scriptDir, err := os.MkdirTemp("", "mgradm-*")
	defer os.RemoveAll(scriptDir)
	if err != nil {
		return nil, utils.Errorf(err, L("failed to create temporary directory"))
	}

	inspector := utils.NewServerInspector(scriptDir)
	if err := inspector.GenerateScript(); err != nil {
		return nil, err
	}

	command := path.Join(utils.InspectContainerDirectory, utils.InspectScriptFilename)

	const podName = "inspector"

	//delete pending pod and then check the node, because in presence of more than a pod GetNode return is wrong
	if err := DeletePod(namespace, podName, ServerFilter); err != nil {
		return nil, utils.Errorf(err, L("cannot delete %s"), podName)
	}

	//this is needed because folder with script needs to be mounted
	nodeName, err := GetNode(namespace, ServerFilter)
	if err != nil {
		return nil, utils.Errorf(err, L("cannot find node running uyuni"))
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
	override, err := GenerateOverrideDeployment(deployData)
	if err != nil {
		return nil, err
	}
	err = RunPod(namespace, podName, ServerFilter, serverImage, pullPolicy, command, override)
	if err != nil {
		return nil, utils.Errorf(err, L("cannot run inspect pod"))
	}

	inspectResult, err := inspector.ReadInspectData()
	if err != nil {
		return nil, utils.Errorf(err, L("cannot inspect data"))
	}

	return inspectResult, err
}
