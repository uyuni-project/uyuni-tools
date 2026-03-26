// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0
//go:build ptf

package kubernetes

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func ptfForKubernetes(_ *types.GlobalFlags,
	_ *kubernetesPTFFlags,
	_ *cobra.Command,
	_ []string,
) error {
	return nil
}
