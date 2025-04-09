//SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

//go:build !nok8s

package kubernetes

import (
	"time"

	adm_utils "github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	"github.com/uyuni-project/uyuni-tools/shared/kubernetes"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/ssl"
	batch "k8s.io/api/batch/v1"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const SetupJobName = "uyuni-setup"

// StartSetupJob creates the job setting up the server.
func StartSetupJob(
	namespace string,
	image string,
	pullPolicy core.PullPolicy,
	pullSecret string,
	mirrorPvName string,
	flags *adm_utils.InstallationFlags,
	fqdn string,
	adminSecret string,
	dbSecret string,
	reportdbSecret string,
	sccSecret string,
	tz string,
) (string, error) {
	job, err := GetSetupJob(
		namespace, image, pullPolicy, pullSecret, mirrorPvName, flags, fqdn,
		adminSecret, dbSecret, reportdbSecret, sccSecret, tz,
	)
	if err != nil {
		return "", err
	}
	return job.ObjectMeta.Name, kubernetes.Apply([]*batch.Job{job}, L("failed to run the setup job"))
}

// GetSetupJob creates the job definition object for the setup.
func GetSetupJob(
	namespace string,
	image string,
	pullPolicy core.PullPolicy,
	pullSecret string,
	mirrorPvName string,
	flags *adm_utils.InstallationFlags,
	fqdn string,
	adminSecret string,
	dbSecret string,
	reportdbSecret string,
	sccSecret string,
	tz string,
) (*batch.Job, error) {
	var maxFailures int32
	timestamp := time.Now().Format("20060102150405")

	template := getServerPodTemplate(image, pullPolicy, tz, pullSecret)

	script, err := adm_utils.GenerateSetupScript(flags, true)
	if err != nil {
		return nil, err
	}

	template.Spec.Containers[0].Name = "setup"
	template.Spec.Containers[0].Command = []string{"sh", "-c", script}
	template.Spec.RestartPolicy = core.RestartPolicyNever

	optional := false

	dbUserEnv := core.EnvVar{Name: "MANAGER_USER", ValueFrom: &core.EnvVarSource{
		SecretKeyRef: &core.SecretKeySelector{
			LocalObjectReference: core.LocalObjectReference{Name: dbSecret},
			Key:                  "username",
			Optional:             &optional,
		},
	}}
	reportdbUserEnv := core.EnvVar{Name: "REPORT_DB_USER", ValueFrom: &core.EnvVarSource{
		SecretKeyRef: &core.SecretKeySelector{
			LocalObjectReference: core.LocalObjectReference{Name: reportdbSecret},
			Key:                  "username",
			Optional:             &optional,
		},
	}}

	envVars := []core.EnvVar{
		{Name: "ADMIN_USER", ValueFrom: &core.EnvVarSource{
			SecretKeyRef: &core.SecretKeySelector{
				LocalObjectReference: core.LocalObjectReference{Name: adminSecret},
				Key:                  "username",
				Optional:             &optional,
			},
		}},
		{Name: "ADMIN_PASS", ValueFrom: &core.EnvVarSource{
			SecretKeyRef: &core.SecretKeySelector{
				LocalObjectReference: core.LocalObjectReference{Name: adminSecret},
				Key:                  "password",
				Optional:             &optional,
			},
		}},
		dbUserEnv,
		{Name: "MANAGER_PASS", ValueFrom: &core.EnvVarSource{
			SecretKeyRef: &core.SecretKeySelector{
				LocalObjectReference: core.LocalObjectReference{Name: dbSecret},
				Key:                  "password",
				Optional:             &optional,
			},
		}},
		reportdbUserEnv,
		{Name: "REPORT_DB_PASS", ValueFrom: &core.EnvVarSource{
			SecretKeyRef: &core.SecretKeySelector{
				LocalObjectReference: core.LocalObjectReference{Name: reportdbSecret},
				Key:                  "password",
				Optional:             &optional,
			},
		}},
		{Name: "REPORT_DB_CA_CERT", Value: ssl.DBCAContainerPath},
		// EXTERNALDB_* variables are not passed yet: only for AWS and it probably doesn't make sense for kubernetes yet.
	}

	// The DB and ReportDB port is expected to be the standard one.
	// When using an external database with a custom port the only solution is to access it using
	// its IP address and a headless service with a custom EndpointSlice.
	// If this is too big a constraint, we'll have to accept the port as a parameter too.
	env := adm_utils.GetSetupEnv(mirrorPvName, flags, fqdn, true)
	for key, value := range env {
		envVars = append(envVars, core.EnvVar{Name: key, Value: value})
	}

	if sccSecret != "" {
		envVars = append(envVars,
			core.EnvVar{Name: "SCC_USER", ValueFrom: &core.EnvVarSource{
				SecretKeyRef: &core.SecretKeySelector{
					LocalObjectReference: core.LocalObjectReference{Name: sccSecret},
					Key:                  "username",
					Optional:             &optional,
				},
			}},
			core.EnvVar{Name: "SCC_PASS", ValueFrom: &core.EnvVarSource{
				SecretKeyRef: &core.SecretKeySelector{
					LocalObjectReference: core.LocalObjectReference{Name: sccSecret},
					Key:                  "password",
					Optional:             &optional,
				},
			}},
		)
	}

	if mirrorPvName != "" {
		envVars = append(envVars, core.EnvVar{Name: "MIRROR_PATH", Value: "/mirror"})
	}
	template.Spec.Containers[0].Env = envVars
	template.Spec.Volumes = append(template.Spec.Volumes,
		kubernetes.CreateConfigVolume("db-ca", kubernetes.DBCAConfigName),
	)

	// Add initContainer waiting for the db and reportdb services to be responding
	template.Spec.InitContainers = append(template.Spec.InitContainers,
		core.Container{
			Name:            "db-waiter",
			Image:           image,
			ImagePullPolicy: pullPolicy,
			Env: []core.EnvVar{
				{Name: "MANAGER_DB_HOST", Value: env["MANAGER_DB_HOST"]},
				{Name: "MANAGER_DB_PORT", Value: env["MANAGER_DB_PORT"]},
				{Name: "MANAGER_DB_NAME", Value: env["MANAGER_DB_NAME"]},
				dbUserEnv,
				{Name: "REPORT_DB_HOST", Value: env["REPORT_DB_HOST"]},
				{Name: "REPORT_DB_PORT", Value: env["REPORT_DB_PORT"]},
				{Name: "REPORT_DB_NAME", Value: env["REPORT_DB_NAME"]},
				reportdbUserEnv,
			},
			Command: []string{
				"sh", "-c",
				`
until pg_isready -U $MANAGER_USER -h $MANAGER_DB_HOST -p $MANAGER_DB_PORT -d $MANAGER_DB_NAME; do
	sleep 60
done

until pg_isready -U $REPORT_DB_USER -h $REPORT_DB_HOST -p $REPORT_DB_PORT -d $REPORT_DB_NAME; do
	sleep 60
done
`,
			},
		},
	)

	job := batch.Job{
		TypeMeta: meta.TypeMeta{Kind: "Job", APIVersion: "batch/v1"},
		ObjectMeta: meta.ObjectMeta{
			Name:      SetupJobName + "-" + timestamp,
			Namespace: namespace,
			Labels:    kubernetes.GetLabels(kubernetes.ServerApp, ""),
		},
		Spec: batch.JobSpec{
			Template:     template,
			BackoffLimit: &maxFailures,
		},
	}

	if pullSecret != "" {
		job.Spec.Template.Spec.ImagePullSecrets = []core.LocalObjectReference{{Name: pullSecret}}
	}

	return &job, nil
}
