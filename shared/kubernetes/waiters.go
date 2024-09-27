// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

//go:build !nok8s

package kubernetes

import (
	"fmt"
	"strings"
	"time"

	"github.com/rs/zerolog"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// WaitForSecret waits for a secret to be available.
func WaitForSecret(namespace string, secret string) {
	for i := 0; ; i++ {
		if err := utils.RunCmd("kubectl", "get", "-n", namespace, "secret", secret); err == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}
}

// WaitForJob waits for a job to be completed before timeout seconds.
//
// If the timeout value is 0 the job will be awaited for for ever.
func WaitForJob(namespace string, name string, timeout int) error {
	for i := 0; ; i++ {
		status, err := jobStatus(namespace, name)
		if err != nil {
			return err
		}
		if status == "error" {
			return fmt.Errorf(
				L("%[1]s job failed, run kubectl logs -n %[2]s -ljob-name=%[1]s for details"),
				name, namespace,
			)
		}
		if status == "success" {
			return nil
		}

		if timeout > 0 && i == timeout {
			return fmt.Errorf(L("%[1]s job failed to complete within %[2]d seconds"), name, timeout)
		}
		time.Sleep(1 * time.Second)
	}
}

func jobStatus(namespace string, name string) (string, error) {
	out, err := utils.RunCmdOutput(
		zerolog.DebugLevel, "kubectl", "get", "job", "-n", namespace, name,
		"-o", "jsonpath={.status.succeeded},{.status.failed}",
	)
	if err != nil {
		return "", utils.Errorf(err, L("failed to get %s job status"), name)
	}
	results := strings.SplitN(strings.TrimSpace(string(out)), ",", 2)
	if len(results) != 2 {
		return "", fmt.Errorf(L("invalid job status response: '%s'"), string(out))
	}
	if results[0] == "1" {
		return "success", nil
	} else if results[1] == "1" {
		return "error", nil
	}
	return "", nil
}
