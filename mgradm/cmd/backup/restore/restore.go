// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package restore

import (
	"errors"
	"fmt"
	"os"
	"path"
	"slices"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/uyuni-project/uyuni-tools/mgradm/cmd/backup/shared"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

var runCmdInput = utils.RunCmdInput
var runCmd = utils.RunCmd

func Restore(
	_ *types.GlobalFlags,
	flags *shared.Flagpole,
	_ *cobra.Command,
	args []string,

) error {
	inputDirectory := args[0]
	printIntro(inputDirectory, flags)
	dryRun := flags.DryRun
	// SanityCheck
	if err := sanityChecks(inputDirectory, flags); err != nil {
		return shared.AbortError(err, false)
	}

	// Gather the list of volumes and images from the backup location
	// Both parses provided flags and the produced list has volumes or images
	// already skipped over if needed.
	volumes, err := gatherVolumesToRestore(inputDirectory, flags)
	if err != nil {
		return shared.AbortError(err, false)
	}
	images, err := gatherImagesToRestore(inputDirectory, flags)
	if err != nil {
		return shared.AbortError(err, false)
	}

	// Restore provided volumes
	// An error with volume restore is considered serious so we abort
	// --continue can be used to skip over already imported images once error
	// is resolved
	if err := restoreVolumes(volumes, flags, dryRun); err != nil {
		return shared.AbortError(err, true)
	}

	// Everything below is not considered a serious error as it can be recreated from
	// defaults, but there may be a data loss
	var hasError error
	if err := restoreImages(images, dryRun); err != nil {
		hasError = err
	}

	// Restore podman config or generate defaults
	if err := restorePodmanConfig(inputDirectory, flags); err != nil {
		hasError = errors.Join(hasError, err)
	}
	// Restore systemd config or generate defaults
	if err := restoreSystemdConfig(inputDirectory, flags); err != nil {
		hasError = errors.Join(hasError, err)
		// TODO: recreate services defaults
	}

	return shared.ReportError(hasError)
}

func printIntro(dir string, flags *shared.Flagpole) {
	log.Debug().Msg("Restoring backup with options:")
	log.Debug().Msgf("input directory: %s", dir)
	log.Debug().Msgf("dry run: %t", flags.DryRun)
	log.Debug().Msgf("skip database: %t", flags.SkipDatabase)
	log.Debug().Msgf("skip config: %t", flags.SkipConfig)
	log.Debug().Msgf("skip restart: %t", flags.NoRestart)
	log.Debug().Msgf("skip images: %t", flags.SkipImages)
	log.Debug().Msgf("skip volumes: %s", flags.SkipVolumes)
	log.Debug().Msgf("extra volumes: %s", flags.ExtraVolumes)
	log.Debug().Msgf("skip existing: %t", flags.SkipExisting)
}

func sanityChecks(inputDirectory string, flags *shared.Flagpole) error {
	if err := shared.SanityChecks(); err != nil {
		return err
	}

	if !utils.FileExists(inputDirectory) {
		return fmt.Errorf(L("input directory %s does not exists"), inputDirectory)
	}

	hostData, err := podman.InspectHost()
	if err != nil {
		return err
	}

	if hostData.HasUyuniServer {
		if flags.ForceRestore {
			log.Warn().Msg(L("Restoring over already initialized server"))
		} else {
			return errors.New(L("server is already initialized. Use force to overwrite"))
		}
	}
	return nil
}

// gatherVolumesToRestore produces a list of volumes to be imported.
// It takes a list from the backup source, checks if volume already exists and if it is
// to be skipped.
// Special `--skipvolume all` handing will cause to return empty list.
func gatherVolumesToRestore(source string, flags *shared.Flagpole) ([]string, error) {
	skipVolumes := flags.SkipVolumes
	if len(skipVolumes) == 1 && skipVolumes[0] == "all" {
		log.Debug().Msg("Skipping restoring of volumes")
		return []string{}, nil
	}

	volumeDir := path.Join(source, "volumes")
	if !utils.FileExists(volumeDir) {
		return []string{}, errors.New(L("No volumes found in the backup"))
	}

	volumes, err := os.ReadDir(volumeDir)
	if err != nil {
		return nil, errors.New(L("Unable to read directory with the volumes"))
	}

	output := []string{}
	for _, v := range volumes {
		if strings.HasSuffix(v.Name(), "sha256sum") {
			// This is checksum file, ignore
			continue
		}
		volName, _ := strings.CutSuffix(v.Name(), ".tar")

		// Skip volumes set as skipvolume option
		if slices.Contains(skipVolumes, volName) {
			log.Info().Msgf(L("Skipping volume %s"), volName)
			continue
		}

		// Skip database volumes if skipdatabase option is used
		if flags.SkipDatabase {
			for _, v := range utils.PgsqlRequiredVolumeMounts {
				if volName == v.Name {
					log.Info().Msgf(L("Skipping database volume %s"), volName)
					continue
				}
			}
		}
		if podman.IsVolumePresent(volName) {
			if flags.SkipExisting {
				log.Info().Msgf(L("Not restoring existing volume %s"), volName)
				continue
			}
			if !flags.ForceRestore {
				return nil, fmt.Errorf(L("Not restoring existing volume %s unless forced"), volName)
			}
			log.Info().Msgf(L("Volume %s will be overwriten"), volName)
		}
		output = append(output, path.Join(volumeDir, v.Name()))
	}
	return output, nil
}

// gatherImagesTorRestore produces a list of images to be imported.
// It checks if images are to be skipped, in which case it returns empty list.
func gatherImagesToRestore(source string, flags *shared.Flagpole) ([]string, error) {
	if flags.SkipImages {
		log.Debug().Msg("Skipping restoring of images")
		return []string{}, nil
	}

	imagesDir := path.Join(source, "images")
	if !utils.FileExists(imagesDir) {
		return []string{}, errors.New(L("No images found in the backup"))
	}
	images, err := os.ReadDir(imagesDir)
	if err != nil {
		return nil, errors.New(L("Unable to read directory with the images"))
	}

	output := []string{}
	for _, image := range images {
		if !strings.HasSuffix(image.Name(), ".tar") {
			continue
		}
		output = append(output, path.Join(imagesDir, image.Name()))
	}
	return output, nil
}

func restoreVolumes(volumes []string, flags *shared.Flagpole, dryRun bool) error {
	var hasError error
	for _, volume := range volumes {
		volName, _ := strings.CutSuffix(volume, ".tar")
		_, volName = path.Split(volName)
		if err := podman.ImportVolume(volName, volume, flags.SkipVerify, dryRun); err != nil {
			hasError = errors.Join(hasError, handleVolumeHacks(volName, err))
		}
	}
	return hasError
}

func restoreImages(images []string, dryRun bool) error {
	var hasErrors error
	for _, image := range images {
		if err := podman.RestoreImage(image, dryRun); err != nil {
			hasErrors = errors.Join(hasErrors, err)
		}
	}
	return hasErrors
}

func restorePodmanConfig(inputDirectory string, flags *shared.Flagpole) error {
	podmanConfigFile := path.Join(inputDirectory, shared.PodmanConfBackupFile)
	if !utils.FileExists(podmanConfigFile) {
		log.Warn().Msg(L("podman config backup not found in the backup location, trying defaults"))
		return defaultPodmanNetwork(flags)
	}

	if !flags.SkipVerify {
		if err := utils.ValidateChecksum(podmanConfigFile); err != nil {
			return errors.Join(err, errors.New(L("Unable to validate podman backup file")))
		}
	}

	return restorePodmanConfiguration(podmanConfigFile, flags)
}

func restoreSystemdConfig(inputDirectory string, flags *shared.Flagpole) error {
	log.Info().Msgf(L("Restoring systemd configuration"))
	systemdConfigFile := path.Join(inputDirectory, shared.SystemdConfBackupFile)
	if !utils.FileExists(systemdConfigFile) {
		return errors.New(L("systemd backup not found in the backup location"))
	}
	if !flags.SkipVerify {
		if err := utils.ValidateChecksum(systemdConfigFile); err != nil {
			return errors.Join(err, errors.New(L("Unable to validate systemd backup file")))
		}
	}

	return restoreSystemdConfiguration(systemdConfigFile, flags)
}
