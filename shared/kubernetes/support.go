// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package kubernetes

import (
	"errors"
	"os"
	"os/exec"
	"path"

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

	configmapFilename, err := fetchConfigMap(dir)
	if err != nil {
		log.Warn().Msg(L("cannot retrieve any configmap. This is expected in no kubernetes host"))
	} else {
		files = append(files, configmapFilename)
	}

	return files, nil
}

func fetchConfigMap(dir string) (string, error) {
	configmapFile, err := os.Create(path.Join(dir, "configmap"))
	if err != nil {
		return "", utils.Errorf(err, L("cannot create %s"), configmapFile.Name())
	}
	defer configmapFile.Close()
	out, err := utils.RunCmdOutput(zerolog.DebugLevel, "kubectl", "get", "configmap", "-o", "yaml")
	if err != nil {
		return "", utils.Errorf(err, L("cannot fetch configmap"))
	}

	_, err = configmapFile.WriteString(string(out))
	if err != nil {
		return "", err
	}
	return configmapFile.Name(), nil
}
