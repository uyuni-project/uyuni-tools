// SPDX-FileCopyrightText: 2023 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package kubernetes

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog"
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

// HelmUninstall runs the helm uninstall command to remove a deployment.
func HelmUninstall(kubeconfig string, deployment string, filter string, dryRun bool) string {
	jsonpath := fmt.Sprintf("jsonpath={.items[?(@.metadata.name==\"%s\")].metadata.namespace}", deployment)
	args := []string{"get", "-A", "deploy", "-o", jsonpath}
	if filter != "" {
		args = append(args, filter)
	}

	out, err := utils.RunCmdOutput(zerolog.DebugLevel, "kubectl", args...)
	if err != nil {
		log.Info().Err(err).Msgf("Failed to find %s's namespace, skipping removal", deployment)
	}
	namespace := string(out)
	if namespace != "" {
		helmArgs := []string{}
		if kubeconfig != "" {
			helmArgs = append(helmArgs, "--kubeconfig", kubeconfig)
		}
		helmArgs = append(helmArgs, "uninstall", "-n", namespace, deployment)

		if dryRun {
			log.Info().Msgf("Would run helm %s", strings.Join(helmArgs, " "))
		} else {
			log.Info().Msgf("Uninstalling %s", deployment)
			message := "Failed to run helm " + strings.Join(helmArgs, " ")
			err := utils.RunCmd("helm", helmArgs...)
			if err != nil {
				log.Fatal().Err(err).Msg(message)
			}
		}
	}
	return namespace
}
