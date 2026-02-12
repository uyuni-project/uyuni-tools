// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package kubernetes

import (
	"bytes"
	"os/exec"

	"github.com/rs/zerolog"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// HasHelmRelease returns whether a helm release is installed or not, even if it failed.
func HasHelmRelease(release string, kubeconfig string) bool {
	if _, err := exec.LookPath("helm"); err == nil {
		args := []string{}
		if kubeconfig != "" {
			args = append(args, "--kubeconfig", kubeconfig)
		}
		args = append(args, "list", "-aAq", "--no-headers", "-f", release)
		out, err := utils.RunCmdOutput(zerolog.TraceLevel, "helm", args...)
		return len(bytes.TrimSpace(out)) != 0 && err == nil
	}
	return false
}
