// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

//go:build !nok8s

package start

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared"
	"github.com/uyuni-project/uyuni-tools/shared/kubernetes"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func kubernetesStart(
	globalFlags *types.GlobalFlags,
	flags *startFlags,
	cmd *cobra.Command,
	args []string,
) error {
	cnx := shared.NewConnection("kubectl", "", kubernetes.ServerFilter)
	namespace, err := cnx.GetNamespace("")
	if err != nil {
		return utils.Errorf(err, L("failed retrieving namespace"))
	}
	return kubernetes.Start(namespace, kubernetes.ServerApp)
}
