// SPDX-FileCopyrightText: 2023 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"os/exec"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/uyuni-project/uyuni-tools/shared"
	"github.com/uyuni-project/uyuni-tools/shared/kubernetes"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func ExecCommand(logLevel zerolog.Level, cnx *shared.Connection, args ...string) error {
	podName, err := cnx.GetPodName()
	if err != nil {
		log.Fatal().Err(err)
	}

	commandArgs := []string{"exec", podName}

	command, err := cnx.GetCommand()
	if err != nil {
		log.Fatal().Err(err)
	}

	if command == "kubectl" {
		commandArgs = append(commandArgs, "-c", "uyuni", "--")
	}

	commandArgs = append(commandArgs, "sh", "-c", strings.Join(args, " "))

	runCmd := exec.Command(command, commandArgs...)
	logger := utils.OutputLogWriter{Logger: log.Logger, LogLevel: logLevel}
	runCmd.Stdout = logger
	runCmd.Stderr = logger
	return runCmd.Run()
}

func RunningImage(cnx *shared.Connection, containerName string) (string, error) {
	command, err := cnx.GetCommand()

	switch command {
	case "podman":
		args := []string{"ps", "--format", "{{.Image}}", "--noheading"}
		image, err := utils.RunCmdOutput(zerolog.DebugLevel, "podman", args...)
		if err != nil {
			return "", err
		}
		return strings.Trim(string(image), "\n"), nil

	case "kubectl":

		//FIXME this will work until containers 0 is uyuni. Then jsonpath should be something like
		// {.items[0].spec.containers[?(@.name=="` + containerName + `")].image but there are problems
		// using RunCmdOutput with an arguments with round brackets
		args := []string{"get", "pods", kubernetes.ServerFilter, "-o", "jsonpath={.items[0].spec.containers[0].image}"}
		image, err := utils.RunCmdOutput(zerolog.DebugLevel, "kubectl", args...)

		log.Info().Msgf("image is: %s", image)
		if err != nil {
			return "", err
		}
		return strings.Trim(string(image), "\n"), nil
	}

	return command, err
}
