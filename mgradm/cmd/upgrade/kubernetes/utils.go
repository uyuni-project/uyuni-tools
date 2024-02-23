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

		if err := kubernetes.RunPgsqlVersionUpgrade(flags.Image, flags.MigrationImage, nodeName, inspectedValues["current_pg_version"], inspectedValues["image_pg_version"]); err != nil {
			return fmt.Errorf("cannot run PostgreSQL version upgrade script: %s", err)
		}
	}

	schemaUpdateRequired := inspectedValues["current_pg_version"] != inspectedValues["image_pg_version"]
	if err := kubernetes.RunPgsqlFinalizeScript(serverImage, flags.Image.PullPolicy, nodeName, schemaUpdateRequired); err != nil {
		return fmt.Errorf("cannot run PostgreSQL version upgrade script: %s", err)
	}

	if err := kubernetes.RunPostUpgradeScript(serverImage, flags.Image.PullPolicy, nodeName); err != nil {
		return fmt.Errorf("cannot run post upgrade script: %s", err)
	}

	err = kubernetes.UyuniUpgrade(serverImage, flags.Image.PullPolicy, &flags.Helm, kubeconfig, fqdn, clusterInfos.Ingress)
	if err != nil {
		return fmt.Errorf("cannot upgrade to image %s: %s", serverImage, err)
	}

	return shared_kubernetes.WaitForDeployment(flags.Helm.Uyuni.Namespace, "uyuni", "uyuni")
}
