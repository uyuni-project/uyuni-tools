// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

//go:build !nok8s

package kubernetes

import (
	"errors"

	"github.com/spf13/cobra"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func ptfForKubernetes(globalFlags *types.GlobalFlags,
	flags *kubernetesInstallFlags,
	cmd *cobra.Command,
	args []string,
) error {
	return errors.New(L("PTF command for kubernetes is not implemented yet"))
}
