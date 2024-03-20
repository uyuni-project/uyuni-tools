// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// Ensure the container image is pulled or pull it if the pull policy allows it.
func PrepareImage(image string, pullPolicy string, args ...string) (string, error) {
	log.Info().Msgf("Ensure image %s is available", image)

	present, err := IsImagePresent(image)
	if err != nil {
		return "", err
	}

	if present && strings.ToLower(pullPolicy) != "always" {
		log.Debug().Msgf("Image %s already present", image)
		return "", nil
	}

	rpmImageFile, err := IsRpmImagePresent(image)
	if err != nil {
		return "", err
	}

	if len(rpmImageFile) > 0 {
		log.Debug().Msgf("Image %s present as RPM. Loading it", image)
		loadedImage, err := loadRpmImage(rpmImageFile)
		if err != nil {
			log.Warn().Msgf("cannot use RPM image for %s:%s", image, err)
			present = false
		} else {
			if strings.ToLower(pullPolicy) == "always" {
				log.Debug().Msg("Ignoring pull policy alway ")
			}

			log.Warn().Msgf("Loading image %s: it's the RPM based image of %s.", strings.TrimSpace(loadedImage), image)
			return loadedImage, nil
		}
	}

	if strings.ToLower(pullPolicy) == "always" {
		log.Debug().Msgf("Pulling image cause pull policy is always %s", image)
		return image, pullImage(image, args...)
	}

	if !present && strings.ToLower(pullPolicy) != "never" {
		log.Debug().Msgf("Pulling image cause is missing and pull policy is not never %s", image)
		return image, pullImage(image, args...)
	}

	return image, fmt.Errorf("image %s is missing and cannot be fetch", image)
}

// GetRpmInfoFromImage return the RPM Image name and the tag, given an image.
func GetRpmInfoFromImage(image string) (rpmImageFile string, tag string) {
	rpmImageFile = strings.ReplaceAll(image, "registry.suse.com/", "")
	rpmImageFile = strings.ReplaceAll(rpmImageFile, "/", "-")
	parts := strings.Split(rpmImageFile, ":")
	tag = "latest"
	if len(parts) > 1 {
		tag = parts[1]
	}
	rpmImageFile = parts[0]
	return rpmImageFile, tag
}

// GetRpmInfoFromImage return the path of an RPM Image.
func GetRpmImage(byteValue []byte, rpmImageFile string, tag string) (string, error) {
	var data types.Metadata
	if err := json.Unmarshal(byteValue, &data); err != nil {
		return "", fmt.Errorf("cannot unmarshal: %s", err)
	}
	fullPathFile := "/usr/share/suse-docker-images/native/" + data.Image.File
	if data.Image.Name == rpmImageFile {
		for _, metadataTag := range data.Image.Tags {
			if metadataTag == tag {
				return fullPathFile, nil
			}
		}
	}
	return "", nil
}

// IsRpmImagePresent return true if the RPM with the provided image is installed.
func IsRpmImagePresent(image string) (string, error) {
	log.Debug().Msgf("looking for RPM image for %s", image)

	rpmImageFile, tag := GetRpmInfoFromImage(image)

	rpmImageDir := "/usr/share/suse-docker-images/native/"
	files, err := os.ReadDir(rpmImageDir)
	if err != nil {
		return "", fmt.Errorf("cannot read directory %s: %s", rpmImageDir, err)
	}

	for _, file := range files {
		if !strings.HasSuffix(file.Name(), "metadata") {
			continue
		}
		fullPathFileName := path.Join(rpmImageDir, file.Name())
		log.Debug().Msgf("parsing %s", fullPathFileName)
		fileHandler, err := os.Open(fullPathFileName)
		if err != nil {
			log.Debug().Msgf("error opening %s: %s", fullPathFileName, err)
			continue
		}
		defer fileHandler.Close()
		byteValue, err := io.ReadAll(fileHandler)
		if err != nil {
			log.Debug().Msgf("error reading %s: %s", fullPathFileName, err)
			continue
		}

		fullPathFile, err := GetRpmImage(byteValue, rpmImageFile, tag)
		if err != nil {
			log.Warn().Msgf("cannot unmarshal %s: %s", fullPathFileName, err)
			return "", err
		}
		if len(fullPathFile) > 0 {
			log.Debug().Msgf("%s match with %s", fullPathFileName, image)
			return fullPathFile, nil
		}
		log.Debug().Msgf("%s does not match with %s", fullPathFileName, image)
	}
	log.Debug().Msgf("image %s does not exists as RPM image", image)
	return "", fmt.Errorf("image %s does not exists as RPM image", image)
}

func loadRpmImage(rpmImageBasePath string) (string, error) {
	out, err := utils.RunCmdOutput(zerolog.DebugLevel, "podman", "load", "--quiet", "--input", rpmImageBasePath)
	if err != nil {
		return "", fmt.Errorf("cannot load image from: %s, continuing trying to pull from the registry: %s", rpmImageBasePath, err)
	}
	loadedImage := strings.ReplaceAll(string(out), "Loaded image: ", "")
	return loadedImage, nil
}

// IsImagePresent return true if the image is present.
func IsImagePresent(image string) (bool, error) {
	out, err := utils.RunCmdOutput(zerolog.DebugLevel, "podman", "images", "--quiet", image)
	if err != nil {
		return false, fmt.Errorf("failed to check if image %s has already been pulled", image)
	}

	if len(bytes.TrimSpace(out)) > 0 {
		return true, nil
	}
	return false, nil
}

// GetPulledImageName returns the fullname of a pulled image.
func GetPulledImageName(image string) (string, error) {
	parts := strings.Split(image, "/")
	imageWithTag := parts[len(parts)-1]
	out, err := utils.RunCmdOutput(zerolog.DebugLevel, "podman", "images", imageWithTag, "--format", "{{.Repository}}")
	if err != nil {
		return "", fmt.Errorf("failed to check if image %s has already been pulled", parts[len(parts)-1])
	}
	return string(bytes.TrimSpace(out)), nil
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
