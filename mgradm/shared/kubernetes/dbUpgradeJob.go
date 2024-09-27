// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

//go:build !nok8s

package kubernetes

import (
	"fmt"

	"github.com/rs/zerolog/log"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"

	"github.com/uyuni-project/uyuni-tools/mgradm/shared/templates"
	"github.com/uyuni-project/uyuni-tools/shared/kubernetes"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
	batch "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// DbUpgradeJobName is the name of the database upgrade job.
const DbUpgradeJobName = "uyuni-db-upgrade"

// StartDbUpgradeJob starts the database upgrade job.
func StartDbUpgradeJob(
	namespace string,
	registry string,
	image types.ImageFlags,
	migrationImage types.ImageFlags,
	oldPgsql string,
	newPgsql string,
) error {
	log.Info().Msgf(L("Upgrading PostgreSQL database from %[1]s to %[2]s…"), oldPgsql, newPgsql)

	var migrationImageUrl string
	var err error
	if migrationImage.Name == "" {
		imageName := fmt.Sprintf("-migration-%s-%s", oldPgsql, newPgsql)
		migrationImageUrl, err = utils.ComputeImage(registry, image.Tag, image, imageName)
	} else {
		migrationImageUrl, err = utils.ComputeImage(registry, image.Tag, migrationImage)
	}
	if err != nil {
		return utils.Errorf(err, L("failed to compute image URL"))
	}

	log.Info().Msgf(L("Using database upgrade image %s"), migrationImageUrl)

	job, err := getDbUpgradeJob(namespace, migrationImageUrl, image.PullPolicy, oldPgsql, newPgsql)
	if err != nil {
		return err
	}

	return kubernetes.Apply([]runtime.Object{job}, L("failed to run the database upgrade job"))
}

func getDbUpgradeJob(
	namespace string,
	image string,
	pullPolicy string,
	oldPgsql string,
	newPgsql string,
) (*batch.Job, error) {
	mounts := []types.VolumeMount{
		{MountPath: "/var/lib/pgsql", Name: "var-pgsql"},
	}

	// Prepare the script
	scriptData := templates.PostgreSQLVersionUpgradeTemplateData{
		OldVersion: oldPgsql,
		NewVersion: newPgsql,
	}

	return kubernetes.GetScriptJob(namespace, DbUpgradeJobName, image, pullPolicy, mounts, scriptData)
}
