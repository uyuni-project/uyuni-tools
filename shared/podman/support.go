// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"errors"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// RunSupportConfigOnProxyHost will run supportconfig command on host machine.
func RunSupportConfigOnHost(dir string) ([]string, error) {
	var files []string
	extensions := []string{"", ".md5"}

	// Run supportconfig on the host if installed
	if _, err := exec.LookPath("supportconfig"); err == nil {
		out, err := utils.RunCmdOutput(zerolog.DebugLevel, "supportconfig")
		if err != nil {
			return []string{}, utils.Errorf(err, L("failed to run supportconfig on the host"))
		}
		tarballPath := utils.GetSupportConfigPath(string(out))

		// Look for the generated supportconfig file
		if tarballPath != "" && utils.FileExists(tarballPath) {
			for _, ext := range extensions {
				files = append(files, tarballPath+ext)
			}
		} else {
			return []string{}, errors.New(L("failed to find host supportconfig tarball from command output"))
		}
	} else {
		log.Warn().Msg(L("supportconfig is not available on the host, skipping it"))
	}

	systemdDump, err := createSystemdDump(dir)
	if err != nil {
		log.Warn().Msg(L("systemd file are not present, skipping them"))
	} else {
		files = append(files, systemdDump)
	}

	containerList, err := runningContainer()
	if err != nil {
		return files, err
	}
	if len(containerList) > 0 {
		for _, container := range containerList {
			inspectDump, err := runPodmanInspectCommand(dir, container)
			if err != nil {
				log.Warn().Msgf(L("cannot podman inspect %s"), container)
			}
			files = append(files, inspectDump)

			bindedFilesDump, err := fetchBindedFileCommand(dir, container)
			if err != nil {
				log.Warn().Msgf(L("cannot fetch the binded files in %s"), container)
			}
			files = append(files, bindedFilesDump)

			logsDump, err := runPodmanLogsCommand(dir, container)
			if err != nil {
				log.Warn().Msgf(L("cannot podman logs %s"), container)
			}
			files = append(files, logsDump)
		}
	}

	return files, nil
}

func createSystemdDump(dir string) (string, error) {
	systemdSupportConfig, err := os.Create(path.Join(dir, "systemd-conf"))
	if err != nil {
		return "", utils.Errorf(err, L("cannot create %s"), systemdSupportConfig.Name())
	}

	out, err := utils.RunCmdOutput(zerolog.DebugLevel, "systemctl", "cat", "uyuni-*")
	if err != nil {
		return "", utils.Errorf(err, L("cannot run systemctl cat uyuni-proxy-pod"))
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
	defer podmanInspectDump.Close()
	if err != nil {
		return "", utils.Errorf(err, L("cannot create %s"), podmanInspectDump)
	}

	out, err := utils.RunCmdOutput(zerolog.DebugLevel, "podman", "inspect", container)
	if err != nil {
		return "", utils.Errorf(err, L("cannot run podman inspect %s"), container)
	}

	_, err = podmanInspectDump.WriteString("====podman inspect " + container + "====\n" + string(out))
	if err != nil {
		return "", err
	}
	return podmanInspectDump.Name(), nil
}

func fetchBindedFileCommand(dir string, container string) (string, error) {
	bindedFilesDump, err := os.Create(path.Join(dir, "binded-files-"+container))
	defer bindedFilesDump.Close()
	if err != nil {
		return "", utils.Errorf(err, L("cannot create %s"), bindedFilesDump)
	}

	_, err = bindedFilesDump.WriteString("====binded files====" + "\n")
	if err != nil {
		return "", err
	}
	out, err := utils.RunCmdOutput(zerolog.DebugLevel, "podman", "inspect", container, "--format", "{{range .Mounts}}{{if eq .Type \"bind\"}} {{.Source}}{{end}}{{end}}")
	if err != nil {
		return "", utils.Errorf(err, L("cannot run podman inspect %s"), container)
	}
	bindedFiles := strings.Split(string(out), " ")

	for _, bindFile := range bindedFiles {
		bindFile = strings.TrimSpace(bindFile)
		if len(bindFile) <= 0 {
			continue
		}
		out, err := utils.RunCmdOutput(zerolog.DebugLevel, "find", bindFile, "-type", "f")
		if err != nil {
			return "", err
		}

		fileList := strings.Split(strings.TrimSpace(string(out)), "\n")
		for _, file := range fileList {
			_, err = bindedFilesDump.WriteString("====" + file + "====" + "\n")
			if err != nil {
				return "", err
			}
			out, err := utils.RunCmdOutput(zerolog.DebugLevel, "cat", file)
			if err != nil {
				return "", err
			}
			_, err = bindedFilesDump.WriteString(string(out) + "\n")
			if err != nil {
				return "", err
			}
		}
	}
	return bindedFilesDump.Name(), nil
}

func runPodmanLogsCommand(dir string, container string) (string, error) {
	podmanLogsDump, err := os.Create(path.Join(dir, "logs-"+container))
	if err != nil {
		return "", utils.Errorf(err, L("cannot create %s"), podmanLogsDump)
	}

	out, err := utils.RunCmdOutput(zerolog.DebugLevel, "podman", "logs", container)
	if err != nil {
		return "", utils.Errorf(err, L("cannot run podman inspect %s"), container)
	}

	_, err = podmanLogsDump.WriteString("====podman logs====\n" + string(out))
	if err != nil {
		return "", err
	}
	return podmanLogsDump.Name(), nil
}

func runningContainer() ([]string, error) {
	containers, err := utils.RunCmdOutput(zerolog.DebugLevel, "podman", "ps", "-a", "--format={{ .Names }}")
	if err != nil {
		return []string{}, err
	}

	containerList := strings.Split(strings.TrimSpace(string(containers)), "\n")

	return containerList, nil
}
