// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/templates"
	"github.com/uyuni-project/uyuni-tools/shared"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// ExecCommand execute commands passed as argument in the current system.
func ExecCommand(logLevel zerolog.Level, cnx *shared.Connection, args ...string) error {
	podName, err := cnx.GetPodName()
	if err != nil {
		return utils.Error(err, L("exec command failed"))
	}

	commandArgs := []string{"exec", podName}

	command, err := cnx.GetCommand()
	if err != nil {
		log.Fatal().Err(err)
	}

	if command == "kubectl" {
		namespace, err := cnx.GetNamespace("")
		if namespace == "" {
			return utils.Error(err, L("failed retrieving namespace"))
		}
		commandArgs = append(commandArgs, "-n", namespace, "-c", "uyuni", "--")
	}

	commandArgs = append(commandArgs, "sh", "-c", strings.Join(args, " "))

	runCmd := exec.Command(command, commandArgs...)
	logger := log.Logger.Level(logLevel)
	runCmd.Stdout = logger
	runCmd.Stderr = logger
	return runCmd.Run()
}

// GeneratePostUpgradeScript generates the script to be run after upgrade.
func GeneratePostUpgradeScript() (string, error) {
	data := templates.PostUpgradeTemplateData{}

	scriptBuilder := new(strings.Builder)
	if err := data.Render(scriptBuilder); err != nil {
		return "", utils.Error(err, L("failed to render database post upgrade script"))
	}
	return scriptBuilder.String(), nil
}

// SanityCheck verifies if an upgrade can be run.
func SanityCheck(inspectedValues *utils.InspectData) error {
	// Skip the uyuni / SUSE Manager release checks if the runningValues is nil.
	if inspectedValues == nil {
		return nil
	}

	isUyuni := inspectedValues.ContainerInspectData.UyuniRelease != ""
	isUyuniImage := inspectedValues.ServerInspectData.UyuniRelease != ""
	isSumaImage := inspectedValues.ServerInspectData.SuseManagerRelease != ""

	if isUyuni && isSumaImage {
		return fmt.Errorf(
			L("currently SUSE Manager %s is installed, instead the image is Uyuni. Upgrade is not supported"),
			inspectedValues.ContainerInspectData.SuseManagerRelease,
		)
	}

	if !isUyuni && isUyuniImage {
		return fmt.Errorf(
			L("currently Uyuni %s is installed, instead the image is SUSE Manager. Upgrade is not supported"),
			inspectedValues.ContainerInspectData.UyuniRelease,
		)
	}

	if isUyuni {
		currentUyuniRelease := inspectedValues.ContainerInspectData.UyuniRelease
		log.Debug().Msgf("Current release is %s", string(currentUyuniRelease))
		if !isUyuniImage {
			return errors.New(L("cannot fetch release from server image"))
		}
		log.Debug().Msgf("Server image release is %s", inspectedValues.ServerInspectData.UyuniRelease)
		if utils.CompareVersion(inspectedValues.ServerInspectData.UyuniRelease, string(currentUyuniRelease)) < 0 {
			return fmt.Errorf(
				L("cannot downgrade from version %[1]s to %[2]s"),
				string(currentUyuniRelease), inspectedValues.ServerInspectData.UyuniRelease,
			)
		}
	} else {
		currentSuseManagerRelease := inspectedValues.ContainerInspectData.SuseManagerRelease
		log.Debug().Msgf("Current release is %s", currentSuseManagerRelease)
		if !isSumaImage {
			return errors.New(L("cannot fetch release from server image"))
		}
		log.Debug().Msgf("Server image release is %s", inspectedValues.ServerInspectData.SuseManagerRelease)
		if utils.CompareVersion(inspectedValues.ServerInspectData.SuseManagerRelease, currentSuseManagerRelease) < 0 {
			return fmt.Errorf(
				L("cannot downgrade from version %[1]s to %[2]s"),
				currentSuseManagerRelease, inspectedValues.ServerInspectData.SuseManagerRelease,
			)
		}
	}
	return nil
}
