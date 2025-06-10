// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

//go:build !nok8s

package kubernetes

import (
	"github.com/rs/zerolog/log"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"

	"github.com/uyuni-project/uyuni-tools/mgradm/shared/templates"
	"github.com/uyuni-project/uyuni-tools/shared/kubernetes"
	batch "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// PostUpgradeJobName is the name of the job apply the database changes after the upgrade.
const PostUpgradeJobName = "uyuni-post-upgrade"

// StartPostUpgradeJob starts the job applying the database changes after the upgrade.
func StartPostUpgradeJob(namespace string, image string, pullPolicy string, pullSecret string) (string, error) {
	log.Info().Msg(L("Performing post upgrade changes…"))

	job, err := getPostUpgradeJob(namespace, image, pullPolicy, pullSecret)
	if err != nil {
		return "", err
	}

	return job.Name, kubernetes.Apply([]runtime.Object{job}, L("failed to run the post upgrade job"))
}

func getPostUpgradeJob(namespace string, image string, pullPolicy string, pullSecret string) (*batch.Job, error) {
	scriptData := templates.PostUpgradeTemplateData{}
	mounts := GetServerMounts()

	return kubernetes.GetScriptJob(namespace, PostUpgradeJobName, image, pullPolicy, pullSecret, mounts, scriptData)
}
