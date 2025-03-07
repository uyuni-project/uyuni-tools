// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package create

import (
	"errors"
	"fmt"
	"os"
	"path"
	"slices"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/uyuni-project/uyuni-tools/mgradm/cmd/backup/shared"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

var systemd podman.Systemd = podman.SystemdImpl{}
var runCmdOutput = utils.RunCmdOutput

func Create(
	_ *types.GlobalFlags,
	flags *shared.Flagpole,
	_ *cobra.Command,
	args []string,
) error {
	dryRun := flags.DryRun
	outputDirectory := args[0]
	printIntro(outputDirectory, flags)

	if err := SanityChecks(outputDirectory, dryRun); err != nil {
		return shared.AbortError(err, false)
	}

	volumesBackupPath := path.Join(outputDirectory, shared.VolumesSubdir)
	imagesBackupPath := path.Join(outputDirectory, shared.ImagesSubdir)

	if err := prepareOuputDirs([]string{outputDirectory, volumesBackupPath, imagesBackupPath}, dryRun); err != nil {
		return shared.AbortError(err, false)
	}

	volumes := gatherVolumesToBackup(flags.ExtraVolumes, flags.SkipVolumes, flags.SkipDatabase)
	images := gatherContainerImagesToBackup(flags.SkipImages)

	if !dryRun {
		if err := shared.StorageCheck(volumes, images, outputDirectory); err != nil {
			return shared.AbortError(err, false)
		}
	}

	// stop service if database is to be backedup. Otherwise do a live backup
	serviceStopped := false
	if !flags.SkipDatabase && !dryRun {
		log.Info().Msg(L("Stopping server service"))
		if err := systemd.StopService(podman.ServerService); err != nil {
			return shared.AbortError(err, false)
		}
		serviceStopped = true
	}

	if err := backupVolumes(volumes, volumesBackupPath, dryRun); err != nil {
		return shared.AbortError(err, true)
	}

	// Remaining backups are not critical, retore can create default values
	// so let's only track if there was an error
	hasError := backupContainerImages(images, imagesBackupPath, dryRun)

	// systemd configuration backup is optional as we have defaults to use
	hasError = errors.Join(hasError, backupSystemdServices(outputDirectory, dryRun))

	// podman configuration backup is optional as we have defaults to use
	hasError = errors.Join(hasError, backupPodmanConfiguration(outputDirectory, dryRun))

	// start service if it was stopped before
	if serviceStopped && !flags.NoRestart && !dryRun {
		log.Info().Msg(L("Restarting server service"))
		hasError = errors.Join(hasError, systemd.StartService(podman.ServerService))
	}

	log.Info().Msgf(L("Backup finished into %s"), outputDirectory)
	return shared.ReportError(hasError)
}

func printIntro(outputDir string, flags *shared.Flagpole) {
	log.Debug().Msg("Creating backup with options:")
	log.Debug().Msgf("output directory: %s", outputDir)
	log.Debug().Msgf("dry run: %t", flags.DryRun)
	log.Debug().Msgf("skip database: %t", flags.SkipDatabase)
	log.Debug().Msgf("skip config: %t", flags.SkipConfig)
	log.Debug().Msgf("skip restart: %t", flags.NoRestart)
	log.Debug().Msgf("skip images: %t", flags.SkipImages)
	log.Debug().Msgf("skip volumes: %s", flags.SkipVolumes)
	log.Debug().Msgf("extra volumes: %s", flags.ExtraVolumes)
}

func prepareOuputDirs(outputDirs []string, dryRun bool) error {
	for _, d := range outputDirs {
		if dryRun {
			log.Info().Msgf(L("Would create '%s' directory"), d)
		} else {
			if err := os.Mkdir(d, 0622); err != nil {
				return fmt.Errorf(L("unable to create target output directory: %w"), err)
			}
		}
	}
	return nil
}

func gatherVolumesToBackup(extraVolumes []string, skipVolumes []string, skipDatabase bool) []string {
	// Construct work volume list
	volumes := extraVolumes

	// Extra handling to skip all, except extra added
	if len(skipVolumes) == 1 && skipVolumes[0] == "all" {
		return volumes
	}

	if skipDatabase {
		skipVolumes = append(skipVolumes, utils.VarPgsqlVolumeMount.Name)
	}
	for _, volume := range utils.ServerVolumeMounts {
		if !slices.Contains(skipVolumes, volume.Name) {
			volumes = append(volumes, volume.Name)
		}
	}
	return volumes
}

func backupVolumes(volumes []string, outputDirectory string, dryRun bool) error {
	log.Info().Msg(L("Backing up container volumes"))
	for _, volume := range volumes {
		log.Debug().Msgf("Backing up %s volume", volume)
		if err := podman.ExportVolume(volume, outputDirectory, dryRun); err != nil {
			return err
		}
	}
	return nil
}

func gatherContainerImagesToBackup(skipImages bool) []string {
	images := []string{}

	if !skipImages {
		for _, service := range utils.UyuniServices {
			images = append(images, service.Image.Name)
		}
	}
	return images
}

func backupContainerImages(images []string, outputDirectory string, dryRun bool) error {
	log.Info().Msg(L("Backing up container images"))
	var hasError error
	for _, image := range images {
		log.Debug().Msgf("Backing up image %s", image)
		if err := podman.ExportImage(image, outputDirectory, dryRun); err != nil {
			log.Warn().Err(err).Msgf(L("Not backing up image %s"), image)
			hasError = errors.Join(hasError, err)
		}
	}
	return hasError
}

func backupSystemdServices(outputDirectory string, dryRun bool) error {
	errorMessage := L("Systemd services and configuration was not backed up")
	log.Info().Msg(L("Backing up Systemd services"))

	if err := exportSystemdConfiguration(outputDirectory, dryRun); err != nil {
		log.Warn().Err(err).Msg(errorMessage)
		return err
	}
	if err := utils.CreateChecksum(path.Join(outputDirectory, shared.SystemdConfBackupFile)); err != nil {
		log.Warn().Err(err).Msg(errorMessage)
		return err
	}
	return nil
}

func backupPodmanConfiguration(outputDirectory string, dryRun bool) error {
	errorMessage := L("Podman configuration was not backed up")
	log.Info().Msg(L("Backing up podman configuration"))
	if err := exportPodmanConfiguration(outputDirectory, dryRun); err != nil {
		log.Warn().Err(err).Msg(errorMessage)
		return err
	}
	if err := utils.CreateChecksum(path.Join(outputDirectory, shared.PodmanConfBackupFile)); err != nil {
		log.Warn().Err(err).Msg(errorMessage)
		return err
	}
	return nil
}

func SanityChecks(outputDirectory string, dryRun bool) error {
	if err := shared.SanityChecks(); err != nil {
		return err
	}

	if utils.FileExists(outputDirectory) {
		if !utils.IsEmptyDirectory(outputDirectory) {
			return fmt.Errorf(L("output directory %s already exists and is not empty"), outputDirectory)
		}
	}

	hostData, err := podman.InspectHost()
	if err != nil {
		return err
	}

	if !hostData.HasUyuniServer {
		return errors.New(L("server is not initialized."))
	}

	return nil
}
