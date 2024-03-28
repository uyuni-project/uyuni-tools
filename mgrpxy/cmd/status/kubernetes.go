// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package status

import (
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/kubernetes"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
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
		return fmt.Errorf(L("failed to discover the cluster type: %s"), err)
	}

	kubeconfig := clusterInfos.GetKubeconfig()
	if !kubernetes.HasHelmRelease("uyuni-proxy", kubeconfig) {
		return errors.New(L("no uyuni-proxy helm release installed on the cluster"))
	}

	namespace, err := kubernetes.FindNamespace("uyuni-proxy", kubeconfig)
	if err != nil {
		return fmt.Errorf(L("failed to find the uyuni-proxy deployment namespace: %s"), err)
	}

	// Is the pod running? Do we have all the replicas?
	status, err := kubernetes.GetDeploymentStatus(namespace, "uyuni-proxy")
	if err != nil {
		return fmt.Errorf(L("failed to get deployment status: %s"), err)
	}
	if status.Replicas != status.ReadyReplicas {
		log.Warn().Msgf(L("Some replicas are not ready: %d / %d"), status.ReadyReplicas, status.Replicas)
	}

	if status.AvailableReplicas == 0 {
		return errors.New(L("the pod is not running"))
	}

	log.Info().Msg(L("Proxy containers up and running"))

	return nil
}
