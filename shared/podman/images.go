// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// Ensure the container image is pulled or pull it if the pull policy allows it.
func PrepareImage(image string, pullPolicy string, args ...string) error {
	log.Info().Msgf("Ensure image %s is available", image)

	needsPull, err := checkImage(image, pullPolicy)
	if err != nil {
		return err
	}

	if needsPull {
		return pullImage(image, args...)
	}
	return nil
}

func calculateRpmImagePath(image string) string {
	imagePrefix := strings.ReplaceAll(image, "registry.suse.com/", "")
	imagePrefix = strings.ReplaceAll(imagePrefix, "/", "-")
	imagePrefix = strings.Split(imagePrefix, ":")[0]

	return "/usr/share/suse-docker-images/native/" + imagePrefix
}

func loadRpmImage(rpmImageBasePath string) error {
	if err := utils.RunCmdStdMapping(zerolog.DebugLevel, "podman", "load", "--input", rpmImageBasePath); err != nil {
		return fmt.Errorf("cannot load image from: %s, continuing trying to pull from the registry: %s", rpmImageBasePath, err)
	}
	return nil
}

func isRpmImagePresent(image string) bool {
	//check if image is available from RPM
	log.Debug().Msgf("Check if RPM based image for %s is present", image)

	rpmImageBasePath := calculateRpmImagePath(image)
	log.Debug().Msgf("Looking for %s", rpmImageBasePath)

	if _, err := os.Stat(rpmImageBasePath); err == nil {
		return true
	}

	return false
}

func checkImage(image string, pullPolicy string) (bool, error) {
	if strings.ToLower(pullPolicy) == "always" {
		return true, nil
	}

	out, err := utils.RunCmdOutput(zerolog.DebugLevel, "podman", "images", "--quiet", image)
	if err != nil {
		return false, fmt.Errorf("failed to check if image %s has already been pulled", image)
	}

	if isRpmImagePresent(image) {
		rpmImageBasePath := calculateRpmImagePath(image)
		if err := loadRpmImage(rpmImageBasePath); err != nil {
			log.Warn().Msgf("cannot use RPM image for %s:%s", image, err)
		}
		return false, nil
	}

	if len(bytes.TrimSpace(out)) == 0 {
		if pullPolicy == "Never" {
			return false, fmt.Errorf("image %s is not available and cannot be pulled due to policy", image)
		}
		return true, nil
	}
	return false, nil
}

func pullImage(image string, args ...string) error {
	log.Info().Msgf("Running podman pull %s", image)
	podmanImageArgs := []string{"pull", image}
	podmanArgs := append(podmanImageArgs, args...)

	loglevel := zerolog.DebugLevel
	if len(args) > 0 {
		loglevel = zerolog.Disabled
		log.Debug().Msg("Additional arguments for pull command will not be shown.")
	}

	return utils.RunCmdStdMapping(loglevel, "podman", podmanArgs...)
}

// ShowAvailableTag  returns the list of avaialable tag for a given image.
func ShowAvailableTag(image string) ([]string, error) {
	log.Info().Msgf("Running podman image search --list-tags %s --format='{{.Tag}}'", image)

	out, err := utils.RunCmdOutput(zerolog.DebugLevel, "podman", "image", "search", "--list-tags", image, "--format='{{.Tag}}'")
	if err != nil {
		return []string{}, fmt.Errorf("cannot find any tag for image %s: %s", image, err)
	}

	tags := strings.Split(string(out), "\n")
	return tags, nil
}
