// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package kubernetes

import (
	"bytes"
	"encoding/json"
	"errors"
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
		return fmt.Errorf(L("failed to %s helm chart %s in namespace %s")+": %s", command, chart, namespace, err)
	}
	return nil
}

// HelmUninstall runs the helm uninstall command to remove a deployment.
func HelmUninstall(kubeconfig string, deployment string, filter string, dryRun bool) (string, error) {
	helmArgs := []string{}
	if kubeconfig != "" {
		helmArgs = append(helmArgs, "--kubeconfig", kubeconfig)
	}

	jsonpath := fmt.Sprintf("jsonpath={.items[?(@.metadata.name==\"%s\")].metadata.namespace}", deployment)
	args := []string{"get", "-A", "deploy", "-o", jsonpath}
	if filter != "" {
		args = append(args, filter)
	}

	out, err := utils.RunCmdOutput(zerolog.DebugLevel, "kubectl", args...)
	if err != nil {
		log.Info().Err(err).Msgf(L("Failed to find %s's namespace, skipping removal"), deployment)
	}

	namespace := string(out)
	if namespace == "" {
		log.Debug().Msgf("Pod is not running, trying to find the namespace using the helm release")
		namespace, err = FindNamespace(deployment, kubeconfig)
		if err != nil {
			log.Info().Err(err).Msgf(L("Cannot guess namespace"))
			return "", nil
		}
	}

	if namespace != "" {
		helmArgs = append(helmArgs, "uninstall", "-n", namespace, deployment)

		if dryRun {
			log.Info().Msgf(L("Would run %s"), "helm "+strings.Join(helmArgs, " "))
		} else {
			log.Info().Msgf(L("Uninstalling %s"), deployment)
			if err := utils.RunCmd("helm", helmArgs...); err != nil {
				return namespace, fmt.Errorf(L("failed to run helm %s: %s"), strings.Join(helmArgs, " "), err)
			}
		}
	}
	return namespace, nil
}

// FindNamespace tries to find the deployment namespace using helm.
func FindNamespace(deployment string, kubeconfig string) (string, error) {
	args := []string{}
	if kubeconfig != "" {
		args = append(args, "--kubeconfig", kubeconfig)
	}
	args = append(args, "list", "-aA", "-f", deployment, "-o", "json")
	out, err := utils.RunCmdOutput(zerolog.DebugLevel, "helm", args...)
	if err != nil {
		return "", fmt.Errorf(L("failed to detect %s's namespace using helm: %s"), deployment, err)
	}
	var data []releaseInfo
	if err = json.Unmarshal(out, &data); err != nil {
		return "", fmt.Errorf(L("helm provided an invalid JSON output: %s"), err)
	}

	if len(data) == 1 {
		return data[0].Namespace, nil
	}
	return "", errors.New(L("found no or more than one deployment"))
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

type releaseInfo struct {
	Namespace string `mapstructure:"namespace"`
}
