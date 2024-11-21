//SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

//go:build !nok8s

package kubernetes

import (
	"strings"
	"time"

	"github.com/uyuni-project/uyuni-tools/mgradm/shared/templates"
	adm_utils "github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	"github.com/uyuni-project/uyuni-tools/shared/kubernetes"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
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
) (string, error) {
	job, err := GetSetupJob(
		namespace, image, pullPolicy, pullSecret, mirrorPvName, flags, fqdn,
		adminSecret, dbSecret, reportdbSecret, sccSecret,
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
) (*batch.Job, error) {
	var maxFailures int32
	timestamp := time.Now().Format("20060102150405")

	template := getServerPodTemplate(image, pullPolicy, flags.TZ, pullSecret)

	script, err := generateSetupScript(flags)
	if err != nil {
		return nil, err
	}

	template.Spec.Containers[0].Name = "setup"
	template.Spec.Containers[0].Command = []string{"sh", "-c", script}
	template.Spec.RestartPolicy = core.RestartPolicyNever

	optional := false

	localHostValues := []string{
		"localhost",
		"127.0.0.1",
		"::1",
		fqdn,
	}

	localDB := "N"
	if utils.Contains(localHostValues, flags.DB.Host) {
		localDB = "Y"
	}

	// The DB and ReportDB port is expected to be the standard one.
	// When using an external database with a custom port the only solution is to access it using
	// its IP address and a headless service with a custom EndpointSlice.
	// If this is too big a constraint, we'll have to accept the port as a parameter too.
	env := []core.EnvVar{
		{Name: "NO_SSL", Value: "Y"},
		{Name: "UYUNI_FQDN", Value: fqdn},
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
		{Name: "MANAGER_USER", ValueFrom: &core.EnvVarSource{
			SecretKeyRef: &core.SecretKeySelector{
				LocalObjectReference: core.LocalObjectReference{Name: dbSecret},
				Key:                  "username",
				Optional:             &optional,
			},
		}},
		{Name: "MANAGER_PASS", ValueFrom: &core.EnvVarSource{
			SecretKeyRef: &core.SecretKeySelector{
				LocalObjectReference: core.LocalObjectReference{Name: dbSecret},
				Key:                  "password",
				Optional:             &optional,
			},
		}},
		{Name: "MANAGER_ADMIN_EMAIL", Value: flags.Email},
		{Name: "MANAGER_MAIL_FROM", Value: flags.EmailFrom},
		{Name: "MANAGER_ENABLE_TFTP", Value: "Y"},
		{Name: "LOCAL_DB", Value: localDB},
		{Name: "MANAGER_DB_NAME", Value: flags.DB.Name},
		{Name: "MANAGER_DB_HOST", Value: flags.DB.Host},
		{Name: "MANAGER_DB_PORT", Value: "5432"},
		{Name: "MANAGER_DB_PROTOCOL", Value: "tcp"},
		{Name: "REPORT_DB_NAME", Value: flags.ReportDB.Name},
		{Name: "REPORT_DB_HOST", Value: flags.ReportDB.Host},
		{Name: "REPORT_DB_PORT", Value: "5432"},
		{Name: "REPORT_DB_USER", ValueFrom: &core.EnvVarSource{
			SecretKeyRef: &core.SecretKeySelector{
				LocalObjectReference: core.LocalObjectReference{Name: reportdbSecret},
				Key:                  "username",
				Optional:             &optional,
			},
		}},
		{Name: "REPORT_DB_PASS", ValueFrom: &core.EnvVarSource{
			SecretKeyRef: &core.SecretKeySelector{
				LocalObjectReference: core.LocalObjectReference{Name: reportdbSecret},
				Key:                  "password",
				Optional:             &optional,
			},
		}},
		{Name: "ISS_PARENT", Value: flags.IssParent},
		{Name: "ACTIVATE_SLP", Value: "N"},
		// TODO EXTERNALDB_* variables are not passed yet: only for AWS and it probably doesn't make sense for kubernetes yet.
	}

	if sccSecret != "" {
		env = append(env,
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
		env = append(env, core.EnvVar{Name: "MIRROR_PATH", Value: "/mirror"})
	}
	template.Spec.Containers[0].Env = env

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

func generateSetupScript(flags *adm_utils.InstallationFlags) (string, error) {
	template := templates.MgrSetupScriptTemplateData{
		DebugJava:      flags.Debug.Java,
		OrgName:        flags.Organization,
		AdminLogin:     "$ADMIN_USER",
		AdminPassword:  "$ADMIN_PASS",
		AdminFirstName: flags.Admin.FirstName,
		AdminLastName:  flags.Admin.LastName,
		AdminEmail:     flags.Admin.Email,
		NoSSL:          true,
	}

	// Prepare the script
	scriptBuilder := new(strings.Builder)
	if err := template.Render(scriptBuilder); err != nil {
		return "", utils.Errorf(err, L("failed to render setup script"))
	}
	return scriptBuilder.String(), nil
}
