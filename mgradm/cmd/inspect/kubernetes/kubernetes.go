// SPDX-FileCopyrightText: 2023 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

//go:build !nok8s

package kubernetes

import "github.com/uyuni-project/uyuni-tools/shared/types"

type kuberntesInspectFlags struct {
	Image types.ImageFlags `mapstructure:",squash"`
	Helm  cmd_utils.HelmFlags
}

//type kubernetesInstallFlags struct {
//	shared.InstallFlags `mapstructure:",squash"`
//}
//
//func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
//
//}
