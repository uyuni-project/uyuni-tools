// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"errors"
	"fmt"
	"os/exec"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
)

// GetSupportConfigPath returns the support config tarball path.
func GetSupportConfigPath(out string) string {
	re := regexp.MustCompile(`/var/log/scc_(.*?)\.txz`)
	return re.FindString(out)
}

// GetSupportConfigFileSaveName returns the support config file name.
func GetSupportConfigFileSaveName() string {
	out, err := RunCmdOutput(zerolog.DebugLevel, "hostname")
	var hostname string
	if err != nil {
		log.Warn().Err(err).Msg(L("Unable to detect hostname, using localhost"))
		hostname = "localhost"
	} else {
		hostname = strings.TrimSpace(string(out))
	}
	now := time.Now()
	return fmt.Sprintf("scc_%s_%s", hostname, now.Format("20060102_1504"))
}

/* CreateSupportConfigTarball will create a tarball in outputFolder with all the supportconfig
 * files generated by host and pods.
 */
func CreateSupportConfigTarball(outputFolder string, files []string) error {
	// Pack it all into a tarball
	log.Info().Msg(L("Preparing the tarball"))

	supportFileName := GetSupportConfigFileSaveName()
	supportFilePath := path.Join(outputFolder, fmt.Sprintf("%s.tar.gz", supportFileName))

	tarball, err := NewTarGz(supportFilePath)
	if err != nil {
		return err
	}

	for _, file := range files {
		if err := tarball.AddFile(file, path.Join(supportFileName, path.Base(file))); err != nil {
			return Errorf(err, L("failed to add %s to tarball"), path.Base(file))
		}
	}
	tarball.Close()
	return nil
}

// GetContainersFromSystemdFiles parse a string of systemdfile and return a list of containers.
func GetContainersFromSystemdFiles(systemdFileList string) []string {
	serviceList := strings.Replace(string(systemdFileList), "/etc/systemd/system/", "", -1)
	containers := strings.Replace(serviceList, ".service", "", -1)

	containerList := strings.Split(strings.TrimSpace(containers), "\n")

	var trimmedContainers []string
	for _, container := range containerList {
		trimmedContainers = append(trimmedContainers, strings.TrimSpace(container))
	}
	return trimmedContainers
}

// RunSupportConfigOnHost will run supportconfig command on host machine.
func RunSupportConfigOnHost(dir string) ([]string, error) {
	var files []string
	extensions := []string{"", ".md5"}

	// Run supportconfig on the host if installed
	if _, err := exec.LookPath("supportconfig"); err == nil {
		out, err := RunCmdOutput(zerolog.DebugLevel, "supportconfig")
		if err != nil {
			return []string{}, Errorf(err, L("failed to run supportconfig on the host"))
		}
		tarballPath := GetSupportConfigPath(string(out))

		// Look for the generated supportconfig file
		if tarballPath != "" && FileExists(tarballPath) {
			for _, ext := range extensions {
				files = append(files, tarballPath+ext)
			}
		} else {
			return []string{}, errors.New(L("failed to find host supportconfig tarball from command output"))
		}
	} else {
		log.Warn().Msg(L("supportconfig is not available on the host, skipping it"))
	}
	return files, nil
}
