// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

const rpmImageDir = "/usr/share/suse-docker-images/native/"

// PrepareImage ensures the container image is pulled or pull it if the pull policy allows it.
//
// Returns the image name to use. Note that it may be changed if the image has been loaded from a local RPM package.
func PrepareImage(authFile string, image string, pullPolicy string, pullEnabled bool) (string, error) {
	if strings.ToLower(pullPolicy) != "always" {
		log.Info().Msgf(L("Ensure image %s is available"), image)

		presentImage, err := IsImagePresent(image)
		if err != nil {
			return image, err
		}

		if len(presentImage) > 0 {
			log.Debug().Msgf("Image %s already present", presentImage)
			return presentImage, nil
		} else {
			log.Debug().Msgf("Image %s is missing", image)
		}
	} else {
		log.Info().Msgf(
			L("Pull Policy is always. Presence of RPM image will be checked and pulled from registry if not present"),
		)
	}

	rpmImageFile := GetRpmImagePath(image)

	if len(rpmImageFile) > 0 {
		log.Debug().Msgf("Image %s present as RPM. Loading it", image)
		loadedImage, err := loadRpmImage(rpmImageFile)
		if err != nil {
			log.Warn().Err(err).Msgf(L("Cannot use RPM image for %s"), image)
		} else {
			log.Info().Msgf(L("Using the %[1]s image loaded from the RPM instead of its online version %[2]s"),
				strings.TrimSpace(loadedImage), image)
			return loadedImage, nil
		}
	} else {
		log.Info().Msgf(L("Cannot find RPM image for %s"), image)
	}

	if strings.ToLower(pullPolicy) != "never" {
		if pullEnabled {
			log.Debug().Msgf("Pulling image %s because it is missing and pull policy is not 'never'", image)
			return image, pullImage(authFile, image)
		} else {
			log.Debug().Msgf("Do not pulling image %s, although the pull policy is not 'never', maybe replicas is zero?", image)
			return image, nil
		}
	}

	return image, fmt.Errorf(L("image %s is missing and cannot be fetched"), image)
}

// GetRpmImageName return the RPM Image name and the tag, given an image.
func GetRpmImageName(image string) (rpmImageFile string, tag string) {
	pattern := regexp.MustCompile(`^https?://|^docker://|^oci://`)
	if pattern.FindStringIndex(image) == nil {
		image = "docker://" + image
	}
	url, err := url.Parse(image)
	if err != nil {
		log.Warn().Msgf(L("Cannot correctly parse image name '%s', local image cannot be used"), image)
		return "", ""
	}
	rpmImageFile = strings.TrimPrefix(url.Path, "/")
	rpmImageFile = strings.ReplaceAll(rpmImageFile, "/", "-")
	parts := strings.Split(rpmImageFile, ":")
	tag = "latest"
	if len(parts) > 1 {
		tag = parts[1]
	}
	rpmImageFile = parts[0]
	return rpmImageFile, tag
}

// BuildRpmImagePath checks the image metadata and returns the RPM Image path.
func BuildRpmImagePath(byteValue []byte, rpmImageFile string, tag string) (string, error) {
	var data types.Metadata
	if err := json.Unmarshal(byteValue, &data); err != nil {
		return "", utils.Errorf(err, L("cannot unmarshal image RPM metadata"))
	}
	fullPathFile := rpmImageDir + data.Image.File
	if data.Image.Name == rpmImageFile {
		for _, metadataTag := range data.Image.Tags {
			if metadataTag == tag {
				return fullPathFile, nil
			}
		}
	}
	return "", nil
}

// GetRpmImagePath return the RPM image path.
func GetRpmImagePath(image string) string {
	log.Debug().Msgf("Looking for installed RPM package containing %s image", image)

	rpmImageFile, tag := GetRpmImageName(image)

	files, err := os.ReadDir(rpmImageDir)
	if err != nil {
		log.Debug().Err(err).Msgf("Cannot read directory %s", rpmImageDir)
		return ""
	}

	for _, file := range files {
		if !strings.HasSuffix(file.Name(), "metadata") {
			continue
		}
		fullPathFileName := path.Join(rpmImageDir, file.Name())
		log.Debug().Msgf("Parsing metadata file %s", fullPathFileName)
		fileHandler, err := os.Open(fullPathFileName)
		if err != nil {
			log.Debug().Err(err).Msgf("Error opening metadata file %s", fullPathFileName)
			continue
		}
		defer fileHandler.Close()
		byteValue, err := io.ReadAll(fileHandler)
		if err != nil {
			log.Debug().Err(err).Msgf("Error reading metadata file %s", fullPathFileName)
			continue
		}

		fullPathFile, err := BuildRpmImagePath(byteValue, rpmImageFile, tag)
		if err != nil {
			log.Warn().Err(err).Msgf(L("Cannot unmarshal metadata file %s"), fullPathFileName)
			return ""
		}
		if len(fullPathFile) > 0 {
			log.Debug().Msgf("%s match with %s", fullPathFileName, image)
			return fullPathFile
		}
		log.Debug().Msgf("%s does not match with %s", fullPathFileName, image)
	}
	log.Debug().Msgf("No installed RPM package containing %s image", image)
	return ""
}

