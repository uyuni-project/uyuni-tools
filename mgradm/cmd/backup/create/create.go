// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package create

import (
	"slices"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	backup "github.com/uyuni-project/uyuni-tools/mgradm/cmd/backup/shared"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

var systemd podman.Systemd = podman.SystemdImpl{}
var runCmdOutput = utils.RunCmdOutput

func Create(
	_ *types.GlobalFlags,
	flags *backup.Flagpole,
	_ *cobra.Command,
	args []string,
) error {
	dryRun := flags.DryRun
	outputDirectory := args[0]

	printIntro(outputDirectory, flags)

	if err := backup.SanityChecks(outputDirectory, dryRun); err != nil {
		return backup.ReportError(err, false)
	}

	volumes := gatherVolumesToBackup(flags.ExtraVolumes, flags.SkipVolumes, flags.SkipDatabase)
	images := gatherContainerImagesToBackup(flags.SkipImages)

	if !dryRun {
		if err := backup.StorageCheck(volumes, images, outputDirectory); err != nil {
			return backup.ReportError(err, false)
		}
	}

	// stop service if database is to be backedup. Otherwise do a live backup
	serviceStopped := false
	if !flags.SkipDatabase && !dryRun {
		log.Info().Msg(L("Stopping server service"))
		if err := systemd.StopService(podman.ServerService); err != nil {
			return backup.ReportError(err, false)
		}
		serviceStopped = true
	}

	if err := backupVolumes(volumes, outputDirectory, dryRun); err != nil {
		return backup.ReportError(err, true)
	}

	if err := backupContainerImages(images, outputDirectory, dryRun); err != nil {
		return backup.ReportError(err, true)
	}

	// systemd configuration backup is optional as we have defaults to use
	backupSystemdServices(outputDirectory, dryRun)

	// podman configuration backup is optional as we have defaults to use
	backupPodmanConfiguration(outputDirectory, dryRun)

	// start service if it was stopped before
	if serviceStopped && !flags.NoRestart && !dryRun {
		log.Info().Msg(L("Restarting server service"))
		if err := systemd.StartService(podman.ServerSalineService); err != nil {
			return backup.ReportError(err, true)
		}
	}

	log.Info().Msgf(L("Backup finished into %s"), outputDirectory)
	return nil
}

func printIntro(outputDir string, flags *backup.Flagpole) {
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
	for _, image := range images {
		log.Debug().Msgf("Backing up image %s", image)
		if err := podman.ExportImage(image, outputDirectory, dryRun); err != nil {
			return err
		}
	}
	return nil
}

func backupSystemdServices(outputDirectory string, dryRun bool) {
	log.Info().Msg(L("Backing up Systemd services"))
	if err := exportSystemdConfiguration(outputDirectory, dryRun); err != nil {
		log.Warn().Err(err).Msg("Systemd services and configuration was not backed up")
	}
}

func backupPodmanConfiguration(outputDirectory string, dryRun bool) {
	log.Info().Msg(L("Backing up podman configuration"))
	if err := exportPodmanConfiguration(outputDirectory, dryRun); err != nil {
		log.Warn().Err(err).Msg("Podman configuration was not backed up")
	}
}
