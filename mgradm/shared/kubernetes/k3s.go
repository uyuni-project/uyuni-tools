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
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
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
func RunPgsqlVersionUpgrade(
	registry string,
	image types.ImageFlags,
	upgradeImage types.ImageFlags,
	namespace string,
	nodeName string,
	oldPgsql string,
	newPgsql string,
) error {
	scriptDir, err := utils.TempDir()
	if err != nil {
		return err
	}
	defer os.RemoveAll(scriptDir)
	if newPgsql > oldPgsql {
		log.Info().Msgf(L("Previous PostgreSQL is %[1]s, new one is %[2]s. Performing a DB version upgrade…"), oldPgsql, newPgsql)

		pgsqlVersionUpgradeContainer := "uyuni-upgrade-pgsql"

		upgradeImageUrl := ""
		if upgradeImage.Name == "" {
			upgradeImageUrl, err = utils.ComputeImage(registry, image.Tag, image, fmt.Sprintf("-migration-%s-%s", oldPgsql, newPgsql))
			if err != nil {
				return utils.Errorf(err, L("failed to compute image URL"))
			}
		} else {
			upgradeImageUrl, err = utils.ComputeImage(registry, image.Tag, upgradeImage)
			if err != nil {
				return utils.Errorf(err, L("failed to compute image URL"))
			}
		}

		log.Info().Msgf(L("Using database upgrade image %s"), upgradeImageUrl)
		pgsqlVersionUpgradeScriptName, err := adm_utils.GeneratePgsqlVersionUpgradeScript(scriptDir, oldPgsql, newPgsql, true)
		if err != nil {
			return utils.Errorf(err, L("cannot generate PostgreSQL database version upgrade script"))
		}

		//delete pending pod and then check the node, because in presence of more than a pod GetNode return is wrong
		if err := kubernetes.DeletePod(namespace, pgsqlVersionUpgradeContainer, kubernetes.ServerFilter); err != nil {
			return utils.Errorf(err, L("cannot delete %s"), pgsqlVersionUpgradeContainer)
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
		overridePgsqlVersioUpgrade, err := kubernetes.GenerateOverrideDeployment(pgsqlVersioUpgradeDeployData)
		if err != nil {
			return err
		}

		err = kubernetes.RunPod(namespace, pgsqlVersionUpgradeContainer, kubernetes.ServerFilter, upgradeImageUrl, image.PullPolicy, "/var/lib/uyuni-tools/"+pgsqlVersionUpgradeScriptName, overridePgsqlVersioUpgrade)
		if err != nil {
			return utils.Errorf(err, L("error running container %s"), pgsqlVersionUpgradeContainer)
		}
	}
	return nil
}

// RunPgsqlFinalizeScript run the script with all the action required to a db after upgrade.
func RunPgsqlFinalizeScript(
	serverImage string, pullPolicy string, namespace string, nodeName string, schemaUpdateRequired bool, migration bool,
) error {
	scriptDir, err := utils.TempDir()
	if err != nil {
		return err
	}
	defer os.RemoveAll(scriptDir)
	pgsqlFinalizeContainer := "uyuni-finalize-pgsql"
	pgsqlFinalizeScriptName, err := adm_utils.GenerateFinalizePostgresScript(scriptDir, true, schemaUpdateRequired, true, migration, true)
	if err != nil {
		return utils.Errorf(err, L("cannot generate PostgreSQL finalization script"))
	}
	//delete pending pod and then check the node, because in presence of more than a pod GetNode return is wrong
	if err := kubernetes.DeletePod(namespace, pgsqlFinalizeContainer, kubernetes.ServerFilter); err != nil {
		return utils.Errorf(err, L("cannot delete %s"), pgsqlFinalizeContainer)
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
					VolumeMounts: append(utils.PgsqlRequiredVolumeMounts,
						types.VolumeMount{MountPath: "/var/lib/uyuni-tools", Name: "var-lib-uyuni-tools"}),
				},
			},
			Volumes: append(utils.PgsqlRequiredVolumes,
				types.Volume{Name: "var-lib-uyuni-tools", HostPath: &types.HostPath{Path: scriptDir, Type: "Directory"}}),
		},
	}
	//transform deploy data in JSON
	overridePgsqlFinalize, err := kubernetes.GenerateOverrideDeployment(pgsqlFinalizeDeployData)
	if err != nil {
		return err
	}
	err = kubernetes.RunPod(namespace, pgsqlFinalizeContainer, kubernetes.ServerFilter, serverImage, pullPolicy, "/var/lib/uyuni-tools/"+pgsqlFinalizeScriptName, overridePgsqlFinalize)
	if err != nil {
		return utils.Errorf(err, L("error running container %s"), pgsqlFinalizeContainer)
	}
	return nil
}

// RunPostUpgradeScript run the script with the changes to apply after the upgrade.
func RunPostUpgradeScript(serverImage string, pullPolicy string, namespace string, nodeName string) error {
	scriptDir, err := utils.TempDir()
	if err != nil {
		return err
	}
	defer os.RemoveAll(scriptDir)
	postUpgradeContainer := "uyuni-post-upgrade"
	postUpgradeScriptName, err := adm_utils.GeneratePostUpgradeScript(scriptDir, "localhost")
	if err != nil {
		return utils.Errorf(err, L("cannot generate PostgreSQL finalization script"))
	}

	//delete pending pod and then check the node, because in presence of more than a pod GetNode return is wrong
	if err := kubernetes.DeletePod(namespace, postUpgradeContainer, kubernetes.ServerFilter); err != nil {
		return utils.Errorf(err, L("cannot delete %s"), postUpgradeContainer)
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
	overridePostUpgrade, err := kubernetes.GenerateOverrideDeployment(postUpgradeDeployData)
	if err != nil {
		return err
	}

	err = kubernetes.RunPod(namespace, postUpgradeContainer, kubernetes.ServerFilter, serverImage, pullPolicy, "/var/lib/uyuni-tools/"+postUpgradeScriptName, overridePostUpgrade)
	if err != nil {
		return utils.Errorf(err, L("error running container %s"), postUpgradeContainer)
	}

	return nil
}
