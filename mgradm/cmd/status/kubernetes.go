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
	_ *types.GlobalFlags,
	_ *statusFlags,
	_ *cobra.Command,
	_ []string,
) error {
	cnx := shared.NewConnection("kubectl", "", kubernetes.ServerFilter)
	namespace, err := cnx.GetNamespace("")
	if err != nil {
		return utils.Errorf(err, L("failed to find the uyuni deployment namespace"))
	}

	// Is the pod running? Do we have all the replicas?
	status, err := kubernetes.GetDeploymentStatus(namespace, kubernetes.ServerApp)
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
	if err := adm_utils.ExecCommand(zerolog.InfoLevel, cnx, "spacewalk-service", "status"); err != nil {
		return utils.Errorf(err, L("failed to run spacewalk-service status"))
	}
	return nil
}
