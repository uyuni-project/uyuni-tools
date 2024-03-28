// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package shared

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/uyuni-project/uyuni-tools/shared"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
)

// CompareVersion compare the server image version and the server deployed  version.
func CompareVersion(imageVersion string, deployedVersion string) int {
	re := regexp.MustCompile(`\((.*?)\)`)
	imageVersionCleaned := strings.ReplaceAll(imageVersion, ".", "")
	imageVersionCleaned = strings.TrimSpace(imageVersionCleaned)
	imageVersionCleaned = re.ReplaceAllString(imageVersionCleaned, "")
	imageVersionInt, _ := strconv.Atoi(imageVersionCleaned)

	deployedVersionCleaned := strings.ReplaceAll(deployedVersion, ".", "")
	deployedVersionCleaned = strings.TrimSpace(deployedVersionCleaned)
	deployedVersionCleaned = re.ReplaceAllString(deployedVersionCleaned, "")
	deployedVersionInt, _ := strconv.Atoi(deployedVersionCleaned)
	return imageVersionInt - deployedVersionInt
}

func isUyuni(cnx *shared.Connection) (bool, error) {
	cnx_args := []string{"/etc/uyuni-release"}
	_, err := cnx.Exec("cat", cnx_args...)
	if err != nil {
		cnx_args := []string{"/etc/susemanager-release"}
		_, err := cnx.Exec("cat", cnx_args...)
		if err != nil {
			return false, errors.New(L("cannot find neither /etc/uyuni-release nor /etc/susemanagere-release"))
		}
		return false, nil
	}
	return true, nil
}

// SanityCheck verifies if an upgrade can be run.
func SanityCheck(cnx *shared.Connection, inspectedValues map[string]string, serverImage string) error {
	isUyuni, err := isUyuni(cnx)
	if err != nil {
		return fmt.Errorf(L("cannot check server release: %s"), err)
	}
	_, isCurrentUyuni := inspectedValues["uyuni_release"]
	_, isCurrentSuma := inspectedValues["suse_manager_release"]

	if isUyuni && isCurrentSuma {
		return fmt.Errorf(L("currently SUSE Manager %s is installed, instead the image is Uyuni. Upgrade is not supported"), inspectedValues["suse_manager_release"])
	}

	if !isUyuni && isCurrentUyuni {
		return fmt.Errorf(L("currently Uyuni %s is installed, instead the image is SUSE Manager. Upgrade is not supported"), inspectedValues["uyuni_release"])
	}

	if isUyuni {
		cnx_args := []string{"s/Uyuni release //g", "/etc/uyuni-release"}
		current_uyuni_release, err := cnx.Exec("sed", cnx_args...)
		if err != nil {
			return fmt.Errorf(L("failed to read current uyuni release: %s"), err)
		}
		log.Debug().Msgf("Current release is %s", string(current_uyuni_release))
		if (len(inspectedValues["uyuni_release"])) <= 0 {
			return fmt.Errorf(L("cannot fetch release from image %s"), serverImage)
		}
		log.Debug().Msgf("Image %s is %s", serverImage, inspectedValues["uyuni_release"])
		if CompareVersion(inspectedValues["uyuni_release"], string(current_uyuni_release)) < 0 {
			return fmt.Errorf(L("cannot downgrade from version %s to %s"), string(current_uyuni_release), inspectedValues["uyuni_release"])
		}
	} else {
		cnx_args := []string{"s/SUSE Manager release //g", "/etc/susemanager-release"}
		current_suse_manager_release, err := cnx.Exec("sed", cnx_args...)
		if err != nil {
			return fmt.Errorf(L("failed to read current susemanager release: %s"), err)
		}
		log.Debug().Msgf("Current release is %s", string(current_suse_manager_release))
		if (len(inspectedValues["suse_manager_release"])) <= 0 {
			return fmt.Errorf(L("cannot fetch release from image %s"), serverImage)
		}
		log.Debug().Msgf("Image %s is %s", serverImage, inspectedValues["suse_manager_release"])
		if CompareVersion(inspectedValues["suse_manager_release"], string(current_suse_manager_release)) < 0 {
			return fmt.Errorf(L("cannot downgrade from version %s to %s"), string(current_suse_manager_release), inspectedValues["suse_manager_release"])
		}
	}

	if (len(inspectedValues["image_pg_version"])) <= 0 {
		return fmt.Errorf(L("cannot fetch postgresql version from %s"), serverImage)
	}
	log.Debug().Msgf("Image %s has PostgreSQL %s", serverImage, inspectedValues["image_pg_version"])
	if (len(inspectedValues["current_pg_version"])) <= 0 {
		return fmt.Errorf(L("posgresql is not installed in the current deployment"))
	}
	log.Debug().Msgf("Current deployment has PostgreSQL %s", inspectedValues["current_pg_version"])

	return nil
}
