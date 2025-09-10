// SPDX-FileCopyrightText: 2025 SUSE LLC
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

// DBUpgradeJobName is the name of the database upgrade job.
const DBUpgradeJobName = "uyuni-db-upgrade"

// StartDBUpgradeJob starts the database upgrade job.
func StartDBUpgradeJob(
	namespace string,
	image types.ImageFlags,
	migrationImage types.ImageFlags,
	pullSecret string,
	oldPgsql string,
	newPgsql string,
) (string, error) {
	log.Info().Msgf(L("Upgrading PostgreSQL database from %[1]s to %[2]s…"), oldPgsql, newPgsql)

	var migrationImageURL string
	var err error
	if migrationImage.Name == "" {
		imageName := fmt.Sprintf("-migration-%s-%s", oldPgsql, newPgsql)
		migrationImageURL, err = utils.ComputeImage(image.Registry.Host, image.Tag, image, imageName)
	} else {
		migrationImageURL, err = utils.ComputeImage(image.Registry.Host, image.Tag, migrationImage)
	}
	if err != nil {
		return "", utils.Error(err, L("failed to compute image URL"))
	}

	log.Info().Msgf(L("Using database upgrade image %s"), migrationImageURL)

	job, err := getDBUpgradeJob(namespace, migrationImageURL, image.PullPolicy, pullSecret, oldPgsql, newPgsql)
	if err != nil {
		return "", err
	}

	return job.ObjectMeta.Name, kubernetes.Apply([]runtime.Object{job}, L("failed to run the database upgrade job"))
}

func getDBUpgradeJob(
	namespace string,
	image string,
	pullPolicy string,
	pullSecret string,
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

	return kubernetes.GetScriptJob(namespace, DBUpgradeJobName, image, pullPolicy, pullSecret, mounts, scriptData)
}
