// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

//go:build !nok8s

package kubernetes

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgradm/cmd/inspect"
	upgrade_shared "github.com/uyuni-project/uyuni-tools/mgradm/cmd/upgrade/shared"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/kubernetes"
	adm_utils "github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	"github.com/uyuni-project/uyuni-tools/shared"
	shared_kubernetes "github.com/uyuni-project/uyuni-tools/shared/kubernetes"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func upgradeKubernetes(
	globalFlags *types.GlobalFlags,
	flags *kubernetesUpgradeFlags,
	cmd *cobra.Command,
	args []string,
) error {
	for _, binary := range []string{"kubectl", "helm"} {
		if _, err := exec.LookPath(binary); err != nil {
			return fmt.Errorf("install %s before running this command", binary)
		}
	}
	cnx := shared.NewConnection("kubectl", "", shared_kubernetes.ServerFilter)

	serverImage, err := utils.ComputeImage(flags.Image.Name, flags.Image.Tag)
	if err != nil {
		return fmt.Errorf("failed to compute image URL")
	}

	inspectedValues, err := inspect.InspectKubernetes(serverImage, flags.Image.PullPolicy)
	if err != nil {
		return err
	}

	err = upgrade_shared.SanityCheck(cnx, inspectedValues, serverImage)
	if err != nil {
		return err
	}

	fqdn, exist := inspectedValues["fqdn"]
	if !exist {
		return fmt.Errorf("inspect function did non return fqdn value")
	}

	clusterInfos := shared_kubernetes.CheckCluster()
	kubeconfig := clusterInfos.GetKubeconfig()

	scriptDir, err := os.MkdirTemp("", "mgradm-*")
	defer os.RemoveAll(scriptDir)
	if err != nil {
		return fmt.Errorf("failed to create temporary directory")
	}

	err = shared_kubernetes.ReplicasTo(shared_kubernetes.ServerFilter, 0)
	if err != nil {
		return fmt.Errorf("cannot set replica to 0: %s", err)
	}

	defer func() {
		err = shared_kubernetes.ReplicasTo(shared_kubernetes.ServerFilter, 1)
	}()
	if inspectedValues["image_pg_version"] > inspectedValues["current_pg_version"] {
		log.Info().Msgf("Previous postgresql is %s, instead new one is %s. Performing a DB migration...", inspectedValues["current_pg_version"], inspectedValues["image_pg_version"])

		migrationContainer := "uyuni-upgrade-pgsql"

		migrationImageUrl := ""
		if flags.MigrationImage.Name == "" {
			migrationImageUrl, err = utils.ComputeImage(flags.Image.Name, flags.Image.Tag, fmt.Sprintf("-migration-%s-%s", inspectedValues["current_pg_version"], inspectedValues["image_pg_version"]))
			if err != nil {
				return fmt.Errorf("failed to compute image URL %s", err)
			}
		} else {
			migrationImageUrl, err = utils.ComputeImage(flags.MigrationImage.Name, flags.Image.Tag)
			if err != nil {
				return fmt.Errorf("failed to compute image URL %s", err)
			}
		}

		log.Info().Msgf("Using migration image %s", migrationImageUrl)
		scriptName, err := adm_utils.GeneratePgMigrationScript(scriptDir, inspectedValues["current_pg_version"], inspectedValues["image_pg_version"], true)
		if err != nil {
			return fmt.Errorf("cannot generate pg migration script: %s", err)
		}

		//delete pending pod and then check the node, because in presence of more than a pod GetNode return is wrong
		out, err := shared_kubernetes.DeletePod(migrationContainer)
		if err != nil {
			return fmt.Errorf("cannot delete %s: %s. Output is: %s", migrationContainer, err, out)
		}

		//this is needed because folder with script needs to be mounted
		nodeName, err := shared_kubernetes.GetNode("uyuni")
		if err != nil {
			return fmt.Errorf("cannot find node for app uyuni %s", err)
		}

		//generate deploy data
		deployData := types.Deployment{
			APIVersion: "v1",
			Spec: &types.Spec{
				RestartPolicy: "Never",
				NodeName:      nodeName,
				Containers: []types.Container{
					{
						Name: migrationContainer,
						VolumeMounts: append(shared_kubernetes.PgsqlRequiredVolumeMounts,
							types.VolumeMount{MountPath: "/var/lib/uyuni-tools", Name: "var-lib-uyuni-tools"}),
					},
				},
				Volumes: append(shared_kubernetes.PgsqlRequiredVolumes,
					types.Volume{Name: "var-lib-uyuni-tools", HostPath: &types.HostPath{Path: scriptDir, Type: "Directory"}}),
			},
		}

		//transform deploy in JSON
		override, err := shared_kubernetes.GenerateOverrideDeployment(deployData)
		if err != nil {
			return err
		}

		err = shared_kubernetes.RunPod(migrationContainer, migrationImageUrl, flags.Image.PullPolicy, "/var/lib/uyuni-tools/"+scriptName, override)
		if err != nil {
			return fmt.Errorf("error running container %s: %s", migrationContainer, err)
		}
	}

	scriptName, err := adm_utils.GenerateFinalizePostgresMigrationScript(scriptDir, true, inspectedValues["current_pg_version"] != inspectedValues["image_pg_version"], true, true, true)
	if err != nil {
		return fmt.Errorf("cannot generate pg migration script: %s", err)
	}

	pgsqlFinalizeContainer := "uyuni-finalize-pgsql"

	//delete pending pod and then check the node, because in presence of more than a pod GetNode return is wrong
	out, err := shared_kubernetes.DeletePod(pgsqlFinalizeContainer)
	if err != nil {
		return fmt.Errorf("cannot delete %s: %s. Output is: %s", pgsqlFinalizeContainer, err, out)
	}

	//this is needed because folder with script needs to be mounted
	nodeName, err := shared_kubernetes.GetNode("uyuni")
	if err != nil {
		return fmt.Errorf("cannot find node for app uyuni %s", err)
	}

	//generate deploy data
	deployData := types.Deployment{
		APIVersion: "v1",
		Spec: &types.Spec{
			RestartPolicy: "Never",
			NodeName:      nodeName,
			Containers: []types.Container{
				{
					Name: pgsqlFinalizeContainer,
					VolumeMounts: append(shared_kubernetes.PgsqlRequiredVolumeMounts,
						types.VolumeMount{MountPath: "/var/lib/uyuni-tools", Name: "var-lib-uyuni-tools"}),
				},
			},
			Volumes: append(shared_kubernetes.PgsqlRequiredVolumes,
				types.Volume{Name: "var-lib-uyuni-tools", HostPath: &types.HostPath{Path: scriptDir, Type: "Directory"}}),
		},
	}
	//transform deploy data in JSON
	override, err := shared_kubernetes.GenerateOverrideDeployment(deployData)
	if err != nil {
		return err
	}
	err = shared_kubernetes.RunPod(pgsqlFinalizeContainer, serverImage, flags.Image.PullPolicy, "/var/lib/uyuni-tools/"+scriptName, override)
	if err != nil {
		return fmt.Errorf("error running container %s: %s", pgsqlFinalizeContainer, err)
	}

	err = kubernetes.UyuniUpgrade(serverImage, flags.Image.PullPolicy, &flags.Helm, kubeconfig, fqdn, clusterInfos.Ingress)
	if err != nil {
		return fmt.Errorf("cannot upgrade to image %s: %s", serverImage, err)
	}

	return shared_kubernetes.WaitForDeployment(flags.Helm.Uyuni.Namespace, "uyuni", "uyuni")
}
