// SPDX-FileCopyrightText: 2025 SUSE LLC
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
		}
		log.Debug().Msgf("Image %s is missing", image)
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
		}
		log.Debug().Msgf("Not pulling image %s, although the pull policy is not 'never', maybe replicas is zero?", image)
		return image, nil
	}

	return image, fmt.Errorf(L("image %s is missing and cannot be fetched"), image)
}

func PrepareImages(
	authFile string,
	image types.ImageFlags,
	pgsqlFlags types.PgsqlFlags,
) (string, string, error) {
	serverImage, err := utils.ComputeImage(image.Registry, utils.DefaultTag, image)
	if err != nil && len(serverImage) > 0 {
		return "", "", utils.Error(err, L("failed to determine image"))
	}

	if len(serverImage) <= 0 {
		log.Debug().Msg("Use deployed image")

		serverImage, err = GetRunningImage(ServerContainerName)
		if err != nil {
			return "", "", utils.Error(err, L("failed to find the image of the currently running server container"))
		}
	}

	pgsqlImage, err := utils.ComputeImage(image.Registry, utils.DefaultTag, pgsqlFlags.Image)
	if err != nil && len(pgsqlImage) > 0 {
		return "", "", utils.Error(err, L("failed to determine pgsql image"))
	}

	if len(pgsqlImage) <= 0 {
		log.Debug().Msg("Use deployed pgsqlimage")

		pgsqlImage, err = GetRunningImage(DBContainerName)
		if err != nil {
			return "", "", utils.Error(err, L("failed to find the image of the currently running db container"))
		}
	}

	preparedServerImage, err := PrepareImage(authFile, serverImage, image.PullPolicy, true)
	if err != nil {
		return preparedServerImage, "", err
	}

	preparedPgsqlImage, err := PrepareImage(authFile, pgsqlImage, image.PullPolicy, true)
	if err != nil {
		return preparedServerImage, preparedPgsqlImage, err
	}

	return preparedServerImage, preparedPgsqlImage, nil
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
	if !utils.FileExists(rpmImageDir) {
		log.Info().Msgf(L("skipping loading image from RPM as %s doesn't exist"), rpmImageDir)
		return ""
	}

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

// ShowAvailableTag returns the list of available tag for a given image.
func ShowAvailableTag(registry string, image types.ImageFlags, authFile string) error {
	log.Info().Msgf(L("Running podman image search --list-tags %s --format={{.Tag}}"), image.Name)

	name, err := utils.ComputeImage(registry, utils.DefaultTag, image)
	if err != nil {
		return err
	}

	args := []string{"image", "search", "--list-tags", name, "--format={{.Tag}}"}
	if authFile != "" {
		args = append(args, "--authfile", authFile)
	}

	out, err := utils.RunCmdOutput(zerolog.DebugLevel, "podman", args...)
	if err != nil {
		return utils.Errorf(err, L("cannot find any tag for image %s"), image)
	}

	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if !strings.HasSuffix(line, ".sig") && !strings.HasSuffix(line, ".att") {
			fmt.Println(line)
		}
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

// HasRemoteImage returns true if the image is available remotely.
//
// The image has to be a full image with registry, path and tag.
func HasRemoteImage(image string) bool {
	out, err := runCmdOutput(zerolog.DebugLevel,
		"podman", "search", "--list-tags", "--format", "{{.Name}}:{{.Tag}}", image,
	)
	if err != nil {
		return false
	}
	imageFinder := regexp.MustCompile("(?Um)^" + image + "$")
	return imageFinder.Match(out)
}

// DeleteImage deletes a podman image based on its name.
// If dryRun is set to true, nothing will be done, only messages logged to explain what would happen.
func DeleteImage(name string, dryRun bool) error {
	exists := imageExists(name)
	if exists {
		if dryRun {
			log.Info().Msgf(L("Would run %s"), "podman image rm "+name)
		} else {
			log.Info().Msgf(L("Run %s"), "podman image rm "+name)
			err := utils.RunCmd("podman", "image", "rm", name)
			if err != nil {
				return utils.Errorf(err, L("Failed to remove image %s"), name)
			}
		}
	}
	return nil
}

// ExportImage saves a podman image based on its name to a specified directory.
// outputDir option expects already existing directory.
// If dryRun is set to true, nothing will be done, only messages logged to explain what would happen.
func ExportImage(name string, outputDir string, dryRun bool) error {
	exists := imageExists(name)
	if exists {
		saveCommand := []string{"podman", "image", "save", "--quiet", "-o", path.Join(outputDir, name+".tar"), name}
		if dryRun {
			log.Info().Msgf(L("Would run %s"), strings.Join(saveCommand, " "))
		} else {
			log.Info().Msgf(L("Run %s"), strings.Join(saveCommand, " "))
			err := utils.RunCmd(saveCommand[0], saveCommand[1:]...)
			if err != nil {
				return utils.Errorf(err, L("Failed to export image %s"), name)
			}
		}
	}
	return nil
}

func imageExists(image string) bool {
	err := utils.RunCmd("podman", "image", "exists", image)
	return err == nil
}

func RestoreImage(imageFile string, dryRun bool) error {
	restoreCommand := []string{"podman", "image", "import", "--quiet", imageFile}
	if dryRun {
		log.Info().Msgf(L("Would run %s"), strings.Join(restoreCommand, " "))
	} else {
		log.Info().Msgf(L("Run %s"), strings.Join(restoreCommand, " "))
		err := utils.RunCmd(restoreCommand[0], restoreCommand[1:]...)
		if err != nil {
			return utils.Errorf(err, L("Failed to restore image %s"), imageFile)
		}
	}
	return nil
}
