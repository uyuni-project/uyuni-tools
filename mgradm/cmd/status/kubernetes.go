// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

//go:build !nok8s

package status

import (
	"errors"
	"fmt"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	adm_utils "github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	"github.com/uyuni-project/uyuni-tools/shared"
	"github.com/uyuni-project/uyuni-tools/shared/kubernetes"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func kubernetesStatus(
	globalFlags *types.GlobalFlags,
	flags *statusFlags,
	cmd *cobra.Command,
	args []string,
) error {
	// Do we have an uyuni helm release?
	clusterInfos, err := kubernetes.CheckCluster()
	if err != nil {
		return fmt.Errorf("failed to discover the cluster type: %s", err)
	}

	kubeconfig := clusterInfos.GetKubeconfig()
	if !kubernetes.HasHelmRelease("uyuni", kubeconfig) {
		return errors.New("no uyuni helm release installed on the cluster")
	}

	namespace, err := kubernetes.FindNamespace("uyuni", kubeconfig)
	if err != nil {
		return fmt.Errorf("failed to find the uyuni deployment namespace: %s", err)
	}

	// Is the pod running? Do we have all the replicas?
	status, err := kubernetes.GetDeploymentStatus(namespace, "uyuni")
	if err != nil {
		return fmt.Errorf("failed to get deployment status: %s", err)
	}
	if status.Replicas != status.ReadyReplicas {
		log.Warn().Msgf("Some replicas are not ready: %d / %d", status.ReadyReplicas, status.Replicas)
	}

	if status.AvailableReplicas == 0 {
		return errors.New("the pod is not running")
	}

	// Are the services running in the container?
	cnx := shared.NewConnection("kubectl", "", kubernetes.ServerFilter)
	if err := adm_utils.ExecCommand(zerolog.InfoLevel, cnx, "spacewalk-service", "status"); err != nil {
		return fmt.Errorf("failed to run spacewalk-service status: %s", err)
	}
	return nil
}
