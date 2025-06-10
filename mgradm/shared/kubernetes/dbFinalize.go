// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

//go:build !nok8s

package kubernetes

import (
	"github.com/rs/zerolog/log"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"

	"github.com/uyuni-project/uyuni-tools/mgradm/shared/templates"
	"github.com/uyuni-project/uyuni-tools/shared/kubernetes"
	batch "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// DBFinalizeJobName is the name of the Database finalization job.
const DBFinalizeJobName = "uyuni-db-finalize"

// StartDBFinalizeJob starts the database finalization job.
func StartDBFinalizeJob(
	namespace string,
	serverImage string,
	pullPolicy string,
	pullSecret string,
	schemaUpdateRequired bool,
	migration bool,
) (string, error) {
	log.Info().Msg(L("Running database finalization, this could be long depending on the size of the database…"))
	job, err := getDBFinalizeJob(namespace, serverImage, pullPolicy, pullSecret, schemaUpdateRequired, migration)
	if err != nil {
		return "", err
	}

	return job.Name, kubernetes.Apply([]runtime.Object{job}, L("failed to run the database finalization job"))
}

func getDBFinalizeJob(
	namespace string,
	image string,
	pullPolicy string,
	pullSecret string,
	schemaUpdateRequired bool,
	migration bool,
) (*batch.Job, error) {
	mounts := []types.VolumeMount{
		{MountPath: "/var/lib/pgsql", Name: "var-pgsql"},
		{MountPath: "/etc/rhn", Name: "etc-rhn"},
	}

	// Prepare the script
	scriptData := templates.FinalizePostgresTemplateData{
		RunReindex:      migration,
		RunSchemaUpdate: schemaUpdateRequired,
		Migration:       migration,
		Kubernetes:      true,
	}

	return kubernetes.GetScriptJob(namespace, DBFinalizeJobName, image, pullPolicy, pullSecret, mounts, scriptData)
}
