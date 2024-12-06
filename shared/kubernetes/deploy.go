// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package kubernetes

import (
	"strconv"
	"strings"

	"github.com/rs/zerolog"
)

// HasDeployment returns true when a deployment matching the kubectl get filter is existing in the namespace.
func HasDeployment(namespace string, filter string) bool {
	out, err := runCmdOutput(zerolog.DebugLevel, "kubectl", "get", "deploy", "-n", namespace, filter, "-o", "name")
	if err == nil && strings.TrimSpace(string(out)) != "" {
		return true
	}
	return false
}

// GetReplicas return the number of replicas of a deployment.
//
// If no such deployment exists, 0 will be returned as if there was a deployment scaled down to 0.
func GetReplicas(namespace string, name string) int {
	out, err := runCmdOutput(zerolog.DebugLevel,
		"kubectl", "get", "deploy", "-n", namespace, name, "-o", "jsonpath={.status.replicas}",
	)
	if err != nil {
		return 0
	}
	replicas, err := strconv.Atoi(strings.TrimSpace(string(out)))
	if err != nil {
		return 0
	}
	return replicas
}
