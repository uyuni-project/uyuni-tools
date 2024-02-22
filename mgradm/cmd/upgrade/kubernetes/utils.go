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

	//this is needed because folder with script needs to be mounted
	//check the node before scaling down
	nodeName, err := shared_kubernetes.GetNode("uyuni")
	if err != nil {
		return fmt.Errorf("cannot find node for app uyuni %s", err)
	}

	err = shared_kubernetes.ReplicasTo(shared_kubernetes.ServerFilter, 0)
	if err != nil {
		return fmt.Errorf("cannot set replica to 0: %s", err)
	}

	defer func() {
		// if something is running, we don't need to set replicas to 1
		if _, err = shared_kubernetes.GetNode("uyuni"); err != nil {
			err = shared_kubernetes.ReplicasTo(shared_kubernetes.ServerFilter, 1)
		}
	}()
	if inspectedValues["image_pg_version"] > inspectedValues["current_pg_version"] {
		log.Info().Msgf("Previous postgresql is %s, instead new one is %s. Performing a DB version upgrade...", inspectedValues["current_pg_version"], inspectedValues["image_pg_version"])

		pgsqlVersionUpgradeContainer := "uyuni-upgrade-pgsql"

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
		pgsqlVersionUpgradeScriptName, err := adm_utils.GeneratePgsqlVersionUpgradeScript(scriptDir, inspectedValues["current_pg_version"], inspectedValues["image_pg_version"], true)
		if err != nil {
			return fmt.Errorf("cannot generate postgresql database version upgrade script: %s", err)
		}

		//delete pending pod and then check the node, because in presence of more than a pod GetNode return is wrong
		if err := shared_kubernetes.DeletePod(pgsqlVersionUpgradeContainer, shared_kubernetes.ServerFilter); err != nil {
			return fmt.Errorf("cannot delete %s: %s", pgsqlVersionUpgradeContainer, err)
		}

		//generate deploy data
		pgsqlVersioUpgradeDeployData := types.Deployment{
			APIVersion: "v1",
			Spec: &types.Spec{
				RestartPolicy: "Never",
				NodeName:      nodeName,
				Containers: []types.Container{
					{
						Name: pgsqlVersionUpgradeContainer,
						VolumeMounts: append(utils.PgsqlRequiredVolumeMounts,
							types.VolumeMount{MountPath: "/var/lib/uyuni-tools", Name: "var-lib-uyuni-tools"}),
					},
				},
				Volumes: append(utils.PgsqlRequiredVolumes,
					types.Volume{Name: "var-lib-uyuni-tools", HostPath: &types.HostPath{Path: scriptDir, Type: "Directory"}}),
			},
		}

		//transform deploy in JSON
		overridePgsqlVersioUpgrade, err := shared_kubernetes.GenerateOverrideDeployment(pgsqlVersioUpgradeDeployData)
		if err != nil {
			return err
		}

		err = shared_kubernetes.RunPod(pgsqlVersionUpgradeContainer, shared_kubernetes.ServerFilter, migrationImageUrl, flags.Image.PullPolicy, "/var/lib/uyuni-tools/"+pgsqlVersionUpgradeScriptName, overridePgsqlVersioUpgrade)
		if err != nil {
			return fmt.Errorf("error running container %s: %s", pgsqlVersionUpgradeContainer, err)
		}
	}

	{
		pgsqlFinalizeContainer := "uyuni-finalize-pgsql"
		pgsqlFinalizeScriptName, err := adm_utils.GenerateFinalizePostgresScript(scriptDir, true, inspectedValues["current_pg_version"] != inspectedValues["image_pg_version"], true, true, true)
		if err != nil {
			return fmt.Errorf("cannot generate psql finalize script: %s", err)
		}
		//delete pending pod and then check the node, because in presence of more than a pod GetNode return is wrong
		if err := shared_kubernetes.DeletePod(pgsqlFinalizeContainer, shared_kubernetes.ServerFilter); err != nil {
			return fmt.Errorf("cannot delete %s: %s", pgsqlFinalizeContainer, err)
		}
		//generate deploy data
		pgsqlFinalizeDeployData := types.Deployment{
			APIVersion: "v1",
			Spec: &types.Spec{
				RestartPolicy: "Never",
				NodeName:      nodeName,
				Containers: []types.Container{
					{
						Name: pgsqlFinalizeContainer,
						VolumeMounts: append(utils.ServerVolumeMounts,
							types.VolumeMount{MountPath: "/var/lib/uyuni-tools", Name: "var-lib-uyuni-tools"}),
					},
				},
				Volumes: append(utils.ServerVolumes,
					types.Volume{Name: "var-lib-uyuni-tools", HostPath: &types.HostPath{Path: scriptDir, Type: "Directory"}}),
			},
		}
		//transform deploy data in JSON
		overridePgsqlFinalize, err := shared_kubernetes.GenerateOverrideDeployment(pgsqlFinalizeDeployData)
		if err != nil {
			return err
		}
		err = shared_kubernetes.RunPod(pgsqlFinalizeContainer, shared_kubernetes.ServerFilter, serverImage, flags.Image.PullPolicy, "/var/lib/uyuni-tools/"+pgsqlFinalizeScriptName, overridePgsqlFinalize)
		if err != nil {
			return fmt.Errorf("error running container %s: %s", pgsqlFinalizeContainer, err)
		}
	}
	{
		postUpgradeContainer := "uyuni-post-upgrade"
		postUpgradeScriptName, err := adm_utils.GeneratePostUpgradeScript(scriptDir, "localhost")
		if err != nil {
			return fmt.Errorf("cannot generate postgresql finalization script %s", err)
		}

		//delete pending pod and then check the node, because in presence of more than a pod GetNode return is wrong
		if err := shared_kubernetes.DeletePod(postUpgradeContainer, shared_kubernetes.ServerFilter); err != nil {
			return fmt.Errorf("cannot delete %s: %s", postUpgradeContainer, err)
		}
		//generate deploy data
		postUpgradeDeployData := types.Deployment{
			APIVersion: "v1",
			Spec: &types.Spec{
				RestartPolicy: "Never",
				NodeName:      nodeName,
				Containers: []types.Container{
					{
						Name: postUpgradeContainer,
						VolumeMounts: append(utils.PgsqlRequiredVolumeMounts,
							types.VolumeMount{MountPath: "/var/lib/uyuni-tools", Name: "var-lib-uyuni-tools"}),
					},
				},
				Volumes: append(utils.PgsqlRequiredVolumes,
					types.Volume{Name: "var-lib-uyuni-tools", HostPath: &types.HostPath{Path: scriptDir, Type: "Directory"}}),
			},
		}
		//transform deploy data in JSON
		overridePostUpgrade, err := shared_kubernetes.GenerateOverrideDeployment(postUpgradeDeployData)
		if err != nil {
			return err
		}

		err = shared_kubernetes.RunPod(postUpgradeContainer, shared_kubernetes.ServerFilter, serverImage, flags.Image.PullPolicy, "/var/lib/uyuni-tools/"+postUpgradeScriptName, overridePostUpgrade)
		if err != nil {
			return fmt.Errorf("error running container %s: %s", postUpgradeContainer, err)
		}
	}

	err = kubernetes.UyuniUpgrade(serverImage, flags.Image.PullPolicy, &flags.Helm, kubeconfig, fqdn, clusterInfos.Ingress)
	if err != nil {
		return fmt.Errorf("cannot upgrade to image %s: %s", serverImage, err)
	}

	return shared_kubernetes.WaitForDeployment(flags.Helm.Uyuni.Namespace, "uyuni", "uyuni")
}
