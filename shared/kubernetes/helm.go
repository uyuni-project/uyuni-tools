// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package kubernetes

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// HelmUpgrade runs the helm upgrade command.
//
// To perform an installation, set the install parameter to true: helm would get the --install parameter.
// If repo is not empty, the --repo parameter will be passed.
// If version is not empty, the --version parameter will be passed.
func HelmUpgrade(kubeconfig string, namespace string, install bool,
	repo string, name string, chart string, version string, args ...string) error {
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
	if err := utils.RunCmdStdMapping(zerolog.DebugLevel, "helm", helmArgs...); err != nil {
		// TODO We cannot use the command variable in the message as that would break localization
		if command == "upgrade" {
			return utils.Errorf(err, L("failed to upgrade helm chart %[1]s in namespace %[2]s"), chart, namespace)
		} else if command == "install" {
			return utils.Errorf(err, L("failed to install helm chart %[1]s in namespace %[2]s"), chart, namespace)
		}
	}
	return nil
}

// HelmUninstall runs the helm uninstall command to remove a deployment.
func HelmUninstall(namespace string, kubeconfig string, deployment string, dryRun bool) error {
	if namespace == "" {
		return fmt.Errorf(L("namespace is required"))
	}

	helmArgs := []string{}
	if kubeconfig != "" {
		helmArgs = append(helmArgs, "--kubeconfig", kubeconfig)
	}
	helmArgs = append(helmArgs, "uninstall", "-n", namespace, deployment)

	if dryRun {
		log.Info().Msgf(L("Would run %s"), "helm "+strings.Join(helmArgs, " "))
	} else {
		log.Info().Msgf(L("Uninstalling %s"), deployment)
		if err := utils.RunCmd("helm", helmArgs...); err != nil {
			return utils.Errorf(err, L("failed to run helm %s"), strings.Join(helmArgs, " "))
		}
	}
	return nil
}

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
