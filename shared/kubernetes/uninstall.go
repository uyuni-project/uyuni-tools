// SPDX-FileCopyrightText: 2023 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package kubernetes

import (
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
)

// Message appended in the uninstall commands for kubernetes.
func UninstallHelp() string {
	return L(`
Note that removing the volumes could also be handled automatically depending on the StorageClass used
when installed on a kubernetes cluster.

For instance on a default K3S install, the local-path-provider storage volumes will
be automatically removed when deleting the deployment even if --purge-volumes argument is not used.`)
}
