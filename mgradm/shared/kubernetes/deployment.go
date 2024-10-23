// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

//go:build !nok8s

package kubernetes

import (
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// ServerDeployName is the name of the server deployment.
const ServerDeployName = "uyuni"

var runCmdOutput = utils.RunCmdOutput

// getRunningServerImage extracts the main server container image from a running deployment.
func getRunningServerImage(namespace string) string {
	out, err := runCmdOutput(
		zerolog.DebugLevel, "kubectl", "get", "deploy", "-n", namespace, ServerDeployName,
		"-o", "jsonpath={.spec.template.spec.containers[0].image}",
	)
	if err != nil {
		// Errors could be that the namespace or deployment doesn't exist, just return no image.
		log.Debug().Err(err).Msg("failed to get the running server container image")
		return ""
	}
	return strings.TrimSpace(string(out))
}
