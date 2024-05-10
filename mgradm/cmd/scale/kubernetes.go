// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

//go:build !nok8s

package scale

import (
	"errors"

	"github.com/spf13/cobra"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func kubernetesScale(
	globalFlags *types.GlobalFlags,
	flags *scaleFlags,
	cmd *cobra.Command,
	args []string,
) error {
	return errors.New(L("kubernetes not supported yet"))
}
