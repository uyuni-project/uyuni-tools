// SPDX-FileCopyrightText: 2023 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package shared

import (
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/uyuni-project/uyuni-tools/shared"
)

func CompareVersion(imageVersion string, deployedVersion string) int {
	imageVersionCleaned := strings.ReplaceAll(imageVersion, ".", "")
	imageVersionCleaned = strings.TrimSpace(imageVersionCleaned)

	imageVersionInt, _ := strconv.Atoi(imageVersionCleaned)

	deployedVersionCleaned := strings.ReplaceAll(deployedVersion, ".", "")

	deployedVersionInt, _ := strconv.Atoi(deployedVersionCleaned)

	return imageVersionInt - deployedVersionInt
}

func SanityCheck(cnx *shared.Connection, inspectedValues map[string]string, serverImage string) error {
	cnx_args := []string{"s/Uyuni release //g", "/etc/uyuni-release"}
	current_uyuni_release, err := cnx.Exec("sed", cnx_args...)

	if err != nil {
		log.Fatal().Err(err).Msg("Failed to read current uyuni release")
	}

	log.Debug().Msgf("Current Uyuni release is %s", string(current_uyuni_release))

	if (len(inspectedValues["uyuni_release"])) <= 0 {
		log.Fatal().Msgf("Cannot fetch release from image %s", serverImage)
	}

	log.Debug().Msgf("Image %s is %s", serverImage, inspectedValues["uyuni_release"])
	if CompareVersion(inspectedValues["uyuni_release"], string(current_uyuni_release)) <= 0 {
		log.Fatal().Msgf("This is not an upgrade, since current Uyuni version is %s and image Uyuni version %s", string(current_uyuni_release), inspectedValues["uyuni_release"])
	}

	if (len(inspectedValues["image_pg_version"])) <= 0 {
		log.Fatal().Msgf("Cannot feth PostgreSQL version from %s", serverImage)
	}
	log.Debug().Msgf("Image %s has PostgreSQL %s", serverImage, inspectedValues["image_pg_version"])

	if (len(inspectedValues["current_pg_version"])) <= 0 {
		log.Fatal().Msgf("PosgreSQL is not installed in the current deploy")
	}
	log.Debug().Msgf("Current deployment has PostgreSQL %s", inspectedValues["current_pg_version"])

	return err
}
