// SPDX-FileCopyrightText: 2024 SUSE LLC
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
func RunSupportConfigOnPodmanHost(dir string) ([]string, error) {
	files, err := utils.RunSupportConfigOnHost(dir)
	if err != nil {
		return files, err
	}

	systemdDump, err := createSystemdDump(dir)
	if err != nil {
		log.Warn().Msg(L("No systemd file to add to the archive"))
	} else {
		files = append(files, systemdDump)
	}

	containerList, err := hostedContainers()
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
		out, err := utils.RunCmdOutput(zerolog.DebugLevel, "find", boundFile, "-type", "f")
		if err != nil {
			return "", err
		}

		fileList := strings.Split(strings.TrimSpace(string(out)), "\n")
		for _, file := range fileList {
			_, err = boundFilesDump.WriteString("====" + file + "====" + "\n")
			if err != nil {
				return "", err
			}
			out, err := utils.RunCmdOutput(zerolog.DebugLevel, "cat", file)
			if err != nil {
				return "", err
			}
			_, err = boundFilesDump.WriteString(string(out) + "\n")
			if err != nil {
				return "", err
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

func hostedContainers() ([]string, error) {
	systemdFiles, err := getSystemdFileList()
	if err != nil {
		return []string{}, err
	}
	servicesList := GetServicesFromSystemdFiles(string(systemdFiles))

	var containerList []string

	for _, service := range servicesList {
		service = strings.Replace(service, ".service", "", -1)
		containerList = append(containerList, strings.Replace(service, "@", "", -1))
	}

	return containerList, nil
}
