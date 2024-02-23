// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package kubernetes

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
	adm_utils "github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	"github.com/uyuni-project/uyuni-tools/shared/kubernetes"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// InstallK3sTraefikConfig installs the K3s Traefik configuration.
func InstallK3sTraefikConfig(debug bool) {
	tcpPorts := []types.PortMap{}
	tcpPorts = append(tcpPorts, utils.TCP_PORTS...)
	if debug {
		tcpPorts = append(tcpPorts, utils.DEBUG_PORTS...)
	}

	kubernetes.InstallK3sTraefikConfig(tcpPorts, utils.UDP_PORTS)
}

// RunPgsqlVersionUpgrade perform a PostgreSQL major upgrade.
func RunPgsqlVersionUpgrade(image types.ImageFlags, migrationImage types.ImageFlags, nodeName string, oldPgsql string, newPgsql string) error {
	scriptDir, err := os.MkdirTemp("", "mgradm-*")
	defer os.RemoveAll(scriptDir)
	if err != nil {
		return fmt.Errorf("failed to create temporary directory")
	}
	if newPgsql > oldPgsql {
		log.Info().Msgf("Previous postgresql is %s, instead new one is %s. Performing a DB version upgrade...", oldPgsql, newPgsql)

		pgsqlVersionUpgradeContainer := "uyuni-upgrade-pgsql"

		migrationImageUrl := ""
		if migrationImage.Name == "" {
			migrationImageUrl, err = utils.ComputeImage(image.Name, image.Tag, fmt.Sprintf("-migration-%s-%s", oldPgsql, newPgsql))
			if err != nil {
				return fmt.Errorf("failed to compute image URL %s", err)
			}
		} else {
			migrationImageUrl, err = utils.ComputeImage(migrationImage.Name, image.Tag)
			if err != nil {
				return fmt.Errorf("failed to compute image URL %s", err)
			}
		}

		log.Info().Msgf("Using migration image %s", migrationImageUrl)
		pgsqlVersionUpgradeScriptName, err := adm_utils.GeneratePgsqlVersionUpgradeScript(scriptDir, oldPgsql, newPgsql, true)
		if err != nil {
			return fmt.Errorf("cannot generate postgresql database version upgrade script: %s", err)
		}

		//delete pending pod and then check the node, because in presence of more than a pod GetNode return is wrong
		if err := kubernetes.DeletePod(pgsqlVersionUpgradeContainer, kubernetes.ServerFilter); err != nil {
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
						VolumeMounts: append(utils.MigrationVolumeMounts,
							types.VolumeMount{MountPath: "/var/lib/uyuni-tools", Name: "var-lib-uyuni-tools"}),
					},
				},
				Volumes: append(utils.MigrationVolumes,
					types.Volume{Name: "var-lib-uyuni-tools", HostPath: &types.HostPath{Path: scriptDir, Type: "Directory"}}),
			},
		}

		//transform deploy in JSON
		overridePgsqlVersioUpgrade, err := kubernetes.GenerateOverrideDeployment(pgsqlVersioUpgradeDeployData)
		if err != nil {
			return err
		}

		err = kubernetes.RunPod(pgsqlVersionUpgradeContainer, kubernetes.ServerFilter, migrationImageUrl, image.PullPolicy, "/var/lib/uyuni-tools/"+pgsqlVersionUpgradeScriptName, overridePgsqlVersioUpgrade)
		if err != nil {
			return fmt.Errorf("error running container %s: %s", pgsqlVersionUpgradeContainer, err)
		}
	}
	return nil
}

// RunPgsqlFinalizeScript run the script with all the action required to a db after upgrade.
func RunPgsqlFinalizeScript(serverImage string, pullPolicy string, nodeName string, schemaUpdateRequired bool) error {
	scriptDir, err := os.MkdirTemp("", "mgradm-*")
	defer os.RemoveAll(scriptDir)
	if err != nil {
		return fmt.Errorf("failed to create temporary directory")
	}
	pgsqlFinalizeContainer := "uyuni-finalize-pgsql"
	pgsqlFinalizeScriptName, err := adm_utils.GenerateFinalizePostgresScript(scriptDir, true, schemaUpdateRequired, true, true, true)
	if err != nil {
		return fmt.Errorf("cannot generate psql finalize script: %s", err)
	}
	//delete pending pod and then check the node, because in presence of more than a pod GetNode return is wrong
	if err := kubernetes.DeletePod(pgsqlFinalizeContainer, kubernetes.ServerFilter); err != nil {
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
					VolumeMounts: append(utils.MigrationVolumeMounts,
						types.VolumeMount{MountPath: "/var/lib/uyuni-tools", Name: "var-lib-uyuni-tools"}),
				},
			},
			Volumes: append(utils.MigrationVolumes,
				types.Volume{Name: "var-lib-uyuni-tools", HostPath: &types.HostPath{Path: scriptDir, Type: "Directory"}}),
		},
	}
	//transform deploy data in JSON
	overridePgsqlFinalize, err := kubernetes.GenerateOverrideDeployment(pgsqlFinalizeDeployData)
	if err != nil {
		return err
	}
	err = kubernetes.RunPod(pgsqlFinalizeContainer, kubernetes.ServerFilter, serverImage, pullPolicy, "/var/lib/uyuni-tools/"+pgsqlFinalizeScriptName, overridePgsqlFinalize)
	if err != nil {
		return fmt.Errorf("error running container %s: %s", pgsqlFinalizeContainer, err)
	}
	return nil
}

// RunPostUpgradeScript run the script with the changes to apply after the upgrade.
func RunPostUpgradeScript(serverImage string, pullPolicy string, nodeName string) error {
	scriptDir, err := os.MkdirTemp("", "mgradm-*")
	defer os.RemoveAll(scriptDir)
	if err != nil {
		return fmt.Errorf("failed to create temporary directory")
	}
	postUpgradeContainer := "uyuni-post-upgrade"
	postUpgradeScriptName, err := adm_utils.GeneratePostUpgradeScript(scriptDir, "localhost")
	if err != nil {
		return fmt.Errorf("cannot generate postgresql finalization script %s", err)
	}

	//delete pending pod and then check the node, because in presence of more than a pod GetNode return is wrong
	if err := kubernetes.DeletePod(postUpgradeContainer, kubernetes.ServerFilter); err != nil {
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
==== BASE ====
					VolumeMounts: append(utils.EtcServerVolumeMounts,
==== BASE ====
						types.VolumeMount{MountPath: "/var/lib/uyuni-tools", Name: "var-lib-uyuni-tools"}),
				},
			},
==== BASE ====
			Volumes: append(utils.EtcServerVolumes,
==== BASE ====
				types.Volume{Name: "var-lib-uyuni-tools", HostPath: &types.HostPath{Path: scriptDir, Type: "Directory"}}),
		},
	}
	//transform deploy data in JSON
	overridePostUpgrade, err := kubernetes.GenerateOverrideDeployment(postUpgradeDeployData)
	if err != nil {
		return err
	}

	err = kubernetes.RunPod(postUpgradeContainer, kubernetes.ServerFilter, serverImage, pullPolicy, "/var/lib/uyuni-tools/"+postUpgradeScriptName, overridePostUpgrade)
	if err != nil {
		return fmt.Errorf("error running container %s: %s", postUpgradeContainer, err)
	}
	return nil
}
