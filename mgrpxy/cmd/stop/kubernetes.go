// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package stop

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared"
	"github.com/uyuni-project/uyuni-tools/shared/kubernetes"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func kubernetesStop(
	_ *types.GlobalFlags,
	_ *stopFlags,
	_ *cobra.Command,
	_ []string,
) error {
	cnx := shared.NewConnection("kubectl", "", kubernetes.ProxyFilter)
	namespace, err := cnx.GetNamespace("")
	if err != nil {
		return utils.Errorf(err, L("failed retrieving namespace"))
	}
	return kubernetes.Stop(namespace, kubernetes.ProxyApp)
}
