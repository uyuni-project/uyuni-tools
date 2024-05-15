// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

//go:build !nok8s

package status

import (
	"errors"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	adm_utils "github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	"github.com/uyuni-project/uyuni-tools/shared"
	"github.com/uyuni-project/uyuni-tools/shared/kubernetes"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
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
		return utils.Errorf(err, L("failed to discover the cluster type"))
	}

	kubeconfig := clusterInfos.GetKubeconfig()
	if !kubernetes.HasHelmRelease("uyuni", kubeconfig) {
		return errors.New(L("no uyuni helm release installed on the cluster"))
	}

	namespace, err := kubernetes.FindNamespace("uyuni", kubeconfig)
	if err != nil {
		return utils.Errorf(err, L("failed to find the uyuni deployment namespace"))
	}

	// Is the pod running? Do we have all the replicas?
	status, err := kubernetes.GetDeploymentStatus(namespace, "uyuni")
	if err != nil {
		return utils.Errorf(err, L("failed to get deployment status"))
	}
	if status.Replicas != status.ReadyReplicas {
		log.Warn().Msgf(L("Some replicas are not ready: %[1]d / %[2]d"), status.ReadyReplicas, status.Replicas)
	}

	if status.AvailableReplicas == 0 {
		return errors.New(L("the pod is not running"))
	}

	// Are the services running in the container?
	cnx := shared.NewConnection("kubectl", "", kubernetes.ServerFilter)
	if err := adm_utils.ExecCommand(zerolog.InfoLevel, cnx, "spacewalk-service", "status"); err != nil {
		return utils.Errorf(err, L("failed to run spacewalk-service status"))
	}
	return nil
}
