// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"os"

	"github.com/spf13/cobra"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func extract(globalFlags *types.GlobalFlags, flags *configFlags, cmd *cobra.Command, args []string) error {
	// Copy the generated file locally
	tmpDir, err := os.MkdirTemp("", "mgrpxy-*")
	if err != nil {
		return utils.Errorf(err, L("failed to create temporary directory"))
	}
	defer os.RemoveAll(tmpDir)

	hostSupportConfigFiles, err := utils.RunSupportConfigOnHost(tmpDir)
	if err != nil {
		return err
	}

	// TODO Get cluster infos in case of kubernetes

	if err := utils.CreateSupportConfigTarball(flags.Output, hostSupportConfigFiles); err != nil {
		return err
	}

	return nil
}

//func getAllProxyCnx(backend string) map[*shared.Connection]bool {
//	cnxs := make(map[*shared.Connection]bool)
//
//	/* this is as hack but it works. We loop for podman proxy container name since the proxy filter for kubernetes
//	* is just one. Storing the context in a unique list, we would have one result if kubernetes and one for each
//	* container for podman.
//	 */
//	for _, container := range podman.ProxyContainerNames {
//		cnx := shared.NewConnection(backend, container, kubernetes.ProxyFilter)
//		if cnxs[cnx] {
//			continue
//		} else {
//			cnxs[cnx] = true
//		}
//	}
//	return cnxs
//}
