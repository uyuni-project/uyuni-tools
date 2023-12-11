// SPDX-FileCopyrightText: 2023 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package kubernetes

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// HelmUpgrade runs the helm upgrade command.
//
// To perform an installation, set the install parameter to true: helm would get the --install parameter.
// If repo is not empty, the --repo parameter will be passed.
// If version is not empty, the --version parameter will be passed.
func HelmUpgrade(kubeconfig string, namespace string, install bool,
	repo string, name string, chart string, version string, args ...string) {

	helmArgs := []string{
		"upgrade",
		"-n", namespace,
		"--create-namespace",
		name,
		chart,
	}
	if kubeconfig != "" {
		helmArgs = append(helmArgs, "--kubeconfig", kubeconfig)
	}

	if repo != "" {
		helmArgs = append(helmArgs, "--repo", repo)
	}
	if version != "" {
		helmArgs = append(helmArgs, "--version", version)
	}
	if install {
		helmArgs = append(helmArgs, "--install")
	}

	helmArgs = append(helmArgs, args...)

	command := "upgrade"
	if install {
		command = "install"
	}
	errorMessage := fmt.Sprintf("Failed to %s helm chart %s in namespace %s", command, chart, namespace)
	if err := utils.RunCmdStdMapping("helm", helmArgs...); err != nil {
		log.Fatal().Err(err).Msg(errorMessage)
	}
}