func loadRpmImage(rpmImageBasePath string) (string, error) {
	out, err := utils.RunCmdOutput(zerolog.DebugLevel, "podman", "load", "--quiet", "--input", rpmImageBasePath)
	if err != nil {
		return "", err
	}
	parseOutput := strings.SplitN(string(out), ":", 2)
	if len(parseOutput) == 2 {
		return strings.TrimSpace(parseOutput[1]), nil
	}
	return "", fmt.Errorf(L("error parsing: %s"), string(out))
}

// IsImagePresent return true if the image is present.
func IsImagePresent(image string) (string, error) {
	log.Debug().Msgf("Checking for %s", image)
	out, err := utils.RunCmdOutput(zerolog.DebugLevel, "podman", "images", "--format={{ .Repository }}", image)
	if err != nil {
		return "", fmt.Errorf(L("failed to check if image %s has already been pulled"), image)
	}

	if len(bytes.TrimSpace(out)) > 0 {
		return image, nil
	}

	splitImage := strings.SplitN(string(image), "/", 2)
	if len(splitImage) < 2 {
		return "", nil
	}
	log.Debug().Msgf("Checking for local image of %s", image)
	out, err = utils.RunCmdOutput(zerolog.DebugLevel, "podman", "images", "--quiet", "localhost/"+splitImage[1])
	if err != nil {
		return "", fmt.Errorf(L("failed to check if image %s has already been pulled"), image)
	}
	if len(bytes.TrimSpace(out)) > 0 {
		return "localhost/" + splitImage[1], nil
	}

	return "", nil
}

// GetPulledImageName returns the fullname of a pulled image.
func GetPulledImageName(image string) (string, error) {
	parts := strings.Split(image, "/")
	imageWithTag := parts[len(parts)-1]
	out, err := utils.RunCmdOutput(zerolog.DebugLevel, "podman", "images", imageWithTag, "--format", "{{.Repository}}")
	if err != nil {
		return "", fmt.Errorf(L("failed to check if image %s has already been pulled"), parts[len(parts)-1])
	}
	return string(bytes.TrimSpace(out)), nil
}

func pullImage(authFile string, image string) error {
	if utils.ContainsUpperCase(image) {
		return fmt.Errorf(L("%s should contains just lower case character, otherwise podman pull would fails"), image)
	}
	log.Info().Msgf(L("Running podman pull %s"), image)
	podmanArgs := []string{"pull", image}

	if authFile != "" {
		podmanArgs = append(podmanArgs, "--authfile", authFile)
	}

	return utils.RunCmdStdMapping(zerolog.DebugLevel, "podman", podmanArgs...)
}

// ShowAvailableTag  returns the list of available tag for a given image.
func ShowAvailableTag(registry string, image types.ImageFlags) error {
	log.Info().Msgf(L("Running podman image search --list-tags %s --format={{.Tag}}"), image.Name)

	name, err := utils.ComputeImage(registry, utils.DefaultTag, image)
	if err != nil {
		return err
	}

	if err := utils.RunCmdStdMapping(
		zerolog.DebugLevel, "podman", "image", "search", "--list-tags", name, "--format={{.Tag}}",
	); err != nil {
		return utils.Errorf(err, L("cannot find any tag for image %s"), image)
	}

	return nil
}

// GetRunningImage given a container name, return the image name.
func GetRunningImage(container string) (string, error) {
	log.Info().Msgf(L("Running podman ps --filter=name=%s --format={{ .Image }}"), container)

	out, err := utils.RunCmdOutput(
		zerolog.DebugLevel, "podman", "ps", fmt.Sprintf("--filter=name=%s", container), "--format='{{ .Image }}'",
	)
	if err != nil {
		return "", utils.Errorf(err, L("cannot find any running image for container %s"), container)
	}

	image := strings.TrimSpace(string(out))
	return image, nil
}
