// SPDX-FileCopyrightText: 2024 SUSE LLC
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

// DbFinalizeJobName is the name of the Database finalization job.
const DbFinalizeJobName = "uyuni-db-finalize"

// StartDbFinalizeJob starts the database finalization job.
func StartDbFinalizeJob(
	namespace string,
	serverImage string,
	pullPolicy string,
	schemaUpdateRequired bool,
	migration bool,
) error {
	log.Info().Msg(L("Running database finalization, this could be long depending on the size of the database…"))
	job, err := getDbFinalizeJob(namespace, serverImage, pullPolicy, schemaUpdateRequired, migration)
	if err != nil {
		return err
	}

	return kubernetes.Apply([]runtime.Object{job}, L("failed to run the database finalization job"))
}

func getDbFinalizeJob(
	namespace string,
	image string,
	pullPolicy string,
	schemaUpdateRequired bool,
	migration bool,
) (*batch.Job, error) {
	mounts := []types.VolumeMount{
		{MountPath: "/var/lib/pgsql", Name: "var-pgsql"},
		{MountPath: "/etc/rhn", Name: "etc-rhn"},
	}

	// Prepare the script
	scriptData := templates.FinalizePostgresTemplateData{
		RunAutotune:     true,
		RunReindex:      true,
		RunSchemaUpdate: schemaUpdateRequired,
		Migration:       migration,
		Kubernetes:      true,
	}

	return kubernetes.GetScriptJob(namespace, DbFinalizeJobName, image, pullPolicy, mounts, scriptData)
}
