// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package create

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
	backup "github.com/uyuni-project/uyuni-tools/mgradm/cmd/backup/shared"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

var systemd podman.Systemd = podman.SystemdImpl{}

func exportSystemdConfiguration(outputDir string, dryRun bool) error {
	filesToBackup := gatherSystemdItems()

	if dryRun {
		log.Info().Msgf(L("Would backup %s"), filesToBackup)
		return nil
	}
	// Create output file
	out, err := os.Create(path.Join(outputDir, backup.SystemdConfBackupFile))
	if err != nil {
		return fmt.Errorf(L("failed to create Systemd backup tarball: %w"), err)
	}
	defer out.Close()

	// Prepare tar buffer
	tw := tar.NewWriter(out)
	defer tw.Close()

	for _, fileToBackup := range filesToBackup {
		f, err := os.Open(fileToBackup)
		if err != nil {
			return err
		}
		fstat, _ := f.Stat()
		h, err := tar.FileInfoHeader(fstat, "")
		if err != nil {
			return err
		}
		// Produced header does not have full path, overwrite it
		h.Name = fileToBackup
		if fstat.IsDir() {
			h.Name += "/"
		}
		if err := tw.WriteHeader(h); err != nil {
			return err
		}
		if fstat.IsDir() {
			continue
		}
		if _, err := io.Copy(tw, f); err != nil {
			return err
		}
	}
	return nil
}

// For each container get service file, service.d and its content.
func gatherSystemdItems() []string {
	result := []string{}

	for _, service := range utils.UyuniServices {
		serviceName, skip := findService(service.Name)
		if skip {
			continue
		}

		servicePath, err := systemd.GetServiceProperty(serviceName, podman.FragmentPath)
		if err != nil {
			log.Debug().Err(err).Msgf("failed to get the path to the %s service file", serviceName)
			// Skipping the dropins since we would likely get a similar error.
			continue
		}
		result = append(result, servicePath)

		// Get the drop in files
		dropIns, err := systemd.GetServiceProperty(serviceName, podman.DropInPaths)
		if err != nil {
			log.Debug().Err(err).Msgf("failed to get the path to the %s service configuration files", serviceName)
		} else {
			dropIns := strings.Split(dropIns, " ")
			result = append(result, filepath.Dir(dropIns[0]))
			result = append(result, dropIns[:]...)
		}
	}
	return result
}

func findService(name string) (serviceName string, skip bool) {
	skip = false
	serviceName = name
	if !systemd.HasService(serviceName) {
		// with optional or more replicas we have service template, check if the service exists at all
		serviceName = name + "@"
		if !systemd.HasService(serviceName) {
			log.Debug().Msgf("No service found for %s, skipping", name)
			skip = true
		}
	}
	return
}
