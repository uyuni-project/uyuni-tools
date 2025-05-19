// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"os"
	"path"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// RunSupportConfigOnPodmanHost will run supportconfig command on podman machine.
func RunSupportConfigOnPodmanHost(systemd Systemd, dir string) ([]string, error) {
	files, err := utils.RunSupportConfigOnHost()
	if err != nil {
		return files, err
	}
	logDump, err := createLogDump(dir)
	if err != nil {
		log.Warn().Msg(L("No logs file on host to add to the archive"))
	} else {
		files = append(files, logDump)
	}

	systemdDump, err := createSystemdDump(dir)
	if err != nil {
		log.Warn().Msg(L("No systemd file to add to the archive"))
	} else {
		files = append(files, systemdDump)
	}

	containerList, err := hostedContainers(systemd)
	if err != nil {
		return files, err
	}
	if len(containerList) > 0 {
		for _, container := range containerList {
			inspectDump, err := runPodmanInspectCommand(dir, container)
			if err != nil {
				log.Warn().Err(err).Msgf(L("failed to run podman inspect %s"), container)
			}
			files = append(files, inspectDump)

			boundFilesDump, err := fetchBoundFileCommand(dir, container)
			if err != nil {
				log.Warn().Err(err).Msgf(L("failed to fetch the config files bound to container %s"), container)
			}
			files = append(files, boundFilesDump)

			logsDump, err := runJournalCtlCommand(dir, container)
			if err != nil {
				log.Warn().Err(err).Msgf(L("failed to run podman logs %s"), container)
			}
			files = append(files, logsDump)
		}
	}

	return files, nil
}

func createLogDump(dir string) (string, error) {
	logConfig, err := os.Create(path.Join(dir, "logs"))
	if err != nil {
		return "", utils.Errorf(err, L("failed to create %s file"), logConfig.Name())
	}

	out, err := utils.RunCmdOutput(zerolog.DebugLevel, "cat", utils.GlobalLogPath)
	if err != nil {
		return "", utils.Errorf(err, L("failed to cat %s"), utils.GlobalLogPath)
	}
	defer logConfig.Close()

	_, err = logConfig.WriteString("====cat " + utils.GlobalLogPath + "====\n" + string(out))
	if err != nil {
		return "", err
	}

	return logConfig.Name(), nil
}

func createSystemdDump(dir string) (string, error) {
	systemdSupportConfig, err := os.Create(path.Join(dir, "systemd-conf"))
	if err != nil {
		return "", utils.Errorf(err, L("failed to create %s file"), systemdSupportConfig.Name())
	}

	out, err := utils.RunCmdOutput(zerolog.DebugLevel, "systemctl", "cat", "uyuni-*")
	if err != nil {
		return "", utils.Errorf(err, L("failed to run systemctl cat uyuni-*"))
	}
	defer systemdSupportConfig.Close()

	_, err = systemdSupportConfig.WriteString("====systemctl cat uyuni-*====\n" + string(out))
	if err != nil {
		return "", err
	}

	return systemdSupportConfig.Name(), nil
}

func runPodmanInspectCommand(dir string, container string) (string, error) {
	podmanInspectDump, err := os.Create(path.Join(dir, "inspect-"+container))
	defer func() {
		if err := podmanInspectDump.Close(); err != nil {
			log.Error().Err(err).Msg(L("failed to close inspect dump file"))
		}
	}()
	if err != nil {
		return "", utils.Errorf(err, L("failed to create %s file"), podmanInspectDump)
	}

	out, err := utils.RunCmdOutput(zerolog.DebugLevel, "podman", "inspect", container)
	if err != nil {
		return "", utils.Errorf(err, L("failed to run podman inspect %s"), container)
	}

	_, err = podmanInspectDump.WriteString("====podman inspect " + container + "====\n" + string(out))
	if err != nil {
		return "", err
	}
	return podmanInspectDump.Name(), nil
}

func fetchBoundFileCommand(dir string, container string) (string, error) {
	boundFilesDump, err := os.Create(path.Join(dir, "bound-files-"+container))
	defer func() {
		if err := boundFilesDump.Close(); err != nil {
			log.Error().Err(err).Msg(L("failed to close bound files"))
		}
	}()
	if err != nil {
		return "", utils.Errorf(err, L("failed to create %s file"), boundFilesDump)
	}

	_, err = boundFilesDump.WriteString("====bound files====" + "\n")
	if err != nil {
		return "", err
	}
	out, err := utils.RunCmdOutput(
		zerolog.DebugLevel, "podman", "inspect", container,
		"--format", "{{range .Mounts}}{{if eq .Type \"bind\"}} {{.Source}}{{end}}{{end}}",
	)
	if err != nil {
		return "", utils.Errorf(err, L("failed to run podman inspect %s"), container)
	}
	boundFiles := strings.Split(string(out), " ")

	for _, boundFile := range boundFiles {
		boundFile = strings.TrimSpace(boundFile)
		if len(boundFile) <= 0 {
			continue
		}
		if stat, err := os.Stat(boundFile); err == nil && stat.Mode().IsRegular() {
			_, err = boundFilesDump.WriteString("====" + boundFile + "====" + "\n")
			if allErrors := utils.JoinErrors(err, utils.CopyFile(boundFile, boundFilesDump)); allErrors != nil {
				return "", allErrors
			}
		}
	}
	return boundFilesDump.Name(), nil
}

func runJournalCtlCommand(dir string, container string) (string, error) {
	journalctlDump, err := os.Create(path.Join(dir, "journalctl-"+container))
	if err != nil {
		return "", utils.Errorf(err, L("failed create %s file"), journalctlDump)
	}

	out, err := utils.RunCmdOutput(zerolog.DebugLevel, "journalctl", "-u", container)
	if err != nil {
		return "", utils.Errorf(err, L("failed to run journalctl -u %s"), container)
	}

	_, err = journalctlDump.WriteString("====journalctl====\n" + string(out))
	if err != nil {
		return "", err
	}
	return journalctlDump.Name(), nil
}

func getSystemdFileList() ([]byte, error) {
	return utils.RunCmdOutput(
		zerolog.DebugLevel, "find", "/etc/systemd/system", "-maxdepth", "1", "-name", "uyuni-*service",
	)
}

func hostedContainers(systemd Systemd) ([]string, error) {
	systemdFiles, err := getSystemdFileList()
	if err != nil {
		return []string{}, err
	}
	servicesList := systemd.GetServicesFromSystemdFiles(string(systemdFiles))

	var containerList []string

	for _, service := range servicesList {
		service = strings.Replace(service, ".service", "", -1)
		containerList = append(containerList, strings.Replace(service, "@", "", -1))
	}

	return containerList, nil
}
