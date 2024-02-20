// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// Ensure the container image is pulled or pull it if the pull policy allows it.
func PrepareImage(image string, pullPolicy string) error {
	log.Info().Msgf("Ensure image %s is available", image)

	needsPull, err := checkImage(image, pullPolicy)
	if err != nil {
		return err
	}

	if needsPull {
		return pullImage(image)
	}
	return nil
}

func checkImage(image string, pullPolicy string) (bool, error) {
	if strings.ToLower(pullPolicy) == "always" {
		return true, nil
	}

	out, err := utils.RunCmdOutput(zerolog.DebugLevel, "podman", "images", "--quiet", image)
	if err != nil {
		return false, fmt.Errorf("failed to check if image %s has already been pulled", image)
	}

	if len(bytes.TrimSpace(out)) == 0 {
		if pullPolicy == "Never" {
			return false, fmt.Errorf("image %s is not available and cannot be pulled due to policy", image)
		}
		return true, nil
	}
	return false, nil
}

func pullImage(image string) error {
	log.Info().Msgf("Running podman pull %s", image)

	return utils.RunCmdStdMapping("podman", "pull", image)
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
